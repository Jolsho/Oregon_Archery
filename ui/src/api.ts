import type { Store } from "@tauri-apps/plugin-store";
import type { EventT } from "./eventT";

type GetRes = {
    events: EventT[]
    token: string
}

export const ROOT_URL   = "http://127.0.0.1:8080";

export async function get_events(cookie_store: Store): Promise<GetRes> {

    let url = "/events";
    const cookie = await cookie_store.get<string>("session_token");
    if (!!cookie) {
        url += `?token=${cookie}`
    }

    let res = await fetch(ROOT_URL + url, { method: "GET" });
    if (!res.ok) {
        console.error(res.text());
        throw res.text();
    } else {
        return await res.json();
    }
}

export async function post_event(cookie_store: Store, event: EventT): Promise<EventT> {
    let url = `/events?title=${encodeURIComponent(event.title)}`;

    const cookie = await cookie_store.get<string>("session_token");
    if (!!cookie) {
        url += `&token=${cookie}`
    }

    try {
        const res = await fetch(ROOT_URL + url, {
            method: "POST",
            body: JSON.stringify(event),
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

export async function delete_event(cookie_store: Store, title: string) {

    let url = ROOT_URL + `/events?title=${title}`;

    const cookie = await cookie_store.get<string>("session_token");
    if (!!cookie) {
        url += `&token=${cookie}`
    }

    let res = await fetch(url, { method: "DELETE" });
    if (!res.ok) {
        const text = await res.text();
        throw new Error(`DELETE /events failed (${res.status}): ${text}`);
    }
}
