import { delete_event, get_events, post_event} from "./api.ts";
import { useEffect,  useState } from "react";
import { Menu } from "./menu.tsx";
import { calculate_leaders, new_event, type EventT } from "./eventT.ts";

import { EventPage } from "./event.tsx";
import { load, Store } from "@tauri-apps/plugin-store";
import { useEventSync, WS_URL } from "./ws.tsx";
export type ConnectionStatus = {
    retryCount: number;
    status: string;
    is_manual: boolean;
};

function App() {

    const [events, setEvents] = useState<EventT[]>([]);
    const [curr_idx, setCurrEvent] = useState<number>(events.length - 1);
    const [menu_open, setMenuOpen] = useState<boolean>(true);

    const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>({
        status: "Offline",
        retryCount: 0,
        is_manual: false,
    });

    const [store, setStore] = useState<Store | null>(null);

    const {send, connect, socket} = useEventSync(setEvents, setConnectionStatus);


    useEffect(() => {
        if (connectionStatus.status == "Connected" && store) {
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
        }
    }, [connectionStatus, store])

    useEffect(() => {
        if (store) {
            get_events(store).then((res) => {
                res.events.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
                res.events.forEach(e => {
                    if (!(e.leaders instanceof Map)) {
                        e.leaders = new Map(Object.entries(e.leaders ?? {}));
                    }
                    calculate_leaders(e);
                });
                setEvents(res.events);
                setCurrEvent(res.events.length - 1);

                console.log("RECIEVED_COOKIE", res.token);
                store.set("session_token", res.token);

                if (!socket && store) {
                    connect(store, WS_URL + `?token=${res.token}`);
                }
            }).catch(e => console.error("GET EVENTS",e));
        }

    }, [store]);

    useEffect(() => {
        if (!store) {
            async function load_store() {
                setStore(await load('store.json'));
            }
            load_store();
            console.log("LOADED STORE");
        }
    }, [])

    function submit_event(ev: EventT) {
        ev.is_persisted = false;
        setEvents((prev) => {
            let evs = [...prev];
            evs[curr_idx] = ev;
            return evs;
        });

        if (!!store) {
            post_event(store, ev)
            .then(() => {
                ev.is_persisted = true;
                setEvents((prev) => {
                    let evs = [...prev];
                    evs[curr_idx] = ev;
                    return evs;
                });
            })
            .catch((e) => console.error(e));
        }
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
                            if (!store) return;
                            delete_event(store, title)
                                .catch((e) => console.error(e));
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
