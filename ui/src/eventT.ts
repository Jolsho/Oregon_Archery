import { DIVISIONS } from "./global";

export type Division = {
    name: string;
    threshold: number;
};

export type Participant = {
    name: string;
    division: number;
    score: number;
    x_count: number;
};
export function plain_member(): Participant {
    return {
        name: "",
        division: 0,
        score: 0,
        x_count: 0,
    }
}

export type Team = {
    name: string;
    score: number;
    x_count: number;
    members: Participant[];
};

export type EventT = {
    title: string;
    scores_per_team: number;
    is_own: boolean;
    is_persisted: boolean;
    divisions: Division[];
    teams: Team[];
    leaders: Map<string, Participant[]>;
    kind: string;
    created_at: Date;
    expires: Date;
}

export function get_division(event: EventT, idx: number): Division {
    if (event.divisions.length > idx) {
        return event.divisions[idx];
    }
    return { name: "N/A", threshold: 0 };
}

export function calculate_leaders(event: EventT): Map<string, Participant[]> {
    let leaders: Map<string, Participant[]> = new Map();

    event.teams.map((team: Team) => {
            for (const member of team.members) {
                const division = event.divisions[member.division];

                if (!leaders.has(division.name)) {
                    leaders.set(division.name, [
                        plain_member(), 
                        plain_member(), 
                        plain_member()
                    ]);
                }
                leaders.get(division.name)?.push(member);
            }
        },
    );

    for (const [division_name, list] of leaders) {
        list.sort(
            (a, b) => b.score - a.score || b.x_count - a.x_count || 
                a.name.localeCompare(b.name),
        );

        leaders.set(division_name, list.slice(0, 3));
    }
    return leaders;
}


export function new_event(title="", teams=[], divisions=DIVISIONS, kind = "OUTDOOR") {
    let e: EventT = {
        title,
        scores_per_team: 6,
        is_own: true,
        is_persisted: false,
        divisions,
        teams,
        leaders: new Map(),
        kind,
        created_at: new Date(),
        expires: new Date(),
    };
    if (teams.length > 0) e.leaders = calculate_leaders(e);
    return e;
}
