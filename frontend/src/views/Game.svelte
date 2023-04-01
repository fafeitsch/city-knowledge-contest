<script lang="ts">
import Button from "../components/Button.svelte";
import store from "../store";
import DefaultLayout from "./DefaultLayout.svelte";
import Map from "./Map.svelte";
import Players from "./Players.svelte";
import WaitingRoom from "./WaitingRoom.svelte";

let gameState = store.get.gameState$;
let players = store.get.players$;

function newGame() {
  store.set.resetGame();
}
</script>

<DefaultLayout>
  <div slot="content-container">
    {#if $gameState === "Waiting"}
      <WaitingRoom />
    {:else if $gameState === "GameEnded"}
      <div
        class="d-flex flex-column justify-content-center gap-3 align-items-center"
      >
        <div class="old-font fs-x-large">Das Spiel ist leider vorbei 🤷</div>
        <Button title="Neues Spiel" on:click="{() => newGame()}" />
      </div>
    {:else}
      <Map />
    {/if}
  </div>

  <Players slot="player-container" players="{$players}" />
</DefaultLayout>
