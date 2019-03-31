package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const nought = "0"
const cross = "X"
const empty = " "

var ErrInvalidMoveAlreadyPlayed = errors.New("type: location already played")

func main() {
	var board = [3][3]string{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	}
	showBoard(board)
	currentPlayer := cross
	finished := false
	moveCount := 0
	for !finished {
		fmt.Printf("It is player %v's go", currentPlayer)
		// get move coordinates
		row, column := getMove(board)
		// update the board
		winner := playMove(&board, row, column, currentPlayer)
		// end if someone won
		if winner {
			fmt.Println(currentPlayer + " has won!!!")
			finished = true
		}
		// end if board is full
		moveCount = moveCount + 1
		if moveCount == 9 && !winner {
			fmt.Println("The board is full and no-one has won")
			finished = true
		}
		// next players turn
		currentPlayer = nextPlayer(currentPlayer)
		// display current state
		showBoard(board)
	}
}

func nextPlayer(currentPlayer string) string {
	if currentPlayer == cross {
		return nought
	} else {
		return cross
	}
}

func playMove(board *[3][3]string, row int, column int, currentPlayer string) bool {
	// update the board
	board[row][column] = currentPlayer
	return checkForWinner(*board, currentPlayer)
}

func checkForWinner(board [3][3]string, currentPlayer string) bool {
	// check to see if someone has won
	// since we do this every turn we only need to check if the player who just moved (currentPlayer) won
	winner := false
	// check for horizontal and vertical lines
	for i := 0; i < 3 && !winner; i++ {
		rowWin := true
		columnWin := true
		for j := 0; j < 3 && (rowWin || columnWin); j++ {
			rowWin = rowWin && board[j][i] == currentPlayer
			columnWin = columnWin && board[i][j] == currentPlayer
		}
		winner = rowWin || columnWin
	}
	if winner {
		return true
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

func getMove(board [3][3]string) (int, int) {
	// get row and column to play
	row := getSingleCoord("row")
	column := getSingleCoord("column")
	// check that the move is okay and ask again if not
	invalidMoveError := checkValidMove(board, row, column)
	if invalidMoveError != nil {
		// there is only one possible invalid move, trying to play where there is already a mark
		fmt.Println("That location has already been played")
		showBoard(board)
		return getMove(board)
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

func checkValidMove(board [3][3]string, row int, column int) error {
	if board[row][column] != empty {
		return ErrInvalidMoveAlreadyPlayed
	}
	return nil
}
