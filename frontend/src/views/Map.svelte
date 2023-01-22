<script lang="ts">
  import L, { LatLng, latLng } from "leaflet";
  import { onMount } from "svelte";
  import { GameState, type Game, type GameResult } from "../store";
  import store from "../store";
  import { combineLatest } from "rxjs";
  import { handleRPCRequest } from "../rpc";
  import Button from "../components/Button.svelte";

  let map;
  let gameState: GameState;
  let countdownValue: string;
  let question: string;
  let game: Game;
  let points: number | undefined;
  let gameResult: GameResult | undefined;

  type AnswerQuestion = {
    points: number;
  };

  onMount(() => {
    combineLatest([
      store.get.gameState$,
      store.get.countdownValue$,
      store.get.question$,
      store.get.game$,
      store.get.gameResult$,
    ]).subscribe(
      ([gameState$, countdownValue$, question$, game$, gameResult$]) => {
        gameState = gameState$;
        countdownValue = countdownValue$;
        question = question$;
        game = game$;
        gameResult = gameResult$;
      }
    );
  });

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
        playerKey: game.playerKey,
        roomKey: game.roomId,
        playerSecret: game.playerSecret,
        guess: [guess.lat, guess.lng],
      },
    }).then((data) => {
      points = data.result.points;
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

<div style="position: relative;">
  {#if gameState === GameState.QuestionCountdown}
    <div class="overlay">{countdownValue}</div>
  {/if}
  <div class="map" style="height:100vh;width:100vw" use:mapAction />
  {#if gameState === GameState.Question}
    <div class="container">
      <div>Suche den Ort {question}</div>
    </div>
  {:else if gameState === GameState.Finished}
    <div class="container">
      <div
        class="d-flex justify-content-spaced align-items-center width-100 p-4"
      >
        <div>Punkte: {points}</div>
        <div>{JSON.stringify(gameResult)}</div>
        {#if gameResult !== undefined}
          <Button on:click={() => advanceGame(game)} title="Weiter" />
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
