<style lang="scss">
.cover-image {
  background: url('../assets/images/skyline.jpeg') no-repeat center center fixed;
  -webkit-background-size: cover;
  -moz-background-size: cover;
  -o-background-size: cover;
  background-size: cover;
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 64px;
}
</style>

<div class="cover-image">
  <slot />
</div>
