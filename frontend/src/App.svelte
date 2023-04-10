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

button {
  background-color: transparent;
  border: none;
  cursor: pointer;
  font-weight: bold;
}

a {
  font-size: smaller;
  font-weight: 600;
  text-decoration: underline;
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
  subscribeToSuccessfullyJoined,
  Topic,
} from './sockets';
import Legal from './components/Legal.svelte';
import rpc from './rpc';

let showDataProtection = false;
let showImprint = false;
let showOsmServices = false;
let legalInfo = rpc.getLegalInformation();
let osmServicesInfo = legalInfo.pipe(
  map(
    (info) =>
      `<h3>Nutzung von OSM-Services</h3>
    <article>Diese Anwendung nutzt zwei Services aus dem OpenStreetMap-Universum. Zum einen Nominatim: Dieser wird
     verwendet, um die Fragen zu stellen und Antworten zu verifizieren. Zum anderen wird ein TileServer verwendet, um die
     Karte darzustellen.</article>
     <article>Die Anwendung nutzt standardmäßig öffentliche Varianten beider Services. Die Entwickler des Spiels weisen aber ausdrücklich
     darauf hin, dass City-Knowledge-Contest mit eigenen Services deployed werden sollte.</article>
     <article>Momentan wird für Nominatim "${info.nominatimServer}" und für den TileServer "${info.tileServer}" genutzt. Der Betreiber
     des Spiels bzw. der Inhaber der URL, auf dem das Spiel läuft, ist dafür verantwortlich, dass beide Services gemäß der jeweiligen
     Benutzungsbedingungen angefragt werden.</article>
`,
  ),
);

let gameState = merge(
  subscribeToSuccessfullyJoined().pipe(
    map((data) => {
      if (data?.started) {
        return 'Question';
      }
      return data ? 'Waiting' : undefined;
    }),
  ),
  subscribeToCountdown().pipe(map((data) => (data ? 'Question' : undefined))),
  subscribeToQuestion().pipe(map((data) => (data ? 'Question' : undefined))),
  subscribeToGameEnded().pipe(map((data) => (data ? 'GameEnded' : undefined))),
).pipe(distinctUntilChanged());
initWebSocket();

function showImprintDialog() {
  showImprint = true;
  showDataProtection = false;
  showOsmServices = false;
}

function showDataProtectionDialog() {
  showImprint = false;
  showDataProtection = true;
  showOsmServices = false;
}

function showOsmServicesDialog() {
  showImprint = false;
  showDataProtection = false;
  showOsmServices = true;
}
</script>

{#if showDataProtection}
  <Legal info="{$legalInfo.dataProtection}" on:closeLegal="{() => (showDataProtection = false)}" />
{/if}
{#if showImprint}
  <Legal info="{$legalInfo.imprint}" on:closeLegal="{() => (showImprint = false)}" />
{/if}
{#if showOsmServices}
  <Legal info="{$osmServicesInfo}" on:closeLegal="{() => (showOsmServices = false)}" />
{/if}
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
    <div>
      <button on:click="{showImprintDialog}">Impressum</button>
      |
      <button on:click="{showDataProtectionDialog}">Datenschutz</button> |
      <button on:click="{showOsmServicesDialog}">OSM-Services</button>
      |
      <a href="https://github.com/fafeitsch/city-knowledge-contest" target="_blank">Projekt auf Github</a>
      |
      <a href="https://www.flaticon.com/free-icons/location-pin" title="location pin icons" target="_blank"
        >Location pin icons created by Smashicons - Flaticon</a
      >
    </div>
  </div>
</footer>
