<script lang="ts">
import Button from "../components/Button.svelte";
import Input from "../components/Input.svelte";
import store from "../store";

let roomKey = "";
let userName = "";
let gameState = store.get.gameState$;

function handleUsernameChange(event: Event) {
  const target = event.target as HTMLInputElement;
  userName = target.value;
}

function handleRoomIdChange(event: Event) {
  const target = event.target as HTMLInputElement;
  roomKey = target.value;
}

function handleOnClick() {
  store.set.username(userName);
}

async function createRoom() {
  await store.methods.createRoom();
}

async function joinRoom() {
  await store.methods.joinRoom(roomKey);
}
</script>

{#if $gameState === "SetupUsername"}
  <div class="d-flex flex-column gap-3 align-items-center">
    <Input
      on:change="{handleUsernameChange}"
      placeholder="Gib deinen Spielername ein"
    />
    <Button on:click="{handleOnClick}" title="Los geht's" />
  </div>
{:else if $gameState === "SetupMap"}
  <div class="d-flex gap-5 align-items-center">
    <div class="d-flex flex-column gap-2 align-items-center">
      <Input
        on:change="{handleRoomIdChange}"
        placeholder="Karten-ID eingeben"
      />
      <Button on:click="{joinRoom}" title="Karte beitreten" />
    </div>
    <div>– oder –</div>
    <Button title="Eine neue Karte erstellen" on:click="{createRoom}" />
  </div>
{/if}
