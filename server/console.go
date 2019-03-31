package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func consoleGame() {
	game := NewGame()
	showBoard(game.Board)
	// repeat till game completed
	for !game.Finished {
		currentPlayer := game.NextToPlay
		fmt.Printf("It is player %v's go", currentPlayer)
		// get move coordinates
		row, column := getMove(*game)
		// update the state
		finished := game.playMove(row, column)
		// if finished display a message
		if finished {
			if game.Winner != "" {
				fmt.Println(currentPlayer + " has won!!!")
			} else {
				fmt.Println("The Board is full and no-one has won")
			}
		}
		showBoard(game.Board)
	}
}

func showBoard(board [3][3]string) {
	fmt.Println(renderBoard(board))
}

func renderBoard(board [3][3]string) string {
	var sb strings.Builder
	sb.WriteString("-------\n")
	for i := 0; i < 3; i++ {
		sb.WriteString("|" + strings.Join(board[i][:], "|") + "|\n")
	}
	sb.WriteString("-------\n")
	return sb.String()
}

func getMove(game Game) (int, int) {
	// get row and column to play
	row := getSingleCoord("row")
	column := getSingleCoord("column")
	// check that the move is okay and ask again if not
	invalidMoveError := game.checkValidMove(row, column)
	if invalidMoveError != nil {
		// it shouldn't be possible to get a game already finished error on the console version
		fmt.Println("That location has already been played")
		showBoard(game.Board)
		return getMove(game)
	}
	return row, column
}

func getSingleCoord(directionName string) int {

	fmt.Println("Please enter the " + directionName + " to place your counter")
	reader := bufio.NewReader(os.Stdin)

	value, _, _ := reader.ReadRune()
	if !(value == '0' || value == '1' || value == '2') {
		fmt.Println("Invalid input. Please enter either 0, 1 or 2")
		return getSingleCoord(directionName)
	}
	return int(value) - '0'
}
