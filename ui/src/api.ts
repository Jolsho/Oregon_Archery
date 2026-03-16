import type { EventT } from "./eventT";
import { useCallback, useEffect, useRef } from "react";

type WSHandlers = {
    onOpen?: () => void;
    onMessage?: (data: any) => void;
    onClose?: (ev: CloseEvent) => void;
};

export function useWebSocket(url: string, handlers: WSHandlers) {
    const wsRef = useRef<WebSocket | null>(null);
    const retryRef = useRef(0);
    const handlersRef = useRef(handlers);

    // keep handlers updated without recreating socket
    useEffect(() => {
        handlersRef.current = handlers;
    }, [handlers]);

    const connect = useCallback(() => {
        const ws = new WebSocket(url);
        wsRef.current = ws;

        ws.onopen = () => {
            retryRef.current = 0;
            handlersRef.current.onOpen?.();
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
            setTimeout(connect, 4000);
        };
    }, [url]);

    useEffect(() => {
        connect();
        return () => wsRef.current?.close();
    }, [connect]);

    const send = useCallback((data: unknown) => {
        const ws = wsRef.current;
        if (!ws || ws.readyState !== WebSocket.OPEN) return;

        ws.send(JSON.stringify(data));
    }, []);

    return { send, socket: wsRef.current };
}

export async function get_events(): Promise<EventT[]> {
    let res = await fetch("/events", {
        method: "GET",
        credentials: "include",
    });
    if (!res.ok) {
        console.error(res.text());
        throw res.text();
    } else {
        return await res.json();
    }
}

export async function post_event(event: EventT): Promise<EventT> {
    const url = `/events?title=${encodeURIComponent(event.title)}`;

    try {
        const res = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(event),
            credentials: "include",
        });

        if (!res.ok) {
            const text = await res.text();

            const err = new Error(`POST /events failed (${res.status}): ${text}`);
            err.name = "http";
            err.cause = res.status;
            throw err;
        }

        return (await res.json()) as EventT;
    } catch (err) {
        throw err;
    }
}

export async function delete_event(title: string) {
    let res = await fetch(`/events?title=${title}`, {
        method: "DELETE",
        credentials: "include",
    });
    if (!res.ok) {
        const text = await res.text();
        throw new Error(`DELETE /events failed (${res.status}): ${text}`);
    }
}
