<script lang="ts">
import { startWith } from 'rxjs';
import { createEventDispatcher } from 'svelte';
import Button from '../../components/Button.svelte';
import { doRpc } from '../../rpc';

type StreetList = {
  FileName: string;
  name: string;
  center: { Lat: number; Lng: number };
};

let dispatch = createEventDispatcher();

let streetLists = doRpc<StreetList[]>('getAvailableStreetLists', {}).pipe(startWith([]));

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
