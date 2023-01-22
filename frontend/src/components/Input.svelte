<script lang="ts">
  export let placeholder: string;
  export let value: string = "";
</script>

<input {placeholder} on:change {value} />

<style lang="scss">
  @import "../styles/variables";

  input {
    padding: 16px 8px;
    width: 100%;
    background-color: $beige;
    border: 3px solid $old-map-lighter;
    font-size: large;
    outline: none;
    transition: 500ms;

    &::placeholder {
      color: $brown-darker;
    }

    &:focus {
      border: 3px solid $old-map-darker;
    }
  }
</style>
