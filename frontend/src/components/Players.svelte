<style lang="scss">
@import '../styles/variables';

.player {
  background-color: $blue-main;
  padding: 16px;
  color: $blue-dark;
  border-radius: 16px;
  color: white;
  font-family: 'LilitaOne';
}

.highlight {
  background-color: $yellow-main;
  color: $blue-dark;
}
</style>

<script lang="ts">
import type { Player } from '../store';
import { fly } from 'svelte/transition';
import { flip } from 'svelte/animate';

export let players: Player[] = [];
export let playerKey = '';
</script>

<div style="align-self: normal" class="d-flex flex-column gap-4 p-4">
  {#each players as player (player.playerKey)}
    <div
      in:fly="{{ x: 200, duration: 1000 }}"
      animate:flip
      id="{player.playerKey}"
      class="d-flex align-items-center justify-content-spaced player {playerKey === player.playerKey
        ? 'highlight'
        : ''}"
    >
      <div>{player.name}</div>
      <div class="d-flex gap-2">
        <div>{!player.points ? 0 : player.points} Punkte</div>
        {#if player.delta}
          <div>+ {player.delta}</div>
        {/if}
      </div>
    </div>
  {/each}
</div>
