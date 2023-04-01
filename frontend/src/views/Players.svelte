<style lang="scss">
@import '../styles/variables';

.players {
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

{#each players as player (player.playerKey)}
  <div in:fly="{{ x: 200, duration: 1000 }}" animate:flip id="{player.playerKey}" class="players">
    <div>{player.name}</div>
    <Badge value="{!player.points ? 0 : player.points}" />
  </div>
{/each}
