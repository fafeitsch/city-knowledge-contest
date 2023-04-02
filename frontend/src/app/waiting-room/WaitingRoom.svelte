<script lang="ts">
import { map, startWith, Subject, switchMap } from 'rxjs';
import AvailableStreetLists from './AvailableStreetLists.svelte';
import Button from '../../components/Button.svelte';
import store from '../../store';
import { type RoomConfiguration, updateRoomConfiguration } from '../../rpc';
import CopyIcon from './CopyIcon.svelte';
import Players from '../../components/Players.svelte';

let room = store.get.room$;
let streetList: string | undefined = undefined;
let gameConfiguration$ = new Subject<RoomConfiguration>();
let errors = gameConfiguration$.pipe(
  switchMap((config) => updateRoomConfiguration(config)),
  startWith(['noConfiguration']),
  map((errors) => errors.length > 0),
);
let players = store.get.players$;

function updateStreetList(event: CustomEvent) {
  streetList = event.detail;
  configureGame();
}

function configureGame() {
  gameConfiguration$.next({ listFileName: streetList });
}

async function startGame() {
  await store.methods.startGame();
}
</script>

<Players players="{$players}" />

<div class="d-flex flex-column align-items-center gap-5">
  <div class="old-font fs-large">Gleich geht das Spiel los …</div>
  <p class="mt-5">Teile den Code, um andere Personen zu diesem Spiel einzuladen:</p>
  <p class="fw-bold p-3 bg-old-map-lighter d-flex align-items-center gap-3">
    {$room.roomKey}
    <CopyIcon
      width="{16}"
      height="{16}"
      className="color-black"
      on:click="{() => {
        navigator.clipboard.writeText($room.roomKey);
      }}"
    />
  </p>
  <AvailableStreetLists on:streetListChanged="{updateStreetList}" />
  <Button title="Spiel starten" on:click="{startGame}" disabled="{$errors}" />
</div>
