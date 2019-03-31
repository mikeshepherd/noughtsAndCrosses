package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {

	game := NewGame()
	store := createGameStoreWithGame(game)

	req := createRequest(t, "GET", "/game/1/get", nil)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/game/{id:[0-9]+}/get", store.getState)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler didn't return ok: got %v", status)
	}

	jsonResponse := &Game{}
	unmarshalResponse(t, rr, jsonResponse)

	expected := Game{
		Board: [3][3]string{
			{empty, empty, empty},
			{empty, empty, empty},
			{empty, empty, empty},
		},
		NextToPlay: cross,
		Finished:   false,
	}
	if *jsonResponse != expected {
		t.Errorf("handler returned incorrect body: got %v wanted %v",
			jsonResponse, expected)
	}
}

func TestCreateAndGet(t *testing.T) {

	router := mux.NewRouter()

	store := &GameStore{}

	router.HandleFunc("/new", store.newGame)
	router.HandleFunc("/game/{id:[0-9]+}/get", store.getState)

	req := createRequest(t, "GET", "/new", nil)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	req = createRequest(t, "GET", "/game/1/get", nil)

	rr = httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler didn't return ok: got %v", status)
	}

	jsonResponse := &Game{}
	unmarshalResponse(t, rr, jsonResponse)

	expected := Game{
		Board: [3][3]string{
			{empty, empty, empty},
			{empty, empty, empty},
			{empty, empty, empty},
		},
		NextToPlay: cross,
		Finished:   false,
	}
	if *jsonResponse != expected {
		t.Errorf("handler returned incorrect body: got %v wanted %v",
			jsonResponse, expected)
	}
}

func TestPlayMoveHandlerBadRequest(t *testing.T) {

	game := NewGame()
	store := createGameStoreWithGame(game)

	req := createRequest(t, "POST", "/game/1/move", game)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/game/{id:[0-9]+}/move", store.play)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler didn't return badRequest: got %v", status)
	}
}

func TestPlayMoveHandlerWrongPlayer(t *testing.T) {

	game := NewGame()
	store := createGameStoreWithGame(game)

	move := &Move{
		Row:    0,
		Column: 0,
		Player: nought,
	}

	req := createRequest(t, "POST", "/game/1/move", move)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/game/{id:[0-9]+}/move", store.play)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler didn't return forbidden: got %v", status)
	}
}

func TestPlayMoveHandlerAlreadyUsedLocation(t *testing.T) {

	game := NewGame()
	game.Board[0][0] = cross
	store := createGameStoreWithGame(game)

	move := &Move{
		Row:    0,
		Column: 0,
		Player: nought,
	}

	req := createRequest(t, "POST", "/game/1/move", move)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/game/{id:[0-9]+}/move", store.play)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler didn't return forbidden: got %v", status)
	}
}

func TestPlayMoveHandlerGameFinished(t *testing.T) {

	game := NewGame()
	game.Finished = true
	store := createGameStoreWithGame(game)

	move := &Move{
		Row:    0,
		Column: 0,
		Player: nought,
	}

	req := createRequest(t, "POST", "/game/1/move", move)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/game/{id:[0-9]+}/move", store.play)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler didn't return forbidden: got %v", status)
	}
}

func TestPlayMoveHandler(t *testing.T) {

	game := NewGame()
	store := createGameStoreWithGame(game)

	move := &Move{
		Row:    0,
		Column: 0,
		Player: cross,
	}

	req := createRequest(t, "POST", "/game/1/move", move)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/game/{id:[0-9]+}/move", store.play)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler didn't return ok: got %v", status)
	}

	jsonResponse := &Game{}
	unmarshalResponse(t, rr, jsonResponse)

	expected := Game{
		Board: [3][3]string{
			{cross, empty, empty},
			{empty, empty, empty},
			{empty, empty, empty},
		},
		NextToPlay: nought,
		Finished:   false,
	}

	if *jsonResponse != expected {
		t.Errorf("handler returned incorrect body: got %v wanted %+v",
			rr.Body.String(), expected)
	}
}

func createGameStoreWithGame(game *Game) *GameStore {

	remoteGame := &RemoteGame{game: game}

	return &GameStore{games: []*RemoteGame{remoteGame}}
}

func createRequest(t *testing.T, method string, url string, v interface{}) *http.Request {
	jsonMove, _ := json.Marshal(v)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonMove))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func unmarshalResponse(t *testing.T, rr *httptest.ResponseRecorder, v interface{}) {
	err := json.Unmarshal(rr.Body.Bytes(), v)

	if err != nil {
		t.Errorf("Could not unmarshal response %v, error %v", rr.Body.String(), err)
	}
}
