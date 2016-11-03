package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/buger/goterm"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const ConfigFile = "config.json"

type Config struct {
	Rows      int    `json:"rows"`
	Columns   int    `json:"columns"`
	Separator string `json:"separator"`
	Alive     string `json:"alive"`
	Dead      string `json:"dead"`
	Seed      int64  `json:"seed,omitempty"`
	Interval  int64  `json:"interval"`
}

var config Config
var board [][]int

func printBoard() {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			var symbol string
			if board[i][j] == 1 {
				symbol = config.Alive
			} else {
				symbol = config.Dead
			}
			goterm.Printf("%s%s", symbol, config.Separator)
		}
		goterm.Println("")
	}
}

func outOfBoundsOrValue(i, j int) int {
	if i < 1 || i > config.Rows-1 || j < 1 || j > config.Columns-1 {
		return 0
	}
	return board[i][j]
}

func neighborCount(i, j int) int {
	var x00, x01, x02,
		x10, x12,
		x20, x21, x22 int
	x00 = outOfBoundsOrValue(i-1, j-1)
	x01 = outOfBoundsOrValue(i-1, j)
	x02 = outOfBoundsOrValue(i-1, j+1)

	x10 = outOfBoundsOrValue(i, j-1)
	x12 = outOfBoundsOrValue(i, j+1)

	x20 = outOfBoundsOrValue(i+1, j-1)
	x21 = outOfBoundsOrValue(i+1, j)
	x22 = outOfBoundsOrValue(i+1, j+1)
	return x00 + x01 + x02 + x10 + x12 + x20 + x21 + x22
}

func aliveCase(count int) int {
	if count < 2 || count > 3 {
		return 0
	}
	return 1
}

func deadCase(count int) int {
	if count == 3 {
		return 1
	}
	return 0
}

func nextFieldValue(i, j int) int {
	var nextValue int
	isAlive := board[i][j] == 1
	count := neighborCount(i, j)
	if isAlive {
		nextValue = aliveCase(count)
	} else {
		nextValue = deadCase(count)
	}
	return nextValue
}

func computeNextBoard() [][]int {
	newBoard := make([][]int, config.Rows)
	for i := 0; i < len(board); i++ {
		newBoard[i] = make([]int, config.Columns)
		for j := 0; j < len(board[i]); j++ {
			newBoard[i][j] = nextFieldValue(i, j)
		}
	}
	return newBoard
}

func initBoard() {
	rand.Seed(config.Seed)
	fmt.Printf("Generating board with seed: %d\n", config.Seed)
	board = make([][]int, config.Rows)
	for i := 0; i < len(board); i++ {
		board[i] = make([]int, config.Columns)
		for j := 0; j < len(board[i]); j++ {
			board[i][j] = rand.Intn(2)
			fmt.Printf("%d%s", board[i][j], config.Separator)
		}
		fmt.Println("")
	}
}

func initConfig() {
	rawConfig, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Using %s\n%s", ConfigFile, rawConfig)
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Fatal(err)
	}
	if config.Seed == 0 {
		fmt.Println("Generating random seed")
		config.Seed = time.Now().UnixNano()
	}
}

func init() {
	fmt.Println("Welcome to the game of life")
	initConfig()
	initBoard()
}

func main() {
	fmt.Println("Press enter to begin")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	goterm.Flush()
	goterm.Clear()
	iteration := 0
	for {
		iteration++
		goterm.MoveCursor(1, 1)
		goterm.Printf("Iteration: %d\n", iteration)
		printBoard()
		goterm.Flush()
		board = computeNextBoard()
		time.Sleep(time.Duration(config.Interval) * time.Millisecond)
	}
}
