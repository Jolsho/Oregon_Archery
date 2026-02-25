import { render_event } from "./event.js";
import { State } from "./state.js";
import { create_elem } from "./utils.js";

export function delete_menu() {
    let menu = document.getElementById("menu");
    !!menu && menu.remove();

    let mini = document.getElementById("menu_minimizer");
    !!mini && mini.remove();
}


let is_open = false;
/**
 * @param {State} state
*/
export function render_menu(state) {
    delete_menu();

    if (!is_open && state.events.length == 0) is_open = true;


    let menu = create_elem("div", state.main, 
        "scrollable", (is_open ? "open" : "close")
    );
    menu.id = "menu";

    let minimizer = create_elem("img", state.main, "scrollable_minimizer", (!is_open) && "flipped");
    minimizer.src = "icons/arrow_back.svg"
    minimizer.id = "menu_minimizer";
    minimizer.addEventListener("click", () => {
        is_open = !is_open;
        if (is_open) {
            menu.classList.replace("close", "open");
            minimizer.classList.remove("flipped");
        } else {
            minimizer.classList.add("flipped");
            menu.classList.replace("open", "close");
        }
    });


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

        is_open = !is_open;
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

            is_open = !is_open;
            minimizer.classList.add("flipped");
            menu.classList.replace("open", "close");
        };

        let container = create_elem("div", event_div, "live_container");
        create_elem("div", container, "live_indicator");

        let title_h2 = create_elem("h2", event_div, "roboto-mono-norm");
        title_h2.textContent = event.title;
    });
}
