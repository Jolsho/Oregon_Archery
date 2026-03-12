import { get_events, run_websocket } from "./src/api.js";
import {render_event} from "./src/event.js"
import { render_menu } from "./src/menu.js";
import { State } from "./src/state.js";

let state = new State();

get_events().then((es) => {
    state.events = es
    state.events.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    render_menu(state);
    render_event(state);
    run_websocket(state);
});
