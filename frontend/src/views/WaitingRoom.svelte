<script lang="ts">
  import { combineLatest } from "rxjs";
  import { onMount } from "svelte";
  import Avatar from "../components/Avatar.svelte";
  import Button from "../components/Button.svelte";
  import store, { GameState } from "../store";

  let players: string[] = [];
  let roomKey: string;
  let playerKey: string;

  onMount(() => {
    combineLatest([
      store.get.players$,
      store.get.playerKey$,
      store.get.roomId$,
    ]).subscribe(([players$, playerKey$, roomKey$]) => {
      players = players$;
      playerKey = playerKey$;
      roomKey = roomKey$;
    });
  });

  function startGame() {
    fetch("http://localhost:23123/rpc", {
      method: "POST",
      body: JSON.stringify({
        method: "updateRoom",
        params: {
          area: [
            [49.795007, 9.892073],
            [49.802597, 9.909925],
            [49.802707, 9.916363],
            [49.803151, 9.925718],
            [49.804646, 9.93619],
            [49.804979, 9.951725],
            [49.795672, 9.967861],
            [49.790186, 9.97035],
            [49.790186, 9.963999],
            [49.777994, 9.954729],
            [49.777384, 9.961081],
            [49.770234, 9.959021],
            [49.771786, 9.951811],
            [49.773338, 9.943657],
            [49.772562, 9.938335],
            [49.775389, 9.932241],
            [49.785531, 9.927092],
            [49.793566, 9.926147],
            [49.790851, 9.907608],
            [49.78941, 9.898338],
            [49.794453, 9.892845],
          ],
          numberOfQuestions: 10,
          playerKey,
          roomKey,
        },
      }),
    })
      .then(() => {
        fetch("http://localhost:23123/rpc", {
          method: "POST",
          body: JSON.stringify({
            method: "startGame",
            params: {
              playerKey,
              roomKey,
            },
          }),
        });
      })
      .then(() => {
        store.set.gameState(GameState.Started);
      });
  }
</script>

<div class="p-3 position-absolute top-0 right-0">
  <div class="d-flex gap-3 align-items-cente4">
    {#each players as player}<Avatar name={player} />{/each}
  </div>
</div>

<div class="d-flex flex-column align-items-center gap-5">
  <div class="old-font fs-large">Gleich geht das Spiel los …</div>
  <Button title="Spiel starten" on:click={startGame} />
</div>
