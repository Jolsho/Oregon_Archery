import { useCallback, useEffect, useRef, type Dispatch, type SetStateAction } from "react";
import type { EventT } from "./eventT";
import type { ConnectionStatus } from "./App";
import { insertSorted } from "./utils";
import type { Store } from "@tauri-apps/plugin-store";

type WSHandlers = {
    onOpen?: () => void;
    onMessage?: (data: any) => void;
    onClose?: (ev: CloseEvent) => void;
};

export const WS_URL = import.meta.env.PROD ? "wss://testohsal.com/ws" : "ws://127.0.0.1:8080/ws";

export type WSReturn = {
    send:       (data: any) => void;
    connect:    (cookie_store: Store, url: string) => Promise<void>;
    socket:     WebSocket | null;
}

export function useWebSocket(handlers: WSHandlers): WSReturn {
    const wsRef = useRef<WebSocket | null>(null);
    const retryRef = useRef(0);
    const handlersRef = useRef(handlers);

    // keep handlers updated without recreating socket
    useEffect(() => {
        handlersRef.current = handlers;
    }, [handlers]);

    const connect = useCallback(async (cookie_store: Store, url: string) => {
        try {
            const ws = new WebSocket(url);
            wsRef.current = ws;

            ws.onopen = () => {
                retryRef.current = 0;
                handlersRef.current.onOpen?.();
                console.log("CONNECTED");
            };

            ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    handlersRef.current.onMessage?.(data);
                } catch {
                    handlersRef.current.onMessage?.(event.data);
                }
            };

            ws.onclose = (ev) => {
                handlersRef.current.onClose?.(ev);
                retryRef.current++;
                setTimeout(() => connect(cookie_store, url), 4000);
            };

            ws.onerror = (err) => {
                console.error("WebSocket error:", err);
            };
        } catch (err) {
            console.error("WebSocket failed to connect:", err);
            retryRef.current++;
            setTimeout(() => connect(cookie_store, url), 4000);
        }
    }, []);

    useEffect(() => {
        return () => wsRef.current?.close();
    }, []);

    const send = useCallback((data: unknown) => {
        const ws = wsRef.current;
        if (!ws || ws.readyState !== WebSocket.OPEN) return;

        ws.send(JSON.stringify(data));
    }, []);

    return { send, connect, socket: wsRef.current };
}




export function useEventSync(
    setEvents: Dispatch<SetStateAction<EventT[]>>,
    setConnectionStatus: Dispatch<SetStateAction<ConnectionStatus>>,
): WSReturn {
    return useWebSocket({
        onOpen: () => {
            console.log("WebSocket connected");

            setConnectionStatus((prev) => ({
                ...prev,
                status: "Connected"
            }));
        },

        onMessage: (data) => {
            switch (data.msg) {
                case "new_event":
                    setEvents((prev) => {
                        const idx = prev.findIndex(
                            (v) => v.title === data.payload.event.title,
                        );
                        if (idx >= 0) {
                            if (prev[idx].is_own) return prev;

                            const events = [...prev];
                            events[idx] = data.payload.event;
                            return events;
                        }

                        const events = [...prev];
                        insertSorted(events, data.payload.event);
                        return events;
                    });
                    break;

                case "delete_event":
                    setEvents((prev) => {
                        const idx = prev.findIndex(
                            (v) => v.title === data.payload.title,
                        );

                        if (idx < 0) return prev;
                        if (prev[idx].is_own) return prev;

                        const events = [...prev];
                        events.splice(idx, 1);
                        return events;
                    });
                    break;
            }
        },

        onClose: (_) => {
            setConnectionStatus((prev) => ({
                ...prev,
                status: "Offline"
            }));
        },
    });
}
