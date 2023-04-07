<style lang="scss">
@import '../../styles/variables';

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
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 150px;
  width: 100vw;
  background-color: $beige;
  font-size: xx-large;
}
</style>

<script lang="ts">
import store from '../../store';
import PartyConfetti from '../../components/PartyConfetti.svelte';
import Button from '../../components/Button.svelte';
import Leaflet from './Leaflet.svelte';
import Players from '../../components/Players.svelte';
import { map, merge, of, Subject, switchMap, tap } from 'rxjs';
import {
  subscribeToCountdown,
  subscribeToQuestion,
  subscribeToQuestionFinished,
  subscribeToSocketTopic,
  Topic,
} from '../../sockets';
import rpc from '../../rpc';

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
let lastResult = guess.pipe(
  switchMap((guess) => {
    if (!guess) {
      return of(undefined);
    }
    return rpc.answerQuestion(guess);
  }),
);
let players = store.get.players$;

function advanceGame() {
  guess.next(undefined);
  rpc.advanceRoom().subscribe();
}

function onAnswerQuestion(event: CustomEvent) {
  guess.next(event.detail);
}
</script>

<div>
  <Players players="{$players}" absolutePosition="true" />
  {#if $countdown}
    <div class="overlay">{$countdown}</div>
  {/if}
  <Leaflet on:answerQuestion="{onAnswerQuestion}" />
  {#if $question && $lastResult === undefined}
    <div class="container">
      <div>Suche den Ort {$question}</div>
    </div>
  {:else if $lastResult !== undefined}
    <div class="container">
      <div class="d-flex justify-content-spaced align-items-center width-100 p-4">
        {#if $lastResult !== 0}
          <div id="party">Richtig 🥳</div>
          <PartyConfetti />
        {:else}
          <div>Leider falsch 🤷</div>
        {/if}
        {#if $gameFinished}
          <Button on:click="{() => advanceGame()}" title="Weiter" />
        {/if}
      </div>
    </div>
  {/if}
</div>
