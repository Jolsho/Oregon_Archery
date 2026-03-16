import type { ConnectionStatus } from "./App.tsx";
import type { EventT } from "./eventT.ts";
import arrow_back from "./assets/arrow_back.svg"
import ohsal_img from "./assets/OHSAL.png"

export type MenuProps = {
    new_event: () => void;
    events: EventT[];
    menuIsOpen: boolean;
    setMenuIsOpen: (is: boolean) => void;
    connectionStatus: ConnectionStatus;
    setCurrentEvent: (idx: number) => void;
};
// Menu Component
export function Menu({ 
    events, new_event,
    menuIsOpen, setMenuIsOpen, 
    connectionStatus, 
    setCurrentEvent 
}: MenuProps) {

    const toggleMenu = () => setMenuIsOpen(!menuIsOpen);

    return (
        <div className="header">
        <img
            className={`scrollable_minimizer ${!menuIsOpen ? "flipped" : ""}`}
            src={arrow_back}
            id="menu_minimizer"
            onClick={toggleMenu}
        />

        <div className="monitor">
            <h1 className="roboto-mono-norm status">{connectionStatus.status}</h1>
        </div>

        <div className={`scrollable ${menuIsOpen ? "open" : "close"}`} id="menu">
            <div className="header_img_container">
                <img className="header_img" src={ohsal_img} />
            </div>

            {/* Create Event Button */}
            <div className="event_container create" onClick={() => {
                new_event();
                toggleMenu();
            }}>
                <h2 className="roboto-mono-norm">Create Event</h2>
            </div>

            {/* Existing Events */}
            {(!events || events.length === 0) ? (
                <div className="event_container no_events">
                    <h2 className="roboto-mono-norm">NO EVENTS</h2>
                </div>
            ) : (
                events.map((event, i) => (
                    <div className="event_container" 
                        key={i}
                        onClick={() => {
                            setCurrentEvent(i);
                            toggleMenu();
                        }}
                    >
                        <div className="live_container">
                            <div className="live_indicator"></div>
                        </div>
                        <h2 className="roboto-mono-norm">{event.title}</h2>
                    </div>
                ))
            )}
        </div>
        </div>
    );
}
