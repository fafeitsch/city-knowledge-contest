<style lang="scss">
@import '../styles/variables';

.legal {
  z-index: 1000;
  position: absolute;
  top: auto;
  left: auto;
  width: 500px;
  overflow: auto;
  background-color: white;
  border-radius: 16px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  border: 3px solid $blue-dark;
  line-height: 1.5;

  :global(h3) {
    color: $blue-dark;
    padding: 0;
    margin: 0;
    font-size: 1.3em;
  }
}
</style>

<script lang="ts">
import Button from './Button.svelte';
import { createEventDispatcher } from 'svelte';

export let info = '';
let dispatch = createEventDispatcher();

function closeLegal() {
  dispatch('closeLegal', {});
}
</script>

<div class="legal">
  {@html info}
  <Button title="Schließen" on:click="{closeLegal}" disabled="{false}" />
</div>
