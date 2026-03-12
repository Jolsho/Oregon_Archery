import { Event, render_event } from "./event.js";
import { update_connection_status, render_menu } from "./menu.js";
import { State } from "./state.js";
import { insertSorted } from "./utils.js";

/**
 * @param {State} state
 *  @returns {Promise<void>}
 */
export async function run_websocket(state) {
    let ws = new WebSocket("/ws");

    ws.onopen = (_event) => {
        console.log("WebSocket connected");

        state.connection_status = "Syncing..."
        update_connection_status(state);

        const msg = "new_event";
        for (const event of state.events) {
            if (!event.is_persisted && event.is_own) {
                ws.send(JSON.stringify({msg, payload: {event}}));
                event.is_persisted = true;
                console.log("SENT", event.title);
            }
        }

        state.connection_status = "Connected";
        update_connection_status(state);
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        switch (data.msg) {
            case "new_event": {
                let idx = state.find_event_idx(data.payload.event.title);
                if (idx >= 0) {
                    // we already processed our own via http.
                    if (state.get_event(idx).is_own) return;

                    state.set_event(data.payload.event, idx);
                } else {
                    insertSorted(state.events, data.payload.event);
                }

                render_menu(state);
                if (state.current_event_idx == idx) {
                    render_event(state);
                }
                return;
            }

            case "delete_event": {
                let idx = state.find_event_idx(data.payload.title);
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

    ws.onclose = (ev) => {
        const retryable =
            ev.code === 1001 || // server restart
            ev.code === 1006 || // network issue
            ev.code >= 4000; // app-defined recoverable errors

        console.log("WebSocket closed", ev.code);

        state.connection_status = retryable
            ? "Reconnecting..."
            : "Export Unsaved Events and Refresh.";
        update_connection_status(state);

        if (retryable) {
            const RETRY_THRESHOLD = 2;
            if (state.retry_count > RETRY_THRESHOLD) {
                state.retry_count = 0;
                state.connection_status = "To Refresh Click -> "
                state.manual_status = true;
                update_connection_status(state)
            } else {
                state.retry_count++;
                setTimeout(() => run_websocket(state), 4000);
            }
        } else {
            console.error("Non-retryable close, manual intervention needed");
        }
    };
}

/**
 *  @returns {Promise<Event[]>}
 */
export async function get_events() {
    let res = await fetch("/events", {
        method: "GET",
        credentials: "include",
    });
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
      err.kind = "http";
      err.status = res.status;
      throw err;
    }

    return await res.json();
  } catch (err) {
    console.error(err);
  }
}

/**
 *  @param {string} title
 *  @returns {Promise<void>}
 */
export async function delete_event(title) {
    let res = await fetch(`/events?title=${title}`, { method: "DELETE", credentials: "include" });
    if (!res.ok) {
        const text = await res.text();
        throw new Error(`DELETE /events failed (${res.status}): ${text}`);
    }
}
