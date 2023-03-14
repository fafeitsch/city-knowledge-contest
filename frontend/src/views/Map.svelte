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

<script lang="ts">
import store from "../store";
import Button from "../components/Button.svelte";
import Leaflet from "../components/Leaflet.svelte";
import PartyConfetti from "../components/PartyConfetti.svelte";

let countdownValue = store.get.countdownValue$;
let question = store.get.question$;
let lastResult = store.get.lastResult$;

async function advanceGame() {
  await store.methods.advanceGame();
}
</script>

<div>
  {#if $countdownValue}
    <div class="overlay">{$countdownValue}</div>
  {/if}
  <Leaflet />
  {#if $question && $lastResult === undefined}
    <div class="container">
      <div>Suche den Ort {$question}</div>
    </div>
  {:else if $lastResult !== undefined}
    <div class="container">
      <div
        class="d-flex justify-content-spaced align-items-center width-100 p-4"
      >
        {#if $lastResult !== 0}
          <div id="party">Richtig 🥳</div>
          <PartyConfetti />
        {:else}
          <div>Leider falsch 🤷</div>
        {/if}
        <Button on:click="{() => advanceGame()}" title="Weiter" />
      </div>
    </div>
  {/if}
</div>
