import store from './store';

let app: HTMLDivElement;

type CreateRoomResult = {
  key: string;
};

type JoinRoomResult = {
  key: string;
  name: string;
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

  form.appendChild(title);
  form.appendChild(addEnterRoomContainer());
  form.appendChild(addCreateRoomContainer());

  app.append(form);
}

function addCreateRoomContainer() {
  function createRoom() {
    fetch('http://localhost:23123/rpc', {
      method: 'POST',
      body: JSON.stringify({
        method: 'createRoom',
      }),
    })
      .then((response) => response.json())
      .then((data: Response<CreateRoomResult>) =>
        store.set.roomId(data.result.key)
      );
  }

  const button = document.createElement('button');
  button.textContent = 'Neue Karte erstellen';
  button.addEventListener('click', () => createRoom());
  return button;
}

function addEnterRoomContainer() {
  function enterRoom(roomId: string) {
    fetch('http://localhost:23123/rpc', {
      method: 'POST',
      body: JSON.stringify({
        method: 'joinRoom',
        params: {
          name: 'foo',
          roomKey: roomId,
        },
      }),
    })
      .then((response) => response.json())
      .then((data: Response<JoinRoomResult>) => {
        store.set.roomId(data.result.key);
        store.set.username(data.result.name);
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
    enterRoom(roomId);
  });
  container.appendChild(input);
  container.appendChild(button);
  return container;
}
