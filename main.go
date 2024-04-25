package main

import (
    "os"
    "fmt"
    "time"
    "math/rand"

    "github.com/gdamore/tcell/v2"
)

var (
    GREEN = tcell.StyleDefault.Foreground(tcell.ColorGreen)
    GREY = tcell.StyleDefault.Foreground(tcell.ColorGrey)
    RED = tcell.StyleDefault.Foreground(tcell.ColorRed)

    WALL = false
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
    over bool
}

func (g *Game) ChangeSnakeDir(d Direction) {
    switch d {
    case Right:
        if g.snake.direction == Left {
            g.snake.direction = Left
        } else {
            g.snake.direction = Right
        }
    case Left:
        if g.snake.direction == Right {
            g.snake.direction = Right
        } else {
            g.snake.direction = Left
        }
    case Down:
        if g.snake.direction == Up {
            g.snake.direction = Up
        } else {
            g.snake.direction = Down
        }
    case Up:
        if g.snake.direction == Down {
            g.snake.direction = Down
        } else {
            g.snake.direction = Up
        }
    }
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

    if g.pos.x >= g.s_cols {
        g.pos.x = 0;
    }
    if g.pos.x < 0 {
        g.pos.x = g.s_cols
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
    char_options := []rune{'%', '#', '=', '$', '!', 'X', '+', '~', ':'}
    g.food.char = char_options[r.Intn(len(char_options))]

    if WALL {
        foodX := randomPosition(g.s_cols)
        if foodX == 0 {
            foodX += 2
        }
        if foodX == g.s_cols - 2 {
            foodX -= 2
        }

        foodY := randomPosition(g.s_rows)
        if foodY == 0 {
            foodY += 1
        }
        if foodY == g.s_rows {
            foodY -= 1
        }

        g.food.cell.x = foodX
        g.food.cell.y = foodY
    } else {
        g.food.cell.x = randomPosition(g.s_cols)
        g.food.cell.y = randomPosition(g.s_rows)
    }
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

    for i := range len(g.snake.body) {
        if WALL {
            if g.pos.x == g.s_cols - 2 || g.pos.x == 0 || g.pos.y == g.s_rows - 1 || g.pos.y == 0 {
                g.over = true
            }
        }

        if g.pos.x == g.snake.body[i].x && g.pos.y == g.snake.body[i].y {
            g.over = true
        }
    }
}

func draw(g *Game, s tcell.Screen) {
    for {
        s.Clear()

        if g.over {
            msg1 := "Game Over"
            msg2 := fmt.Sprintf("%v food collected", len(g.snake.body) - 1)
            msg3 := "Press [ESC] to exit"
            msg4 := "Press [r] to try again"

            for i, r := range msg1 {
                s.SetContent((g.s_cols / 2) - (len(msg1) / 2) + i, g.s_rows / 2 - 2, r, nil, RED.Bold(true))
            }

            for i, r := range msg2 {
                s.SetContent((g.s_cols / 2) - (len(msg2) / 2) + i, g.s_rows / 2 - 1, r, nil, tcell.StyleDefault)
            }

            for i, r := range msg3 {
                s.SetContent((g.s_cols / 2) - (len(msg3) / 2) + i, g.s_rows / 2 + 1, r, nil, tcell.StyleDefault)
            }

            for i, r := range msg4 {
                s.SetContent((g.s_cols / 2) - (len(msg4) / 2) + i, g.s_rows / 2 + 2, r, nil, tcell.StyleDefault)
            }
        } else {
            g.MoveSnake()
            g.CheckCollison()

            if WALL {
                for i := 0; i < g.s_cols - 1; i++ {
                    for j := 0; j < g.s_rows; j++ {
                        if j == 0 || j == g.s_rows - 1 || i == 0 || i == g.s_cols - 2 {
                            s.SetContent(i, j, '░', nil, tcell.StyleDefault)
                            s.SetContent(i + 1, j, '░', nil, tcell.StyleDefault)
                        }
                    }
                }
            }

            for i := range len(g.snake.body) {
                s.SetContent(g.snake.body[i].x, g.snake.body[i].y, g.snake.char, nil, g.snake.body[i].color)
                s.SetContent(g.snake.body[i].x + 1, g.snake.body[i].y, g.snake.char, nil, g.snake.body[i].color)
            }

            s.SetContent(g.food.cell.x, g.food.cell.y, g.food.char, nil, g.food.cell.color)
            s.SetContent(g.food.cell.x + 1, g.food.cell.y, g.food.char, nil, g.food.cell.color)

        }

        s.Show()

        time.Sleep(time.Second / time.Duration(60 / 6))
    }
}

func usage(program string) {
    fmt.Printf ("USAGE: %s <options>\n", program)
    fmt.Println("To control the snake, use the arrow keys or vim movements")
    fmt.Println("")
    fmt.Println("Options:")
    fmt.Println("   -h | --help: show program usage")
    fmt.Println("   -w | --wall: you loose if snake hist the wall")
}

func main() {
    args := os.Args

    if len(args) > 1 {
        switch args[1] {
        case "-h", "--help":
            usage(args[0])
            os.Exit(0)
        case "-w", "--wall":
            WALL = true
        default:
            fmt.Fprintf(os.Stderr, "Incorrect argument: %v\n\n", args[1])
            usage(args[0])
            os.Exit(1)
        }
    }

    s, err := tcell.NewScreen()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(1)
    }

    if err := s.Init(); err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(1)
    }

    cols, rows := s.Size()

    if cols % 2 != 0 {
        cols--
    }

    g := Game {
        snake: Snake {
            char: '█',
            direction: Right,
            body: []Cell {
                {x: 0, y: 0, color: GREEN},
            },
        },
        food: Food {
            char: '*',
            cell: Cell {
                x: 20, y: 5,
                color: RED.Bold(true),
            },
        },
        pos: struct { x, y int }{
            x: 2,
            y: randomPosition(rows - 1) + 1,
        },
        s_cols: cols,
        s_rows: rows,
        over: false,
    }

    s.Clear()

    go draw(&g, s)

    for {
        switch ev := s.PollEvent().(type) {
        case *tcell.EventResize:
            s.Sync()
            cols, rows := s.Size()

            if cols % 2 != 0 {
                cols--
            }

            g.s_cols = cols
            g.s_rows = rows
        case *tcell.EventKey:
            switch ev.Key() {
            case tcell.KeyEscape:
                s.Fini()
                os.Exit(0)
            case tcell.KeyRight:
                g.ChangeSnakeDir(Right)
            case tcell.KeyLeft:
                g.ChangeSnakeDir(Left)
            case tcell.KeyDown:
                g.ChangeSnakeDir(Down)
            case tcell.KeyUp:
                g.ChangeSnakeDir(Up)
            case tcell.KeyRune:
                switch ev.Rune() {
                case 'l':
                    g.ChangeSnakeDir(Right)
                case 'h':
                    g.ChangeSnakeDir(Left)
                case 'j':
                    g.ChangeSnakeDir(Down)
                case 'k':
                    g.ChangeSnakeDir(Up)
                case 'r':
                    if g.over {
                        g.snake.body = []Cell {
                            {x: 0, y: 0, color: GREEN},
                        }
                        g.pos = struct { x, y int }{
                            x: 2,
                            y: randomPosition(g.s_rows - 1) + 1,
                        }
                        g.snake.direction = Right
                        g.over = false
                    }
                }
            }
        }
    }
}
