package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

type coord struct {
	row, col int
}

type Model struct {
	simulate bool
	grid     [][]int
	cursor   coord
}

func New() Model {
	m := new(Model)

	m.simulate = false

	grid := make([][]int, 12)
	for i := range grid {
		grid[i] = make([]int, 48)
	}
	m.grid = grid

	cursor := coord{5, 5}
	m.cursor = cursor

	return *m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor.row > 0 {
				m.cursor.row--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor.row < len(m.grid) - 1 {
				m.cursor.row++
			}

		// The "left" and "h" keys move the cursor up
		case "left", "h":
			if m.cursor.col > 0 {
				m.cursor.col--
			}

		// The "right" and "l" keys move the cursor down
		case "right", "l":
			if m.cursor.col < len(m.grid[0]) - 1 {
				m.cursor.col++
			}
		// the selected state for the item that the cursor is pointing at.
		case " ":
			row, col := m.cursor.row, m.cursor.col
			if m.grid[row][col] == 1 {
				m.grid[row][col] = 0
			} else {
				m.grid[row][col] = 1
			}
		case "enter":
			if m.simulate {
				m.simulate = false
			} else {
				m.simulate = true
			}
			return m, tick()
		}

	case TickMsg:
		return m.Mutate()
	}

	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/3, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Mutate() (tea.Model, tea.Cmd) {
	if !m.simulate {
		return m, nil
	}
	state := make([][]int, len(m.grid))
	for row := range state {
		state[row] = make([]int, len(m.grid[0]))
	}

	for row := range m.grid {
		for col := range m.grid[row] {
			neighbors := []coord{
				{row-1, col-1},
				{row-1, col},
				{row-1, col+1},
				{row  , col-1},
				{row  , col+1},
				{row+1, col-1},
				{row+1, col},
				{row+1, col+1},
			}

			var count int

			for _, neighbor := range neighbors {
				if neighbor.row < 0 {
					neighbor.row = len(m.grid) - 1
				}
				if neighbor.row > len(m.grid) - 1 {
					neighbor.row = 0
				}
				if neighbor.col < 0 {
					neighbor.col = len(m.grid[0]) - 1
				}
				if neighbor.col > len(m.grid[0]) - 1 {
					neighbor.col = 0
				}
				count += m.grid[neighbor.row][neighbor.col]
			}
			if count == 3 {
				// either dead or alive, the cell is or turns live
				state[row][col] = 1
			} else if m.grid[row][col] == 1 && count == 2 {
				state[row][col] = 1
			} else {
				state[row][col] = 0
			}
		}
	}

	m.grid = state
	return m, tick()
}

func (m Model) View() string {
	var result string

	// Iterate over our choices
	cursor := m.cursor
	for r := range m.grid {
		for c := range m.grid[r] {
			if r == cursor.row && c == cursor.col {
				result += "X"
			} else if m.grid[r][c] == 1 {
				result += "1"
			} else {
				result += "0"
			}
		}
		result += "\n"
	}

	return result 
}

func main() {
	program := tea.NewProgram(New())
	log.Fatal(program.Start())
}
