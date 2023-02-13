<div class="page-container">
  <div class="player-container">
    <slot name="player-container" />
  </div>
  <div class="content-container">
    <slot name="content-container" />
  </div>
</div>

<style lang="scss">
  .page-container {
    width: 100%;
    height: 100%;
    position: relative;

    .content-container {
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .player-container {
      z-index: 1000;
      margin: 32px;
      position: absolute;
      top: 0;
      right: 0;
      z-index: 1000;
      width: 300px;
      display: flex;
      flex-direction: column;
      gap: 16px;
      align-items: flex-start;
      padding: 16px;
    }
  }
</style>
