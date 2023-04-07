<script lang="ts">
import { createEventDispatcher } from 'svelte';
import Button from '../../components/Button.svelte';
import rpc from '../../rpc';
import { startWith } from 'rxjs';

let dispatch = createEventDispatcher();

let streetLists = rpc.getStreetLists().pipe(startWith([]));

function updateStreetList(streetList: string) {
  dispatch('streetListChanged', streetList);
}
</script>

Welche Karte möchtest du spielen?
<div class="d-flex align-items-center gap-3">
  {#each $streetLists as streetList}
    <Button on:click="{() => updateStreetList(streetList.FileName)}" title="{streetList.name}" />
  {/each}
</div>
