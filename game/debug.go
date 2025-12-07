// debug.go provides a simple logging system
package game

import (
	"fmt"
	"time"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var defaultLogger = &Logger{textFormat: nil}

// Log level to Pico-8 color
// Debug 12
// Info 7
// Warn 9
// Error 8
var LogLevelToColor map[int]color.Color = map[int]color.Color{
	0: color.RGBA{41, 173, 255, 255},
	1: color.RGBA{255, 241, 232, 255},
	2: color.RGBA{255, 163, 0, 255},
	3: color.RGBA{255, 0, 77, 255},
}

type LogEntry struct {
	timestamp int64
	msg       string
	col       color.Color
}

type Logger struct {
	textFormat *text.GoTextFaceSource
	history    []LogEntry
}

func InitLogger(textFormat *text.GoTextFaceSource) {
	defaultLogger = &Logger{
		textFormat: textFormat,
		history:    make([]LogEntry, 0),
	}
}

func (l *Logger) Print(target *ebiten.Image, msg string, x, y float64, col color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(target, msg, &text.GoTextFace{
		Source: l.textFormat,
		Size:   5,
	}, op)
}

func (l *Logger) Render(target *ebiten.Image, x, y float64, h int) {
	lines := min(len(l.history)-1, h-1)
	for i := lines; i >= 0; i -= 1 {
		l.Print(target, l.history[i].msg, x, y+float64((lines-i)*7), l.history[i].col)
	}

	// trim history to max size
	if len(l.history) > h {
		l.history = l.history[:h]
	}
}

func (l *Logger) NewEntry(level int, msg string) {
	// insert the log entry at the top of the history
	l.history = append([]LogEntry{{
		timestamp: time.Now().Unix(),
		msg:       msg,
		col:       LogLevelToColor[level],
	}}, l.history...)
}

func (l *Logger) Debug(s string) {
	l.NewEntry(0, s)
}

func (l *Logger) Info(s string) {
	l.NewEntry(1, s)
}

// alias to Info
func (l *Logger) Log(s string) {
	l.NewEntry(1, s)
}

func (l *Logger) Warn(s string) {
	l.NewEntry(2, s)
}

func (l *Logger) Error(s string) {
	l.NewEntry(3, s)
}

func Log(msg string) {
	defaultLogger.Log(msg)
}

func LogDebug(msg string) {
	defaultLogger.Debug(msg)
}

func LogInfo(msg string) {
	defaultLogger.Info(msg)
}

func LogWarn(msg string) {
	defaultLogger.Warn(msg)
}

func LogError(msg string) {
	defaultLogger.Error(msg)
}

func RenderLogger(target *ebiten.Image, x, y float64, h int) {
	defaultLogger.Render(target, x, y, h)
}

type LogData struct {
	Key   string
	Value any
}

func PrintData(target *ebiten.Image, x, y float64, data []LogData) {
	max_length := 0
	for _, d := range data {
		if len(d.Key) > max_length {
			max_length = len(d.Key)
		}
	}
	// print the data using the default logger in the format:
	//     a short key: value
	// a very long key: value
	for i, d := range data {
		defaultLogger.Print(target, fmt.Sprintf("%s: %v", d.Key, d.Value), x+float64((max_length-len(d.Key))*6), y+float64(i*7), LogLevelToColor[0])
	}
}
