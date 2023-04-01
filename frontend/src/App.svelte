<script lang="ts">
import Login from './views/Login.svelte';
import store from './store';
import { filter } from 'rxjs';
import Game from './views/Game.svelte';
import { environment } from './environment';

let gameState = store.get.gameState$;

function initWebSocket() {
  store.get.room$.pipe(filter((room) => !!room)).subscribe((room) => {
    const websocket = new WebSocket(environment[import.meta.env.MODE].wsUrl + room.roomKey + '/' + room.playerKey);
    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.topic === 'successfullyJoined') {
        console.log(data.payload.players);
        store.set.players(data.payload.players);
      } else if (data.topic === 'playerJoined') {
        store.set.addPlayer(data.payload);
      } else if (data.topic === 'questionCountdown') {
        store.set.countdownValue(data.payload.followUps + 1);
      } else if (data.topic === 'question') {
        store.set.countdownValue(undefined);
        store.set.question(data.payload.find);
      } else if (data.topic === 'questionFinished') {
        store.set.gameResult(data.payload);
        store.set.updatePlayerRanking(data.payload);
      } else if (data.topic === 'gameEnded') {
        store.set.game(undefined);
      }
    };
  });
}

initWebSocket();
</script>

{#if $gameState === 'SetupMap'}
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
    <a class="color-black" href="https://www.flaticon.com/free-icons/location-pin" title="location pin icons"
      >Location pin icons created by Smashicons - Flaticon</a
    >
  </div>
</footer>
