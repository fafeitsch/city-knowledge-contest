<script lang="ts">
import Players from '../../components/Players.svelte';
import { filter, map, merge, tap } from 'rxjs';
import {
  subscribeToCountdown,
  subscribeToQuestion,
  subscribeToQuestionFinished,
  subscribeToRoomUpdated,
  subscribeToSuccessfullyJoined,
} from '../../sockets.ts';
import { store } from '../../store.ts';
import rpc from '../../rpc.js';

let players = store.get.players$.pipe(filter((players) => !!players));
let room = store.get.room$.pipe(filter((room) => !!room));
let currentQuestion = merge(
  subscribeToCountdown().pipe(
    filter((countdown) => !!countdown),
    map((countdown) => countdown.questionNumber),
  ),
  subscribeToQuestion().pipe(
    filter((question) => !!question),
    map((question) => question.questionNumber),
  ),
  subscribeToQuestionFinished().pipe(
    filter((result) => !!result),
    map((result) => result.questionNumber),
  ),
).pipe(map((question) => question + 1));
let totalQuestions = merge(
  subscribeToSuccessfullyJoined().pipe(
    filter((result) => !!result),
    map((result) => result.options),
  ),
  subscribeToRoomUpdated().pipe(
    filter((result) => !!result),
    map((result) => result),
  ),
).pipe(map((options) => options.numberOfQuestions));

function kickPlayer({ detail }: CustomEvent) {
  rpc.kickPlayer(detail).subscribe();
}
</script>

<div class="d-flex flex-column">
  {#if $currentQuestion !== undefined}
    <span class="pl-4 pr-4 pt-4 d-flex justify-content-end color-grey"
      >Frage {$currentQuestion} von {$totalQuestions}</span
    >
  {/if}
  <Players
    players="{$players}"
    playerKey="{$room.playerKey}"
    enableKick="{true}"
    on:kickPlayer="{(event) => kickPlayer(event)}"
  />
</div>
