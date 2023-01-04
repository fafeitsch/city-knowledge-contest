import './styles.scss';

import { initRoom } from './room';
import { initForm } from './login';
import store from './store';

let app: HTMLDivElement;

export function initMain() {
  app = document.getElementById('app')! as HTMLDivElement;

  store.get.roomId$.subscribe((roomId) => {
    app.childNodes.forEach((child) => {
      app.removeChild(child);
    });
    if (roomId) {
      initRoom();
    } else {
      initForm();
    }
  });
}

initMain();
