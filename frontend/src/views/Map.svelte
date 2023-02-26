<script lang="ts">
  import { GameState, type Game } from "../store";
  import store from "../store";
  import { handleRPCRequest } from "../rpc";
  import Button from "../components/Button.svelte";
  import Leaflet from "../components/Leaflet.svelte";

  let gameState = store.get.gameState$;
  let countdownValue = store.get.countdownValue$;
  let question = store.get.question$;
  let game = store.get.game$;
  let gameResult = store.get.gameResult$;

  let currentResult: number = 0;
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
</script>

<div>
  {#if $gameState === GameState.QuestionCountdown}
    <div class="overlay">{$countdownValue}</div>
  {/if}
  <Leaflet
    game={$game}
    solution={$gameResult ? $gameResult.solution : undefined}
    bind:currentResult
  />
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
        {#if $gameResult !== undefined}
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
