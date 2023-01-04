import './styles.scss';

import { initForm } from './login';
import store from './store';
import { combineLatest } from 'rxjs';

let app: HTMLDivElement;

export function initMain() {
  app = document.getElementById('app')! as HTMLDivElement;

  store.get.gameState$.subscribe((state) => {
    app.childNodes.forEach((child) => {
      app.removeChild(child);
    });

    if (state === 'waiting') {
      const waitingRoom = document.createElement('waiting-room');
      app.appendChild(waitingRoom);
    } else {
      initForm();
    }
  });
}

function initWebSocket() {
  combineLatest([
    store.get.roomId$,
    store.get.playerKey$,
    store.get.gameState$,
  ]).subscribe(([roomId, playerKey, gameState]) => {
    if (gameState !== 'waiting') {
      return;
    }
    const websocket = new WebSocket(
      'ws://127.0.0.1:23123/ws/' + roomId + '/' + playerKey
    );
    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.topic === 'successfullyJoined') {
        store.set.players(data.payload.players);
      } else if (data.topic === 'playerJoined') {
        store.set.addPlayer(data.payload.name);
      }
    };
  });
}

initMain();
initWebSocket();
