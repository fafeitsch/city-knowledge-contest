<style lang="scss">
@import '../styles/variables';

.absolute-position {
  z-index: 1000;
  margin: 32px;
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
  max-height: 400px;
  overflow: auto;
  background-color: white;
  border-radius: 16px;
  padding: 16px;
}

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

export let players: Player[];
export let absolutePosition = false;
export let playerKey = '';
</script>

<div style="align-self: normal" class="d-flex flex-column gap-4 p-4" class:absolute-position="{absolutePosition}">
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
      <div>{!player.points ? 0 : player.points} Punkte</div>
    </div>
  {/each}
</div>
