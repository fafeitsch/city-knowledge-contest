<style lang="scss">
@import '../styles/variables';

.card {
  background-color: $blue-light-transparent;
  border-radius: 16px;
  padding: 16px;
  gap: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>

<div class="card">
  <slot />
</div>
