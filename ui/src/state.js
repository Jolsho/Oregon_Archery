import { Event } from "./event.js";
import { create_elem } from "./utils.js";

export class State {
    constructor() {
        /** @type {Event[]} */
        this.events = [];

        /** @type {number} */
        this.current_event_idx = 0;

        /** @type {HTMLElement} */
        this.main = document.getElementById("main");

        /** @type {HTMLElement} */
        this.header = create_elem("div", this.main, "header");

        /** @type {boolean} */
        this.menu_is_open = true;

        /** @type {string} */
        this.connection_status = "Connected";

        /** @type {number} */
        this.retry_count = 0;

        /** @type {boolean} */
        this.manual_status = false;
    };

    /** 
     * @param {number} idx 
     * @returns {Event | null} 
    */
    get_event(idx=-1) {
        if (idx == -1) idx = this.current_event_idx;
        if (this.events.length > idx) {
            return this.events[idx]
        } 
        return null;
    };


    /** 
     * @param {Event} ev 
     * @param {number} idx 
    */
    set_event(ev, idx=-1) {
        if (idx == -1) idx = this.current_event_idx;
        if (this.events.length > idx) {
            this.events[idx] = ev;
        }
    }

    /** @param {number} idx */
    new_current_event(idx=-1) {
        if (idx < 0) {
            // filter out events that are currently being created.
            // prevents nesting of new/unfilled events.
            this.events = this.events.filter(ev => 
                !!ev.title && ev.divisions.length > 0
            );

            this.events.push(new Event());
            idx = this.events.length - 1;
        }
        this.current_event_idx = Math.min(idx, this.events.length - 1);
    }


    /** 
     * @param {string} title 
     * @returns {boolean} 
    */
    is_unique_title(title) {
        let idx = this.events.findIndex((v) => v.title == title);
        return idx == -1 || idx == this.current_event_idx;
    }

    remove_event(idx=-1) {
        if (idx == -1) idx = this.current_event_idx;

        if (this.events.length > this.current_event_idx) {
            this.events.splice(this.current_event_idx, 1);
            this.current_event_idx = Math.max(this.events.length - 1, 0);
        } 
    }

    /** 
     * @param {string} title 
     * @returns {number} 
    */
    find_event_idx(title) {
        let idx = this.events.findIndex((v) => v.title == title);
        return idx;
    }
};
