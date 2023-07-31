<style lang="scss">
@import '../../styles/variables';

.select {
  padding: 16px 8px;
  width: 100%;
  border: 3px solid $gray-main;
  background-color: white;
  border-radius: 8px;
  font-size: large;
  outline: none;
  transition: 500ms;
}

.placeholder {
  color: $gray-dark;
}
</style>

<script lang="ts">
import { createEventDispatcher } from 'svelte';
import rpc from '../../rpc';
import { startWith } from 'rxjs';

export let selectedStreetList = '';
let dispatch = createEventDispatcher();
let streetLists = rpc.getStreetLists().pipe(startWith([]));

function updateStreetList(event: any) {
  dispatch('streetListChanged', event.target.value);
}
</script>

<select class="select" value="{selectedStreetList}" on:change="{updateStreetList}" data-testid="select-streetlist">
  {#each $streetLists as streetList}
    <option class="placeholder" value="" disabled selected hidden data-testid="streetlist-option-initial-selection"
      >Straßenkarte auswählen</option
    >
    <option value="{streetList.FileName}" data-testid="streetlist-option">{streetList.name}</option>
  {/each}
</select>
