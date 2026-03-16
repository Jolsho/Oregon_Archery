import { new_event, type EventT } from "./eventT";

export type State = {
    events: EventT[];
    current_event_idx: number;
    menu_is_open: boolean;
    connection_status: string;
    retry_count: number;
    manual_status: boolean;
};


export function new_state(): State {
    return ({
        events: [],
        current_event_idx: 0,
        menu_is_open: true,
        connection_status: "Connected",
        retry_count: 0,
        manual_status: false,
    });

};

export function get_event(state: State, idx=-1): EventT | null {
    if (idx == -1) idx = state.current_event_idx;
    if (state.events.length > idx) {
        return state.events[idx]
    } 
    return null;
};

export function set_event(state: State, ev: EventT, idx=-1) {
    if (idx == -1) idx = state.current_event_idx;
    if (state.events.length > idx) {
        state.events[idx] = ev;
        state.events.sort((a, b) => 
            b.created_at.getTime() - 
            a.created_at.getTime()
        );

        if (idx = state.current_event_idx) {
            state.current_event_idx = state.events.findIndex(e => e == ev);
        }
    }
}

export function new_current_event(state: State, idx=-1) {
    if (idx < 0) {
        // filter out events that are currently being created.
        // prevents nesting of new/unfilled events.
        state.events = state.events.filter(ev => 
            !!ev.title && ev.divisions.length > 0
        );

        state.events.push(new_event());
        idx = state.events.length - 1;
    }
    state.current_event_idx = Math.min(idx, state.events.length - 1);
}

export function is_unique_title(state: State, title: string): Boolean {
    let idx = state.events.findIndex((v) => v.title == title);
    return idx == -1 || idx == state.current_event_idx;
}

export function remove_event(state: State, idx=-1) {
    if (idx == -1) idx = state.current_event_idx;

    if (state.events.length > state.current_event_idx) {
        state.events.splice(state.current_event_idx, 1);
        state.current_event_idx = Math.max(state.events.length - 1, 0);
    } 
}

export function find_event_idx(state: State, title: string): number {
    let idx = state.events.findIndex((v) => v.title == title);
    return idx;
}
