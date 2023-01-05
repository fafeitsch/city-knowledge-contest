<script lang="ts">
  import store from "./store";

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

  function handleUsernameChange(event: Event) {
    const target = event.target as HTMLInputElement;
    username = target.value;
  }

  function handleRoomIdChange(event: Event) {
    const target = event.target as HTMLInputElement;
    roomId = target.value;
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

{#if username.length === 0}
  <div class="d-flex flex-column gap-3">
    <input
      on:change={handleUsernameChange}
      placeholder="Gib deinen Spielername ein"
    />
    <!-- <button>Neue Karte erstellen</button> -->
  </div>
{:else}
  <div class="d-flex gap-3">
    <div class="d-flex flex-column gap-2">
      <input on:change={handleRoomIdChange} placeholder="Karten-ID eingeben" />
      <button on:click={joinRoom}>Karte beitreten</button>
    </div>
    <button on:click={createRoom}>Erstelle eine neue Karte</button>
  </div>
{/if}
