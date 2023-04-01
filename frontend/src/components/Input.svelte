<style lang="scss">
@import '../styles/variables';

.input-container {
  position: relative;
  display: flex;
}

input {
  padding: 16px 8px;
  width: 100%;
  background-color: $beige;
  border: 3px solid $old-map-lighter;
  font-size: large;
  outline: none;
  transition: 500ms;

  &:focus {
    border: 3px solid $old-map-darker;
  }
}

.label {
  color: $brown-darker;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 9px;
  right: 0;
  border: 3px solid transparent;
  background-color: transparent;
  display: flex;
  align-items: center;
  width: 100%;
  pointer-events: none;
}

input:focus + .label .placeholder,
.placeholder:not(div[data-value=''] .placeholder) {
  font-size: 0.8rem;
  transform: translate(0, -130%);
}

.placeholder {
  transform: translate(0);
  transition: transform 0.15s ease-out, font-size 0.15s ease-out, background-color 0.2s ease-out, color 0.15s ease-out;
}
</style>

<script lang="ts">
import { createEventDispatcher } from 'svelte';

const dispatch = createEventDispatcher();

function dispatchInputEvent(event) {
  dispatch('input', event.target.value);
}

export let placeholder: string;
export let value: string = '';
</script>

<div class="input-container" data-value="{value}">
  <input on:change on:input="{dispatchInputEvent}" autocomplete="off" bind:value="{value}" />
  <label class="label">
    <div class="placeholder">{placeholder}</div>
  </label>
</div>
