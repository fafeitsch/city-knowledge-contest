<style lang="scss">
.kick-button {
  border: 1px lightgray solid;
  background-color: lightgray;
  border-radius: 12px;
  min-width: 24px;
  height: 24px;
  cursor: pointer;
  font-family: 'LilitaOne', sans-serif;
}
</style>

<script>
import { createEventDispatcher } from 'svelte';

let dispatch = createEventDispatcher();
let caption = '🗙';
let state = 'initial';

function buttonClicked() {
  if (state === 'initial') {
    caption = 'Kick?';
    state === 'pending';
    setTimeout(() => {
      state = 'armed';
      caption = 'Kick!';
    }, 1000);
    setTimeout(() => {
      state = 'initial';
      caption = '🗙';
    }, 2000);
  } else if (state === 'armed') {
    dispatch('confirm', undefined);
  }
}
</script>

<button class="kick-button" on:click="{buttonClicked}" data-testid="kick-player-button">{caption}</button>
