<script lang="ts">
import L, { Icon, latLng, type LatLng, LeafletMouseEvent, type Map, Marker, TileLayer } from 'leaflet';
import { delay, filter, map, merge, of, Subject, switchMap, take, takeUntil, tap } from 'rxjs';
import { createEventDispatcher, onDestroy, onMount } from 'svelte';
import { environment } from '../../environment';
import img from '../../assets/images/pin.png';
import { subscribeToQuestionFinished, subscribeToRoomUpdated, subscribeToSuccessfullyJoined } from '../../sockets';
import store from '../../store';

export let disabled = false;

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
});

onDestroy(() => {
  destroy$.next(undefined);
  destroy$.complete();
  mapContainer.remove();
  mapContainer = null;
});

function createMap() {
  const leafletMap = L.map('mapContainer').setView(latLng(50, 10), 5);
  let layer: TileLayer | undefined = undefined;
  store.get.room$.pipe(takeUntil(destroy$)).subscribe((room) => {
    if (layer) {
      layer.removeFrom(leafletMap);
    }
    console.log('setting up tile layer');
    layer = L.tileLayer(environment[import.meta.env.MODE].tileUrl.replace('roomKey', room.roomKey), {
      attribution: `&copy;<a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a>`,
    }).addTo(leafletMap);
  });
  subscribeToRoomUpdated()
    .pipe(
      switchMap((data) => {
        if (!data) {
          return subscribeToSuccessfullyJoined().pipe(
            map((payload) => payload.options),
            take(1),
          );
        }
        return of(data);
      }),
      takeUntil(destroy$),
      filter((data) => !!data),
    )
    .subscribe((config) => {
      leafletMap.invalidateSize();
      if (config?.boundingBox) {
        leafletMap.setMaxBounds(config.boundingBox);
        leafletMap.options.maxBoundsViscosity = 1;
      }
      leafletMap.setMinZoom(config.minZoom);
      leafletMap.setMaxZoom(config.maxZoom);
      leafletMap.flyTo({ lat: config.center[0], lng: config.center[1] }, config.maxZoom / 2 + config.minZoom / 2);
    });
  leafletMap.addEventListener('click', (e) => onMapClicked(e));

  return leafletMap;
}

function onMapClicked(event: LeafletMouseEvent) {
  if (disabled) {
    return;
  }
  dispatch('mapClicked', [event.latlng.lat, event.latlng.lng]);
}
</script>

<div id="mapContainer" class="map full-viewheight full-viewwidth"></div>
