<script lang="ts">
import { filter, map, merge, Observable, startWith, Subject, switchMap, tap } from 'rxjs';
import AvailableStreetLists from './AvailableStreetLists.svelte';
import Button from '../../components/Button.svelte';
import store from '../../store';
import rpc, { type RoomConfiguration, type RoomConfigurationResult } from '../../rpc';
import CopyIcon from './CopyIcon.svelte';
import Players from '../../components/Players.svelte';
import CoverImage from '../../components/CoverImage.svelte';
import Card from '../../components/Card.svelte';
import Input from '../../components/Input.svelte';
import { subscribeToJoined, subscribeToRoomUpdated } from '../../sockets';

const decimalRegex = /^\d+$/;
let room = store.get.room$;
let remoteConfiguration: Observable<RoomConfigurationResult | undefined> = merge(
  subscribeToRoomUpdated(),
  subscribeToJoined().pipe(map((payload) => payload?.options)),
).pipe(filter((config) => !!config));
let players = store.get.players$;

function updateStreetList(event: CustomEvent, config: RoomConfiguration) {
  const newConfig: RoomConfiguration = {
    maxAnswerTimeSec: config.maxAnswerTimeSec,
    listFileName: event.detail,
    numberOfQuestions: config.numberOfQuestions,
  };
  configureGame(newConfig);
}

function updateNumberOfQuestions(event: CustomEvent, config: RoomConfiguration) {
  let text = event.detail;
  if (!decimalRegex.test(text)) {
    text = text.replaceAll(/[^0-9]/g, '');
  }
  const newConfig: RoomConfiguration = {
    maxAnswerTimeSec: config.maxAnswerTimeSec,
    listFileName: config.listFileName,
    numberOfQuestions: parseInt(text, 10),
  };
  configureGame(newConfig);
}

function updateMaxAnswerTimeSec(event: CustomEvent, config: RoomConfiguration) {
  let text = event.detail;
  if (!decimalRegex.test(text)) {
    text = text.replaceAll(/[^0-9]/g, '');
  }
  const newConfig: RoomConfiguration = {
    maxAnswerTimeSec: parseInt(text, 10),
    listFileName: config.listFileName,
    numberOfQuestions: config.numberOfQuestions,
  };
  configureGame(newConfig);
}

function configureGame(config: RoomConfiguration) {
  rpc.updateRoomConfiguration(config).subscribe();
}

function startGame() {
  rpc.startGame().subscribe();
}
</script>

<Players playerKey="{$room.playerKey}" players="{$players}" absolutePosition="{true}" />

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
    <AvailableStreetLists
      on:streetListChanged="{(event) => updateStreetList(event, $remoteConfiguration)}"
      selectedStreetList="{$remoteConfiguration?.listFileName || ''}"
    />
    <div class="d-flex gap-4">
      <Input
        placeholder="Anzahl der Fragen"
        on:input="{(event) => updateNumberOfQuestions(event, $remoteConfiguration)}"
        value="{$remoteConfiguration?.numberOfQuestions || ''}"
        type="number"
      />
      <Input
        placeholder="Sekunden pro Frage"
        on:input="{(event) => updateMaxAnswerTimeSec(event, $remoteConfiguration)}"
        value="{$remoteConfiguration?.maxAnswerTimeSec || ''}"
        type="number"
      />
    </div>
    <Button
      title="Spiel starten"
      on:click="{startGame}"
      disabled="{!$remoteConfiguration || $remoteConfiguration.errors.length > 0}"
    />
  </Card>
</CoverImage>
