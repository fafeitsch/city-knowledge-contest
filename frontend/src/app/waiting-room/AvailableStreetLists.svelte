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

let dispatch = createEventDispatcher();
let streetLists = rpc.getStreetLists().pipe(startWith([]));
let selectedStreetList = '';

function updateStreetList(streetList: string) {
  dispatch('streetListChanged', streetList);
}
</script>

<select
  placeholder="Strassenkarte auswählen"
  class="select"
  bind:value="{selectedStreetList}"
  on:change="{() => updateStreetList(selectedStreetList)}"
>
  {#each $streetLists as streetList}
    <option class="placeholder" value="" disabled selected hidden>Straßenkarte auswählen</option>
    <option value="{streetList.FileName}">{streetList.name}</option>
  {/each}
</select>
