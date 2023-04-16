import store from './store';
import { BehaviorSubject, Observable } from 'rxjs';
import { environment } from './environment';
import type { RoomConfiguration } from './rpc';

export enum Topic {
  'roomUpdated' = 'roomUpdated',
  'question' = 'question',
  'successfullyJoined' = 'successfullyJoined',
  'questionCountdown' = 'questionCountdown',
  'gameEnded' = 'gameEnded',
  'questionFinished' = 'questionFinished',
}

let websocket: WebSocket | undefined = undefined;
let subscriptions: Record<string, BehaviorSubject<any>> = {};

export function initWebSocket() {
  store.get.room$.subscribe((room) => {
    Object.values(subscriptions).forEach((subscription) => subscription.next(undefined));
    if (websocket !== undefined) {
      websocket.close();
    }
    if (!room) {
      return;
    }
    websocket = new WebSocket(environment[import.meta.env.MODE].wsUrl + room.roomKey + '/' + room.playerKey);
    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (subscriptions[data.topic]) {
        subscriptions[data.topic].next(data.payload);
      }
      if (data.topic === 'successfullyJoined') {
        store.set.players(data.payload.players);
      } else if (data.topic === 'playerJoined') {
        store.set.addPlayer(data.payload);
      } else if (data.topic === 'questionFinished') {
        store.set.updatePlayerRanking(data.payload.points);
      } else if (data.topic === 'playerAnswered') {
        store.set.updatePlayerDelta(data.payload);
      } else if (data.topic === 'question' || data.topic === 'questionCountdown') {
        store.set.removePlayerDelta();
      }
    };
  });
}

export function subscribeToQuestion(): Observable<{ find: string }> {
  return subscribeToSocketTopic<{ find: string }>(Topic.question);
}

export function subscribeToCountdown(): Observable<{ followUps: number }> {
  return subscribeToSocketTopic<{ followUps: number }>(Topic.questionCountdown);
}

export function subscribeToJoined(): Observable<{ options: RoomConfigurationResult }> {
  return subscribeToSocketTopic<{ options: RoomConfigurationResult }>(Topic.successfullyJoined);
}

export function subscribeToRoomUpdated(): Observable<RoomConfigurationResult> {
  return subscribeToSocketTopic<RoomConfigurationResult>(Topic.roomUpdated);
}

export function subscribeToSuccessfullyJoined(): Observable<SuccessfulJoinResult> {
  return subscribeToSocketTopic(Topic.successfullyJoined);
}

export type SuccessfulJoinResult = { options: RoomConfigurationResult; started: boolean };

export type RoomConfigurationResult = RoomConfiguration & {
  errors: string[];
  center: [number, number];
  minZoom: number;
  maxZoom: number;
  boundingBox?: [[number, number], [number, number]];
};

export type GameResult = {
  delta: Record<string, number>;
  points: Record<string, number>;
  question: string;
  solution: {
    lat: number;
    lon: number;
  };
};

export function subscribeToQuestionFinished(): Observable<GameResult> {
  return subscribeToSocketTopic<GameResult>(Topic.questionFinished);
}

export function subscribeToGameEnded(): Observable<unknown> {
  return subscribeToSocketTopic<unknown>(Topic.gameEnded);
}

export function subscribeToSocketTopic<T>(topic: Topic): BehaviorSubject<T | undefined> {
  if (!subscriptions[topic]) {
    subscriptions[topic] = new BehaviorSubject(undefined);
  }
  return subscriptions[topic];
}
