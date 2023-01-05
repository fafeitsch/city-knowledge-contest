<script lang="ts">
  import { onMount } from "svelte";
  import store from "./store";

  let players: string[] = [];
  onMount(() => {
    store.get.players$.subscribe((value) => {
      players = value;
    });
  });
</script>

<div>Waiting Room</div>
<ul>
  {#each players as player}
    <li>{player}</li>
  {/each}
</ul>
