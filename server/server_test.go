package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {

	game := NewGame()
	store := RemoteGame{game: game}

	req := createRequest(t, "GET", "/get", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(store.getState)

	handler.ServeHTTP(rr, req)

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
	store := RemoteGame{game: game}

	req := createRequest(t, "POST", "/play", game)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(store.play)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler didn't return badRequest: got %v", status)
	}
}

func TestPlayMoveHandlerWrongPlayer(t *testing.T) {

	game := NewGame()
	store := RemoteGame{game: game}

	move := &Move{
		Row:    0,
		Column: 0,
		Player: nought,
	}

	req := createRequest(t, "POST", "/play", move)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(store.play)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler didn't return forbidden: got %v", status)
	}
}

func TestPlayMoveHandlerAlreadyUsedLocation(t *testing.T) {

	game := NewGame()
	game.Board[0][0] = cross
	store := RemoteGame{game: game}

	move := &Move{
		Row:    0,
		Column: 0,
		Player: nought,
	}

	req := createRequest(t, "POST", "/play", move)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(store.play)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler didn't return forbidden: got %v", status)
	}
}

func TestPlayMoveHandlerGameFinished(t *testing.T) {

	game := NewGame()
	game.Finished = true
	store := RemoteGame{game: game}

	move := &Move{
		Row:    0,
		Column: 0,
		Player: nought,
	}

	req := createRequest(t, "POST", "/play", move)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(store.play)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler didn't return forbidden: got %v", status)
	}
}

func TestPlayMoveHandler(t *testing.T) {

	game := NewGame()
	store := RemoteGame{game: game}

	move := &Move{
		Row:    0,
		Column: 0,
		Player: cross,
	}
	req := createRequest(t, "POST", "/play", move)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(store.play)

	handler.ServeHTTP(rr, req)

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
