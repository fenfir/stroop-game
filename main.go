package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type colorWord struct {
	Color color.RGBA
	Text  string
	Value int
}

var colorNames []string
var colors []color.RGBA

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	colorNames = []string{
		"RED",
		"WHITE",
		"BLUE",
		"GREEN",
		"YELLOW",
	}

	colors = []color.RGBA{
		colornames.Red,
		colornames.White,
		colornames.Blue,
		colornames.Green,
		colornames.Yellow,
	}

	pixelgl.Run(run)
}

// DrawType indicates the type of redraw that should happen
type DrawType int

const (
	// NewWord specifies that a new word should be drawn
	NewWord int = iota
)

const (
	// MaxTimeout is the maximum number of timeouts until the game exits
	MaxTimeout = 5
	// TimeoutLength is the number of seconds before a timeout occurs
	TimeoutLength = 5
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Stroop Effect!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Panic(err)
	}

	atlas := text.NewAtlas(basicfont.Face7x13,
		text.ASCII,
	)

	drawChan := make(chan int)
	inputChan := make(chan int)

	colorText := text.New(pixel.V(100, 500), atlas)
	messageText := text.New(pixel.V(100, 50), atlas)
	timerText := text.New(pixel.V(100, 100), atlas)
	scoreText := text.New(pixel.V(700, 700), atlas)

	incorrectCount := 0
	score := 0

	timer := NewSecondsTimer(TimeoutLength)
	word := drawNewWord(colorText)

	fps := time.Tick(time.Second / 120)
	for !win.Closed() {
		win.Clear(colornames.Black)

		if win.JustPressed(pixelgl.Key1) {
			go func() {
				inputChan <- 1
			}()
		} else if win.JustPressed(pixelgl.Key2) {
			go func() {
				inputChan <- 2
			}()
		} else if win.JustPressed(pixelgl.Key3) {
			go func() {
				inputChan <- 3
			}()
		} else if win.JustPressed(pixelgl.Key4) {
			go func() {
				inputChan <- 4
			}()
		} else if win.JustPressed(pixelgl.Key5) {
			go func() {
				inputChan <- 5
			}()
		}

		select {
		case _ = <-drawChan:
			log.Printf("New Word\n")
			timer = NewSecondsTimer(TimeoutLength)
			word = drawNewWord(colorText)
		case _ = <-timer.Timer.C:
			log.Printf("Timeout\n")
			setMessageText(messageText, "Timeout!", colornames.Blue)
			incorrectCount++
			go func() {
				drawChan <- NewWord
			}()
		case c := <-inputChan:
			log.Printf("Input received: %d\n", c)
			if c == word.Value {
				log.Println("Correct input")
				score++
				setMessageText(messageText, "Correct!", colornames.Green)
				go func() {
					drawChan <- NewWord
				}()
			} else {
				setMessageText(messageText, "Wrong!", colornames.Red)
			}
		default:
		}

		scoreMessage(scoreText, score, incorrectCount)
		timerMessage(timerText, timer.Remaining())

		colorText.Draw(win, pixel.IM.Scaled(colorText.Orig, 16))
		messageText.Draw(win, pixel.IM.Scaled(messageText.Orig, 2))
		timerText.Draw(win, pixel.IM.Scaled(timerText.Orig, 2))
		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 2))

		if incorrectCount >= MaxTimeout {
			win.SetClosed(true)
		}

		win.Update()
		<-fps
	}
}

func randomColorWord() *colorWord {
	randColorName := rand.Int() % len(colorNames)
	randColor := rand.Int() % len(colors)

	return &colorWord{
		Color: colors[randColor],
		Text:  colorNames[randColorName],
		Value: randColor + 1,
	}
}

func scoreMessage(scoreText *text.Text, score int, wrong int) {
	scoreText.Clear()
	scoreText.Color = colornames.Chartreuse
	fmt.Fprintf(scoreText, "%d words correct\n%d/5 words incorrect", score, wrong)
}

func timerMessage(timerText *text.Text, timeRemaining time.Duration) {
	timerText.Clear()
	timerText.Color = colornames.Red
	fmt.Fprintf(timerText, "%1.1f seconds remaining", timeRemaining.Seconds())
}

func drawNewWord(colorText *text.Text) *colorWord {
	colorText.Clear()
	randColorWord := randomColorWord()
	colorText.Color = randColorWord.Color
	fmt.Fprintln(colorText, randColorWord.Text)

	return randColorWord
}

func setMessageText(messageText *text.Text, text string, color color.RGBA) {
	messageText.Clear()
	messageText.Color = color
	fmt.Fprintf(messageText, text)
}
