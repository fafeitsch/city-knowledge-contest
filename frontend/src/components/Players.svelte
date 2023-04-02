<style lang="scss">
@import '../styles/variables';

.absolute-position {
  z-index: 1000;
  margin: 32px;
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
}

.player {
  background-color: rgba($brown, 0.7);
  padding: 16px;
  color: $beige;
}
</style>

<script lang="ts">
import Badge from '../components/Badge.svelte';
import type { Player } from '../store';
import { fly } from 'svelte/transition';
import { flip } from 'svelte/animate';

export let players: Player[];
export let absolutePosition = false;
</script>

<div style="align-self: normal" class="d-flex flex-column gap-4 p-4" class:absolute-position="{absolutePosition}">
  {#each players as player (player.playerKey)}
    <div
      in:fly="{{ x: 200, duration: 1000 }}"
      animate:flip
      id="{player.playerKey}"
      class="d-flex align-items-center justify-content-spaced player"
    >
      <div>{player.name}</div>
      <Badge value="{!player.points ? 0 : player.points}" />
    </div>
  {/each}
</div>
