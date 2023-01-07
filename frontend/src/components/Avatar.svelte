<script lang="ts">
  export let name: string;

  function getRandomColor() {
    var letters = "0123456789ABCDEF";
    var color = "#";
    for (var i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  }
</script>

<div style="--color:{getRandomColor()};">
  {name.toUpperCase().slice(0, 1)}
</div>

<style lang="scss">
  div {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--color);
    font-size: large;
    font-weight: 600;
    color: black;
  }
</style>
