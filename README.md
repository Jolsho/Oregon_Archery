# SERVER

`./main.go`
Initialize the state, which is just a list of Events.
Then the Networker which holds active WebSocket Connections along with a logger, and keys for session cookies.
These two things are parameters to the handlers.


`./handlers/events.go && ./handlers/ws.go` 
There are 3 handlers, the first of which just serves static UI files.
The Second one is a classic REST API where users can create, update, and delete events.
The third upgrades the HTTP connection to a WebSocket connection where we can send live updates.


`./network/middleware.go` 
There is then a middleware which looks for a session cookie. 
If it doesn't exist then it creates one and sends it to the client.
This session cookie is used to establish ownership of events.
That way not anyone could edit an event.

`./networker/cookie.go`
The cookie is a digital signature using the Networkers key.
The data being signed is a hash of a random nonce, and an expiration.
This is just simply to first ensure uniqueness, and then obviously to have lifetimes on cookies.

# UI
`./ui/index.html && ./ui/index.js`
These are the entry point to the ui, and are what is servered at the root of the website.
The html document preloads source files, fonts, andn icons.
The index.js file immediately calls a function called `get_events()`.
This function and others like it(`post_event()`, `delete_event()`) are located in `./ui/src/api.js`.
These are the functions that handle interactions with the server.
Anyway once events are recieved the render functions are called `render_menu()`, and `render_event()`.
The prior handles rendering the menu in the top left, and the latter renders everything else.
These two functions are called sequentially everywhere which is not the most efficient way to do it but its simple.

# NGINX
There is a lot of config in this section.
This establishes a reverse proxy for the go server.
It acts as a TLS boundary, so that data between the client and it are encrypted.
However, once the data is within this system it is now plain text and immediately forwarded to the go server.
Given such, this is where we configure TLS certifications which is automated via CertBot.

There are two `server` blocks in `./nginx/prod/default.conf` and the first is just for certification.
It can also upgrade to https, but its predomantly to prove ownership of the domain so as to obtain a TLS cert.
Then the main `server` block points directly to the go server.
It has a rate limit of 20 packets in a short amount of time. 
It also forwards http headers so the go server can see where the request originates from.

There is then `./nginx/fail2ban` which holds config files for a process that watches nginx logs and bans bad actors.
So for example in the go server if someone tries to access an event they don't own the status code 401 is sent back.
This process will see that and if the same client does that 3 times in 60 seconds they get banned for a few minutes.
The idea here is to just prevent probes/scrapers and bad actors from having unlimited attempts at exploiting the server.

# BASH
`renewal_ssl.sh && w ./start.sh`
Bash is a programming language for the Bourne Again Shell.
The first is ran as a cron job periodically to make sure the TLS certificate does not expire.
The second is a deployment script designed to automate building and running the docker containers for this server.
It also ensures that all the shared volumes between the containers are created before they are launched.
It copies config files that are not useful inside the project directory like fail2ban files.
This script is especially useful for testing, as you can simply run `./start.sh` to start and `./start.sh stop` to stop.

