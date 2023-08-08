<style lang="scss">
@import '../styles/variables';

.player {
  background-color: $blue-main;
  padding: 16px;
  border-radius: 16px;
  color: white;
  font-family: 'LilitaOne', sans-serif;
}

.highlight {
  background-color: $yellow-main;
  color: $blue-dark;
}

.kick-button {
  position: absolute;
  right: -12px;
  top: -12px;
}

.container {
  position: relative;
}
</style>

<script lang="ts">
import type { Player } from '../store';
import { fly } from 'svelte/transition';
import { flip } from 'svelte/animate';
import KickButton from './KickButton.svelte';
import { createEventDispatcher } from 'svelte';

export let players: Player[] = [];
export let playerKey = '';
export let enableKick = false;

let dispatch = createEventDispatcher();

function kickPlayer(player: Player) {
  dispatch('kickPlayer', player);
}
</script>

<div style="align-self: normal" class="d-flex flex-column gap-4 p-4">
  {#each players as player (player.playerKey)}
    <div
      class="d-flex align-items-center container"
      id="{player.playerKey}"
      in:fly="{{ x: 200, duration: 1000 }}"
      animate:flip
    >
      <div
        class=" flex-grow-1 d-flex w-100 gap-4 align-items-center player {playerKey === player.playerKey
          ? 'highlight'
          : ''}"
        data-testid="player-list-entry"
      >
        <span class="flex-grow-1 flex-shrink-1 ellipsis">{player.name}</span>
        <div class="d-flex flex-shrink-0">
          {!player.points ? 0 : player.points} Punkte
          {#if player.delta}
            + {player.delta}
          {/if}
        </div>
      </div>
      {#if playerKey !== player.playerKey && enableKick}
        <div class="kick-button">
          <KickButton on:confirm="{() => kickPlayer(player)}" />
        </div>
      {/if}
    </div>
  {/each}
</div>
