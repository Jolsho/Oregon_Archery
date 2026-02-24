import { get_events, run_websocket } from "./api.js";
import {Event, render_event} from "./event.js"
import { render } from "./initial.js";

/** @type {Event[]} */
let events = [];

get_events().then((es) => {
    events = es
    render(events);
    render_event(events, 0);
    run_websocket(events);
});

