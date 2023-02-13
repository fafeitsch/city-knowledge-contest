<script lang="ts">
  import store, { GameState } from "../store";
  import DefaultLayout from "./DefaultLayout.svelte";
  import Map from "./Map.svelte";
  import Players from "./Players.svelte";
  import WaitingRoom from "./WaitingRoom.svelte";

  let gameState = store.get.gameState$;
  let players = store.get.players$;
</script>

<DefaultLayout>
  <div slot="content-container">
    {#if $gameState === GameState.Waiting}
      <WaitingRoom />
    {:else}
      <Map />
    {/if}
  </div>

  <Players slot="player-container" players={$players} />
</DefaultLayout>
