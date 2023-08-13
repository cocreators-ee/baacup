# Baacup - Backup your savegames

Baacup monitors the savegame folders for your games for changes, and creates backups of your
savegames. Many of us have run into a situation where you've spent a good amount of time with a
game, mostly relying on quicksaves, autosaves, or a few save slots, only to find out later that
you've ended up with some kind of issue.

These could be game breaking bugs, random encounters that you can't beat, or other unexpected
events. Then you go and look at your old saves to go back to and find out that you either have no
saves old enough to avoid the issue, or the only saves you have are way too old for you to want to
go back to them.

Baacup was made to save your savegames, to avoid these kinds of situations ruining your enjoyment of
the game.

Baacup supports Windows, Linux, and MacOS, and relies on community-driven rules to know what files
to back up, and provides an interface for you to easily restore back to a previous version.

## Data files

The data files for Baacup will be stored under the appropriate `BASE_PATH` depending on the
platform:

- For \*nix: $HOME/Baacup
- For Windows: My Documents\Baacup

### Rules

Stored in `{BASE_PATH}/rules/{game}-{variant}.yaml`, e.g. `Baldurs Gate 2 - Steam.yaml`.

Contains the following:

```yaml
name: Baldur's Gate 2
platforms:
  windows:
    executable: *\\bg2.exe
    savegames:
      - %LOCALAPPDATA%\\SomePublisher\\BG2\\quicksave.sav
      - %LOCALAPPDATA%\\SomePublisher\\BG2\\autosave.sav
      - %LOCALAPPDATA%\\SomePublisher\\BG2\\saves\\*.sav
  linux:
    executable: */bg2.bin
    savegames:
      - $HOME/.local/savegames/BG2/*.sav
  macos:
    executable: */bin/bg2
    savegames:
      - $HOME/Save Games/Baldur's Gate 2/*.macsav
```

### Config

Stored in `{BASE_PATH}/config.yaml`.

Contains the following:

```yaml
disabled_rules:
  - {game}-{variant}.yaml
backups:
  keep_saves: 50
  max_mb_per_game: 500
compaction:
  compact_after_days: 365
  keep_saves: 5
rules_last_updated: 2023-04-01T11:22:33
rules_autoupdate: true
```

### Backups

For all the rules from above, the results of the backups shall be put to
`{BASE_PATH}/backups/{game}-{variant}/`. Each file backed up will be named
`{original_filename_before_ext}-{date}-{timestamp}.{ext}`, and will be accompanied by a metadata
file with the same name but extension replaced with `yaml`.

The `yaml` metadata will look like this:

```yaml
source: /full/path/to/file/source.sav
backup_time: RFC 3339 timestamp
sha256_hash: SHA256 hash of the file contents
```

## Development

Built with [Wails](https://wails.io/). You will need Go 1.20+, Node 18+, and Pnpm 8.6.0+ to work on
this code.

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite
development server that will provide very fast hot reload of your frontend changes. If you want to
develop in a browser and have access to your Go methods, there is also a dev server that runs on
http://localhost:34115. Connect to this in your browser, and you can call your Go code from
devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## License

The code is released under the GPL v3 license. Details in the [LICENSE.md](./LICENSE.md) file.

# Financial support

This project has been made possible thanks to [Cocreators](https://cocreators.ee) and
[Lietu](https://lietu.net). You can help us continue our open source work by supporting us on
[Buy me a coffee](https://www.buymeacoffee.com/cocreators).

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/cocreators)
