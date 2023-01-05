<script lang="ts">
  import { combineLatest } from "rxjs";
  import Counter from "./lib/Counter.svelte";
  import Login from "./Login.svelte";
  import store from "./store";
  import WaitingRoom from "./WaitingRoom.svelte";

  let gameState = "";
  store.get.gameState$.subscribe((state) => {
    gameState = state;
  });

  function initWebSocket() {
    combineLatest([
      store.get.roomId$,
      store.get.playerKey$,
      store.get.gameState$,
    ]).subscribe(([roomId, playerKey, gameState]) => {
      if (gameState !== "waiting") {
        return;
      }
      const websocket = new WebSocket(
        "ws://127.0.0.1:23123/ws/" + roomId + "/" + playerKey
      );
      websocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.topic === "successfullyJoined") {
          store.set.players(data.payload.players);
        } else if (data.topic === "playerJoined") {
          store.set.addPlayer(data.payload.name);
        }
      };
    });
  }
  initWebSocket();
</script>

<main>
  <h1>
    {#if gameState === "inital"}
      <Login />
    {:else}
      <WaitingRoom />
    {/if}
  </h1>
</main>
