package main

import (
	"testing"
)

func TestPlayMove(t *testing.T) {
	game := Game{board: [3][3]string{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	},
		nextToPlay: cross,
	}
	game.playMove(1, 1)
	if game.board[1][1] != cross {
		t.Errorf("Board not updated after move correctly")
	}
	game.nextPlayer()
	game.playMove(1, 2)
	if game.board[1][2] != nought {
		t.Errorf("Board not updated after move correctly")
	}
}

func TestWinnerCheck(t *testing.T) {
	for _, player := range []string{nought, cross} {
		winningBoards := [][3][3]string{
			{
				{player, player, player},
				{empty, empty, empty},
				{empty, empty, empty},
			},
			{
				{empty, empty, empty},
				{player, player, player},
				{empty, empty, empty},
			},
			{
				{empty, empty, empty},
				{empty, empty, empty},
				{player, player, player},
			},
			{
				{player, empty, empty},
				{player, empty, empty},
				{player, empty, empty},
			},
			{
				{empty, player, empty},
				{empty, player, empty},
				{empty, player, empty},
			},
			{
				{empty, empty, player},
				{empty, empty, player},
				{empty, empty, player},
			},
			{
				{player, empty, empty},
				{empty, player, empty},
				{empty, empty, player},
			},
			{
				{empty, empty, player},
				{empty, player, empty},
				{player, empty, empty},
			},
		}

		for _, board := range winningBoards {
			game := Game{board: board, nextToPlay: player}
			winner := game.checkForWinner()
			if !winner {
				t.Errorf("Board should be a winner but wasn't\n" + renderBoard(board))
			}
		}
	}
}

func TestNotWinnerCheck(t *testing.T) {
	notWinningBoards := [][3][3]string{
		{
			{cross, empty, empty},
			{empty, nought, empty},
			{empty, empty, empty},
		},
		{
			{cross, cross, empty},
			{empty, nought, empty},
			{empty, empty, empty},
		},
		{
			{cross, cross, nought},
			{empty, empty, empty},
			{empty, empty, empty},
		},
		{
			{cross, cross, empty},
			{nought, empty, empty},
			{empty, empty, empty},
		},
		{
			{cross, cross, nought},
			{empty, nought, cross},
			{empty, empty, empty},
		},
		{
			{cross, cross, nought},
			{empty, empty, cross},
			{empty, empty, nought},
		},
		{
			{cross, cross, nought},
			{empty, nought, cross},
			{empty, nought, cross},
		},
		{
			{cross, cross, nought},
			{nought, nought, cross},
			{empty, empty, cross},
		},
		{
			{cross, cross, nought},
			{empty, empty, cross},
			{nought, nought, cross},
		},
		{
			{cross, cross, nought},
			{nought, nought, cross},
			{empty, nought, cross},
		},
	}

	for _, board := range notWinningBoards {
		game := Game{board: board, nextToPlay: cross}
		crossWinner := game.checkForWinner()
		game = Game{board: board, nextToPlay: nought}
		noughtWinner := game.checkForWinner()
		if crossWinner || noughtWinner {
			t.Errorf("Board should not be a winner but was\n" + renderBoard(board))
		}
	}
}

func TestValidMoves(t *testing.T) {
	game := Game{board: [3][3]string{
		{cross, empty, empty},
		{empty, nought, empty},
		{empty, empty, empty},
	}}

	invalidMoveErr := game.checkValidMove(0, 1)
	if invalidMoveErr != nil {
		t.Errorf("Move should be valid but wasn't")
	}

	invalidMoveErr = game.checkValidMove( 1, 0)
	if invalidMoveErr != nil {
		t.Errorf("Move should be valid but wasn't")
	}
}

func TestInvalidMoves(t *testing.T) {
	game := Game{board: [3][3]string{
		{cross, empty, empty},
		{empty, nought, empty},
		{empty, empty, empty},
	}}

	invalidMoveErr := game.checkValidMove(0, 0)
	if invalidMoveErr == nil {
		t.Errorf("Move should be invalid but wasn't")
	}

	invalidMoveErr = game.checkValidMove( 1, 1)
	if invalidMoveErr == nil {
		t.Errorf("Move should be invalid but wasn't")
	}
}

func TestNextPlayer(t *testing.T) {
	game := Game{nextToPlay:cross}
	game.nextPlayer()
	if game.nextToPlay != nought {
		t.Errorf("The player after cross should be nought")
	}
	game.nextPlayer()
	if game.nextToPlay != cross {
		t.Errorf("The player after nought should be cross")
	}
}
