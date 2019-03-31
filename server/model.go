package main

import (
	"errors"
	"log"
)

const nought = "0"
const cross = "X"
const empty = " "

type Game struct {
	Board      [3][3]string `json:""`
	NextToPlay string       `json:""`
	moveCount  int
	Winner     string `json:",omitempty"`
	Finished   bool   `json:""`
}

var ErrInvalidMoveAlreadyPlayed = errors.New("type: location already played")
var ErrGameFinished = errors.New("type: game Finished")

func NewGame() *Game {
	return &Game{Board: [3][3]string{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	},
		moveCount:  0,
		NextToPlay: cross,
	}
}

func (game *Game) nextPlayer() {
	if game.NextToPlay == cross {
		game.NextToPlay = nought
	} else {
		game.NextToPlay = cross
	}
}

func (game *Game) playMove(row int, column int) bool {
	currentPlayer := game.NextToPlay
	// update the board
	game.Board[row][column] = currentPlayer
	game.moveCount = game.moveCount + 1
	log.Printf("Move %d was player %v at %d,%d", game.moveCount, currentPlayer, row, column)
	// check to see if the game is over
	winner := game.checkForWinner()
	game.nextPlayer()
	if winner {
		log.Printf("Winner found.")
		game.Winner = currentPlayer
		game.Finished = true
	}
	if game.moveCount == 9 {
		log.Printf("Board filled with no Winner.")
		game.Finished = true
	}
	return game.Finished
}

func (game Game) checkForWinner() bool {
	// check to see if someone has won
	// since we do this every turn we only need to check if the player who just moved (currentPlayer) won
	board := &game.Board
	currentPlayer := game.NextToPlay
	for i := 0; i < 3; i++ {
		rowWin := true
		columnWin := true
		for j := 0; j < 3 && (rowWin || columnWin); j++ {
			rowWin = rowWin && board[j][i] == currentPlayer
			columnWin = columnWin && board[i][j] == currentPlayer
		}
		if rowWin || columnWin {
			return true
		}
	}
	// check for diagonal, there's only two fixed possibilites, both of which go through the centre
	// it's easiest just to hard code these
	diagonalOneWin := false
	diagonalTwoWin := false
	if board[1][1] == currentPlayer {
		diagonalOneWin = board[0][0] == currentPlayer && board[2][2] == currentPlayer
		diagonalTwoWin = board[0][2] == currentPlayer && board[2][0] == currentPlayer
	}
	return diagonalOneWin || diagonalTwoWin
}

func (game Game) checkValidMove(row int, column int) error {
	if game.Board[row][column] != empty {
		return ErrInvalidMoveAlreadyPlayed
	}
	if game.Finished {
		return ErrGameFinished
	}
	return nil
}
