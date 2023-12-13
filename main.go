package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"

	"github.com/chuckhoffmann/wired-logic/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	simulation             *simulation.Simulation
	simulationImage        *image.Paletted
	backgroundImage        *ebiten.Image
	wireImages             []*ebiten.Image
	width                  int
	height                 int
	scale                  int
	cursorx                int
	cursory                int
	mouse_cursorx          int
	mouse_cursory          int
	leftMouseButtonPressed bool
	simulationPaused       bool
	key_debounce           int
}

func (g *Game) Update() error {
	// Update the simulation
	var newSimulation *simulation.Simulation
	if g.simulationPaused {
		newSimulation = g.simulation
	} else {
		newSimulation = g.simulation.Step()
	}
	// Update the background image
	newSimulation.Draw(g.simulationImage)
	g.backgroundImage = ebiten.NewImageFromImage(g.simulationImage)
	wires := g.simulation.Circuit().Wires()
	for i, wire := range wires {
		oldCharge := g.simulation.State(wire).Charge()
		charge := newSimulation.State(wire).Charge()
		if oldCharge == charge {
			continue
		}
		position := wire.Bounds().Min
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(position.X), float64(position.Y))
		r, gr, b, a := g.simulationImage.Palette[charge+1].RGBA()
		op.ColorM.Scale(float64(r)/0xFFFF, float64(gr)/0xFFFF, float64(b)/0xFFFF, float64(a)/0xFFFF)
		g.backgroundImage.DrawImage(g.wireImages[i], op)
	}
	g.simulation = newSimulation
	//
	// Handle keyboard and mouse inputs
	g.handleKeyboard()
	g.handleMouse()

	return nil
}

func (g *Game) handleKeyboard() error {
	// Handle various keyboard inputs
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("closing game")
	}
	if ebiten.IsKeyPressed((ebiten.KeySpace)) {

		flipPixel(g.cursorx, g.cursory, g.simulationImage)
		g.reloadSimulation()
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		// Debounce the keypress
		if g.key_debounce > 0{
			g.key_debounce --
			return nil
		}
		g.key_debounce = 10
		// Pause the simulation
		g.simulationPaused = !g.simulationPaused
		drawPoweredDown(g.simulationImage)
		g.reloadSimulation()
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		saveImage(g.simulationImage, "test.gif")
	}

	switch {
	// Handle the arrow keys since up and down can be pressed at the same time.
	// if both are pressed, the evaluation order is W, ArrowUp, S, ArrowDown
	case ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp):
		if g.cursory > 0 {
			g.cursory--
		}
	case ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown):
		if g.cursory < g.height-1 {
			g.cursory++
		}
	default:
		// Do nothing

	}

	switch {
	// Handle the arrow keys since left and right can be pressed at the same time.
	// if both are pressed, the evaluation order is A, ArrowLeft, D, ArrowRight
	case ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft):
		if g.cursorx > 0 {
			g.cursorx--
		}
	case ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight):
		if g.cursorx < g.width-1 {
			g.cursorx++
		}
	}

	return nil
}

func (g *Game) handleMouse() {
	// Handle mouse inputs
	// Get the mouse position. If the mouse has moved, update the cursor position.
	x, y := ebiten.CursorPosition()
	if x != g.mouse_cursorx || y != g.mouse_cursory {
		// The mouse has moved
		g.mouse_cursorx = x
		g.mouse_cursory = y
		g.cursorx = g.mouse_cursorx
		g.cursory = g.mouse_cursory
	} else {
		// The mouse has not moved
		g.mouse_cursorx = x
		g.mouse_cursory = y
	}
	// Check if the left mouse button is pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.leftMouseButtonPressed {
			// The left mouse button has just been pressed
			flipPixel(g.cursorx, g.cursory, g.simulationImage)
			g.reloadSimulation()
			g.leftMouseButtonPressed = true
		} else {
			// The left mouse button is being held down
			// If the cursor has moved (i.e. through arrow keys) flip the pixel
			if g.cursorx != g.mouse_cursorx || g.cursory != g.mouse_cursory {
				flipPixel(g.cursorx, g.cursory, g.simulationImage)
				g.reloadSimulation()
			}
		}

	} else {
		if g.leftMouseButtonPressed {
			// The left mouse button has just been released
			g.leftMouseButtonPressed = false

		}
	}
}

func (g *Game) reloadSimulation() {
	g.simulation = simulation.New(g.simulationImage)
	g.simulation.Draw(g.simulationImage)

	for _, img := range g.wireImages {
		img.Dispose()
	}

	wires := g.simulation.Circuit().Wires()
	g.wireImages = make([]*ebiten.Image, len(wires))
	for i, wire := range wires {
		img := drawMask(wire)
		g.wireImages[i] = ebiten.NewImageFromImage(img)
	}

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.backgroundImage, nil)
	vector.DrawFilledRect(screen, float32(g.cursorx), float32(g.cursory), float32(1), float32(1), color.RGBA{128, 128, 128, 128}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func main() {
	// Get the command line arguments
	var speed, scale, width, height int
	flag.IntVar(&speed, "speed", 15, "simulation steps per second")
	flag.IntVar(&scale, "scale", 12, "pixel scale factor")
	flag.IntVar(&width, "width", 64, "width of the simulation")
	flag.IntVar(&height, "height", 64, "height of the simulation")
	flag.Parse()
	flag.Args()
	// Create a new game struct and set the properties
	game := &Game{
		width:   width,
		height:  height,
		scale:   scale,
		cursorx: width / 2,
		cursory: height / 2,
	}
	ebiten.SetTPS(speed)
	// Create a new blank simulation image
	p := color.Palette{
		color.Black,
		color.RGBA{0x88, 0, 0, 0xFF},
		color.RGBA{0xFF, 0, 0, 0xFF},
		color.RGBA{0xFF, 0x22, 0, 0xFF},
		color.RGBA{0xFF, 0x44, 0, 0xFF},
		color.RGBA{0xFF, 0x66, 0, 0xFF},
		color.RGBA{0xFF, 0x88, 0, 0xFF},
		color.RGBA{0xFF, 0xAA, 0, 0xFF},
	}
	game.simulationImage = image.NewPaletted(image.Rect(0, 0, game.width, game.height), p)
	game.simulation = simulation.New(game.simulationImage)
	game.backgroundImage = ebiten.NewImageFromImage(game.simulationImage)
	game.reloadSimulation()

	ebiten.SetWindowSize(game.width*game.scale, game.height*game.scale)
	ebiten.SetWindowTitle("Wired Logic Sandbox")
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}
	fmt.Println("All done!")
}

func init() {
	fmt.Println("Starting the program...")
}

func saveImage(img *image.Paletted, filename string) {
	// Create a new file
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	// Encode the image as a gif
	if err := gif.Encode(f, img, nil); err != nil {
		fmt.Println(err)
	}
}

func flipPixel(x, y int, img *image.Paletted) {
	// Flip the pixel at the given coordinates
	// if index is 0 set it to 1
	// else set it to 0
	col := img.ColorIndexAt(x, y)
	// Flip the pixel
	if col == 0 {
		img.SetColorIndex(x, y, 1)
	} else {
		img.SetColorIndex(x, y, 0)
	}
}

func drawMask(wire *simulation.Wire) image.Image {
	// Draw a mask for the wire.
	bounds := image.Rect(0, 0, wire.Bounds().Dx(), wire.Bounds().Dy())
	bounds = bounds.Union(image.Rect(0, 0, 4, 4))
	position := wire.Bounds().Min
	img := image.NewRGBA(bounds)
	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	for _, pixel := range wire.Pixels() {
		img.SetRGBA(pixel.X-position.X, pixel.Y-position.Y, white)
	}
	return img
}

func drawPoweredDown(img *image.Paletted) {
	// Draw a powered down simulation
	// Loop through the pixels and set
	// any pixel that has an index > 1
	// and less than or equal to 7 to 1.
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			col := img.ColorIndexAt(x, y)
			if col > 1 && col <= 7 {
				img.SetColorIndex(x, y, 1)
			}
		}
	}

}
