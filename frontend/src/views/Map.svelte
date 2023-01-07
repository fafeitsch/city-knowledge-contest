<script lang="ts">
  import L, { latLng } from "leaflet";
  import { onMount } from "svelte";
  import { GameState } from "../store";
  import store from "../store";
  import { combineLatest } from "rxjs";

  let map;
  let gameState: GameState;
  let countdownValue: string;

  onMount(() => {
    combineLatest([store.get.gameState$, store.get.countdownValue$]).subscribe(
      ([gameState$, countdownValue$]) => {
        gameState = gameState$;
        countdownValue = countdownValue$;
      }
    );
  });

  function createMap(container) {
    let m = L.map(container).setView(latLng(50, 10), 5);
    L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: `&copy;<a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>,
	        &copy;<a href="https://carto.com/attributions" target="_blank">CARTO</a>`,
      subdomains: "abcd",
      maxZoom: 20,
    }).addTo(m);

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
</div>

<style lang="scss">
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
</style>
