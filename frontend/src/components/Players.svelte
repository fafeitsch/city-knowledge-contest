<style lang="scss">
@import '../styles/variables';

.player-container {
  z-index: 1000;
  margin: 32px;
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: flex-start;
  padding: 16px;
}

.player {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: rgba($brown, 0.7);
  padding: 16px;
  width: 100%;
  color: $beige;
}
</style>

<script lang="ts">
import Badge from '../components/Badge.svelte';
import type { Player } from '../store';
import { fly } from 'svelte/transition';
import { flip } from 'svelte/animate';

export let players: Player[];
</script>

<div class="player-container">
  {#each players as player (player.playerKey)}
    <div in:fly="{{ x: 200, duration: 1000 }}" animate:flip id="{player.playerKey}" class="player">
      <div>{player.name}</div>
      <Badge value="{!player.points ? 0 : player.points}" />
    </div>
  {/each}
</div>
