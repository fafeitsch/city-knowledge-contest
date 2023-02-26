<script lang="ts">
  import L, { latLng, Marker, type LatLng, type Map } from "leaflet";
  import { filter } from "rxjs";
  import { onMount } from "svelte";
  import { environment } from "../environment";
  import { handleRPCRequest } from "../rpc";
  import store, { GameState, type Game } from "../store";

  export let solution: { lat: number; lon: number } | undefined = undefined;
  export let currentResult: number;
  export let game: Game;

  let mapContainer: Map;

  type AnswerQuestion = {
    points: number;
  };

  onMount(() => {
    mapContainer = createMap();
    store.get.gameResult$
      .pipe(filter((result) => !!result))
      .subscribe((value) => {
        const marker = new Marker({
          lat: value.solution[0],
          lng: value.solution[1],
        });
        mapContainer
          .flyTo(
            {
              lat: value.solution[0],
              lng: value.solution[1],
            },
            18
          )
          .on("moveend", () => {
            marker.addTo(mapContainer);
          });
      });
    return {
      destroy: () => {
        mapContainer.remove();
        mapContainer = null;
      },
    };
  });

  function answerQuestion(guess: LatLng) {
    handleRPCRequest<
      {
        playerKey: string;
        roomKey: string;
        playerSecret: string;
        guess: Array<number>;
      },
      AnswerQuestion
    >({
      method: "answerQuestion",
      params: {
        playerKey: game.playerKey,
        roomKey: game.roomId,
        playerSecret: game.playerSecret,
        guess: [guess.lat, guess.lng],
      },
    }).then((data) => {
      currentResult = data.result.points;
      store.set.gameState(GameState.Finished);
    });
  }

  if (solution !== undefined) {
    console.log(solution);
  }

  function createMap() {
    const map = L.map("mapContainer").setView(latLng(50, 10), 5);

    L.tileLayer(environment[import.meta.env.MODE].tileUrl, {
      attribution: `&copy;<a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>,
	        &copy;<a href="https://carto.com/attributions" target="_blank">CARTO</a>`,
      subdomains: "abcd",
      maxZoom: 20,
    }).addTo(map);

    map.flyTo({ lat: 49.79465390310462, lng: 9.929384801847446 }, 16);

    map.addEventListener("click", (e) => answerQuestion(e.latlng));

    return map;
  }
</script>

<div id="mapContainer" class="map full-viewheight full-viewwidth" />
