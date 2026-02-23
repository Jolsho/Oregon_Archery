# API

## STRUCTURES

    type Member struct {
        Name     string `json:"name"`
        Division int    `json:"division"`
        Score    int    `json:"score"`
        XCount   int    `json:"x_count"`
    }

    type Team struct {
        Name    string   `json:"name"`
        Score   int      `json:"score"`
        XCount  int      `json:"x_count"`
        Members []Member `json:"members"`
    }

    type Event struct {
        Title     string          `json:"title"`
        IsOwn	  bool			  `json:"is_own"`
        Leaders   map[string]int  `json:"leaders"`
        Divisions []string        `json:"divisions"`
        Teams     map[string]Team `json:"teams"`
        Secret 	  string 		   `json:"-"`
    }

## ENDPOINTS

    /events
        "ALL" {
            make sure the user has a secret value in a cookie
        }

        "GET" {
            return Event list json encoded
        }

        "PUT/POST" {
            decode JSON Event from body

            if event is in Events:
                check if user owns that event(Cookies)
                update that entry

            else:
                set event.secret = cookie.value
                set event.divisions
                set event.is_own = true
                append event to events
            
            return event json encoded 
        }
        "DELETE" {
            get event_title from url     

            if event exists and user owns event remove it
        }
