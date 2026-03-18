import { useEffect, useMemo, useState } from "react";
import { calculate_leaders, type Division, type EventT, type Team } from "./eventT";
import { TeamBoard } from "./teamboard";
import { MAX_NAME_LENGTH } from "./global";
import garbage_svg from "./assets/garbage.svg"
import submit_svg from "./assets/submit.svg"
import edit_svg from "./assets/edit.svg"
import add_svg from "./assets/add.svg"

type EditPanelProps = {
    submit_event: () => void;
    kind: string;
    setKind: (kind: string) => void;
    scores_per_team: number;
    set_scores_per_team: (spt: number) => void;
    divisions: Division[];
    setDivisions: (divs: Division[]) => void;
};

function EditEventPanel({
    submit_event,
    kind, setKind,
    scores_per_team, set_scores_per_team,
    divisions, setDivisions,
}: EditPanelProps) {
    const KINDS = ["OUTDOOR", "INDOOR"];

    return (
    <div className="edit_event_panel">
        <div className="kind_container">
            {KINDS.map((k) => (
                <label className="roboto-mono-norm" key={"kind" + k}>{k}
                    <input type="radio" name="event-kind" 
                        value={k} 
                        checked={kind === k}
                        onChange={(e) => setKind(e.target.value)}
                    />
                </label>
            ))}
        </div>
        {kind === "OUTDOOR" ? (
            <label className="roboto-mono-norm">
                {"Scores Per Team"}
                <input className="roboto-mono-norm scores_per_team" 
                    type="number" value={scores_per_team}
                    onInput={(e) => set_scores_per_team(Number(e.currentTarget.value)) }
                    onKeyDown={(e) => {
                        if (e.key == "Enter") submit_event();
                    }}
                />
            </label>
        ) : (
            divisions.map((divi, i) => (
                <label className="roboto-mono-norm right_label" key={i + "div_label"} >
                    {divi.name}
                    <input className="scores_per_team" type="number" value={divi.threshold} 
                        onChange={(e) => {
                            const value = e.target.value;
                            const divs = [...divisions];
                            divs[i].threshold = !value ? 0 : Number(value);
                            setDivisions(divs);
                        }}
                        onKeyDown={(e) => {
                            if (e.key == "Enter") submit_event();
                        }}
                    />
                </label>
            ))
        )}
    </div>
    );
}

type LeaderBoardProps = {
    event: EventT;
};

function LeaderBoards({ event }: LeaderBoardProps) {

    const leaders = useMemo(() => {
        return calculate_leaders(event);
    }, [event])

    const team_order: Team[] = useMemo(() => {

        let order: Team[] = [];

        event?.teams.forEach((t) => {
            let team = {...t, score: 0, x_count: 0};
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
                    const threshold = event.divisions[mem.division].threshold;
                    if (mem.score >= threshold) team.score++;
                }
            }
            order.push(team);
        });

        if (order.length > 9) order = order.slice(0, 9);
        order.sort((a, b) => b.score - a.score || b.x_count - a.x_count);

        return order;

    }, [event]);

    return (
    <div className="leaderboard_grid">
    {/* DIVISION LEADER BOARDS */}
    {event.divisions.map((div, i) => {
        const ls = leaders.get(div.name);
        if (!ls) return;

        return (
        <div className="leaderboard" key={div.name + i}>
            <div className="leader">
                <h1 className="roboto-mono-norm">
                    {div.name.split(" ")[1]}
                </h1>
                <h2 className="roboto-mono-norm side_value left">
                    {div.name.split(" ")[0][0]}
                </h2>
                <h2 className="roboto-mono-norm side_value">
                    V:{div.threshold}
                </h2>
            </div>
            {ls.map((leader, i) => {
                let name = leader.name;
                if (name.length > MAX_NAME_LENGTH) {
                    let names = name.split(" ");

                    if (names[0].length > 10) {
                        names[0] = names[0].slice(0, 8);
                        names[0] += ".";
                    }
                    names[1] = names[1].slice(0, MAX_NAME_LENGTH - names[0].length - 1);
                    name = names[0] + " " + names[1] + ".";
                }

                return (
                <div className="leader" key={"leader" + leader.name + i}>
                    {!!name &&
                    <>
                    <h2 className="roboto-mono-norm">{name}</h2>
                    <hr/>
                    <h3 className="roboto-mono-norm">{leader.score}</h3>
                    <hr/>
                    <h3 className="roboto-mono-norm">{leader.x_count}</h3>
                    <h4 className="placement roboto-mono-norm">{i + 1}</h4>
                    </>
                    }
                </div>
                )
            })}
        </div>
        );
    })}

    {/* TEAM LEADER BOARD */}
    {team_order.length > 0 &&
    <div className="leaderboard">
        <div className="leader">
            <h1 className="roboto-mono-norm">TEAM</h1>
        </div>
        {team_order.map((team, i) => {
            let name = team.name;
            if (team.name.length > MAX_NAME_LENGTH) {
                let names = team.name.split(" ");

                if (names[0].length > 10) {
                    names[0] = names[0].slice(0, 10);
                    names[0] += ".";
                }
                name = names[0] + " " + names[1][0].toUpperCase() + ".";
            }
            return (
            <div className="leader" key={"team_leader" + team.name + i}>
                <h2 className="roboto-mono-norm">{name}</h2>
                <hr/>
                <h3 className="roboto-mono-norm">{team.score}</h3>
                <hr/>
                <h3 className="roboto-mono-norm">
                    {(event.kind == "OUTDOOR") ? team.x_count : 100 - i}
                </h3>
                <h4 className="roboto-mono-norm placement">{i + 1}</h4>
            </div>
            );
        })}
    </div>
    }
    </div> 
    );
}

export type EventPageProps = {
    events: EventT[];
    idx: number;
    post_event: (event: EventT) => void;
    remove_event: () => void;
    maluable?: boolean;
};

export function EventPage({
    events, idx, 
    post_event, remove_event
}: EventPageProps) {

    const event = events[idx];

    const [kind,            setKind]    = useState(event.kind || "");
    const [scores_per_team, setScores]  = useState(event.scores_per_team);
    const [divisions,       setDivisions] = useState(event.divisions);

    const [teams,   setTeams]   = useState<Team[]>(event?.teams || [])
    const [is_new,  setIsNew]   = useState(!event?.title)
    const [title,   setTitle]   = useState(event?.title || "")
    const [error,   setError]   = useState("")
    const [maluable, setMaluable] = useState(!event?.title)

    useEffect(() => {
        if (event) {
            setTeams(event.teams);
            setIsNew(!event.title);
            setMaluable(!event.title);
            setTitle(event.title);
            setKind(event.kind);
            setScores(event.scores_per_team);
            setDivisions(event.divisions);
            setError("");
        }
    }, [idx])

    const submit_event = (ev?: EventT) => {
        const found = events.findIndex(e => e.title == title);
        if (!title || (found !== -1 && found !== idx)) {
            setError("Event Needs Unique Title");
            return;
        }

        if (!ev) ev = {...event, teams };

        ev.teams = ev.teams.filter(t => !!t.name);

        post_event({...ev, title, kind, divisions, scores_per_team });
        setMaluable(false);
        setIsNew(false);
    };


    return (
    <div 
        id="current_event_page"
        className="event_feed_scrollable" 
    > 
        <div className="event_feed_container" > 
            <div className="title_container" > 
                {!!error && (
                    <h3 className="error roboto-mono-norm" id="event_err_msg"> 
                        {error}
                    </h3>
                )}

                {maluable ? (
                    <>
                    <img src={garbage_svg}
                        className="team_action_btn"
                        onClick={remove_event}
                    />
                    <input 
                        className="roboto-mono-norm title" 
                        value={title}
                        placeholder="EVENT TITLE"
                        onInput={(e) => setTitle(e.currentTarget.value)}
                        onKeyDown={(e) => e.key == "Enter" && submit_event()}
                        autoFocus={true}
                    />
                    <img src={submit_svg}
                        className="team_action_btn"
                        onClick={() => submit_event()}
                    />
                    <EditEventPanel 
                        submit_event={submit_event} 
                        kind={kind} setKind={setKind}
                        scores_per_team={scores_per_team}
                        set_scores_per_team={setScores}
                        divisions={divisions}
                        setDivisions={setDivisions}
                    />
                    </>
                ) : (
                    (event.is_own) ? (
                    <>
                        <div/>
                        <h1 className="roboto-mono-norm title" > {event.title} </h1>
                        <img src={edit_svg}
                            className="team_action_btn"
                            onClick={() => setMaluable(true)}
                        />
                    </>
                    ) : (
                    <> 
                    <div/> 
                    <h1 className="roboto-mono-norm title" > {event.title} </h1>
                    <div/> 
                    </>
                    )
                )}
            </div>

            <LeaderBoards event={event} />

            <div className="team_boards">
                {teams?.map((team, i) => (
                    <TeamBoard
                        key={i + team.name}
                        isown={event.is_own}
                        divisions={event.divisions}
                        team={team}
                        submit_team={(t) => {
                            let ts = [...teams];
                            ts[i] = t;
                            setTeams(ts);
                            submit_event({...event, teams: ts});
                        }}
                        remove_team={() => {
                            let ts = [...teams];
                            ts.splice(i, 1);
                            setTeams(ts);
                            submit_event({...event, teams: ts});
                        }}
                    />
                ))}

                {!is_new && event.is_own && 
                    <div className="team_scoreboard add_team_container">
                        <img src={add_svg}
                            className="add_team"
                            onClick={() => {
                                setTeams(prev => (
                                    [...prev, {
                                            name: "", 
                                            score: 0, 
                                            x_count: 0,
                                            members: []
                                    }]
                                ));
                            }}
                        >
                        </img>
                        <div className="vert"/>
                        <div className="horz"/>
                    </div>
                }
            </div>
        </div>
    </div>
    );
}
