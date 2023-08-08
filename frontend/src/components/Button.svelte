<style lang="scss">
@import '../styles/variables';

button {
  background-color: $yellow-main;
  color: $blue-dark;
  font-family: 'LilitaOne', sans-serif;
  padding: 16px 16px;
  text-align: center;
  font-size: large;
  text-transform: uppercase;
  border: none;
  border-radius: 8px;
  transition: 0.3s;
  &:hover {
    cursor: pointer;
    background-color: $blue-dark;
    color: $yellow-main;
  }
}

.disabled {
  background-color: $gray-light;
  color: $gray-dark;
  &:hover {
    cursor: not-allowed;
    background-color: $gray-light;
    color: $gray-dark;
  }
}
</style>

<!-- Polygon shaped button inspired by https://www.codecraftsman.us/create-a-hexagonal-button-in-css/ -->
<script lang="ts">
export let title: string;
export let disabled: boolean = false;
export let e2eTestId: string = 'button';
</script>

<button class="{disabled ? 'disabled' : ''}" on:click disabled="{disabled}" data-testid="{e2eTestId}">{title}</button>
