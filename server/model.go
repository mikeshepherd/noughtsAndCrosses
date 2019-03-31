package main

import "errors"

const nought = "0"
const cross = "X"
const empty = " "

type Game struct {
	board      [3][3]string
	nextToPlay string
	moveCount  int
	winner     string
	finished   bool
}

var ErrInvalidMoveAlreadyPlayed = errors.New("type: location already played")
var ErrGameFinished = errors.New("type: game finished")

func NewGame() *Game {
	return &Game{board: [3][3]string{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	},
		moveCount:  0,
		nextToPlay: cross,
	}
}

func (game *Game) nextPlayer() {
	if game.nextToPlay == cross {
		game.nextToPlay = nought
	} else {
		game.nextToPlay = cross
	}
}

func (game *Game) playMove(row int, column int) bool {
	// update the board
	game.board[row][column] = game.nextToPlay
	// increase move count
	game.moveCount = game.moveCount + 1
	winner := game.checkForWinner()
	// check to see if the game is over
	if winner {
		game.winner = game.nextToPlay
		game.finished = true
	}
	if game.moveCount == 9 {
		game.finished = true
	}
	return game.finished
}

func (game Game) checkForWinner() bool {
	// check to see if someone has won
	// since we do this every turn we only need to check if the player who just moved (currentPlayer) won
	board := &game.board
	currentPlayer := game.nextToPlay
	// check for horizontal and vertical lines
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
	if game.board[row][column] != empty {
		return ErrInvalidMoveAlreadyPlayed
	}
	if game.finished {
		return ErrGameFinished
	}
	return nil
}
