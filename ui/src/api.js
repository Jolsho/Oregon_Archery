import { Event, render_event } from "./event.js";
import { State } from "./state.js";

/**
 * @param {State} state 
 *  @returns {Promise<void>}
 */
export async function run_websocket(state) {
    let ws = new WebSocket("/ws");

    ws.onopen = (_event) => {
        console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        switch (data.msg) {
        case "new_event": {
            let idx = state.find_event_idx(data.event.title);
            if (idx > 0) {

                // we already processed our own via http.
                if (state.get_event(idx).is_own) return;

                state.set_event(data.event, idx);

            } else {
                state.events.push(data.event);
            }

            render_menu(state);
            if (state.current_event_idx == idx) {
                render_event(state);
            }
            return;
        }

        case "delete_event": {
            let idx = state.find_event_idx(data.title);
            if (idx < 0) return;

            // we already processed our own via http.
            if (state.get_event(idx).is_own) return;
            state.remove_event(idx);

            render_menu(state);
            if (state.current_event_idx == idx) {
                render_event(state);
            }
            return;
        }
        }
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
