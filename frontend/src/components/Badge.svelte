<script lang="ts">
  export let value: string | number;
</script>

<div class="badge">
  {value}
</div>

<style lang="scss">
  @import "../styles/_variables";

  .badge {
    padding: 8px;
    background-color: $brown-darker;
    color: $beige;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
    border-radius: 8px;
    width: 24px;
    display: flex;
    justify-content: center;
  }
</style>
