import { useEffect, useState, useRef } from 'react';
import useAuthStore from './useAuth';

export function useLiveMetrics() {
  const [metrics, setMetrics] = useState<any>(null);
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const { accessToken } = useAuthStore();

  useEffect(() => {
    if (!accessToken) return;

    let reconnectInterval: NodeJS.Timeout;

    const connect = () => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = process.env.NEXT_PUBLIC_API_URL?.replace(/^http(s?):\/\//, '') || window.location.host;
      const wsUrl = `${protocol}//${host}/ws/v1/stream?token=${accessToken}`;

      wsRef.current = new WebSocket(wsUrl);

      wsRef.current.onopen = () => {
        setIsConnected(true);
        console.log('Telemetry WebSocket connected');
      };

      wsRef.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          setMetrics(data);
        } catch (err) {
          console.error('Failed to parse websocket message', err);
        }
      };

      wsRef.current.onclose = () => {
        setIsConnected(false);
        console.log('Telemetry WebSocket disconnected. Attempting to reconnect...');
        reconnectInterval = setTimeout(connect, 3000);
      };

      wsRef.current.onerror = (err) => {
        console.error('Telemetry WebSocket error', err);
        wsRef.current?.close();
      };
    };

    connect();

    return () => {
      if (reconnectInterval) clearTimeout(reconnectInterval);
      if (wsRef.current) {
        wsRef.current.onclose = null;
        wsRef.current.close();
      }
    };
  }, [accessToken]);

  return { metrics, isConnected };
}
