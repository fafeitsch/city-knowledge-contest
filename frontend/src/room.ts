let app: HTMLDivElement;

export function initRoom() {
  app = document.getElementById('app')! as HTMLDivElement;

  const room = document.createElement('div');
  room.innerText = `Let's play`;

  app.append(room);
}
