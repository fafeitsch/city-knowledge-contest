<style lang="scss">
@import '../../styles/variables';

.game-panel {
  z-index: 1000;
  margin: 32px;
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
  max-height: 400px;
  background-color: white;
  border-radius: 16px;
  padding: 8px;
  overflow-x: hidden;
  overflow-y: auto;
}

.overlay {
  z-index: 10000;
  height: 100vh;
  width: 100vw;
  position: absolute;
  top: 0;
  left: 0;
  background-color: rgba($color: #000000, $alpha: 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 128px;
}

.container {
  z-index: 800;
  position: absolute;
  bottom: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80px;
  left: 50%;
  transform: translateX(-50%);
}

.card {
  background-color: white;
  border-radius: 16px;
  width: 30%;
  padding: 16px;
  font-size: xx-large;
  font-family: 'LilitaOne', sans-serif;
  text-align: center;
}
</style>

<script lang="ts">
import PartyConfetti from '../../components/PartyConfetti.svelte';
import Button from '../../components/Button.svelte';
import Leaflet from './Leaflet.svelte';
import { delay, map, merge, Observable, of, startWith, Subject, switchMap, tap, zip } from 'rxjs';
import { subscribeToCountdown, subscribeToQuestion, subscribeToQuestionFinished } from '../../sockets';
import rpc from '../../rpc';
import GamePanel from './GamePanel.svelte';
import LeaveButton from '../../components/LeaveButton.svelte';

let countdown = merge(
  subscribeToQuestion().pipe(map(() => undefined)),
  subscribeToCountdown().pipe(map((data) => (data ? data.followUps + 1 : undefined))),
);
let question = merge(countdown.pipe(map(() => undefined)), subscribeToQuestion().pipe(map((data) => data?.find)));
let gameFinished = merge(
  countdown.pipe(map(() => undefined)),
  question.pipe(map(() => undefined)),
  subscribeToQuestionFinished().pipe(map((data) => (data ? 'questionFinished' : undefined))),
);

let guess = new Subject<[number, number] | undefined>();
let lastResult: Observable<number | undefined | 'pending'> = merge(
  guess.pipe(
    switchMap((guess) => {
      if (!guess) {
        return of(undefined);
      }
      return rpc.answerQuestion(guess).pipe(startWith('pending'));
    }),
  ),
  question.pipe(map(() => undefined)),
);

function advanceGame() {
  guess.next(undefined);
  rpc.advanceRoom().subscribe();
}

function onAnswerQuestion(event: CustomEvent) {
  guess.next(event.detail);
}
</script>

<div>
  <div class="game-panel">
    <GamePanel />
  </div>
  {#if $countdown}
    <div class="overlay" data-testid="countdown-overlay">{$countdown}</div>
  {/if}
  <Leaflet on:mapClicked="{onAnswerQuestion}" disabled="{$lastResult !== undefined}" />
  <LeaveButton />
  {#if $question && !$gameFinished && $lastResult === undefined}
    <div class="card container" data-testid="question-card">Suche den Ort {$question}</div>
  {:else if $gameFinished || $lastResult !== undefined}
    <div class="container card d-flex flex-column gap-3" data-testid="game-state-card">
      {#if $lastResult > 0}
        Richtig
      {:else if $lastResult === 0}
        Falsch
      {:else if $lastResult === 'pending'}
        Antwort wird geprüft …
      {/if}
      {#if $gameFinished}
        <Button on:click="{() => advanceGame()}" title="Weiter" e2eTestId="proceed-game-button" />
      {/if}
    </div>
  {/if}
</div>
{#if $lastResult > 0}
  <PartyConfetti />
{/if}
