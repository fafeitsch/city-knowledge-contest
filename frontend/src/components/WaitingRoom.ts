import store from '../store';

export class WaitingRoom extends HTMLElement {
  constructor() {
    super();
  }
  connectedCallback() {
    const wrapper = document.createElement('div');
    wrapper.innerHTML = '<div>Lets play</div>';
    const playerList = document.createElement('div');
    playerList.setAttribute('id', 'player-list');
    store.get.players$.subscribe((players) => {
      while (playerList.firstChild) {
        playerList.removeChild(playerList.firstChild);
      }
      players.forEach((player) => {
        const playerNode = document.createElement('div');
        playerNode.innerText = player;
        playerList.appendChild(playerNode);
      });
    });
    wrapper.appendChild(playerList);
    super.appendChild(wrapper);
  }
}

customElements.define('waiting-room', WaitingRoom);
