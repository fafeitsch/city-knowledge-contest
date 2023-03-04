<script lang="ts">
  import { onMount } from "svelte";
  import { handleRPCRequest } from "../rpc";
  import store from "../store";
  import {
    defaultRoomSeetings,
    type UpdateRoomParams,
  } from "../views/WaitingRoom.svelte";
  import Button from "./Button.svelte";

  type StreetList = {
    FileName: string;
    name: string;
    center: { Lat: number; Lng: number };
  };

  let streetLists: StreetList[] = [];
  let game = store.get.game$;

  onMount(() => {
    handleRPCRequest<{}, StreetList[]>({
      method: "getAvailableStreetLists",
      params: {},
    }).then((data) => (streetLists = data.result));
  });

  function updateStreetList(streetList: string) {
    handleRPCRequest<UpdateRoomParams, {}>({
      method: "updateRoom",
      params: {
        ...defaultRoomSeetings,
        listFileName: streetList,
        playerKey: $game.playerKey,
        roomKey: $game.roomId,
        playerSecret: $game.playerSecret,
      },
    }).then((data) => console.log(data.result));
  }
</script>

Welche Karte möchtest du spielen?
<div class="d-flex align-items-center gap-3">
  {#each streetLists as streetList}<Button
      on:click={() => updateStreetList(streetList.FileName)}
      title={streetList.name}
    />{/each}
</div>
