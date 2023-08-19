<script lang="ts">
</script>

<div class="titlebar">
  <slot />
</div>

<style lang="scss">
  @import "../setup";

  .titlebar {
    width: inherit;
    max-width: inherit;

    display: flex;
    flex-direction: row;
    align-items: center;
    background: $color-primary-4;
    margin-left: -1rem;
    padding: 16px;

    margin-bottom: $spacing-md;
  }
</style>
