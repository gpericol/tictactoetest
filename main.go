package main

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"syscall/js"
)

const (
	empty    = " "
	player   = "X"
	computer = "O"
)

var board [3][3]string
var code string

func initBoard() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = empty
		}
	}
}

func getCode() string {
	str := code + "EOFuserArenaStateread"

	hash := sha256.Sum256([]byte(str))
	str = hex.EncodeToString(hash[:])

	hash = sha256.Sum256([]byte(str))
	str = hex.EncodeToString(hash[:])

	hash = sha256.Sum256([]byte(str))
	str = hex.EncodeToString(hash[:])

	return str
}

func getBoard(this js.Value, p []js.Value) interface{} {
	// return a matrix of the board as js.Value
	retArr := js.Global().Get("Array").New()
	for i := 0; i < 3; i++ {
		retArr.SetIndex(i, js.Global().Get("Array").New())
		for j := 0; j < 3; j++ {
			retArr.Index(i).SetIndex(j, board[i][j])
		}
	}
	return retArr
}

func checkWin(player string) bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true
		}
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true
		}
	}
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true
	}
	return false
}

func isBoardFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == empty {
				return false
			}
		}
	}
	return true
}

func computerTurn() {
	var freeSpaces [][2]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == empty {
				freeSpaces = append(freeSpaces, [2]int{i, j})
			}
		}
	}
	space := freeSpaces[rand.Intn(len(freeSpaces))]
	board[space[0]][space[1]] = computer
}

func playTurn(this js.Value, p []js.Value) interface{} {
	if code == "" {
		return "-"
	}
	// verify input that are two parameters, they are int and from 0 to 2
	if len(p) != 2 {
		return "-"
	}
	if p[0].Type() != js.TypeNumber || p[1].Type() != js.TypeNumber {
		return "-"
	}
	if p[0].Int() < 0 || p[0].Int() > 2 || p[1].Int() < 0 || p[1].Int() > 2 {
		return "-"
	}

	x := p[0].Int()
	y := p[1].Int()

	// check if the space is empty
	if board[x][y] != empty {
		return "-"
	}

	board[x][y] = player
	if checkWin(player) {
		return getCode()
	}
	if isBoardFull() {
		return "draw"
	}
	computerTurn()
	if checkWin(computer) {
		return "loss"
	}
	if isBoardFull() {
		return "draw"
	}
	return "+"
}

func resetBoard(this js.Value, p []js.Value) interface{} {
	initBoard()
	return nil
}

func setCode(this js.Value, p []js.Value) interface{} {
	if len(p) != 1 {
		return nil
	}
	code = p[0].String()
	return nil
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	initBoard()
	js.Global().Set("setCode", js.FuncOf(setCode))
	js.Global().Set("resetBoard", js.FuncOf(resetBoard))
	js.Global().Set("getBoard", js.FuncOf(getBoard))
	js.Global().Set("playTurn", js.FuncOf(playTurn))

	<-c
}
