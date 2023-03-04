<script lang="ts" context="module">
  export type UpdateRoomParams = {
    listFileName: string;
    numberOfQuestions: number;
    playerKey: string;
    playerSecret: string;
    roomKey: string;
    maxAnswerTimeSec: number;
  };

  export const defaultRoomSeetings = {
    numberOfQuestions: 10,
    listFileName: "wuerzburg.json",
    maxAnswerTimeSec: 600,
  };
</script>

<script lang="ts">
  import AvailableStreetLists from "../components/AvailableStreetLists.svelte";
  import Button from "../components/Button.svelte";
  import CopyIcon from "../components/CopyIcon.svelte";
  import {handleRPCRequest} from "../rpc";
  import store, {GameState} from "../store";
  import {map} from 'rxjs';

  let game = store.get.game$;
  let errors = store.get.errors$.pipe(map(errors => errors.length > 0));

  async function startGame() {
    await store.methods.startGame()
  }
</script>

<div class="d-flex flex-column align-items-center gap-5">
  <div class="old-font fs-large">Gleich geht das Spiel los …</div>
  <p class="mt-5">
    Teile den Code, um andere Personen zu diesem Spiel einzuladen:
  </p>
  <p class="fw-bold p-3 bg-old-map-lighter d-flex align-items-center gap-3">
    {$game.roomId}
    <CopyIcon
      width={16}
      height={16}
      className="color-black"
      on:click={() => {
        navigator.clipboard.writeText($game.roomId);
      }}
    />
  </p>
  <AvailableStreetLists/>
  <Button title="Spiel starten" on:click={startGame} disabled="{$errors}"/>
</div>
