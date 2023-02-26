<script lang="ts">
  import Button from "../components/Button.svelte";
  import { handleRPCRequest } from "../rpc";
  import store, { GameState, type Game, type Player } from "../store";

  let game = store.get.game$;

  function startGame() {
    handleRPCRequest<
      {
        listFileName: string;
        numberOfQuestions: number;
        playerKey: string;
        playerSecret: string;
        roomKey: string;
        maxAnswerTimeSec: number;
      },
      null
    >({
      method: "updateRoom",
      params: {
        listFileName: "wuerzburg.json",
        numberOfQuestions: 10,
        playerKey: $game.playerKey,
        roomKey: $game.roomId,
        playerSecret: $game.playerSecret,
        maxAnswerTimeSec: 600,
      },
    }).then(() => {
      handleRPCRequest<
        { playerKey: string; playerSecret: string; roomKey: string },
        null
      >({
        method: "startGame",
        params: {
          playerKey: $game.playerKey,
          roomKey: $game.roomId,
          playerSecret: $game.playerSecret,
        },
      }).then(() => {
        store.set.gameState(GameState.Started);
      });
    });
  }
</script>

<div class="d-flex flex-column align-items-center gap-5">
  <div class="old-font fs-large">Gleich geht das Spiel los …</div>
  <Button title="Spiel starten" on:click={startGame} />
</div>
