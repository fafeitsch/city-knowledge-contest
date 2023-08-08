<script lang="ts">
import Button from '../../components/Button.svelte';
import Players from '../../components/Players.svelte';
import CoverImage from '../../components/CoverImage.svelte';
import Card from '../../components/Card.svelte';
import { store } from '../../store';

let players = store.get.players$;
let room = store.get.room$;

function newGame() {
  window.location = window.location.protocol + '//' + window.location.host;
  store.set.resetGame();
}
</script>

<CoverImage>
  <h1 data-testid="game-over-title">Das Spiel ist leider vorbei.</h1>
  <Card>
    <div style="width: 400px">
      <Players playerKey="{$room.playerKey}" players="{$players}" />
    </div>
    <Button title="Neues Spiel starten" on:click="{() => newGame()}" e2eTestId="start-a-new-game-button" />
  </Card>
</CoverImage>
