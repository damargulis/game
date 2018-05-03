package game

import (
	"bufio"
	"fmt"
	"github.com/damargulis/game/interfaces"
	"os"
	"strconv"
	"strings"
)

func readInts(prompt string) []int {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	ints := strings.Split(strings.TrimSpace(text), ",")
	nums := make([]int, len(ints))
	for i, num := range ints {
		x, _ := strconv.Atoi(num)
		nums[i] = x
	}
	return nums
}

func isInside(g game.Game, row, col int) bool {
	maxRow, maxCol := g.GetBoardDimensions()
	return row >= 0 && row < maxRow && col >= 0 && col < maxCol
}
