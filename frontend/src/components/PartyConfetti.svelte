<style>
.confetti {
  position: fixed;
  top: -50px;
  left: 0;
  height: 100vh;
  width: 100vw;
  display: flex;
  justify-content: center;
  overflow: hidden;
  pointer-events: none;
  z-index: 700;
}
</style>

<script lang="ts">
import { Confetti } from 'svelte-confetti';
</script>

<div class="confetti">
  <Confetti
    x="{[-5, 5]}"
    y="{[0, 0.1]}"
    delay="{[500, 2000]}"
    infinite
    duration="5000"
    amount="200"
    fallDistance="100vh"
  />
</div>
