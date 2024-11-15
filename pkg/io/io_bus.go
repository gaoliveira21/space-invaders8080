package io

type AudioPlayer interface {
	Play(soundType byte)
}

type IOBus struct {
	//Read 1
	//BIT	0	coin (0 when active)
	//1	P2 start button
	//2	P1 start button
	//3	?
	//4	P1 shoot button
	//5	P1 joystick left
	//6	P1 joystick right
	//7	?
	input1 byte

	//Read 2
	//BIT	0,1	dipswitch number of lives (0:3,1:4,2:5,3:6)
	//2	tilt 'button'
	//3	dipswitch bonus life at 1:1000,0:1500
	//4	P2 shoot button
	//5	P2 joystick left
	//6	P2 joystick right
	//7	dipswitch coin info 1:off,0:on
	input2 byte

	shiftH byte
	shiftL byte
	offset byte

	audioPlayer AudioPlayer
}

func NewIOBus(ap AudioPlayer) *IOBus {
	return &IOBus{
		input1:      0x00,
		input2:      0x00,
		shiftH:      0x00,
		shiftL:      0x00,
		audioPlayer: ap,
	}
}

func (io *IOBus) Read(port byte) byte {
	switch port {
	case 0x1:
		return io.input1
	case 0x2:
		return io.input2
	case 0x3:
		shift := (uint16(io.shiftH)<<8 | uint16(io.shiftL))
		return byte(shift >> (8 - uint16(io.offset)))
	default:
		return 0
	}
}

func (io *IOBus) Write(port byte, A byte) {
	switch port {
	case 0x2:
		io.offset = A & 0x7
	case 0x3:
		io.audioHandler(A, 0x1, UFORepeatsSound)
		io.audioHandler(A, 0x2, ShotSound)
		io.audioHandler(A, 0x4, ExplosionSound)
		io.audioHandler(A, 0x8, InvaderDieSound)
	case 0x4:
		io.shiftL = io.shiftH
		io.shiftH = A
	case 0x5:
		io.audioHandler(A, 0x1, FleetMovement1Sound)
		io.audioHandler(A, 0x2, FleetMovement2Sound)
		io.audioHandler(A, 0x4, FleetMovement3Sound)
		io.audioHandler(A, 0x8, FleetMovement4Sound)
		io.audioHandler(A, 0x10, UFOHitSound)
	}
}

func (io *IOBus) OnInput(port uint8, bit uint8, pressed bool) {
	if port == 1 {
		io.setInput(&io.input1, bit, pressed)
	}

	if port == 2 {
		io.setInput(&io.input2, bit, pressed)
	}
}

func (io *IOBus) setInput(inputPtr *byte, bit uint8, pressed bool) {
	if pressed {
		*inputPtr |= 1 << bit
	} else {
		*inputPtr &= ^(1 << bit)
	}
}

func (io *IOBus) audioHandler(soundType, bitMask, soundId byte) {
	if io.audioPlayer == nil {
		return
	}

	if soundType&bitMask != 0 {
		io.audioPlayer.Play(soundId)
	}
}
