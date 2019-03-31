package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type GameStore struct {
	sync.Mutex
	games []*RemoteGame
}

type RemoteGame struct {
	sync.Mutex
	game *Game
}

type GameId struct {
	Id int `json:""`
}

type Move struct {
	Row    int    `json:""`
	Column int    `json:""`
	Player string `json:""`
}

var ErrGameNotFound = errors.New("type: game not found")

func (gameStore *GameStore) play(w http.ResponseWriter, req *http.Request) {
	log.Printf("received move request %v", req)
	// lock the game store so we don't get problems
	gameStore.Lock()
	// get the game that is being played
	remoteGame, err := gameStore.getGame(req)

	if err != nil {
		log.Printf("Move submitted for non-existant game %v", req)
		http.Error(w, "Game not found", 404)
		// don't forget to unlock the global store
		gameStore.Unlock()
		return
	}

	// lock the game
	remoteGame.Lock()
	defer remoteGame.Unlock()
	// now that we have the specific game and locked it we can release the global store
	gameStore.Unlock()

	game := remoteGame.game
	var m Move
	if req.Body == nil {
		log.Printf("invalid move request %v", req)
		http.Error(w, "Please send a request body", 400)
		return
	}

	// read the json body
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields() // don't succeed if we get spurious fields
	err = decoder.Decode(&m)
	if err != nil {
		log.Printf("could not decode request %v", req)
		http.Error(w, err.Error(), 400)
		return
	}
	log.Printf("Got decoded move %v", m)

	// check the right player is making the move
	if m.Player != game.NextToPlay {
		log.Printf("incorrect player tried to move %v", req)
		http.Error(w, "Incorrect player tried to move", 403)
		return
	}

	// check the move is valid
	invalidMoveError := game.checkValidMove(m.Row, m.Column)
	if invalidMoveError != nil {
		log.Printf("invalid move %v, %v", invalidMoveError, req)
		http.Error(w, invalidMoveError.Error(), 403)
		return
	}

	// update the game state
	game.playMove(m.Row, m.Column)

	// return the new state
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (gameStore *GameStore) getState(w http.ResponseWriter, req *http.Request) {

	// lock the store
	gameStore.Lock()
	defer gameStore.Unlock()
	remoteGame, err := gameStore.getGame(req)

	// return the game state
	if err != nil {
		log.Printf("Status requested for non-existant game, %v", req)
		http.Error(w, "Game not found", 404)
		return
	}

	// lock the game
	remoteGame.Lock()
	defer remoteGame.Unlock()

	err = json.NewEncoder(w).Encode(remoteGame.game)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (gameStore *GameStore) newGame(w http.ResponseWriter, req *http.Request) {
	game := NewGame()
	remoteGame := RemoteGame{game: game}

	gameStore.Lock()
	defer gameStore.Unlock()

	newGames := append(gameStore.games, &remoteGame)

	gameStore.games = newGames

	newGameId := GameId{Id: len(gameStore.games)}

	err := json.NewEncoder(w).Encode(newGameId)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (gameStore *GameStore) getGame(req *http.Request) (*RemoteGame, error) {

	// this is a helper function so locking the mutex should be handled by the caller
	vars := mux.Vars(req)
	gameId, _ := strconv.Atoi(vars["id"])

	if gameId > len(gameStore.games) {
		return nil, ErrGameNotFound
	}

	return gameStore.games[gameId-1], nil
}

func startServer(port int) {
	gameStore := GameStore{}

	r := mux.NewRouter()
	r.HandleFunc("/games/{id:[0-9]+}", gameStore.getState).Methods("GET")
	r.HandleFunc("/games/{id:[0-9]+}/move", gameStore.play).Methods("POST")
	r.HandleFunc("/games", gameStore.newGame).Methods("POST")
	http.Handle("/", r)

	log.Printf("Going to listen on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}
