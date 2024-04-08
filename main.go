package main

import (
    "os"
    "time"
    "math/rand"
    "github.com/gdamore/tcell/v2"
)

var (
    GREEN = tcell.StyleDefault.Foreground(tcell.ColorGreen)
    GREY = tcell.StyleDefault.Foreground(tcell.ColorGrey)
    RED = tcell.StyleDefault.Foreground(tcell.ColorRed)
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Direction int

const (
    Right Direction = iota
    Down
    Left
    Up
)

type Cell struct {
    x, y int
    color tcell.Style
}

type Snake struct {
    direction Direction
    body []Cell
    char rune
    lenght int
}

type Food struct {
    cell Cell
    char rune
}

type Game struct {
    snake Snake
    food Food
    pos struct {
        x, y int
    }
    s_cols int
    s_rows int
}

func (g *Game) MoveSnake() {
    if len(g.snake.body) > 1 {
        for i := len(g.snake.body) - 1; i > 0; i-- {
            g.snake.body[i].x = g.snake.body[i-1].x
            g.snake.body[i].y = g.snake.body[i-1].y
        }
    }

    g.snake.body[0].x = g.pos.x
    g.snake.body[0].y = g.pos.y

    switch g.snake.direction {
    case Right:
        g.pos.x += 2
    case Down:
        g.pos.y += 1
    case Left:
        g.pos.x -= 2
    case Up:
        g.pos.y -= 1
    }

    if g.pos.x > g.s_cols {
        g.pos.x = 0;
    }
    if g.pos.x < 0 {
        g.pos.x = g.s_cols - 1
    }
    if g.pos.y >= g.s_rows {
        g.pos.y = 0;
    }
    if g.pos.y < 0 {
        g.pos.y = g.s_rows - 1
    }
}

func randomPosition(limit int) int {
    result := r.Intn(limit)
    if result % 2 != 0 {
        result--
    }
    return result
}

func (g *Game) NewFood() {
    char_options := []rune{'*', '%', '#', '=', '$', '!', 'X', '+', '~'}
    g.food.char = char_options[r.Intn(len(char_options))]
    g.food.cell.x = randomPosition(g.s_cols)
    g.food.cell.y = randomPosition(g.s_rows)
}

func (g *Game) CheckCollison() {
    if g.pos.x == g.food.cell.x && g.pos.y == g.food.cell.y {
        g.snake.body = append(
            g.snake.body,
            Cell {
                x: g.snake.body[len(g.snake.body) - 1].x,
                y: g.snake.body[len(g.snake.body) - 1].y,
                color: GREY,
            },
        )

        g.NewFood()
    }
}

func draw(g *Game, s tcell.Screen) {
    for {
        g.MoveSnake()
        g.CheckCollison()

        s.Clear()

        for i := range len(g.snake.body) {
            s.SetContent(g.snake.body[i].x, g.snake.body[i].y, g.snake.char, nil, g.snake.body[i].color)
            s.SetContent(g.snake.body[i].x + 1, g.snake.body[i].y, g.snake.char, nil, g.snake.body[i].color)
        }

        s.SetContent(g.food.cell.x, g.food.cell.y, g.food.char, nil, g.food.cell.color)
        s.SetContent(g.food.cell.x + 1, g.food.cell.y, g.food.char, nil, g.food.cell.color)

        s.Show()

        time.Sleep(time.Second / time.Duration(60 / 6))
    }
}

func main() {
    s, _ := tcell.NewScreen()
    s.Init()

    cols, rows := s.Size()

    g := Game {
        snake: Snake {
            char: '█',
            lenght: 5,
            direction: Right,
            body: []Cell {
                {x: 0, y: 0, color: GREEN},
            },
        },
        food: Food {
            char: '*',
            cell: Cell {
                x: 20, y: 5,
                color: RED,
            },
        },
        pos: struct { x, y int }{
            x: randomPosition(cols),
            y: randomPosition(rows),
        },
        s_cols: cols,
        s_rows: rows,
    }

    s.Clear()

    go draw(&g, s)

    for {
        switch ev := s.PollEvent().(type) {
        case *tcell.EventResize:
            s.Sync()
        case *tcell.EventKey:
            switch ev.Key() {
            case tcell.KeyEscape:
                s.Fini()
                os.Exit(0)
            case tcell.KeyRune:
                switch ev.Rune() {
                case 'l':
                    if g.snake.direction == Left {
                        g.snake.direction = Left
                    } else {
                        g.snake.direction = Right
                    }
                case 'h':
                    if g.snake.direction == Right {
                        g.snake.direction = Right
                    } else {
                        g.snake.direction = Left
                    }
                case 'j':
                    if g.snake.direction == Up {
                        g.snake.direction = Up
                    } else {
                        g.snake.direction = Down
                    }
                case 'k':
                    if g.snake.direction == Down {
                        g.snake.direction = Down
                    } else {
                        g.snake.direction = Up
                    }
                }
            }
        }
    }
}
