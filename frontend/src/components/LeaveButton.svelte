<style lang="scss">
@import '../styles/variables';

.leave-button {
  z-index: 999;
  position: absolute;
  top: 0;
  left: 32px;
  margin: 16px;
  padding: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50px;
  width: 200px;
}
</style>

<script lang="ts">
import store from '../store';
import Button from '../components/Button.svelte';
import rpc from '../rpc';

function leaveGame() {
  rpc.leaveGame().subscribe((data) => {
    window.location = window.location.protocol + '//' + window.location.host;
    store.set.resetGame();
  });
}
</script>

<div class="leave-button">
  <Button on:click="{leaveGame}" title="Raum verlassen" e2eTestId="leave-game-button" />
</div>
