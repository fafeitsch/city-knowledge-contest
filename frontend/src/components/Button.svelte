<style lang="scss">
@import '../styles/variables';

div {
  filter: drop-shadow(1px 1px 3px $brown);
}

button {
  background-color: $brown-darker;
  display: inline-block;
  clip-path: polygon(0 50%, 24px 100%, calc(100% - 24px) 100%, 100% 50%, calc(100% - 24px) 0, 24px 0);
  padding: 16px 64px;
  text-align: center;
  font-size: large;
  text-transform: uppercase;
  color: $old-map-lighter;
  border: none;
  font-family: 'MedievalSharp';
  background-image: radial-gradient(
    200% 70% at center 20%,
    $brown -30%,
    $brown 49%,
    $brown-darker 50%,
    $brown-darker 150%
  );
  background-repeat: no-repeat;
  transition: background-position-y ease-in-out 300ms;
  &:hover {
    cursor: pointer;
    background-position-y: -50px;
  }
}

.disabled {
  background-color: $old-map-lighter;
  color: darken($old-map-lighter, 20%);
  cursor: auto !important;
  background-image: none;
}
</style>

<!-- Polygon shaped button inspired by https://www.codecraftsman.us/create-a-hexagonal-button-in-css/ -->
<script lang="ts">
export let title: string;
export let disabled: boolean;
</script>

<div>
  <button class="{disabled ? 'disabled' : ''}" on:click disabled="{disabled}">{title}</button>
</div>
