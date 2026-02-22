import { Event } from "./event.js";

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
