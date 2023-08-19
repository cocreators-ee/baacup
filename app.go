package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/gobwas/glob"
	"github.com/goccy/go-yaml"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const logMaxLength = 20

func newConfig() *Config {
	return &Config{
		DisabledRules: []string{},
		Backups: &BackupConfig{
			KeepSaves:    50,
			MaxMBPerGame: 500,
		},
		Compaction: &CompactionConfig{
			KeepSaves:        5,
			CompactAfterDays: 180,
		},
		RulesLastUpdated: time.Time{},
		RulesAutoUpdate:  true,
	}
}

// newApp creates a new App application struct
func newApp() *App {
	exit := make(chan bool)
	a := &App{
		Config:         newConfig(),
		Rules:          map[string]ActiveRule{},
		ActiveMonitors: []Monitor{},
		ProcessList:    []string{},
		BasePath:       getBasePath(),
		Errors:         []string{},
		Events:         []string{},
		exit:           exit,
	}

	return a
}

func (a *App) ensurePath(createPath string) {
	err := os.MkdirAll(createPath, 0o700)
	if err != nil {
		if err != os.ErrExist {
			a.ReportError(fmt.Errorf("could not create %s", createPath))
		}
	}
}

func (a *App) getConfigPath() string {
	return filepath.Join(a.BasePath, "config.yaml")
}

func (a *App) getBackupsPath() string {
	return filepath.Join(a.BasePath, "backups")
}

func (a *App) getRulesPath() string {
	return filepath.Join(a.BasePath, "rules")
}

// LoadConfig loads or reload the application configuration
func (a *App) LoadConfig() {
	// Configure default config
	a.Config = newConfig()
	configPath := a.getConfigPath()
	contents, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No configuration file exists yet, make sure we save the defaults
			a.SaveConfig()
			return
		}

		panic(err)
	}

	err = yaml.Unmarshal(contents, a.Config)
	if err != nil {
		a.ReportError(err)
		a.ReportError(fmt.Errorf("error parsing config from %s", configPath))
	}

	a.Config.PathSeparator = a.GetPathSeparator()

	wailsRuntime.EventsEmit(a.ctx, "configUpdated", a.Config)
}

// SaveConfig saves the application configuration
func (a *App) SaveConfig() {
	configPath := a.getConfigPath()

	// Dump config into YAML
	data, err := yaml.Marshal(a.Config)
	if err != nil {
		a.ReportError(err)
		a.ReportError(fmt.Errorf("error writing config to %s", configPath))
		return
	}

	err = os.WriteFile(configPath, data, 0o600)
	if err != nil {
		a.ReportError(err)
		a.ReportError(fmt.Errorf("error writing config to %s", configPath))
		return
	}
}

func (a *App) preprocessRules(rules map[string]ActiveRule) map[string]ActiveRule {
	newRules := map[string]ActiveRule{}
	for key, rule := range rules {
		// Expand ${HOME} etc. in the "Savegames" list
		savegames := []string{}
		for _, v := range rule.Platform.Savegames {
			savegames = append(savegames, os.ExpandEnv(v))
		}

		executableGlob, err := glob.Compile(rule.Platform.Executable)
		if err != nil {
			a.ReportError(fmt.Errorf("failed to parse pattern %s from %s.yaml", rule.Platform.Executable, rule.RuleFilename))
			continue
		}

		newRules[key] = ActiveRule{
			RuleFilename: rule.RuleFilename,
			Name:         rule.Name,
			Platform: RulePlatform{
				Executable: rule.Platform.Executable,
				Savegames:  savegames,
			},
			executableGlob: executableGlob,
		}
	}
	return newRules
}

// LoadRules loads our rules/*.yaml configuration files
func (a *App) LoadRules() {
	activeRules := a.loadRules()
	if len(activeRules) == 0 {
		// Seems like we need to update rules from the internet and try again
		if a.Config.RulesAutoUpdate {
			a.DownloadRules()
			activeRules = a.loadRules()
		}
	}

	a.Rules = a.preprocessRules(activeRules)
	wailsRuntime.EventsEmit(a.ctx, "rulesUpdated", a.Rules)
	a.AddEvent(fmt.Sprintf("Loaded %d rules", len(a.Rules)))
}

func (a *App) loadRules() map[string]ActiveRule {
	platform := getPlatform()

	rulesPath := a.getRulesPath()
	rulesFiles, err := filepath.Glob(filepath.Join(rulesPath, "*.yaml"))
	if err != nil {
		panic(err)
	}

	activeRules := map[string]ActiveRule{}
	for _, rulesFile := range rulesFiles {
		rule := &Rule{}

		name := strings.TrimSuffix(filepath.Base(rulesFile), ".yaml")

		contents, err := os.ReadFile(rulesFile)
		if err != nil {
			a.ReportError(err)
			a.ReportError(fmt.Errorf("error parsing rules from %s", rulesFile))
			continue
		}

		err = yaml.Unmarshal(contents, rule)
		if err != nil {
			a.ReportError(err)
			a.ReportError(fmt.Errorf("error parsing rules from %s", rulesFile))
			continue
		}

		rulePlatform, ok := rule.Platforms[platform]
		if !ok {
			a.ReportError(fmt.Errorf("%s does not support %s, skipping", rulesFile, platform))
			continue
		}

		activeRules[name] = ActiveRule{
			RuleFilename: name,
			Name:         rule.Name,
			Platform:     rulePlatform,
		}
	}

	return activeRules
}

// DownloadRules downloads the community-driven rules files off GitHub
func (a *App) DownloadRules() {
	// TODO: Download off github cocreators-ee/baacup-rules
	a.AddEvent("Should download rules off of github.com/cocreators-ee/baacup-rules")
}

func (a *App) runMonitor() {
	pollMonitors := time.NewTicker(time.Second)
	pollRules := time.NewTicker(time.Second * 15)

	for {
		select {
		case <-a.exit:
			pollMonitors.Stop()
			pollRules.Stop()
			return

		case <-pollMonitors.C:
			a.checkMonitors()

		case <-pollRules.C:
			a.CheckRules()
		}
	}
}

func (a *App) checkMonitors() {
	for _, monitor := range a.ActiveMonitors {
		newFiles := a.findNewFiles(monitor.Path, a.Backups[monitor.RuleFilename])
		for _, newFile := range newFiles {
			a.backupFile(monitor.RuleFilename, newFile)
		}
	}
}

func (a *App) findNewFiles(filePath string, backups []BackupMetadata) []string {
	var newFiles []string

	files, err := filepath.Glob(filePath)
	if err != nil {
		a.ReportError(err)
		return newFiles
	}

	for _, f := range files {
		needsBackup := true
		metas := listMetadataBySource(backups, f)
		for _, meta := range metas {
			if !a.fileModifiedAfter(f, meta.LastModified) {
				needsBackup = false
			}
		}

		if needsBackup {
			a.AddEvent(fmt.Sprintf("%s needs backup", f))
			newFiles = append(newFiles, f)
		}
	}

	return newFiles
}

func (a *App) fileModifiedAfter(file string, lastModified time.Time) bool {
	s, err := os.Lstat(file)
	if err != nil {
		a.ReportError(err)
		return false
	}
	return s.ModTime().After(lastModified)
}

func listMetadataBySource(backups []BackupMetadata, source string) []BackupMetadata {
	matches := []BackupMetadata{}
	for _, bm := range backups {
		if bm.Source == source {
			matches = append(matches, bm)
		}
	}

	return matches
}

func (a *App) backupFile(ruleFilename string, sourcePath string) {
	rule := a.Rules[ruleFilename]

	stat, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		if os.IsPermission(err) {
			a.ReportError(err)
			return
		}
	}

	meta := BackupMetadata{
		Filename:     "",
		Source:       sourcePath,
		BackupTime:   time.Now(),
		LastModified: stat.ModTime(),
	}

	meta, err = a.makeBackup(ruleFilename, meta)
	if err != nil {
		a.ReportError(err)
		return
	}

	// Now that we've successfully backed it up, load the metadata in memory
	a.Backups[ruleFilename] = append(a.Backups[ruleFilename], meta)

	// Then limit it to max length
	if len(a.Backups[ruleFilename]) > a.Config.Backups.KeepSaves {
		sort.SliceStable(a.Backups[ruleFilename], func(li, ri int) bool {
			l := a.Backups[ruleFilename][li]
			r := a.Backups[ruleFilename][ri]
			return l.BackupTime.After(r.BackupTime)
		})

		extra := a.Backups[ruleFilename][a.Config.Backups.KeepSaves:]
		for _, meta := range extra {
			a.DeleteBackup(ruleFilename, meta.Filename)
		}
	}

	wailsRuntime.EventsEmit(a.ctx, "backupsUpdated", a.Backups)
	a.AddEvent(fmt.Sprintf("Backed up %s savegame %s", rule.Name, filepath.Base(sourcePath)))
}

// DeleteBackup deletes a specific backup
func (a *App) DeleteBackup(ruleFilename string, filename string) {
	// Figure out filenames
	backupPath := filepath.Join(a.getBackupsPath(), ruleFilename)
	baseNoExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	metaFilename := fmt.Sprintf("%s.baacup.yaml", baseNoExt)

	filePath := filepath.Join(backupPath, filename)
	metaPath := filepath.Join(backupPath, metaFilename)

	a.tryDeleteFile(filePath)
	a.tryDeleteFile(metaPath)

	backups := []BackupMetadata{}
	for _, meta := range a.Backups[ruleFilename] {
		if meta.Filename != filename {
			backups = append(backups, meta)
		}
	}
	a.Backups[ruleFilename] = backups

	rule := a.Rules[ruleFilename]
	a.AddEvent(fmt.Sprintf("Removed old backup for %s, %s", rule.Name, filename))
}

func (a *App) tryDeleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		a.ReportError(err)
	}
}

func (a *App) makeBackup(ruleFilename string, meta BackupMetadata) (BackupMetadata, error) {
	// Figure out filenames
	backupPath := filepath.Join(a.getBackupsPath(), ruleFilename)
	ext := filepath.Ext(meta.Source)
	baseNoExt := strings.TrimSuffix(filepath.Base(meta.Source), ext)
	timestamp := meta.BackupTime.Format("2006-01-02T150405.000")

	backupFilename := fmt.Sprintf("%s-%s%s", baseNoExt, timestamp, ext)
	metaFilename := fmt.Sprintf("%s-%s.baacup.yaml", baseNoExt, timestamp)

	// Make sure we update the metadata with this new filename
	meta.Filename = backupFilename

	// Ensure the destination path exists
	err := os.MkdirAll(backupPath, 0o700)
	if err != nil {
		return meta, err
	}

	// Write metadata file
	metaFile := filepath.Join(backupPath, metaFilename)
	data, err := yaml.Marshal(meta)
	if err != nil {
		a.ReportError(fmt.Errorf("error writing backup metadata to %s", metaFile))
		return meta, err
	}

	err = os.WriteFile(metaFile, data, 0o600)
	if err != nil {
		a.ReportError(fmt.Errorf("error writing backup metadata to %s", metaFile))
		return meta, err
	}

	// Make the copy of the file
	backupDst := filepath.Join(backupPath, backupFilename)
	err = copyFile(meta.Source, backupDst)

	return meta, err
}

func copyFile(src, dst string) error {
	// Verify source file
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("cannot copy %s to %s, %s is not a regular file", src, dst, src)
	}

	// Open the source
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		err := source.Close()
		if err != nil {
			fmt.Printf("Error closing %s: %s", src, err)
		}
	}()

	// Handle the destination
	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		err := destination.Close()
		if err != nil {
			fmt.Printf("Error closing %s: %s", dst, err)
		}
	}()

	// Do the copy
	_, err = io.Copy(destination, source)

	return err
}

// RestoreBackup restores a selected backup
func (a *App) RestoreBackup(ruleFilename string, filename string) bool {
	metadata := a.findBackupMetadataForRestore(ruleFilename, filename)
	if metadata.Source == "" {
		// For some reason couldn't find the metadata
		return false
	}

	dst := metadata.Source
	src := filepath.Join(a.getBackupsPath(), ruleFilename, filename)

	err := copyFile(src, dst)
	if err != nil {
		return false
	}

	err = os.Chtimes(dst, metadata.LastModified, metadata.LastModified)
	if err != nil {
		a.ReportError(err)
		return false
	}

	rule := a.Rules[ruleFilename]
	a.AddEvent(fmt.Sprintf("Restored %s %s to %s", rule.Name, filename, dst))

	return true
}

func existsInList(items []string, item string) bool {
	for _, val := range items {
		if val == item {
			return true
		}
	}
	return false
}

func (a *App) pollProcessList() {
	procs, err := process.Processes()
	if err != nil {
		a.ReportError(err)
		return
	}

	processList := []string{}
	for _, proc := range procs {
		exe, err := proc.Exe()
		if err != nil {
			if !os.IsPermission(err) && !os.IsNotExist(err) {
				a.ReportError(err)
			}
			continue
		}

		if !existsInList(processList, exe) {
			processList = append(processList, exe)
			// wailsRuntime.LogPrint(a.ctx, exe)
		}
	}

	a.ProcessList = processList
}

// CheckRules checks which of the rules are currently matching running processes and we should monitor for
func (a *App) CheckRules() {
	a.pollProcessList()

	// Add all savegame paths for all running games to monitoring
	var monitors []Monitor
	backups := map[string][]BackupMetadata{}
	rules := map[string]ActiveRule{}

	for key, rule := range a.Rules {
		backups[rule.RuleFilename] = a.findBackupMetadata(rule.RuleFilename)
		if a.IsRunning(rule.executableGlob) {
			var ruleMonitors []Monitor
			for _, savePath := range rule.Platform.Savegames {
				ruleMonitors = append(ruleMonitors, Monitor{
					Path:         savePath,
					RuleFilename: rule.RuleFilename,
				})
			}

			monitors = append(monitors, ruleMonitors...)
			rules[key] = rule

			a.AddEvent(fmt.Sprintf("Detected %s as running", rule.Name))
		}
	}

	// TODO: Check for things no longer running

	a.Backups = backups
	a.ActiveMonitors = monitors
	a.ActiveRules = rules

	wailsRuntime.EventsEmit(a.ctx, "backupsUpdated", a.Backups)
	wailsRuntime.EventsEmit(a.ctx, "activeMonitorsUpdated", a.ActiveMonitors)
	wailsRuntime.EventsEmit(a.ctx, "activeRulesUpdated", a.ActiveRules)
}

func (a *App) findBackupMetadata(ruleFilename string) []BackupMetadata {
	backupPath := filepath.Join(a.getBackupsPath(), ruleFilename)
	backups := []BackupMetadata{}

	// Read directory contents
	files, err := os.ReadDir(backupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return backups
		}
	}

	for _, file := range files {
		if file.IsDir() {
			// Some idiot put extra folders here
			continue
		}

		filename := file.Name()
		filePath := filepath.Join(backupPath, filename)

		// We're only interested in metadata
		if !strings.HasSuffix(filePath, ".baacup.yaml") {
			continue
		}

		contents, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				// Someone just deleted this as we were reading it, no big deal
				continue
			} else {
				a.ReportError(err)
				continue
			}
		}

		meta := BackupMetadata{}
		err = yaml.Unmarshal(contents, &meta)
		if err != nil {
			a.ReportError(err)
			a.ReportError(fmt.Errorf("error parsing metadata from %s", filePath))
		}

		base := strings.TrimSuffix(filename, ".baacup.yaml")
		ext := filepath.Ext(meta.Source)
		backupFilename := fmt.Sprintf("%s%s", base, ext)
		meta.Filename = backupFilename

		backups = append(backups, meta)
	}

	return backups
}

func (a *App) findBackupMetadataForRestore(ruleFilename string, filename string) BackupMetadata {
	for _, meta := range a.Backups[ruleFilename] {
		if meta.Filename == filename {
			return meta
		}
	}

	a.ReportError(fmt.Errorf("tried to restore %s but couldn't find its metadata", filename))

	return BackupMetadata{}
}

// IsRunning checks if this process is currently running
func (a *App) IsRunning(exePattern glob.Glob) bool {
	for _, exeFullPath := range a.ProcessList {
		match := exePattern.Match(exeFullPath)

		if match {
			return true
		}
	}
	return false
}

// AddEvent adds an event to the event log
func (a *App) AddEvent(msg string) {
	events := append([]string{msg}, a.Events...)

	last := len(events)
	if last >= logMaxLength {
		last = logMaxLength
	}

	a.Events = events[:last]
	wailsRuntime.EventsEmit(a.ctx, "eventsUpdated", a.Events)
}

// ReportError reports an error to the UI
func (a *App) ReportError(err error) {
	errors := append([]string{err.Error()}, a.Errors...)

	last := len(errors)
	if last >= logMaxLength {
		last = logMaxLength
	}

	a.Errors = errors[:last]
	wailsRuntime.EventsEmit(a.ctx, "errorsUpdated", a.Errors)
}

// GetRules returns the rules currently loaded
func (a *App) GetRules() map[string]ActiveRule {
	return a.Rules
}

// GetConfig returns the application configuration
func (a *App) GetConfig() Config {
	return *a.Config
}

// GetActiveMonitors returns a list of things we're monitoring
func (a *App) GetActiveMonitors() []Monitor {
	return a.ActiveMonitors
}

// GetActiveBackups returns a list of backups of the currently active games
func (a *App) GetActiveBackups() map[string][]BackupMetadata {
	return a.Backups
}

// GetActiveRules returns a list of the rules that are currently active
func (a *App) GetActiveRules() map[string]ActiveRule {
	return a.ActiveRules
}

// GetErrors returns the list of latest errors
func (a *App) GetErrors() []string {
	return a.Errors
}

// GetEvents returns a list of the latest events
func (a *App) GetEvents() []string {
	return a.Events
}

// GetPathSeparator returns current path separator on the host system
func (a *App) GetPathSeparator() string {
	return string(os.PathSeparator)
}

// Wails events etc.
func (a *App) menu() *menu.Menu {
	appMenu := menu.NewMenu()

	fileMenu := appMenu.AddSubmenu("File")

	fileMenu.AddText("Reload rules", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {
		a.LoadRules()
		a.CheckRules()
	})

	fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		wailsRuntime.Quit(a.ctx)
	})

	if runtime.GOOS == "darwin" {
		// on macOS platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcuts
		appMenu.Append(menu.EditMenu())
	}

	return appMenu
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Ensure our paths exist
	a.ensurePath(a.getBackupsPath())
	a.ensurePath(a.getRulesPath())

	// Try to load config and rules
	a.LoadConfig()
	a.LoadRules()
	a.CheckRules()

	// Start the monitor
	go a.runMonitor()
}

func (a *App) shutdown(ctx context.Context) {
	a.exit <- true
}
