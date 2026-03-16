import { delete_event, get_events, post_event, useWebSocket } from "./api.ts";
import { useEffect, useState, type Dispatch, type SetStateAction } from "react";
import { Menu } from "./menu.tsx";
import { calculate_leaders, new_event, type EventT } from "./eventT.ts";

import { EventPage } from "./event.tsx";
import { insertSorted } from "./utils.ts";
export type ConnectionStatus = {
    retryCount: number;
    status: string;
    is_manual: boolean;
};

export function useEventSync(
    setEvents: Dispatch<SetStateAction<EventT[]>>,
    setConnectionStatus: Dispatch<SetStateAction<ConnectionStatus>>,
) {
    return useWebSocket("/ws", {
        onOpen: () => {
            console.log("WebSocket connected");

            setConnectionStatus((prev) => ({
                ...prev,
                status: "Syncing...",
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
                            (v) => v.title === data.payload.event.title,
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
                status: "Reconnecting..."
            }));
        },
    });
}

function App() {

    const [events, setEvents] = useState<EventT[]>([]);
    const [curr_idx, setCurrEvent] = useState<number>(events.length - 1);
    const [menu_open, setMenuOpen] = useState<boolean>(true);

    const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>({
        status: "Connected",
        retryCount: 0,
        is_manual: false,
    });
    const {send} = useEventSync(setEvents, setConnectionStatus);


    useEffect(() => {
        if (connectionStatus.status == "Syncing...") {
            setEvents(prev => {
                const msg = "new_event";
                let evs = [...prev];
                evs.forEach(event => {
                    if (!event.is_persisted && event.is_own) {
                        send({msg, payload: {event}});
                        event.is_persisted = true;
                        calculate_leaders(event);
                        console.log("SENT", event.title);
                    }
                });
                return evs;
            })
            setConnectionStatus({...connectionStatus, status: "Connected"});
        }
    }, [connectionStatus])


    useEffect(() => {
        get_events().then((es) => {
            es.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
            es.forEach(e => {
                if (!(e.leaders instanceof Map)) {
                    e.leaders = new Map(Object.entries(e.leaders ?? {}));
                }
                calculate_leaders(e);
            });
            setEvents(es);
            setCurrEvent(es.length - 1);
        });
    }, []);

    async function submit_event(ev: EventT): Promise<string> {
        let msg = "";
        try {
            await post_event(ev);
            ev.is_persisted = true;
        } catch (e) {
            if (e instanceof Error) msg = e.message;
            ev.is_persisted = false;
        }

        setEvents((prev) => {
            let evs = [...prev];
            evs[curr_idx] = ev;
            return evs;
        });

        return msg;
    }

    return (
        <>
            <Menu
                events={events}
                new_event={() => {
                    let evs = [...events, new_event()];
                    setEvents(evs);
                    setCurrEvent(evs.length - 1);
                }}
                menuIsOpen={menu_open}
                setMenuIsOpen={setMenuOpen}
                connectionStatus={connectionStatus}
                setCurrentEvent={setCurrEvent}
            />
            {(events.length > curr_idx && curr_idx > -1) ? (
                <EventPage
                    events={events}
                    idx={curr_idx}
                    post_event={submit_event}
                    remove_event={() => {
                        let evs = [...events];
                        const title = evs[curr_idx].title;
                        evs.splice(curr_idx, 1);
                        setEvents(evs);
                        setCurrEvent(curr_idx - 1);
                        if (!!title) {
                            delete_event(title).catch((e) => console.error(e));
                        }
                    }}
                />
            ) : (
                <> </>
            )}
        </>
    );
}

export default App;
