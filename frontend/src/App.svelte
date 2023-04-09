<style lang="scss">
.page-container {
  width: 100%;
  height: 100%;
  position: relative;

  .content-container {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>

<script lang="ts">
import { distinctUntilChanged, map, merge } from 'rxjs';
import Login from './app/login/Login.svelte';
import WaitingRoom from './app/waiting-room/WaitingRoom.svelte';
import Map from './app/map/Map.svelte';
import EndScreen from './app/end-screen/EndScreen.svelte';
import {
  initWebSocket,
  subscribeToCountdown,
  subscribeToGameEnded,
  subscribeToQuestion,
  subscribeToSocketTopic,
  Topic,
} from './sockets';

let gameState = merge(
  subscribeToSocketTopic(Topic.successfullyJoined).pipe(map((data) => (data ? 'Waiting' : undefined))),
  subscribeToCountdown().pipe(map((data) => (data ? 'Question' : undefined))),
  subscribeToQuestion().pipe(map((data) => (data ? 'Question' : undefined))),
  subscribeToGameEnded().pipe(map((data) => (data ? 'GameEnded' : undefined))),
).pipe(distinctUntilChanged());
initWebSocket();
</script>

<div class="page-container">
  <div class="content-container">
    {#if $gameState === undefined}
      <Login />
    {:else if $gameState === 'Waiting'}
      <WaitingRoom />
    {:else if $gameState === 'GameEnded'}
      <EndScreen />
    {:else}
      <Map />
    {/if}
  </div>
</div>

<footer>
  <div class="p-3 d-flex flex-column align-items-center gap-3">
    <div>Fancy Footer | License</div>
    <a class="color-black" href="https://www.flaticon.com/free-icons/location-pin" title="location pin icons"
      >Location pin icons created by Smashicons - Flaticon</a
    >
  </div>
</footer>
