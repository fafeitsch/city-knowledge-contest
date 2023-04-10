const isSsl = window.location.protocol === 'https:';
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
