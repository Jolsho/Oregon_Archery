# LANGUAGES / TOOLS

## BACKEND
- Python
  - Django -- http/server
  - Pandas -- data/sorting

## FRONTEND
- TS || JS
- HTML/CSS

---
# API

## STRUCTURES

    type ParticipantID = (team_name, name)
    
    class PARTICIPANT {
        name: STRING,
        team: STRING,
        division: STRING,
        score: int(default 0),
        x_hits: int(default 0)
    }

    class TEAM {
        points: int
        members: []PARTICIPANT -- ORDERED by SCORE/XHITS DESCENDING
    }

    class LEADERBOARD: {
        divisions: DICT(division: []PARTICIPANTID )
        -- ORDERED by SCORE/XHITS DESCENDING
            -- keep list limited to top 3 or 5 or something
    }

    class EVENT {
        creator: IP_ADDRESS,
        title: STRING,
        city: STRING,
        date: STRING(dd/mm/yyyy),
        team_points: DICT(team_name: int)
        leaderboard: LEADERBOARD
    }

## ENDPOINTS

    "/events":
        GET:
            parse url, and if it has "event_title" return that specific event(EVENT).
            else return all EVENTS name & city & date ONLY!

        POST/PUT:
            body is a json encoded EVENT STRUCT where each field is the field name(lowercase)
            usings the title field check if it already exists 

            if it does NOT and the method is POST:
                add it to the Events Dictionary.
            else if method == PUT and it IP matches creator:
                edit the event in place
            else:
                reject and return error

            - IMPORTANT!!
                make sure to capture the IP address of the person creating the event.
                We will use that to allow people to edit the event, and add scores/teams/etc.



    "/events/team"
        POST/PUT:
            grab event "title" from url
            parse json encoded TEAM STRUCT from body and insert into EVENT->TEAMS DICT
                again check if it exists first and if it does send back and error indicating that

            MAKE SURE TO CHECK IP AGAINST EVENT CREATOR



    "/events/team/member"
        POST/PUT:
            same idea as above but this time we use name as the key and insert into TEAM->MEMBERS
            This time you need to insert and shift the LIST to maintain order
            Also "division" field will come with name, so just set those( team, name, division )

            MAKE SURE TO CHECK IP AGAINST EVENT CREATOR



    "/events/score"
        POST/PUT:
            given "team_name", "name", "score", "x_hits"
                find the participant and set score and xhits.
                from there go to leaderboard and see if they made it.
                if they did insert it and potentially bump others.

            MAKE SURE TO CHECK IP AGAINST EVENT CREATOR

