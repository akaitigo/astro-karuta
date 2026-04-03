"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type { WSMessage, WSMessageType } from "@/types/game";

export type ConnectionStatus = "connecting" | "connected" | "disconnected";

type MessageHandler = (message: WSMessage) => void;

interface UseWebSocketOptions {
  url: string;
  onMessage?: MessageHandler;
  autoConnect?: boolean;
  maxReconnectAttempts?: number;
}

interface UseWebSocketReturn {
  status: ConnectionStatus;
  send: (type: WSMessageType, payload: unknown) => void;
  connect: () => void;
  disconnect: () => void;
}

const BASE_DELAY_MS = 1000;
const MAX_DELAY_MS = 30000;

function getBackoffDelay(attempt: number): number {
  const delay = Math.min(BASE_DELAY_MS * Math.pow(2, attempt), MAX_DELAY_MS);
  // Add jitter (0-25% of delay)
  return delay + Math.random() * delay * 0.25;
}

export function useWebSocket({
  url,
  onMessage,
  autoConnect = true,
  maxReconnectAttempts = 5,
}: UseWebSocketOptions): UseWebSocketReturn {
  const [status, setStatus] = useState<ConnectionStatus>("disconnected");
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttemptRef = useRef(0);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const onMessageRef = useRef(onMessage);
  const intentionalCloseRef = useRef(false);
  const urlRef = useRef(url);
  const maxReconnectRef = useRef(maxReconnectAttempts);

  // Keep refs current
  useEffect(() => {
    onMessageRef.current = onMessage;
  }, [onMessage]);

  useEffect(() => {
    urlRef.current = url;
  }, [url]);

  useEffect(() => {
    maxReconnectRef.current = maxReconnectAttempts;
  }, [maxReconnectAttempts]);

  const clearReconnectTimer = useCallback(() => {
    if (reconnectTimerRef.current !== null) {
      clearTimeout(reconnectTimerRef.current);
      reconnectTimerRef.current = null;
    }
  }, []);

  // Use a ref-based connect to avoid circular dependency in useCallback
  const connectImplRef = useRef<() => void>(() => {});

  connectImplRef.current = () => {
    // Close existing connection
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    intentionalCloseRef.current = false;
    reconnectAttemptRef.current = 0;
    setStatus("connecting");

    const ws = new WebSocket(urlRef.current);
    wsRef.current = ws;

    ws.onopen = () => {
      setStatus("connected");
      reconnectAttemptRef.current = 0;
    };

    ws.onmessage = (event: MessageEvent) => {
      try {
        const data: unknown = JSON.parse(String(event.data));
        if (
          typeof data === "object" &&
          data !== null &&
          "type" in data &&
          "payload" in data
        ) {
          const message = data as WSMessage;
          onMessageRef.current?.(message);
        }
      } catch {
        // Ignore malformed messages
      }
    };

    ws.onclose = () => {
      wsRef.current = null;
      setStatus("disconnected");

      if (
        !intentionalCloseRef.current &&
        reconnectAttemptRef.current < maxReconnectRef.current
      ) {
        const delay = getBackoffDelay(reconnectAttemptRef.current);
        reconnectAttemptRef.current += 1;
        reconnectTimerRef.current = setTimeout(() => {
          connectImplRef.current();
        }, delay);
      }
    };

    ws.onerror = () => {
      // onclose will fire after onerror, reconnection is handled there
    };
  };

  const connect = useCallback(() => {
    connectImplRef.current();
  }, []);

  const disconnect = useCallback(() => {
    intentionalCloseRef.current = true;
    clearReconnectTimer();
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setStatus("disconnected");
  }, [clearReconnectTimer]);

  const send = useCallback(
    (type: WSMessageType, payload: unknown) => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type, payload }));
      }
    },
    [],
  );

  useEffect(() => {
    if (autoConnect) {
      connect();
    }
    return () => {
      intentionalCloseRef.current = true;
      clearReconnectTimer();
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return { status, send, connect, disconnect };
}
