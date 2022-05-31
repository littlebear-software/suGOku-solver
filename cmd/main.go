package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type coordinate [2]int
type row [9]int
type col [9]int
type square [9]coordinate
type board [9][9]int

var boardState int = 0

const (
	ROW int = 0
	COL int = 1
)

var squares [3][3]square = [3][3]square{
	{
		square{
			{0, 0},
			{0, 1},
			{0, 2},
			{1, 0},
			{1, 1},
			{1, 2},
			{2, 0},
			{2, 1},
			{2, 2},
		},
		square{
			{0, 3},
			{0, 4},
			{0, 5},
			{1, 3},
			{1, 4},
			{1, 5},
			{2, 3},
			{2, 4},
			{2, 5},
		},
		square{
			{0, 6},
			{0, 7},
			{0, 8},
			{1, 6},
			{1, 7},
			{1, 8},
			{2, 6},
			{2, 7},
			{2, 8},
		},
	},
	{
		square{
			{3, 0},
			{3, 1},
			{3, 2},
			{4, 0},
			{4, 1},
			{4, 2},
			{5, 0},
			{5, 1},
			{5, 2},
		},
		square{
			{3, 3},
			{3, 4},
			{3, 5},
			{4, 3},
			{4, 4},
			{4, 5},
			{5, 3},
			{5, 4},
			{5, 5},
		},
		square{
			{3, 6},
			{3, 7},
			{3, 8},
			{4, 6},
			{4, 7},
			{4, 8},
			{5, 6},
			{5, 7},
			{5, 8},
		},
	},
	{
		square{
			{6, 0},
			{6, 1},
			{6, 2},
			{7, 0},
			{7, 1},
			{7, 2},
			{8, 0},
			{8, 1},
			{8, 2},
		},
		square{
			{6, 3},
			{6, 4},
			{6, 5},
			{7, 3},
			{7, 4},
			{7, 5},
			{8, 3},
			{8, 4},
			{8, 5},
		},
		square{
			{6, 6},
			{6, 7},
			{6, 8},
			{7, 6},
			{7, 7},
			{7, 8},
			{8, 6},
			{8, 7},
			{8, 8},
		},
	},
}

func main() {
	// initialize board
	var board [9][9]int
	var err error

	boardstring := os.Args[1]
	board = readBoard(boardstring)
	// board, err = generateBoard()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("board state: original")
	draw(board)

	solution, valid := solve(board)
	if valid {
		draw(solution)
	} else {
		fmt.Println("failed to solve..........")
	}
}

// generate a new puzzle
func generateBoard() (board, error) {
	var err error
	var board struct {
		Cells board `json:"board"`
	}
	var data []byte
	url := "https://sugoku.herokuapp.com/board?difficulty=random"
	if resp, err := http.Get(url); err != nil {
		return board.Cells, err
	} else {
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			return board.Cells, err
		}
		if resp.StatusCode != http.StatusOK {
			return board.Cells, errors.New(string(data))
		}
		err = json.Unmarshal(data, &board)
	}

	return board.Cells, err
}

func readBoard(boardstring string) (b board) {
	boardchars := strings.Split(boardstring, "")
	row := 0
	col := 0
	for _, char := range boardchars {
		val, _ := strconv.Atoi(char)
		b[row][col] = val
		col++
		if col > 8 {
			col = 0
			row++
		}
	}
	return b
}

// draw the board
func draw(board board) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")

	for _, row := range board {
		rowvals := []string{}
		for _, col := range row {
			val := fmt.Sprint(col)
			if col == 0 {
				val = " "
			}
			rowvals = append(rowvals, val)
		}
		table.Append(rowvals)
	}
	table.Render()
}

// pick the next empty spot
func pickASpot(board board) coordinate {
	for row, cols := range board {
		for col, value := range cols {
			if value == 0 {
				return coordinate{row, col}
			}
		}
	}
	return coordinate{-1, -1}
}

// choose possible numbers
func chooseNextPossible(board board, spot coordinate) []int {
	var options []int
	// check for possibility of numbers min-9
	for choice := 1; choice < 10; choice++ {
		if checkRow(board, spot, choice)+
			checkCol(board, spot, choice)+
			checkSquare(board, spot, choice) == 3 {
			options = append(options, choice)
		}
	}
	return options
}

// check for any conflicts in row
func checkRow(board board, spot coordinate, value int) int {
	for _, col := range board[spot[ROW]] {
		if value == col {
			return -1
		}
	}
	return 1
}

// check for any conflicts in column
func checkCol(board board, spot coordinate, value int) int {
	for _, row := range board {
		if row[spot[COL]] == value {
			return -1
		}
	}
	return 1
}

// check for conflicts in square
func checkSquare(board board, spot coordinate, value int) int {
	var vertical, horizontal int
	var square square
	switch {
	case spot[ROW] < 3:
		vertical = 0
	case spot[ROW] > 5:
		vertical = 2
	default:
		vertical = 1
	}
	switch {
	case spot[COL] < 3:
		horizontal = 0
	case spot[COL] > 5:
		horizontal = 2
	default:
		horizontal = 1
	}
	square = squares[vertical][horizontal]
	for _, coord := range square {
		if board[coord[ROW]][coord[COL]] == value {
			return -1
		}
	}
	return 1
}

func solve(b board) (board, bool) {
	/*
		steps:
			1. pick a spot
			2. check for possible numbers
			3. for each possible update board state
			4. for each board state not invalid, repeat 1-3
	*/

	boardState++
	var nextSpot coordinate
	var noSpot = coordinate{-1, -1}
	var boardStates []board

	nextSpot = pickASpot(b)
	possibleValues := chooseNextPossible(b, nextSpot)
	if nextSpot == noSpot || len(possibleValues) == 0 {
		return b, false
	}
	for _, possibleValue := range possibleValues {
		newBoard := b
		newBoard[nextSpot[ROW]][nextSpot[COL]] = possibleValue
		boardStates = append(boardStates, newBoard)
		fmt.Println("board generation:", boardState)
		draw(newBoard)
	}

	for _, bs := range boardStates {
		solution, valid := solve(bs)
		if valid {
			draw(bs)
			b = solution
		}
	}

	return b, true
}
