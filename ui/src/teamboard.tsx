import { type KeyboardEvent, useState, useEffect } from "react";
import { plain_member, type Division, type Participant, type Team } from "./eventT"
import { MASCOT_NAMES, MAX_NAME_LENGTH, NA_DIV } from "./global.ts";

import garbage_svg from "./assets/garbage.svg"
import edit_svg from "./assets/edit.svg"
import add_svg from "./assets/add.svg"
import submit_svg from "./assets/submit.svg"
import close_svg from "./assets/close.svg"
import { random_uint32 } from "./utils.ts";

type MemberProps = {
    divisions: Division[];
    member: Participant;
};
function Member({ 
    divisions, 
    member, 
}: MemberProps) {
    let name = member.name;
    if (name.length > MAX_NAME_LENGTH) {
        let names = name.split(" ");

        if (names[0].length > 10) {
            names[0] = names[0].slice(0, 8);
            names[0] += ".";
        }
        names[1] = names[1].slice(0, MAX_NAME_LENGTH - names[0].length - 1);
        name = names[0] + " " + names[1] + ".";
    }

    let div = member.division < divisions.length ? divisions[member.division] : NA_DIV;
    let names = div.name.split(" ");
    let d = names[1].slice(0, 3).toUpperCase();
    let gender =  names[0][0].toUpperCase();
    let div_name = `${gender} ${d}`;

    return (
    <>
        <p className="name"> {name} </p>
        <hr/>
        <p className="division">
            <span>{div_name}</span>
        </p>
        <hr/>
        <p className="number"> {member.score} </p>
        <hr/>
        <p className="number"> {member.x_count} </p>
    </>
    );
};


type MutMemberProps = {
    divisions: Division[];
    member: Participant;
    add_member: () => void;
    update_member: (member: Participant) => void;
    remove_member: () => void;
    focus: boolean;
    submit: () => void;
};
function Mut_Member({ 
    divisions, 
    member, 
    add_member, 
    update_member, 
    remove_member, 
    focus, 
    submit 
}: MutMemberProps) {

    const [menu_visible, set_menu_visiblity] = useState(false);

    const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) : boolean => {
        if (e.shiftKey && e.key === "Enter") {
            add_member();
        } else if (e.shiftKey && e.key === "Backspace") {
            remove_member();
        } else if (e.key === "Enter") {
            submit();
        } else {
            return false;
        }
        return true;
    };
    let names = divisions[member.division].name.split(" ");
    let d = names[1].slice(0,3).toUpperCase();
    let gender =  names[0][0].toUpperCase();

    return (
        <>
        <img src={close_svg} onClick={remove_member}/>

        <input 
            className="roboto-mono-norm name"
            placeholder="Name" 
            value={member.name} 
            autoFocus={focus} 
            onKeyDown={(e) => handleKeyDown(e)}
            onChange={(e) => update_member({...member, 
                name: e.currentTarget.value
            })}
        />

        <hr/>

        <p className="division maluable_division"
            onClick={() => set_menu_visiblity(!menu_visible)}
        >

            {menu_visible && 
                <div className="division_menu">
                    {divisions.map((div, i) => {
                        let names = div.name.split(" ");
                        let d = names[1].toUpperCase();
                        let gender =  names[0][0].toUpperCase();
                        const text = `${gender} ${d}`;

                        return (
                        <div className="div_option roboto-mono-norm"
                            onClick={(e) => {
                                e.stopPropagation();
                                update_member({
                                    ...member,
                                    division: i,
                                });
                                set_menu_visiblity(false);
                            }}
                        >{text}</div>
                        );

                    })}
                </div>
            }
            <span>
            {`${gender} ${d}`}
            </span>
        </p>

        <hr/>

        <input 
            className="roboto-mono-norm"
            placeholder="Score" 
            value={member.score} 
            onKeyDown={(e) => handleKeyDown(e)}
            onChange={(e) => update_member({...member, 
                score: Number(e.currentTarget.value)
            })}
        />

        <hr/>

        <input 
            className="roboto-mono-norm"
            placeholder="Xs" 
            value={member.x_count} 
            onKeyDown={(e) => handleKeyDown(e)}
            onChange={(e) => update_member({...member, 
                    x_count: Number(e.currentTarget.value)
            })}
        />
    </>
    );
}


type TeamBoardProps = {
    isown: boolean;
    divisions: Division[];
    team: Team;
    remove_team: () => void;
    submit_team: (team: Team) => void;
};

export function TeamBoard({ isown, divisions, team, submit_team, remove_team}: TeamBoardProps) {
    const [members, setMembers] = useState(team.members);
    const [name, setName] = useState(team.name );
    const [is_mut, setMut] = useState(!team.name);

    useEffect(() => {
        if (is_mut && members.length == 0) {
            setMembers([plain_member()]);
        }
    },[is_mut])

    return (
        <div className="team_scoreboard">
            <div className="team_header">
                {is_mut ? (
                    <>
                    <img 
                        className="team_action_btn" 
                        src={garbage_svg}
                        onClick={remove_team}
                    />
                    <input
                        className="roboto-mono-norm"
                        type="text"
                        value={name}
                        placeholder="TEAM NAME"
                        onChange={(e) => setName(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.shiftKey && e.key == "Enter") {
                                    setMembers(prev => [...prev, plain_member()])
                            } else if (e.key == "Enter") {
                                setMut(false);
                                let mems = members.filter(m => !!m.name);
                                submit_team({...team, 
                                    name: !!name ? name : MASCOT_NAMES[random_uint32() % MASCOT_NAMES.length], 
                                    members: mems
                                });
                                setMembers(mems);
                            }
                        }}
                        autoFocus={!name}
                    />
                    <img 
                        className="team_action_btn" 
                        src={submit_svg}
                        onClick={() => {
                            setMut(false);
                            let mems = members.filter(m => !!m.name);
                            submit_team({...team, 
                                name: !!name ? name : MASCOT_NAMES[random_uint32() % MASCOT_NAMES.length], 
                                members: mems
                            });
                            setMembers(mems);
                        }}
                    />
                    </>
                ) : (
                    <>
                    {isown ? (
                        <img 
                            className="team_action_btn" 
                            src={edit_svg}
                            onClick={() => setMut(true) }
                        />
                        ) : (
                        <div/>
                        ) 
                    }
                        <h2 className="roboto-mono-norm">{name}</h2>
                        <div/>
                    </>
                )}
            </div>
            {members.map((m, i) => (
                <div key={i} className="participant_entry roboto-mono-norm">
                    {is_mut ?
                        <Mut_Member
                            divisions={divisions}
                            member={m}
                            add_member={() =>
                                setMembers([...members, { 
                                    name: "", 
                                    division: 0, 
                                    score: 0, 
                                    x_count: 0 
                                }])
                            }
                            update_member={(m) => {
                                const newMembers = [...members];
                                newMembers[i] = m;
                                setMembers(newMembers);
                            }}
                            remove_member={() => {
                                const newMembers = [...members];
                                newMembers.splice(i, 1);
                                setMembers(newMembers);
                            }}
                            focus={(!!name && i == members.length - 1)}
                            submit={() => {
                                let mems = members.filter(m => !!m.name);
                                submit_team({...team, 
                                    name: !!name ? name : MASCOT_NAMES[random_uint32() % MASCOT_NAMES.length], 
                                    members: mems
                                });
                                setMembers(mems);
                                setMut(false);
                            }}
                        />
                        :
                        <Member
                            divisions={divisions}
                            member={m}
                        />
                    }
                </div>
            ))}
            {is_mut && 
                <div className="adder_entry">
                    <img src={add_svg}
                        onClick={() =>
                            setMembers([...members, { 
                                name: "", 
                                division: 0, 
                                score: 0, 
                                x_count: 0 
                            }])
                        }
                    />
                </div>
            }
        </div>
    );
}
