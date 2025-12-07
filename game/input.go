// input.go provides a centralized system to handle mouse and keyboard inputs
package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button struct {
	WasDown      bool
	IsDown       bool
	Pressed      bool
	Released     bool
	DoubleTapped bool

	HoldFor     int
	ReleasedFor int

	DoubleTapLength int

	checkFn func() bool
}

func NewButton(checkFn func() bool) *Button {
	return &Button{
		checkFn: checkFn,
	}
}

func (b *Button) Update() {
	b.IsDown = b.checkFn()
	b.Pressed = b.IsDown && !b.WasDown
	b.Released = !b.IsDown && b.WasDown
	b.DoubleTapped = b.Pressed && b.ReleasedFor > 0 && b.ReleasedFor < b.DoubleTapLength

	if b.IsDown {
		b.HoldFor++
		b.ReleasedFor = 0
	} else {
		b.ReleasedFor++
		b.HoldFor = 0
	}

	b.WasDown = b.IsDown
}

func (b *Button) Reset() {
	b.WasDown = false
	b.IsDown = false
	b.Pressed = false
	b.Released = false
	b.DoubleTapped = false
	b.HoldFor = 0
	b.ReleasedFor = 0
}

func NewKeyboardButton(key ebiten.Key) *Button {
	return NewButton(func() bool {
		return ebiten.IsKeyPressed(key)
	})
}

func NewMouseButton(button ebiten.MouseButton) *Button {
	return NewButton(func() bool {
		return ebiten.IsMouseButtonPressed(button)
	})
}

type Cursor struct {
	PosX      int
	Posy      int
	LeftDown  bool
	RightDown bool
}

var (
	Cursors         = make([]Cursor, 10)
	NumberOfCursors = 1

	MouseX     int
	MouseY     int
	MousePrevX int
	MousePrevY int

	MousePos       = V2()
	MouseWheelUp   bool
	MouseWheelDown bool

	MouseActive bool
	TouchActive bool

	MouseLeft  = NewMouseButton(ebiten.MouseButtonLeft)
	MouseRight = NewMouseButton(ebiten.MouseButtonRight)

	KeyEsc = NewKeyboardButton(ebiten.KeyEscape)

	KeyF1  = NewKeyboardButton(ebiten.KeyF1)
	KeyF2  = NewKeyboardButton(ebiten.KeyF2)
	KeyF3  = NewKeyboardButton(ebiten.KeyF3)
	KeyF4  = NewKeyboardButton(ebiten.KeyF4)
	KeyF12 = NewKeyboardButton(ebiten.KeyF12)

	KeyN1 = NewKeyboardButton(ebiten.Key1)
	KeyN2 = NewKeyboardButton(ebiten.Key2)
	KeyN3 = NewKeyboardButton(ebiten.Key3)
	KeyN4 = NewKeyboardButton(ebiten.Key4)
	KeyN5 = NewKeyboardButton(ebiten.Key5)
	KeyN6 = NewKeyboardButton(ebiten.Key6)
	KeyN7 = NewKeyboardButton(ebiten.Key7)
	KeyN8 = NewKeyboardButton(ebiten.Key8)
	KeyN9 = NewKeyboardButton(ebiten.Key9)

	KeyBracketLeft  = NewKeyboardButton(ebiten.KeyBracketLeft)
	KeyBracketRight = NewKeyboardButton(ebiten.KeyBracketRight)

	KeyB = NewKeyboardButton(ebiten.KeyB)
	KeyD = NewKeyboardButton(ebiten.KeyD)
	KeyL = NewKeyboardButton(ebiten.KeyL)
)

func UpdateInputs() {
	_, dy := ebiten.Wheel()
	MouseWheelUp = dy > 0
	MouseWheelDown = dy < 0

	x, y := ebiten.CursorPosition()
	if x < 0 {
		x = 0
	}

	if x >= WorldWidth {
		x = WorldWidth - 1
	}

	if y < 0 {
		y = 0
	}

	if y >= WorldHeight {
		y = WorldHeight - 1
	}

	MouseLeft.Update()
	MouseRight.Update()

	MousePrevX = MouseX
	MousePrevY = MouseY

	// set mouse to active if any movement or button / wheel action is detected
	if MouseX != MousePrevX || MouseY != MousePrevY || MouseWheelUp || MouseWheelDown || MouseLeft.Pressed || MouseRight.Pressed {
		MouseActive = true
	}

	MouseX = x
	MouseY = y
	MousePos.Set(float32(x), float32(y))

	if MouseActive {
		NumberOfCursors = 1
		Cursors[0].PosX = MouseX
		Cursors[0].Posy = MouseY
		Cursors[0].LeftDown = MouseLeft.IsDown
		Cursors[0].RightDown = MouseRight.IsDown
	} else {
		touchIds := []ebiten.TouchID{}
		touchIds = ebiten.AppendTouchIDs(touchIds)
		NumberOfCursors = min(len(touchIds), 10)
		for i := 0; i < NumberOfCursors; i++ {
			posX, posY := ebiten.TouchPosition(touchIds[i])
			Cursors[i].PosX = posX
			Cursors[i].Posy = posY
			Cursors[i].LeftDown = !inpututil.IsTouchJustReleased(touchIds[i])
			Cursors[i].RightDown = false
		}
	}

	KeyEsc.Update()
	KeyF1.Update()
	KeyF2.Update()
	KeyF3.Update()
	KeyF4.Update()
	KeyF12.Update()

	KeyN1.Update()
	KeyN2.Update()
	KeyN3.Update()
	KeyN4.Update()
	KeyN5.Update()
	KeyN6.Update()
	KeyN7.Update()
	KeyN8.Update()
	KeyN9.Update()

	KeyBracketLeft.Update()
	KeyBracketRight.Update()

	KeyB.Update()
	KeyD.Update()
	KeyL.Update()
}
