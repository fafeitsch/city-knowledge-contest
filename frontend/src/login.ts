import { take } from 'rxjs';
import store from './store';

let app: HTMLDivElement;

type RoomResult = {
  roomKey: string;
  playerKey: string;
};

type Response<T> = {
  jsonrpc: string;
  result: T;
  id: null;
};

export function initForm() {
  app = document.getElementById('app')! as HTMLDivElement;

  const form = document.createElement('div');
  form.classList.add('d-flex', 'flex-column', 'gap-4', 'align-items-center');

  const title = document.createElement('h1');
  title.innerText = 'Guten Tag';

  const usernameInput = document.createElement('input');
  usernameInput.placeholder = 'Gib deinen Spielernamen ein …';
  usernameInput.addEventListener('blur', (event) => {
    store.set.username((event.target as any).value);
  });

  form.appendChild(title);
  form.appendChild(usernameInput);
  form.appendChild(addJoinRoomContainer());
  form.appendChild(addCreateRoomContainer());

  app.append(form);
}

function addCreateRoomContainer() {
  function createRoom(username: string) {
    fetch('http://localhost:23123/rpc', {
      method: 'POST',
      body: JSON.stringify({
        method: 'createRoom',
        params: {
          name: username,
        },
      }),
    })
      .then((response) => response.json())
      .then((data: Response<RoomResult>) => {
        store.set.roomId(data.result.roomKey);
        store.set.playerKey(data.result.playerKey);
        store.set.gameState('waiting');
      });
  }

  const button = document.createElement('button');
  button.textContent = 'Neue Karte erstellen';
  button.addEventListener('click', () => {
    store.get.username$
      .pipe(take(1))
      .subscribe((value) => createRoom(value ?? ''));
  });
  return button;
}

function addJoinRoomContainer() {
  function joinRoom(roomId: string, username: string) {
    fetch('http://localhost:23123/rpc', {
      method: 'POST',
      body: JSON.stringify({
        method: 'joinRoom',
        params: {
          name: username,
          roomKey: roomId,
        },
      }),
    })
      .then((response) => response.json())
      .then((data: Response<RoomResult>) => {
        store.set.roomId(roomId);
        store.set.playerKey(data.result.playerKey);
        store.set.gameState('waiting');
      });
  }

  const container = document.createElement('div');
  container.classList.add('d-flex', 'flex-column', 'gap-2');
  const input = document.createElement('input');
  input.placeholder = 'Karten-ID eingeben …';
  const button = document.createElement('button');
  button.textContent = 'Karte beitreten';
  button.addEventListener('click', () => {
    const roomId = input.value;
    store.get.username$
      .pipe(take(1))
      .subscribe((value) => joinRoom(roomId, value ?? ''));
  });
  container.appendChild(input);
  container.appendChild(button);
  return container;
}
