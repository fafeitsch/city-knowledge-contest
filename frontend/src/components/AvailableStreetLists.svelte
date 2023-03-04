<script lang="ts">
  import { onMount } from "svelte";
  import { handleRPCRequest } from "../rpc";
  import store from "../store";
  import Button from "./Button.svelte";

  type StreetList = {
    FileName: string;
    name: string;
    center: { Lat: number; Lng: number };
  };

  let streetLists: StreetList[] = [];

  onMount(() => {
    handleRPCRequest<{}, StreetList[]>({
      method: "getAvailableStreetLists",
      params: {},
    }).then((data) => (streetLists = data.result));
  });

  function updateStreetList(streetList: string) {
      store.set.streetList(streetList);
  }
</script>

Welche Karte möchtest du spielen?
<div class="d-flex align-items-center gap-3">
  {#each streetLists as streetList}<Button
      on:click={() => updateStreetList(streetList.FileName)}
      title={streetList.name}
    />{/each}
</div>
