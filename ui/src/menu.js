import { run_websocket } from "./api.js";
import { render_event } from "./event.js";
import { State } from "./state.js";
import { create_elem } from "./utils.js";

export function delete_menu() {
    let menu = document.getElementById("menu");
    !!menu && menu.remove();

    let mini = document.getElementById("menu_minimizer");
    !!mini && mini.remove();

    let monitor = document.getElementById("connection_monitor");
    !!monitor && monitor.remove();

    let refresh = document.getElementById("status_refresh");
    !!refresh && refresh.remove();
}

/** @param {State} state */
function render_connection_monitor(state) {
    let monitor = create_elem("div", state.header, "monitor");
    monitor.id = "connection_monitor";

    let status = create_elem("h1", monitor, "roboto-mono-norm", "status");
    status.id = "connection_status";

    update_connection_status(state);
}

/** @param {State} state */
export function update_connection_status(state) {
    let status = document.getElementById("connection_status");
    status.textContent = state.connection_status;

    if (state.manual_status) {
        let refresh = document.getElementById("status_refresh");
        if (!refresh) refresh = create_elem("img", state.header, "status_refresh");

        refresh.id = "status_refresh";
        refresh.src = "icons/refresh.svg"
        refresh.addEventListener("click", () => {
            state.manual_status = false;
            run_websocket(state);
            refresh.remove();
        });
    }
}

/** @param {State} state */
export function render_menu(state) {
    delete_menu();

    if (!state.menu_is_open && state.events.length == 0) state.menu_is_open = true;

    let menu = create_elem("div", state.header, 
        "scrollable", (state.menu_is_open ? "open" : "close")
    );
    menu.id = "menu";


    let minimizer = create_elem("img", state.header, "scrollable_minimizer", (!state.menu_is_open) && "flipped");
    minimizer.src = "icons/arrow_back.svg"
    minimizer.id = "menu_minimizer";
    minimizer.addEventListener("click", () => {
        state.menu_is_open = !state.menu_is_open;
        if (state.menu_is_open) {
            menu.classList.replace("close", "open");
            minimizer.classList.remove("flipped");
        } else {
            minimizer.classList.add("flipped");
            menu.classList.replace("open", "close");
        }
    });

    render_connection_monitor(state);

    let img_container = create_elem("div", menu, "header_img_container");
    let img = create_elem("img", img_container, "header_img");
    img.src = "imgs/OHSAL.png";


    // CREATE EVENT
    let event_div = create_elem("div", menu, "event_container", "create");
    event_div.onclick = () => {
        let current = document.getElementById("current_event_page");
        if (!!current) state.main.removeChild(current);

        state.new_current_event();
        render_event(state);

        state.menu_is_open = !state.menu_is_open;
        minimizer.classList.add("flipped");
        menu.classList.replace("open", "close");
    };

    let title_h2 = create_elem("h2", event_div, "roboto-mono-norm");
    title_h2.textContent = "Create Event";

    if (!state.events || state.events.length == 0) {
        let event_div = create_elem("div", menu, "event_container", "no_events");
        let title_h2 = create_elem("h2", event_div, "roboto-mono-norm");
        title_h2.textContent = "NO EVENTS";
        return;
    }

    state.events.forEach((event, i) => {
        let event_div = create_elem("div", menu, "event_container");
        event_div.onclick = () => {
            let current = document.getElementById("current_event_page")
            if (!!current) state.main.removeChild(current);

            
            state.new_current_event(i);
            render_event(state);

            state.menu_is_open = !state.menu_is_open;
            minimizer.classList.add("flipped");
            menu.classList.replace("open", "close");
        };

        let container = create_elem("div", event_div, "live_container");
        create_elem("div", container, "live_indicator");

        let title_h2 = create_elem("h2", event_div, "roboto-mono-norm");
        title_h2.textContent = event.title;
    });
}
