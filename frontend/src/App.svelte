<script lang="ts">
  import { combineLatest } from "rxjs";
  import Login from "./views/Login.svelte";
  import store from "./store";
  import WaitingRoom from "./views/WaitingRoom.svelte";

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
  {#if gameState === "inital"}
    <div class="d-flex flex-column gap-4 align-items-center">
      <h1 class="old-font">City Knowledge Contest</h1>
      <p class="mb-5">Wer findet die gesuchten Orte am schnellsten?</p>
      <Login />
    </div>
  {:else}
    <WaitingRoom />
  {/if}
</main>

<footer>
  <div class="p-3">Fancy Footer | License</div>
</footer>
