<script lang="ts">
import { onMount } from 'svelte';
import store from '../../store';
import Input from '../../components/Input.svelte';
import Button from '../../components/Button.svelte';
import rpc from '../../rpc';

let roomKey = '';
let userName = '';
let userNameSet = false;

onMount(() => {
  userName = localStorage.getItem('userName');
  userNameSet = !!userName;
});

function handleUsernameChange(newName: CustomEvent<string>) {
  userName = newName.detail;
  localStorage.setItem('userName', userName);
}

function handleRoomIdChange(event: CustomEvent<string>) {
  roomKey = event.detail;
}

function createRoom() {
  rpc.createRoom(userName).subscribe((data) => store.set.game(data));
}

async function joinRoom() {
  rpc.joinRoom(roomKey, userName).subscribe((data) => store.set.game(data));
}
</script>

<div class="d-flex flex-column gap-3 align-items-center">
  <Input on:input="{handleUsernameChange}" placeholder="Spielername" value="{userName}" />
</div>
<div class="d-flex gap-5 align-items-center">
  <div class="d-flex flex-column gap-2 align-items-center">
    <Input on:input="{handleRoomIdChange}" placeholder="Karten-ID" />
    <Button on:click="{joinRoom}" title="Karte beitreten" disabled="{!userName || !roomKey}" />
  </div>
  <div>– oder –</div>
  <Button title="Eine neue Karte erstellen" on:click="{createRoom}" disabled="{!userName}" />
</div>