import { delete_event, post_event } from "./api.js";
import { render } from "./initial.js";
import { create_elem, random_uint32 } from "./utils.js";

/**
 * @typedef {Object} Division
 * @property {string} name
 * @property {number} threshold
 */

/**
 *  @typedef {Object} Participant
 *  @property {string} name
 *  @property {number} division
 *  @property {number} score
 *  @property {number} x_count
 */

/**
 *  @typedef {Object} Team
 *  @property {string} name
 *  @property {number} points
 *  @property {Participant[]} members
 */

const MAX_NAME_LENGTH = 15;
const MASCOT_NAMES = [
    "Grizzlies",
    "Raptors",
    "Stallions",
    "Thunder",
    "Vipers",
    "Phoenix",
    "Wolverines",
    "Falcons",
    "Titans",
    "Sharks",
    "Cougars",
    "Mustangs",
    "Panthers",
    "Hurricanes",
    "Bulls",
    "Eagles",
    "Spartans",
    "Bears",
    "Knights",
    "Wolves",
];

export class Event {
    constructor(title = "", teams = [], divisions = [], kind = "OUTDOOR") {
        /** @type {string} */
        this.title = title;

        /** @type {number} */
        this.scores_per_team = 6;

        /** @type {boolean} */
        this.is_own = false;

        /** @type {Division[]} */
        this.divisions = divisions;

        /** @type {Map<string, Team>} */
        this.teams = new Map(teams);

        /** @type {Map<string, Participant[]>} */
        this.leaders = new Map();
        if (teams.length > 0) this.calculate_leaders();

        /** @type {string[]} */
        this.kind = kind;
    }
}

function get_division(event, idx) {
    if (event.divisions.length > idx) {
        return event.divisions[idx];
    }
    return { name: "N/A", threshold: 0 };
}

/** @param {Event} event */
function calculate_leaders(event) {
    event.leaders = new Map();

    Object.values(event.teams).map(
        /** @param {Team} team */
        (team) => {
            for (const member of team.members) {
                const division = get_division(event, member.division);

                if (!event.leaders.has(division.name)) {
                    event.leaders.set(division.name, []);
                }
                event.leaders.get(division.name).push(member);
            }
        },
    );

    for (const [division_name, list] of event.leaders) {
        list.sort(
            (a, b) => b.score - a.score || b.x_count - a.x_count || b.name - a.name,
        );

        event.leaders.set(division_name, list.slice(0, 3));
    }
}

/**
 * @param {HTMLElement} element
 * @param {Team} team
 * @param {Event} event
 * @param {string} teamboard_id
 * @param {HTMLElement} entry
 * @param {Participant} participant
 * @param {(string) => void} input_callback
 */
function team_input_field(
    element,
    team,
    event,
    teamboard_id,
    entry,
    participant,
    input_callback,
) {
    element.addEventListener("input", (e) => {
        input_callback(e.target.value);
    });

    element.addEventListener("keydown", (e) => {
        let board = document.getElementById(teamboard_id);
        if (e.shiftKey && e.key == "Enter") {
            team.members.push({
                name: "",
                division: 0,
                score: 0,
                x_count: 0,
            });

            let parent = board.parentNode;

            parent.insertBefore(render_teamboard(event, team, true), board);
            parent.removeChild(board);
        } else if (e.shiftKey && e.key == "Backspace") {
            board.removeChild(entry);
            team.members = team.members.filter((p) => p !== participant);
        } else if (e.key == "Enter") {
            submit_team(event, team, teamboard_id);
        }
    });
}

/**
 * @param {Event} event
 * @param {Team} team
 * @param {string} board_id
 */
function submit_team(event, team, board_id) {
    if (!team.name) {
        team.name = MASCOT_NAMES[random_uint32() % 20];
    }

    team.members = team.members.filter((member) => !!member.name);

    event.teams[team.name] = team;

    post_event(event)
        .then((ev) => {
            event = ev;
            let board = document.getElementById(board_id);
            let parent = board.parentNode;

            parent.insertBefore(render_teamboard(event, team, false), board);
            parent.removeChild(board);

            let leader_grid = document.getElementsByClassName("leaderboard_grid")[0];
            leader_grid.innerHTML = "";
            calculate_leaders(event);
            render_leaderboard(leader_grid, event);
        })
        .catch((e) => console.error(e));
}

/**
 * @param {Event} event
 * @param {Team | null} team
 * @param {Boolean} is_maluable
 * @returns {HTMLElement}
 */
function render_teamboard(event, team, is_maluable = false) {
    const team_existed = !!team;
    team = !team_existed
        ? {
            name: "",
            members: [{ name: "", division: 0, score: 0, x_count: 0 }],
            score: 0,
            x_count: 0,
        }
        : team;

    if (is_maluable || !team_existed)
        return render_teamboard_mut(event, team, team_existed);

    team.members.sort(
        (a, b) => b.score - a.score || b.x_count - a.x_count || b.name - a.name,
    );

    let teamboard = create_elem("div", null, "team_scoreboard");
    teamboard.id = random_uint32();

    let team_header = create_elem("div", teamboard, "team_header");

    let team_name_div = create_elem("h2", teamboard, "roboto-mono-norm");
    if (team.name.length > MAX_NAME_LENGTH) {
        let names = team.name.split(" ");

        if (names[0].length > 10) {
            names[0] = names[0].slice(0, 10);
            names[0] += ".";
        }
        team_name_div.textContent =
            names[0] + " " + names[1][0].toUpperCase() + ".";
    } else {
        team_name_div.textContent = team.name;
    }
    team_header.appendChild(team_name_div);

    team_header.insertBefore(document.createElement("img"), team_name_div);

    if (event.is_own) {
        let edit_team_btn = create_elem("img", team_header, "team_action_btn");
        edit_team_btn.src = "icons/edit.svg";
        edit_team_btn.addEventListener("click", () => {
            let board = document.getElementById(teamboard.id);
            let parent = board.parentNode;

            parent.insertBefore(render_teamboard(event, team, true), board);
            parent.removeChild(board);
        });
    } else {
        team_header.appendChild(document.createElement("img"));
    }

    team.members.forEach((part) => {
        let entry = create_elem(
            "div",
            teamboard,
            "participant_entry",
            "roboto-mono-norm",
        );

        // NAME
        let name = create_elem("p", entry, "name");
        if (part.name.length > MAX_NAME_LENGTH) {
            let names = part.name.split(" ");

            if (names[0].length > 10) {
                names[0] = names[0].slice(0, 8);
                names[0] += ".";
            }
            names[1] = names[1].slice(0, MAX_NAME_LENGTH - names[0].length - 1);
            name.textContent = names[0] + " " + names[1] + ".";
        } else {
            name.textContent = part.name;
        }

        let seperator = create_elem("hr", entry);

        // DIVISION
        let division = create_elem("p", entry, "division");
        let label = create_elem("span", division);
        label.textContent = get_division(event, part.division)
            .name.slice(0, 3)
            .toUpperCase();

        seperator = create_elem("hr", entry);

        let score = create_elem("p", entry, "number");
        score.textContent = part.score;

        seperator = create_elem("hr", entry);

        let x_count = create_elem("p", entry, "number");
        x_count.textContent = part.x_count;
    });

    return teamboard;
}

/**
 * @param {Event} event
 * @param {Team | null} team
 * @param {Boolean} team_existed
 * @returns {HTMLElement}
 */
function render_teamboard_mut(event, team, team_existed) {
    let teamboard = create_elem("div", null, "team_scoreboard");
    teamboard.id = random_uint32();

    let team_header = create_elem("div", teamboard, "team_header");

    let team_name_div = create_elem("input", teamboard, "roboto-mono-norm");
    team_name_div.type = "text";
    team_name_div.placeholder = "TEAM NAME";

    if (team_existed) {
        team_name_div.value = team.name;
    } else {
        requestAnimationFrame(() => team_name_div.focus());
    }
    team_name_div.addEventListener("input", (e) => {
        team.name = e.target.value;
    });
    team_header.appendChild(team_name_div);

    let remove_team_btn = document.createElement("img");
    remove_team_btn.src = "icons/garbage.svg";
    remove_team_btn.classList.add("team_action_btn", "left");
    remove_team_btn.addEventListener("click", () => {
        delete event.teams[team.name];
        post_event(event)
            .then((ev) => {
                event = ev;
                let board = document.getElementById(teamboard.id);
                board.parentNode.removeChild(board);
                delete event.teams[team.name];

                let leader_grid =
                    document.getElementsByClassName("leaderboard_grid")[0];
                leader_grid.innerHTML = "";
                calculate_leaders(event);
                render_leaderboard(leader_grid, event);
            })
            .catch((e) => console.error(e));
    });
    team_header.insertBefore(remove_team_btn, team_name_div);

    let submit_team_btn = document.createElement("img");
    submit_team_btn.src = "icons/submit.svg";
    submit_team_btn.classList.add("team_action_btn");
    team_header.appendChild(submit_team_btn);

    submit_team_btn.addEventListener("click", () =>
        submit_team(event, team, teamboard.id),
    );

    team.members.forEach((part, idx) => {
        let entry = create_elem(
            "div",
            teamboard,
            "participant_entry",
            "roboto-mono-norm",
        );

        let remove_part = create_elem("img", entry);
        remove_part.src = "icons/close.svg";
        remove_part.addEventListener("click", () => {
            teamboard.removeChild(entry);
            team.members = team.members.filter((p) => p !== part);
        });

        let name = create_elem("input", entry, "name");
        if (!!part.name) {
            name.value = part.name;
        } else {
            name.placeholder = "Name";
            if (!!team.name) {
                requestAnimationFrame(() => name.focus());
            }
        }
        team_input_field(
            name,
            team,
            event,
            teamboard.id,
            entry,
            part,
            (v) => (team.members[idx].name = v),
        );

        let seperator = create_elem("hr", entry);

        let division = create_elem("p", entry, "division", "maluable_division");
        let label = create_elem("span", division);
        label.textContent = get_division(event, part.division)
            .name.slice(0, 3)
            .toUpperCase();

        let menu = create_elem("div", division, "division_menu");

        event?.divisions.forEach((div, i) => {
            let option = create_elem("div", menu, "div_option", "roboto-mono-norm");

            option.textContent = div.name;
            option.addEventListener("click", (e) => {
                e.stopPropagation(); // prevents reopening immediately
                label.textContent = div.name.slice(0, 3).toUpperCase(); // only update the label
                team.members[idx].division = i;
                menu.classList.remove("visible");
            });
        });

        division.addEventListener("click", () => {
            if (menu.classList.contains("visible")) {
                menu.classList.remove("visible");
            } else {
                menu.classList.add("visible");
            }
        });

        seperator = create_elem("hr", entry);

        let score = create_elem("input", entry, "number");
        !!part.score ? (score.value = part.score) : (score.placeholder = "Score");

        team_input_field(
            score,
            team,
            event,
            teamboard.id,
            entry,
            part,
            (v) => (team.members[idx].score = Number(v)),
        );

        seperator = create_elem("hr", entry);

        let x_count = create_elem("input", entry, "number");
        !!part.x_count
            ? (x_count.value = part.x_count)
            : (x_count.placeholder = "Xs");

        team_input_field(
            x_count,
            team,
            event,
            teamboard.id,
            entry,
            part,
            (v) => (team.members[idx].x_count = Number(v)),
        );
    });

    let adder_entry = create_elem("div", teamboard, "adder_entry");

    let adder = create_elem("img", adder_entry);
    adder.src = "icons/add.svg";
    adder.addEventListener("click", () => {
        team.members.push({
            name: "",
            division: 0,
            score: 0,
            x_count: 0,
        });

        let board = document.getElementById(teamboard.id);
        let parent = board.parentNode;

        parent.insertBefore(render_teamboard(event, team, true), board);
        parent.removeChild(board);
    });

    return teamboard;
}

/**
 * @param {HTMLElement} container
 * @param {Event} event
 */
function render_leaderboard(container, event) {
    for (const div of event?.divisions) {
        let leaders = event.leaders.get(div.name);

        let box = create_elem("div", null, "leaderboard");

        let entry = create_elem("div", box, "leader");

        let header = create_elem("h1", entry, "roboto-mono-norm");
        header.textContent = div.name;

        if (!!leaders) {
            leaders.forEach((leader, i) => {
                let entry = create_elem("div", box, "leader");

                let name = create_elem("h2", entry, "roboto-mono-norm");
                if (leader.name.length > MAX_NAME_LENGTH) {
                    let names = leader.name.split(" ");

                    if (names[0].length > 10) {
                        names[0] = names[0].slice(0, 8);
                        names[0] += ".";
                    }
                    names[1] = names[1].slice(0, MAX_NAME_LENGTH - names[0].length - 1);
                    name.textContent = names[0] + " " + names[1] + ".";
                } else {
                    name.textContent = leader.name;
                }

                let seperator = create_elem("hr", entry);

                let score = create_elem("h3", entry, "roboto-mono-norm");
                score.textContent = leader.score;

                seperator = create_elem("hr", entry);

                let xs = create_elem("h3", entry, "roboto-mono-norm");
                xs.textContent = leader.x_count;

                let place = create_elem("h4", entry, "placement", "roboto-mono-norm");
                place.textContent = `${i + 1}`;
            });
            container.appendChild(box);
        }
    }
    render_team_leaderboard(container, event);
}

/**
 * @param {HTMLElement} container
 * @param {Event} event
 */
function render_team_leaderboard(container, event) {
    let team_order = [];

    Object.values(event.teams).map(
        /** @param {Team} team*/
        (team) => {
            team.score = 0;
            team.x_count = 0;
            if (event.kind == "OUTDOOR") {
                for (
                    let i = 0;
                    i < Math.min(team.members.length, event.scores_per_team);
                    i++
                ) {
                    team.score += team.members[i].score;
                    team.x_count += team.members[i].x_count;
                }
            } else {
                for (const mem of team.members) {
                    const threshold = get_division(event, mem.division).threshold;
                    if (mem.score >= threshold) team.score++;
                }
            }
            team_order.push(team);
        },
    );

    team_order.sort((a, b) => b.score - a.score || b.x_count - a.x_count);

    // I THINK THEY WANT TO SEE ALL TEAMS...
    // team_order = team_order.slice(0,3);

    if (team_order.length == 0) return;

    let box = create_elem("div", container, "leaderboard");

    let entry = create_elem("div", box, "leader");

    let header = create_elem("h1", entry, "roboto-mono-norm");
    header.textContent = "TEAM";

    team_order.forEach((team, i) => {
        let entry = create_elem("div", box, "leader");

        let name = create_elem("h2", entry, "roboto-mono-norm");
        if (team.name.length > MAX_NAME_LENGTH) {
            let names = team.name.split(" ");

            if (names[0].length > 10) {
                names[0] = names[0].slice(0, 10);
                names[0] += ".";
            }
            name.textContent = names[0] + " " + names[1][0].toUpperCase() + ".";
        } else {
            name.textContent = team.name;
        }

        let seperator = create_elem("hr", entry);

        let score = create_elem("h3", entry, "roboto-mono-norm");
        score.textContent = team.score;

        seperator = create_elem("hr", entry);

        let xs = create_elem("h3", entry, "roboto-mono-norm");
        xs.textContent = (event.kind == "OUTDOOR") ? team.x_count : 100 - i;

        let place = create_elem("h4", entry, "placement", "roboto-mono-norm");
        place.textContent = `${i + 1}`;
    });
}

function submit_event(
    events, main, idx, parent, title
) {
    let event = events[idx];
    if (!!event.title) {
        let existing_idx = events.findIndex((v) => v.title == event.title);
        if (existing_idx == -1 || existing_idx == idx) {
            post_event(event)
                .then((ev) => {
                    events[idx] = ev;
                    main.innerHTML = "";
                    render(events);
                    render_event(events, idx, main);
                })
                .catch((e) => console.error(e));
            return;
        }
        let err_msg = document.getElementById("event_err_msg");
        if (!err_msg) {
            err_msg = create_elem("h3", parent, "error", "roboto-mono-norm");
            err_msg.id = "event_err_msg";
        }
        err_msg.textContent = "An event with that name already exists.";

    } else {
        let err_msg = document.getElementById("event_err_msg");
        if (!err_msg) {
            err_msg = create_elem("h3", parent, "error", "roboto-mono-norm");
            err_msg.id = "event_err_msg";
        }
        err_msg.textContent = "Events must have titles.";

        requestAnimationFrame(() => title.focus());
    }
}

/**
 * @param {Event[]} events
 * @param {Number} idx
 * @param {HTMLElement | null} main
 * @param {Boolean} is_maluable
 */
export function render_event(events, idx, main = null, is_maluable = false) {
    if (events.length <= idx) return;
    let event = events[idx];

    const is_new = !event.title;

    if (!is_maluable && is_new) is_maluable = true;

    if (is_maluable && !is_new && !event.is_own) is_maluable = false;

    if (!main) main = document.getElementById("main");

    let scrollable = create_elem("div", main, "event_feed_scrollable");
    scrollable.id = "current_event_page";

    let container = create_elem("div", scrollable, "event_feed_container");

    let title_container = create_elem("div", container, "title_container");

    let title = create_elem(
        is_maluable ? "input" : "h1",
        title_container,
        "roboto-mono-norm",
        "title",
    );
    title.textContent = event.title;

    if (is_maluable) {
        if (!is_new) title.value = event.title;

        title.placeholder = "Event Title";
        title.addEventListener("input", (e) => {
            event.title = e.target.value;
        });
        title.addEventListener("keydown", (e) => {
            if (e.key == "Enter") {
                submit_event(events, main, idx, title_container, title);
            }
        });

        let remove_event_btn = create_elem("img", null, "team_action_btn");
        remove_event_btn.src = "icons/garbage.svg";
        remove_event_btn.addEventListener("click", () => {
            delete_event(event.title).then(() => {
                main.innerHTML = "";
                events.splice(idx, 1);
                if (events.length > 0) render_event(events, events.length - 1, main);
                render(events);
            });
        });
        title_container.insertBefore(remove_event_btn, title);

        let submit_team_btn = create_elem(
            "img",
            title_container,
            "team_action_btn",
        );
        submit_team_btn.src = "icons/submit.svg";
        submit_team_btn.addEventListener("click", () => {
            submit_event(events, main, idx, title_container, title);
        });
        requestAnimationFrame(() => title.focus());

        let edit_event_panel = create_elem(
            "div", title_container, "edit_event_panel",
        );
        let kind_container = create_elem("div", edit_event_panel, "kind_container");
        const KINDS = ["OUTDOOR", "INDOOR"];
        let scoresPerTeamLabel = null;

        function renderKindExtras(kind) {
            // Clear previous UI
            if (scoresPerTeamLabel) {
                scoresPerTeamLabel.remove();
                scoresPerTeamLabel = null;
            }


            if (kind === "OUTDOOR") {
                scoresPerTeamLabel = create_elem(
                    "label", edit_event_panel, "roboto-mono-norm"
                );
                scoresPerTeamLabel.textContent = "Scores Per Team";

                let scoresPerTeamInput = create_elem(
                    "input", scoresPerTeamLabel, "roboto-mono-norm", "scores_per_team"
                );
                scoresPerTeamInput.id = "OUTDOOR_SCORES_PER_TEAM";
                scoresPerTeamInput.type = "number";
                scoresPerTeamInput.value = event.scores_per_team;

                scoresPerTeamInput.addEventListener("input", (e) => {
                    event.scores_per_team = Number(e.target.value);
                });

                scoresPerTeamInput.addEventListener("keydown", (e) => {
                    if (e.key == "Enter") {
                        submit_event(events, main, idx, title_container, title);
                    }
                });
            } else if (kind === "INDOOR") {
            }
        }

        KINDS.forEach((kind, i) => {
            const label = create_elem("label", kind_container, "roboto-mono-norm");
            label.textContent = kind;

            const input = create_elem("input", label);
            input.type = "radio";
            input.name = "event-kind";
            input.value = kind;

            if (i === 0) {
                input.checked = true;
                event.kind = kind;
                renderKindExtras(kind);
            }

            input.addEventListener("change", (e) => {
                if (!e.target.checked) return;

                event.kind = e.target.value;
                renderKindExtras(event.kind);
            });
        });
    } else if (event.is_own) {
        title_container.insertBefore(document.createElement("div"), title);

        let edit_btn = create_elem("img", title_container, "team_action_btn");
        edit_btn.src = "icons/edit.svg";
        edit_btn.addEventListener("click", () => {
            main.removeChild(scrollable);
            render_event(events, idx, main, true);
        });
    } else {
        title_container.insertBefore(document.createElement("div"), title);
        title_container.appendChild(document.createElement("div"));
    }

    let leader_container = create_elem("div", container, "leaderboard_grid");

    calculate_leaders(event);
    render_leaderboard(leader_container, event);

    let team_boards = create_elem("div", container, "team_boards");

    Object.values(event.teams).map((team) =>
        team_boards.appendChild(render_teamboard(event, team)),
    );

    if (!is_new && event.is_own) {
        let add_cont = create_elem(
            "div",
            team_boards,
            "team_scoreboard",
            "add_team_container",
        );

        let add = create_elem("img", add_cont, "add_team");
        add.src = "icons/add.svg";
        add.addEventListener("click", () => {
            team_boards.insertBefore(render_teamboard(event, null, true), add_cont);
        });

        create_elem("div", add, "vert");
        create_elem("div", add, "horz");
    }
}
