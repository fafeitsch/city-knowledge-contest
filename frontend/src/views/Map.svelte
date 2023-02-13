<script lang="ts">
  import L, { LatLng, latLng } from "leaflet";
  import { GameState, type Game } from "../store";
  import store from "../store";
  import { handleRPCRequest } from "../rpc";
  import Button from "../components/Button.svelte";

  let map;
  let gameState = store.get.gameState$;
  let countdownValue = store.get.countdownValue$;
  let question = store.get.question$;
  let game = store.get.game$;
  let gameResult = store.get.gameResult$;

  let currentResult: number = 0;

  type AnswerQuestion = {
    points: number;
  };

  function answerQuestion(guess: LatLng) {
    handleRPCRequest<
      {
        playerKey: string;
        roomKey: string;
        playerSecret: string;
        guess: Array<number>;
      },
      AnswerQuestion
    >({
      method: "answerQuestion",
      params: {
        playerKey: $game.playerKey,
        roomKey: $game.roomId,
        playerSecret: $game.playerSecret,
        guess: [guess.lat, guess.lng],
      },
    }).then((data) => {
      currentResult = data.result.points;
      store.set.gameState(GameState.Finished);
    });
  }

  function advanceGame(game: Game) {
    handleRPCRequest<
      { playerKey: string; playerSecret: string; roomKey: string },
      {}
    >({
      method: "advanceGame",
      params: {
        playerKey: game.playerKey,
        playerSecret: game.playerSecret,
        roomKey: game.roomId,
      },
    }).then((data) => console.log(data));
  }

  function createMap(container) {
    let m = L.map(container).setView(latLng(50, 10), 5);

    L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: `&copy;<a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>,
	        &copy;<a href="https://carto.com/attributions" target="_blank">CARTO</a>`,
      subdomains: "abcd",
      maxZoom: 20,
    }).addTo(m);

    m.addEventListener("click", (e) => answerQuestion(e.latlng));

    return m;
  }

  function mapAction(container) {
    map = createMap(container);

    return {
      destroy: () => {
        map.remove();
        map = null;
      },
    };
  }
</script>

<div>
  {#if $gameState === GameState.QuestionCountdown}
    <div class="overlay">{$countdownValue}</div>
  {/if}
  <div class="map full-viewheight full-viewwidth" use:mapAction />
  {#if $gameState === GameState.Question}
    <div class="container">
      <div>Suche den Ort {$question}</div>
    </div>
  {:else if $gameState === GameState.Finished}
    <div class="container">
      <div
        class="d-flex justify-content-spaced align-items-center width-100 p-4"
      >
        {#if currentResult !== 0}
          <div>Richtig 🥳</div>
        {:else}
          <div>Leider falsch 🤷</div>
        {/if}
        {#if gameResult !== undefined}
          <Button on:click={() => advanceGame($game)} title="Weiter" />
        {/if}
      </div>
    </div>
  {/if}
</div>

<style lang="scss">
  @import "../styles/variables";

  .overlay {
    z-index: 10000;
    height: 100vh;
    width: 100vw;
    position: absolute;
    top: 0;
    left: 0;
    background-color: rgba($color: #000000, $alpha: 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 128px;
  }

  .container {
    z-index: 999;
    position: absolute;
    bottom: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 150px;
    width: 100vw;
    background-color: $beige;
    font-size: xx-large;
  }
</style>
