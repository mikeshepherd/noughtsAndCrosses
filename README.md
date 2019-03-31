# Backend Code Challenge

This is a simple Noughts & Crosses API server.

The API supports multiple concurrent games.

### Building
Go 1.11 is required to build the service.
Run `go build` to build and `go test` to test.

### Playing

To start the service run the executable built by `go build`, passing the server command to it e.g. `./noughtsandcrosses server`.  
Without `server` the executable will provide a simple console version of the game instead.  
The port the server listens on can be configured by adding `-port XXXX` to the command.

To start a new game POST an empty request to `http://localhost:8000/game/new`.
The response will be a JSON object containing the id for the new game e.g. `{"Id":1}`

To see the state of a game GET `http://localhost:8000/game/{id}`.
The response will be a JSON object of the game state.  
E.g. `{"Board":[[" ","X"," "],["0","X","0"],[" "," "," "]],"NextToPlay":"X","Finished":false}`

To make a move POST a JSON move object to `http://localhost:8000/game/{id}/move`.
The move object contains the row and column to play and the player name (X or 0).  
E.g. `{"Row":1, "Column": 1, "Player": "X"}`


## Design

Initially I built a simple console client to allow me to focus on getting the mechanics of the game correct 
and to more easily interact with the model as I developed it without requiring a carefully
orchestrated list of http calls.  
Once this once done I could extract the gameplay mechanics into a separate place and add a http server to interact with it,
 keeping the logic for the game itself independent of the API service or console interface.
The API is very simple, requiring only two endpoints, one to get the state and one to make a move.  
A third was added later to support multiple games at once.
The game state is stored in a global-ish variable shared between every instance of the http handlers,
 due to the shared nature this requires mutexes to avoid race conditions, possible cheating and getting stale state.