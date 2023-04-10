<script lang="ts">
import L, { Icon, latLng, type LatLng, type Map, Marker } from 'leaflet';
import { filter, map, Subject, takeUntil } from 'rxjs';
import { createEventDispatcher, onDestroy, onMount } from 'svelte';
import { environment } from '../../environment';
import img from '../../assets/images/pin.png';
import { subscribeToQuestionFinished, subscribeToRoomUpdated, subscribeToSocketTopic, Topic } from '../../sockets';

let mapContainer: Map;

const markerIcon = new Icon({
  iconUrl: img,
  iconSize: [50, 50],
  iconAnchor: [25, 50],
});

let marker: Marker | undefined = undefined;
let destroy$ = new Subject<void>();
let dispatch = createEventDispatcher();

onMount(() => {
  mapContainer = createMap();
  subscribeToQuestionFinished()
    .pipe(
      takeUntil(destroy$),
      filter((result) => !!result),
    )
    .subscribe((value) => {
      if (marker !== undefined) {
        marker.removeFrom(mapContainer);
      }
      marker = new Marker(
        {
          lat: value.solution[0],
          lng: value.solution[1],
        },
        { icon: markerIcon },
      );
      mapContainer
        .flyTo(
          {
            lat: value.solution[0],
            lng: value.solution[1],
          },
          18,
        )
        .on('moveend', () => {
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

onDestroy(() => {
  destroy$.next(undefined);
  destroy$.complete();
});

function createMap() {
  const map = L.map('mapContainer').setView(latLng(50, 10), 5);

  L.tileLayer(environment[import.meta.env.MODE].tileUrl, {
    attribution: `&copy;<a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>,
	        &copy;<a href="https://carto.com/attributions" target="_blank">CARTO</a>`,
    subdomains: 'abcd',
    maxZoom: 20,
  }).addTo(map);

  subscribeToRoomUpdated()
    .pipe(takeUntil(destroy$))
    .subscribe((config) => {
      map.flyTo({ lat: config.center[0], lng: config.center[1] }, 16);
    });
  map.addEventListener('click', (e) => dispatch('answerQuestion', [e.latlng.lat, e.latlng.lng]));

  return map;
}
</script>

<div id="mapContainer" class="map full-viewheight full-viewwidth"></div>
