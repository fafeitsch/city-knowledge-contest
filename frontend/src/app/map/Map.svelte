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
  z-index: 999;
  position: absolute;
  bottom: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 150px;
  width: 100vw;
}

.card {
  background-color: white;
  border-radius: 16px;
  width: 30%;
  padding: 16px;
  font-size: xx-large;
  font-family: 'LilitaOne';
  text-align: center;
}
</style>

<script lang="ts">
import store from '../../store';
import PartyConfetti from '../../components/PartyConfetti.svelte';
import Button from '../../components/Button.svelte';
import Leaflet from './Leaflet.svelte';
import { map, merge, of, Subject, switchMap, tap, zip } from 'rxjs';
import { subscribeToCountdown, subscribeToQuestion, subscribeToQuestionFinished } from '../../sockets';
import rpc from '../../rpc';
import GamePanel from './GamePanel.svelte';

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
let lastResult = merge(
  guess.pipe(
    switchMap((guess) => {
      if (!guess) {
        return of(undefined);
      }
      return rpc.answerQuestion(guess);
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
    <div class="overlay">{$countdown}</div>
  {/if}
  <Leaflet on:mapClicked="{onAnswerQuestion}" disabled="{$lastResult !== undefined}" />
  {#if $question && !$gameFinished && $lastResult === undefined}
    <div class="container">
      <div class="card">Suche den Ort {$question}</div>
    </div>
  {:else if $gameFinished || $lastResult !== undefined}
    <div class="container">
      <div class="card">
        <div class="mb-3">
          {#if $lastResult}
            <div>Richtig</div>
            <PartyConfetti />
          {:else}
            <div>Falsch</div>
          {/if}
        </div>
        {#if $gameFinished}
          <Button on:click="{() => advanceGame()}" title="Weiter" />
        {/if}
      </div>
    </div>
  {/if}
</div>
