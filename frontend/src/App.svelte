<script lang="ts">
  import Login from "./views/Login.svelte";
  import store, { GameState } from "./store";
  import { combineLatest } from "rxjs";
  import Game from "./views/Game.svelte";
  import { environment } from "./environment";

  let gameState: GameState;
  store.get.gameState$.subscribe((state) => {
    gameState = state;
  });

  function initWebSocket() {
    combineLatest([store.get.game$, store.get.gameState$]).subscribe(
      ([game, gameState]) => {
        if (gameState !== GameState.Waiting) {
          return;
        }
        const websocket = new WebSocket(
          environment[import.meta.env.MODE].wsUrl +
            game.roomId +
            "/" +
            game.playerKey
        );
        websocket.onmessage = (event) => {
          const data = JSON.parse(event.data);
          if (data.topic === "successfullyJoined") {
            store.set.players(data.payload.players);
          } else if (data.topic === "playerJoined") {
            store.set.addPlayer(data.payload);
          } else if (data.topic === "gameStarted") {
            store.set.gameState(GameState.Started);
          } else if (data.topic === "questionCountdown") {
            store.set.gameState(GameState.QuestionCountdown);
            store.set.countdownValue(data.payload.followUps);
            store.set.gameResult(undefined);
          } else if (data.topic === "question") {
            store.set.gameState(GameState.Question);
            store.set.question(data.payload.find);
          } else if (data.topic === "questionFinished") {
            store.set.gameState(GameState.Finished);
            store.set.gameResult(data.payload);
            store.set.updatePlayerRanking(data.payload);
          } else if (data.topic === "gameEnded") {
            store.set.gameState(GameState.GameEnded);
          }
        };
      }
    );
  }
  initWebSocket();
</script>

{#if gameState === GameState.SetupUsername || gameState === GameState.SetupMap}
  <div class="d-flex flex-column gap-4 align-items-center">
    <h1 class="old-font">City Knowledge Contest</h1>
    <p class="mb-5 fs-large">Wer findet die gesuchten Orte am schnellsten?</p>
    <Login />
  </div>
{:else}
  <Game />
{/if}

<footer>
  <div class="p-3 d-flex flex-column align-items-center gap-3">
    <div>Fancy Footer | License</div>
    <a
      class="color-black"
      href="https://www.flaticon.com/free-icons/location-pin"
      title="location pin icons"
      >Location pin icons created by Smashicons - Flaticon</a
    >
  </div>
</footer>
