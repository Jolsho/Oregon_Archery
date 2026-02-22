import { render_event, Event } from "./event.js";
import { create_elem } from "./utils.js";


let is_open = false;
/**
 * @param {Event[]} events
*/
export function render(events) {
    let main = document.getElementById("main");
    if (!is_open && events.length == 0) is_open = true;


    let feed = create_elem("div", main, 
        "scrollable", (is_open ? "open" : "close")
    );
    feed.id = "feed";

    let minimizer = create_elem("img", main, "scrollable_minimizer", (!is_open) && "flipped");
    minimizer.src = "icons/arrow_back.svg"
    minimizer.addEventListener("click", () => {
        is_open = !is_open;
        if (is_open) {
            feed.classList.replace("close", "open");
            minimizer.classList.remove("flipped");
        } else {
            minimizer.classList.add("flipped");
            feed.classList.replace("open", "close");
        }
    });


    let img_container = create_elem("div", feed, "header_img_container");
    let img = create_elem("img", img_container, "header_img");
    img.src = "OHSAL.png";


    // CREATE EVENT
    let event_div = create_elem("div", feed, "event_container", "create");
    event_div.onclick = () => {
        let current = document.getElementById("current_event_page");
        if (!!current) main.removeChild(current);

        events.push(new Event())
        render_event(events, events.length - 1, main);

        is_open = !is_open;
        minimizer.classList.add("flipped");
        feed.classList.replace("open", "close");
    };

    let title_h2 = create_elem("h2", event_div, "roboto-mono-norm");
    title_h2.textContent = "Create Event";

    if (!events || events.length == 0) {
        let event_div = create_elem("div", feed, "event_container", "no_events");
        let title_h2 = create_elem("h2", event_div, "roboto-mono-norm");
        title_h2.textContent = "NO EVENTS";
        return;
    }

    events.forEach((event, i) => {
        let event_div = create_elem("div", feed, "event_container");
        event_div.onclick = () => {
            let current = document.getElementById("current_event_page")
            if (!!current) main.removeChild(current);

            render_event(events, i, main);

            is_open = !is_open;
            minimizer.classList.add("flipped");
            feed.classList.replace("open", "close");
        };

        let container = create_elem("div", event_div, "live_container");
        create_elem("div", container, "live_indicator");

        let title_h2 = create_elem("h2", event_div, "roboto-mono-norm");
        title_h2.textContent = event.title;
    });
}
