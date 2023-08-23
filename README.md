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

![Screenshot of Baacup](./screenshot.png?raw=true "Screenshot of Baacup")

## Status

The software is currently in an early stage.

- It looks ugly
- The core logic seems to kinda work
- It can make backups based on rules
- It can maybe restore backups without data loss
- It can read the config
- It supports the `backups.keep_saves` and `backups.max_mb_per_game` options
- It does not support any of the other configuration options
- There is no way to configure the application from the GUI
- There has been very little testing in general
- Data loss is possible, though we try to avoid it
- You have to manually download the rules from
  [cocreators-ee/baacup-rules](https://github.com/cocreators-ee/baacup-rules) or create them
  yourself

## Data files

The data files for Baacup will be stored under the appropriate `BASE_PATH` depending on the
platform:

- For \*nix: `$HOME/Baacup`
- For Windows: `My Documents\Baacup`

### Rules

Stored in `{BASE_PATH}/rules/{game}-{variant}.yaml`, e.g. `Baldurs Gate 2 - Steam.yaml`.

Contains the following:

```yaml
name: Baldur's Gate 2
issues: Optional explanation of any issues with these rules.
platforms:
  windows:
    executable: *\\bg2.exe
    savegames:
      - ${LOCALAPPDATA}\\SomePublisher\\BG2\\quicksave.sav
      - ${LOCALAPPDATA}\\SomePublisher\\BG2\\autosave.sav
      - ${LOCALAPPDATA}\\SomePublisher\\BG2\\saves\\*.sav
  linux:
    executable: */bg2.bin
    savegames:
      - ${HOME}/.local/savegames/BG2/*.sav
  macos:
    executable: */bin/bg2
    savegames:
      - ${HOME}/Save Games/Baldur's Gate 2/*.macsav
```

You can create these files manually if you want, but we'd prefer you then contribute them to
[cocreators-ee/baacup-rules](https://github.com/cocreators-ee/baacup-rules) for the rest of the
community to benefit from them as well.

If you want to manually get the community-published rules, you can
[download them from cocreators-ee/baacup-rules](https://github.com/cocreators-ee/baacup-rules/archive/refs/heads/main.zip)
and then copy the contents of the `rules` folder in the archive to the "rules" path under the data
file directory as explained above.

### Config

Stored in `{BASE_PATH}/config.yaml`.

Contains the following:

```yaml
disabled_rules:
  - "{game}-{variant}.yaml"
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
last_modified: RFC 3339 timestamp
```

## Development

Built with [Wails](https://wails.io/) and [Svelte](https://svelte.dev). You will need the following
installed:

- [Wails](https://wails.io/docs/gettingstarted/installation)
- [Go 1.20+](https://go.dev/dl/)
- [Node 18+ (likely LTS)](https://nodejs.org/en)
- [Pnpm 8.6.0+](https://pnpm.io/installation)
- [go-pre-commit](https://github.com/lietu/go-pre-commit#using-the-hooks)
- [pre-commit](https://pre-commit.com/#install)

## Design

Before implementation we'd like to have a design. Current draft is on
[Figma](https://www.figma.com/file/7UrzEb3GO1o4jJ7i1WauEO/Baacup?type=design&node-id=11%3A93).

## Live Development

To run in live development mode, run `wails dev -loglevel Info` in the project directory. This will
run a Vite development server that will provide very fast hot reload of your frontend changes. If
you want to develop in a browser and have access to your Go methods, there is also a dev server that
runs on http://localhost:34115. Connect to this in your browser, and you can call your Go code from
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
