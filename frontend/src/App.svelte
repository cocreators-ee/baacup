<script lang="ts">
  import { Column, Grid, Loading, Row } from "carbon-components-svelte"

  import { hash } from "./router"
  import Game from "./routes/Game.svelte"
  import Home from "./routes/Home.svelte"
  import { activeRuleStore, configStore, ruleStore } from "./state"
  import { formatDateTime, formatNumber } from "./utils.js"

  const routes = {
    "": Home,
    Game: Game,
  }

  $: view = routes[$hash.split("/")[0]]
</script>

<main>
  {#if $configStore === undefined}
    <Loading />
  {:else}
    <Grid>
      <Row noGutter>
        <Column md={3} noGutter>
          <div class="sidebar">
            <div class="logo">
              <a href="#"> Baacup </a>
            </div>
            <div>
              <h2>Now playing</h2>
              <div class="now-playing">
                {#if Object.keys($activeRuleStore).length > 0}
                  {#each Object.keys($activeRuleStore) as rule}
                    {@const activeRule = $activeRuleStore[rule]}
                    <a href={`#Game/${rule}`}>{activeRule.name}</a>
                  {/each}
                {:else}
                  <p>No games detected...</p>
                {/if}
              </div>
            </div>
            <div class="grow">&nbsp;</div>
            <div class="footer">
              <div>
                <p>Rules last updated</p>
                <p>{formatDateTime($configStore.rulesLastUpdated)}</p>
                <p>{formatNumber(Object.keys($ruleStore).length)} rules loaded</p>
              </div>
              <div>
                <em>Baacup: Backup your savegames</em>
              </div>
            </div>
          </div>
        </Column>
        <Column md={5} noGutter>
          <div class="main">
            <svelte:component this={view} />
          </div>
        </Column>
      </Row>
    </Grid>
  {/if}
</main>

<style lang="scss">
  @import "setup";

  :global(main > .bx--grid) {
    padding: 0;
  }

  :global(main > .bx--grid > .bx--row) {
    margin: 0;
  }

  .main {
    display: flex;
    flex-direction: column;
  }

  .sidebar {
    position: fixed;
    width: inherit;
    max-width: inherit;

    display: flex;
    padding: 32px 16px 32px 32px;
    flex-direction: column;
    align-items: flex-start;
    gap: 32px;
    flex-shrink: 0;
    margin: 0 0 0 -1rem;

    //background: $color-secondary-1-4;
    height: 100vh;

    background: conic-gradient(from 270deg at 100% -0%, $color-secondary-1-4 0deg, #3f0c23 310deg);

    .logo {
      a {
        color: $color-complement-1;
        text-decoration: none;
      }
      text-align: center;
      text-shadow: 0px 0px 4px lighten($color-complement-1, 30%);
      font-size: 48px;
      font-style: normal;
      font-weight: 700;
      line-height: normal;
      letter-spacing: 0.96px;
      width: 100%;
    }

    h2 {
      font-size: $spacing-2xl;
      font-style: normal;
      font-weight: 700;
      line-height: normal;
      letter-spacing: 0.64px;
    }

    > div {
      display: flex;
      flex-direction: column;
      gap: $spacing-md;

      &.grow {
        flex-grow: 1;
      }

      &.footer {
        width: 100%;
        text-align: center;
        letter-spacing: 0.32px;

        p,
        em {
          font-size: 14px;
        }

        em {
          text-align: center;
          font-style: italic;
          font-weight: 700;
          line-height: normal;
          letter-spacing: 0.32px;
        }
      }
    }

    .now-playing {
      padding-left: $spacing-lg;
      gap: $spacing-md;

      p {
        font-style: italic;
      }
    }
  }
</style>
