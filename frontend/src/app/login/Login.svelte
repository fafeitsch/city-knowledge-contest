<style lang="scss">
.header-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.save-name-notice {
  width: 400px;
  line-height: 1.5;
}
</style>

<script lang="ts">
import { onMount } from 'svelte';
import store from '../../store';
import Input from '../../components/Input.svelte';
import Button from '../../components/Button.svelte';
import rpc from '../../rpc';
import CoverImage from '../../components/CoverImage.svelte';
import Card from '../../components/Card.svelte';
import { Subject, catchError, map, of, switchMap, tap } from 'rxjs';

let roomKey = '';
let userName = '';
let userNameSet = false;

let joinTrigger = new Subject<void>();
let joinError = joinTrigger.pipe(
  switchMap(() =>
    rpc.joinRoom(roomKey, userName).pipe(
      tap((data) => store.set.game(data)),
      map(() => undefined),
      catchError((err) => of(err)),
    ),
  ),
);

onMount(() => {
  userName = localStorage.getItem('userName');
  userNameSet = !!userName;

  const pathname = window.location.pathname;
  roomKey = pathname.substring(pathname.lastIndexOf('/') + 1);
});

function handleUsernameChange(newName: CustomEvent<string>) {
  userName = newName.detail;
  localStorage.setItem('userName', userName);
}

function createRoom() {
  rpc.createRoom(userName).subscribe((data) => {
    store.set.game(data);
    window.history.replaceState(null, '', window.location + 'room/' + data.roomKey);
  });
}

async function joinRoom() {
  joinTrigger.next(undefined);
}
</script>

<CoverImage>
  <div class="header-container">
    <h1>City Knowledge Contest</h1>
    <div>Wer findet die gesuchten Orte am schnellsten?</div>
  </div>
  <Card>
    <span class="save-name-notice"
      >Mit der Eingabe eines Pseudonyms erklärst du dich damit einverstanden, dass dieses in deinem Browser gespeichert
      wird. So musst du es beim nächsten Spiel nicht erneut eingeben.</span
    >
    <Input
      on:input="{handleUsernameChange}"
      placeholder="Wie heißt du?"
      value="{userName}"
      e2eTestId="user-name-input"
    />
    {#if roomKey.length < 1}
      <Button
        title="Neue Karte erstellen"
        on:click="{createRoom}"
        disabled="{!userName}"
        e2eTestId="create-room-button"
      />
    {:else}
      <Button
        on:click="{joinRoom}"
        title="Karte beitreten"
        disabled="{!userName || !roomKey}"
        e2eTestId="join-room-button"
      />
    {/if}
    {#if $joinError}
      <div>Du konntest dem Raum nicht beitreten. Prüfe den Raum-ID oder erstelle einen neuen Raum</div>
      <Button
        title="Neue Karte erstellen"
        on:click="{createRoom}"
        disabled="{!userName}"
        e2eTestId="create-room-button"
      />
    {/if}
  </Card>
</CoverImage>
