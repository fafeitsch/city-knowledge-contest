<script lang="ts">
  import Button from "../components/Button.svelte";
  import Input from "../components/Input.svelte";
  import store from "../store";

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
  let enterUsername = false;
  let roomId = "";

  function handleUsernameChange(event: Event) {
    const target = event.target as HTMLInputElement;
    username = target.value;
  }

  function handleRoomIdChange(event: Event) {
    const target = event.target as HTMLInputElement;
    roomId = target.value;
  }

  function handleOnClick(event: Event) {
    if (username.length > 0) {
      enterUsername = true;
    }
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
        store.set.gameState("waiting");
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
        store.set.gameState("waiting");
      });
  }
</script>

{#if !enterUsername}
  <div class="d-flex flex-column gap-3 align-items-center">
    <Input
      on:change={handleUsernameChange}
      placeholder="Gib deinen Spielername ein"
    />
    <Button on:click={handleOnClick} title="Los geht's" />
  </div>
{:else}
  <div class="d-flex gap-5 align-items-center">
    <div class="d-flex flex-column gap-2 align-items-center">
      <Input on:change={handleRoomIdChange} placeholder="Karten-ID eingeben" />
      <Button on:click={joinRoom} title="Karte beitreten" />
    </div>
    <div>– oder –</div>
    <Button title="Erstelle eine neue Karte" on:click={createRoom} />
  </div>
{/if}
