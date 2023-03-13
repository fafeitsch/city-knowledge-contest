<script lang="ts" context="module">
</script>

<script lang="ts">
  import AvailableStreetLists from "../components/AvailableStreetLists.svelte";
  import Button from "../components/Button.svelte";
  import CopyIcon from "../components/CopyIcon.svelte";
  import store from "../store";
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
    {$game.roomKey}
    <CopyIcon
      width={16}
      height={16}
      className="color-black"
      on:click={() => {
        navigator.clipboard.writeText($game.roomKey);
      }}
    />
  </p>
  <AvailableStreetLists/>
  <Button title="Spiel starten" on:click={startGame} disabled="{$errors}"/>
</div>
