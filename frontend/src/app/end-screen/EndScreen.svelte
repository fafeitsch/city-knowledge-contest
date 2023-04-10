<script lang="ts">
import Button from '../../components/Button.svelte';
import store from '../../store';
import Players from '../../components/Players.svelte';
import { get } from 'svelte/store';
import App from '../../App.svelte';

let players = store.get.players$;
let room = store.get.room$;

function newGame() {
  window.location = window.location.protocol + '//' + window.location.host;
  store.set.resetGame();
}
</script>

<div class="d-flex flex-column gap-3 align-items-center">
  <div class="old-font fs-x-large">Das Spiel ist leider vorbei.</div>
  <Players playerKey="{$room.playerKey}" players="{$players}" />
  <Button title="Neues Spiel" on:click="{() => newGame()}" />
</div>
