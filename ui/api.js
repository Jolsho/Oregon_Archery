import { Event } from "./event.js";

/**
 * @param {Event[]} events 
 *  @returns {Promise<void>}
 */
export async function run_websocket(events) {
    let ws = new WebSocket("/ws");

    ws.onopen = (event) => {
        console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        events.onMessage?.(data);
    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
        console.log("WebSocket closed");
    };
}

/**
 *  @returns {Promise<Event[]>}
 */
export async function get_events() {
    let res = await fetch("/events");
    if (!res.ok) {
        console.error(res.error);
        throw res.error;
    } else {
        return await res.json();
    }
}

/**
 *  @param {Event} event
 *  @returns {Promise<Event>}
 */
export async function post_event(event) {
    const url = `/events?title=${encodeURIComponent(event.title)}`;

    const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(event),
    });

    if (!res.ok) {
        const text = await res.text();
        throw new Error(`POST /events failed (${res.status}): ${text}`);
    }
    return await res.json();
}

/**
 *  @param {string} title
 *  @returns {Promise<void>}
 */
export async function delete_event(title) {
    let res = await fetch(`/events?title=${title}`,
        { method: "DELETE" },
    );
    if (!res.ok) {
        const text = await res.text();
        throw new Error(`DELETE /events failed (${res.status}): ${text}`);
    }
}
