<script lang="ts">
  import Title from "$lib/Title.svelte"

  import {
    backupStore,
    configStore,
    errorStore,
    eventStore,
    monistorStore,
    ruleStore,
  } from "../state"
</script>

<Title>
  <div class="titlebar">
    <div>
      <h2>Baacup - Backup your savegames</h2>
    </div>
  </div>
</Title>

<article>
  <section>
    <h2>This page is intentionally pretty ugly at this stage.</h2>
    <p>Just click on the game name on the left when one shows up.</p>
    <p>You can also edit the rule files, and then go File -> Reload rules or press Ctrl+R.</p>
  </section>

  <section>
    <h2>Supported games</h2>
    {#if $ruleStore}
      <ul>
        {#each Object.keys($ruleStore) as key}
          {@const rule = $ruleStore[key]}
          <li><a href={`#Game/${key}`}>{rule.name}</a></li>
        {/each}
      </ul>
    {:else}
      <p>No games supported. Please check configuration.</p>
    {/if}
  </section>

  <section>
    <h2>Events</h2>
    <ul>
      {#each $eventStore as event}
        <li>{event}</li>
      {/each}
    </ul>
    {#if $errorStore.length > 0}
      <h2>Errors</h2>
      <ul>
        {#each $errorStore as error}
          <li>{error}</li>
        {/each}
      </ul>
    {:else}
      <h2>No errors</h2>
    {/if}
  </section>

  <section>
    <h2>Rules</h2>
    <pre>{JSON.stringify($ruleStore, null, 2)}</pre>
  </section>

  <section>
    <h2>Config</h2>
    <pre>{JSON.stringify($configStore, null, 2)}</pre>
  </section>

  <section>
    <h2>Monitors</h2>
    <pre>{JSON.stringify($monistorStore, null, 2)}</pre>
  </section>

  <section>
    <h2>Backups</h2>
    <pre>{JSON.stringify($backupStore, null, 2)}</pre>
  </section>
</article>

<style lang="scss">
  @import "../setup";

  article {
    display: flex;
    flex-direction: column;
    gap: $spacing-xl;

    section {
      h2 {
        margin-bottom: $spacing-md;
      }
    }
  }

  ul {
    text-align: left;
  }

  pre {
    white-space: pre-wrap;
  }

  .titlebar {
    display: flex;
    flex-direction: row;
    flex-grow: 1;
    align-items: center;
    gap: 8px;
    justify-content: space-between;

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
    }
  }
</style>
