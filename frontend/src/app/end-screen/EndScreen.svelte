<script lang="ts">
import Button from '../../components/Button.svelte';
import store from '../../store';
import Players from '../../components/Players.svelte';
import { get } from 'svelte/store';
import App from '../../App.svelte';
import CoverImage from '../../components/CoverImage.svelte';
import Card from '../../components/Card.svelte';

let players = store.get.players$;
let room = store.get.room$;

function newGame() {
  window.location = window.location.protocol + '//' + window.location.host;
  store.set.resetGame();
}
</script>

<CoverImage>
  <h1>Das Spiel ist leider vorbei.</h1>
  <Card>
    <div style="width: 400px">
      <Players playerKey="{$room.playerKey}" players="{$players}" />
    </div>
    <Button title="Neues Spiel starten" on:click="{() => newGame()}" />
  </Card>
</CoverImage>
