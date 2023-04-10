const isSsl = window.location.protocol === 'https:';
console.log(window.location.host, window.location, window.location.hostname, window.location.port);
export const environment: Record<string, { wsUrl: string; apiUrl: string; tileUrl: string }> = {
  development: {
    wsUrl: 'ws://localhost:23123/ws/',
    apiUrl: 'http://localhost:23123/rpc',
    tileUrl: 'http://localhost:23123/tile/{z}/{x}/{y}',
  },
  production: {
    wsUrl: (isSsl ? 'wss://' : 'ws://') + window.location.host + `/ws/`,
    apiUrl: '/rpc',
    tileUrl: '/tile/{z}/{x}/{y}',
  },
};
