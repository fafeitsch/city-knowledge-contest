<script lang="ts">
import { map, startWith, Subject, switchMap } from 'rxjs';
import AvailableStreetLists from './AvailableStreetLists.svelte';
import Button from '../../components/Button.svelte';
import store from '../../store';
import rpc, { type RoomConfiguration } from '../../rpc';
import CopyIcon from './CopyIcon.svelte';
import Players from '../../components/Players.svelte';
import CoverImage from '../../components/CoverImage.svelte';
import Card from '../../components/Card.svelte';
import Input from '../../components/Input.svelte';

let room = store.get.room$;
let streetList: string | undefined = undefined;
let gameConfiguration$ = new Subject<RoomConfiguration>();
let errors = gameConfiguration$.pipe(
  switchMap((config) => rpc.updateRoomConfiguration(config)),
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

function startGame() {
  rpc.startGame().subscribe();
}
</script>

<Players players="{$players}" absolutePosition="{true}" />

<CoverImage>
  <h1>Gleich geht's los…</h1>
  <Card>
    <p>Teile den Link, um andere Personen zu diesem Spiel einzuladen:</p>
    <p class="fw-bold d-flex align-items-center gap-3">
      {window.location}
      <CopyIcon
        width="{16}"
        height="{16}"
        className="color-black"
        on:click="{() => {
          navigator.clipboard.writeText(window.location.toString());
        }}"
      />
    </p>
  </Card>
  <Card>
    <AvailableStreetLists on:streetListChanged="{updateStreetList}" />
    <div class="d-flex gap-4">
      <Input placeholder="Anzahl der Fragen" />
      <Input placeholder="Zeit" />
    </div>
    <Button title="Spiel starten" on:click="{startGame}" disabled="{$errors}" />
  </Card>
</CoverImage>
