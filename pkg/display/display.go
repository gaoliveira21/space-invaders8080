package display

import "github.com/veandco/go-sdl2/sdl"

var (
	width  = int32(224)
	height = int32(256)
	scale  = int32(2)
)

var window *sdl.Window

func Init() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	win, err := sdl.CreateWindow("Space Invaders", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width*scale, height*scale, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	window = win
}

func Draw(vram []byte) {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(nil, 0)
	if err != nil {
		panic(err)
	}

	for i, b := range vram {
		x := i / 32
		for bit := 0; bit < 8; bit++ {
			y := (i%32)*8 + bit

			pixel := b & (0x1 << bit)

			color := uint32(0x00000000)
			if pixel > 0 {
				color = 0xffffffff
			}

			rwidth := 1 * scale
			rheight := 1 * scale

			surface.FillRect(&sdl.Rect{
				X: int32(x) * rwidth,
				Y: (height * scale) - int32(y)*rheight,
				W: rwidth,
				H: rheight,
			}, color)
		}
	}

	err = window.UpdateSurface()
	if err != nil {
		panic(err)
	}
}

func Destroy() {
	sdl.Quit()
	window.Destroy()
}
