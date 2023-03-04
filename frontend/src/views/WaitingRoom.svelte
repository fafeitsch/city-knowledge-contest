<script lang="ts" context="module">
  export type UpdateRoomParams = {
    listFileName: string;
    numberOfQuestions: number;
    playerKey: string;
    playerSecret: string;
    roomKey: string;
    maxAnswerTimeSec: number;
  };

  export const defaultRoomSeetings = {
    numberOfQuestions: 1,
    listFileName: "wuerzburg.json",
    maxAnswerTimeSec: 600,
  };
</script>

<script lang="ts">
  import AvailableStreetLists from "../components/AvailableStreetLists.svelte";
  import Button from "../components/Button.svelte";
  import CopyIcon from "../components/CopyIcon.svelte";
  import { handleRPCRequest } from "../rpc";
  import store, { GameState } from "../store";

  let game = store.get.game$;

  function startGame() {
    handleRPCRequest<UpdateRoomParams, null>({
      method: "updateRoom",
      params: {
        ...defaultRoomSeetings,
        playerKey: $game.playerKey,
        roomKey: $game.roomId,
        playerSecret: $game.playerSecret,
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
  <p class="mt-5">
    Teile den Code, um andere Personen zu diesem Spiel einzuladen:
  </p>
  <p class="fw-bold p-3 bg-old-map-lighter d-flex align-items-center gap-3">
    {$game.roomId}<CopyIcon
      width={16}
      height={16}
      className="color-black"
      on:click={() => {
        navigator.clipboard.writeText($game.roomId);
      }}
    />
  </p>
  <AvailableStreetLists />
</div>
