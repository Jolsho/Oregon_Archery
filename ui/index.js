import { get_events, run_websocket } from "./src/api.js";
import {render_event} from "./src/event.js"
import { render_menu } from "./src/menu.js";
import { State } from "./src/state.js";

let state = new State();

get_events().then((es) => {
    state.events = es
    render_menu(state);
    render_event(state);
    run_websocket(state);
});
