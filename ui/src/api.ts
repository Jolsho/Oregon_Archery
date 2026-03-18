import type { Store } from "@tauri-apps/plugin-store";
import type { EventT } from "./eventT";

type GetRes = {
    events: EventT[]
    token: string
}

const DEFAULT_TIMEOUT   = 5000; // timeout in milliseconds
export const ROOT_URL = import.meta.env.PROD ? "https://testohsal.com" : "http://127.0.0.1:8080";

export async function get_events(
    cookie_store: Store,
    timeout = DEFAULT_TIMEOUT
): Promise<GetRes> {

    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), timeout);

    let url = "/events";
    try {
        const cookie = await cookie_store.get<string>("session_token");
        if (!!cookie) {
            url += `?token=${cookie}`
        }

        let res = await fetch(ROOT_URL + url, { 
            method: "GET", 
            signal: controller.signal
        });
        if (!res.ok) {
            console.error(res.text());
            throw res.text();
        } 

        return await res.json();

    } catch (err) {
        // distinguish timeout from other errors
        if ((err as any).name === "AbortError") {
            throw new Error(`POST /events timed out after ${timeout}ms`);
        }
        throw err;
    } finally {
        clearTimeout(id);
    }
}

export async function post_event(
    cookie_store: Store,
    event: EventT,
    timeout = DEFAULT_TIMEOUT
): Promise<EventT> {
    let url = `/events?title=${encodeURIComponent(event.title)}`;

    const cookie = await cookie_store.get<string>("session_token");
    if (cookie) {
        url += `&token=${cookie}`;
    }

    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), timeout);

    try {
        const res = await fetch(ROOT_URL + url, {
            method: "POST",
            body: JSON.stringify(event),
            signal: controller.signal,
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
        // distinguish timeout from other errors
        if ((err as any).name === "AbortError") {
            throw new Error(`POST /events timed out after ${timeout}ms`);
        }
        throw err;
    } finally {
        clearTimeout(id);
    }
}

export async function delete_event(
    cookie_store: Store, 
    title: string,
    timeout = DEFAULT_TIMEOUT
) {

    let url = ROOT_URL + `/events?title=${title}`;

    const cookie = await cookie_store.get<string>("session_token");
    if (!!cookie) {
        url += `&token=${cookie}`
    }

    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), timeout);

    try {
        let res = await fetch(url, { 
            method: "DELETE", 
            signal: controller.signal
        });
        if (!res.ok) {
            const text = await res.text();
            throw new Error(`DELETE /events failed (${res.status}): ${text}`);
        }
    } catch (err) {
        // distinguish timeout from other errors
        if ((err as any).name === "AbortError") {
            throw new Error(`POST /events timed out after ${timeout}ms`);
        }
        throw err;
    } finally {
        clearTimeout(id);
    }
}
