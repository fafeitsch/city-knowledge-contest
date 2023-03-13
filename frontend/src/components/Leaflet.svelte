<script lang="ts">
  import L, {Icon, latLng, type LatLng, type Map, Marker} from "leaflet";
  import {filter} from "rxjs";
  import {onMount} from "svelte";
  import {environment} from "../environment";
  import {handleRPCRequest} from "../rpc";
  import store, {type Game, GameState} from "../store";
  import img from "../assets/images/pin.png";

  export let game: Game;

  let mapContainer: Map;

  type AnswerQuestion = {
    points: number;
  };

  const markerIcon = new Icon({
    iconUrl: img,
    iconSize: [50, 50],
    iconAnchor: [25, 50],
  });

  let marker: Marker;

  onMount(() => {
    mapContainer = createMap();
    store.get.gameResult$
      .pipe(filter((result) => !!result))
      .subscribe((value) => {
        if (marker !== undefined) {
          marker.removeFrom(mapContainer);
        }
        marker = new Marker(
          {
            lat: value.solution[0],
            lng: value.solution[1],
          },
          {icon: markerIcon}
        );
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
    handleRPCRequest<AnswerQuestion>(
      "answerQuestion",
      {
        playerKey: game.playerKey,
        roomKey: game.roomKey,
        playerSecret: game.playerSecret,
        guess: [guess.lat, guess.lng],
      },
    ).then(() => {
      store.set.gameState(GameState.Finished);
    });
  }

  function createMap() {
    const map = L.map("mapContainer").setView(latLng(50, 10), 5);

    L.tileLayer(environment[import.meta.env.MODE].tileUrl, {
      attribution: `&copy;<a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>,
	        &copy;<a href="https://carto.com/attributions" target="_blank">CARTO</a>`,
      subdomains: "abcd",
      maxZoom: 20,
    }).addTo(map);

    map.flyTo({lat: 49.79465390310462, lng: 9.929384801847446}, 16);

    map.addEventListener("click", (e) => answerQuestion(e.latlng));

    return map;
  }
</script>

<div id="mapContainer" class="map full-viewheight full-viewwidth"/>
