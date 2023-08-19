<script lang="ts">
  import { Button, Loading } from "carbon-components-svelte"
  import Checkmark from "carbon-icons-svelte/lib/Checkmark.svelte"
  import Restart from "carbon-icons-svelte/lib/Restart.svelte"

  import Title from "$lib/Title.svelte"

  import { RestoreBackup } from "../../wailsjs/go/main/App"
  import { hash } from "../router"
  import {
    type ActiveRule,
    type BackupMetadata,
    backupStore,
    configStore,
    ruleStore,
  } from "../state"
  import { formatDateTime, formatNumber } from "../utils.js"

  let game: string = ""
  let rule: ActiveRule = undefined
  let backups: BackupMetadata[] = []
  let restored: string = undefined
  let clearRestoreTimeout = undefined
  let pathSeparator = "/"

  const SHOW_SUCCESS_MS = 1_500

  function basename(path: string): string {
    return path.split(pathSeparator).pop()
  }

  async function restore(backup: BackupMetadata) {
    if (restored === backup.filename) {
      return
    }

    console.log("Want to restore", game, backup.filename)
    const result = await RestoreBackup(game, backup.filename)
    if (result) {
      restored = backup.filename
      if (clearRestoreTimeout) {
        clearTimeout(clearRestoreTimeout)
      }
      clearRestoreTimeout = setTimeout(() => {
        restored = undefined
      }, SHOW_SUCCESS_MS)
    }
  }

  $: {
    pathSeparator = $configStore.pathSeparator
    game = $hash.split("/")[1]
    rule = $ruleStore[game]
    backups = $backupStore[game]

    if (backups) {
      backups.sort((a, b) => {
        return b.backupTime.localeCompare(a.backupTime)
      })
    }

    /*
    // For testing a longer list
    backups = [
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
      {
        filename: "quicksav.json",
        source: "/home/lietu/.local/share/godot/app_userdata/Toleo/quicksave.json",
        backupTime: new Date().toISOString(),
        lastModified: new Date().toISOString(),
      },
    ]
     */
  }
</script>

{#if !rule || !backups}
  <Loading />
{:else}
  <Title>
    <div class="titlebar">
      <div class="game">
        <h2>Game:</h2>
        <h3 title={rule.name}>{rule.name}</h3>
      </div>
      <div>
        <h2>Backups:</h2>
        <h3>{formatNumber(backups.length)}</h3>
      </div>
    </div>
  </Title>

  <article>
    <section class="monitoring">
      <h2>Monitoring</h2>
      <ul class="monitors">
        {#each rule.platform.savegames as monitor}
          <!-- // @formatter:off -->
          {@const parts = monitor.split(pathSeparator)}
          <li>
            {#each parts as part, i}{part}{i === parts.length - 1 ? "" : pathSeparator}<wbr
              />{/each}
          </li>
          <!-- // @formatter:on -->
        {/each}
      </ul>
    </section>

    <section class="backups">
      <h2>Backups</h2>
      {#if backups.length === 0}
        <p>No backups yet...</p>
      {:else}
        {#each backups as backup, i}
          {@const ts = formatDateTime(backup.backupTime).split(" ")}
          {#if i > 0}
            <div class="separator" />
          {/if}
          <div class="backup">
            <div class="name" title={backup.filename}>{basename(backup.source)}</div>
            <div class="end">
              <div class="timestamp" title={backup.filename}>
                {ts[0]}<br />
                {ts[1]}
              </div>
              <div class="actions">
                <Button
                  size="small"
                  kind={restored === backup.filename ? "secondary" : "primary"}
                  icon={restored === backup.filename ? Checkmark : Restart}
                  on:click={() => restore(backup).then(() => {})}
                >
                  Restore
                </Button>
              </div>
            </div>
          </div>
        {/each}
      {/if}
    </section>
  </article>
{/if}

<style lang="scss">
  @import "../setup";

  .titlebar {
    display: flex;
    flex-direction: row;
    flex-grow: 1;
    align-items: center;
    gap: 8px;
    justify-content: space-between;
    max-width: 100%;
    width: 100%;

    h2,
    h3 {
      font-size: 16px;
      font-weight: 400;
    }

    h3 {
      color: $color-secondary-1-1;
      font-weight: 700;
      letter-spacing: 0.32px;
    }

    div {
      display: flex;
      flex-direction: row;
      gap: 16px;
      min-width: 0;

      &.game {
        flex-shrink: 10000;

        h3 {
          white-space: nowrap;
          text-overflow: ellipsis;
          overflow: hidden;
        }
      }
    }
  }

  article {
    display: flex;
    flex-direction: column;
    gap: $spacing-md;
  }

  section {
    display: flex;
    flex-direction: column;
    gap: $spacing-xs;

    h2 {
      font-size: $spacing-lg;
      font-style: normal;
      font-weight: 700;
      line-height: normal;
      letter-spacing: 0.4px;
    }

    &.backups {
      display: flex;
      flex-direction: column;
      width: 100%;

      .backup {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: $spacing-md;
        padding: $spacing-xs;
        font-size: 16px;
        justify-content: space-between;
        width: 100%;

        .name,
        .timestamp {
          color: $color-secondary-2-1;
        }

        .name {
          flex-shrink: 100000;
          white-space: nowrap;
          text-overflow: ellipsis;
          overflow: hidden;
        }

        .end {
          display: flex;
          flex-direction: row;
          gap: $spacing-md;
        }

        .fill {
          flex-grow: 1;
          min-width: 0;
        }

        .timestamp {
          text-align: right;
          font-size: 14px;
        }
      }

      .separator {
        align-self: center;
        height: 1px;
        width: 85%;
        background: #2f5aa4;
      }
    }
  }

  ul {
    text-align: left;
    margin: 0 $spacing-xs;
    color: $color-secondary-2-1;
    font-weight: 400;
    font-size: 16px;
    word-break: break-word;
    overflow-wrap: break-word;

    &.monitors {
      li:not(:last-child) {
        margin-bottom: $spacing-xs;
      }
    }
  }

  pre {
    white-space: pre-wrap;
  }
</style>
