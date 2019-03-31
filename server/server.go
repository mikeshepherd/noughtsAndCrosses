package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type RemoteGame struct {
	sync.Mutex
	game *Game
}

type Move struct {
	Row    int    `json:""`
	Column int    `json:""`
	Player string `json:""`
}

func (remoteGame *RemoteGame) play(w http.ResponseWriter, req *http.Request) {
	log.Printf("received move request %v", req)
	// lock the game so we don't get problems
	remoteGame.Lock()
	defer remoteGame.Unlock()

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
	err := decoder.Decode(&m)
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

func (remoteGame RemoteGame) getState(w http.ResponseWriter, req *http.Request) {

	// lock the game so we don't get problems
	remoteGame.Lock()
	defer remoteGame.Unlock()

	// return the game state
	err := json.NewEncoder(w).Encode(remoteGame.game)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func startServer(port int) {
	game := NewGame()
	store := RemoteGame{game: game}

	r := mux.NewRouter()
	r.HandleFunc("/get", store.getState).Methods("GET")
	r.HandleFunc("/move", store.play).Methods("POST")
	http.Handle("/", r)

	log.Printf("Going to listen on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}
