<script lang="ts">
  import { onMount } from "svelte";
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import store, { GameState } from "../store";

  type RoomResult = {
    roomKey: string;
    playerKey: string;
  };

  type Response<T> = {
    jsonrpc: string;
    result: T;
    id: null;
  };

  let username = "";
  let roomId = "";
  let gameState: GameState;

  onMount(() => {
    store.get.gameState$.subscribe((value) => {
      gameState = value;
    });
  });

  function handleUsernameChange(event: Event) {
    const target = event.target as HTMLInputElement;
    username = target.value;
  }

  function handleRoomIdChange(event: Event) {
    const target = event.target as HTMLInputElement;
    roomId = target.value;
  }

  function handleOnClick() {
    store.set.gameState(GameState.SetupMap);
  }

  function createRoom() {
    fetch("http://localhost:23123/rpc", {
      method: "POST",
      body: JSON.stringify({
        method: "createRoom",
        params: {
          name: username,
        },
      }),
    })
      .then((response) => response.json())
      .then((data: Response<RoomResult>) => {
        store.set.roomId(data.result.roomKey);
        store.set.playerKey(data.result.playerKey);
        store.set.gameState(GameState.Waiting);
      });
  }

  function joinRoom() {
    fetch("http://localhost:23123/rpc", {
      method: "POST",
      body: JSON.stringify({
        method: "joinRoom",
        params: {
          name: username,
          roomKey: roomId,
        },
      }),
    })
      .then((response) => response.json())
      .then((data: Response<RoomResult>) => {
        store.set.roomId(roomId);
        store.set.playerKey(data.result.playerKey);
        store.set.gameState(GameState.Waiting);
      });
  }
</script>

{#if gameState === GameState.SetupUsername}
  <div class="d-flex flex-column gap-3 align-items-center">
    <Input
      on:change={handleUsernameChange}
      placeholder="Gib deinen Spielername ein"
    />
    <Button on:click={handleOnClick} title="Los geht's" />
  </div>
{:else if gameState === GameState.SetupMap}
  <div class="d-flex gap-5 align-items-center">
    <div class="d-flex flex-column gap-2 align-items-center">
      <Input on:change={handleRoomIdChange} placeholder="Karten-ID eingeben" />
      <Button on:click={joinRoom} title="Karte beitreten" />
    </div>
    <div>– oder –</div>
    <Button title="Eine neue Karte erstellen" on:click={createRoom} />
  </div>
{/if}
