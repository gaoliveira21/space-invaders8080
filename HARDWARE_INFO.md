# Space Invaders Hardware Info

Space Invaders, (C) Taito 1978, Midway 1979

CPU: Intel 8080 @ 2MHz (CPU similar to the (newer) Zilog Z80)

Interrupts: $cf (RST 8) at the start of vblank, $d7 (RST $10) at the end of vblank.

Video: 256(x)*224(y) @ 60Hz, vertical monitor. Colours are simulated with a
plastic transparent overlay and a background picture.
Video hardware is very simple: 7168 bytes 1bpp bitmap (32 bytes per scanline).

Sound: SN76477 and samples.

Memory map:

ROM

- $0000-$07ff:    invaders.h
- $0800-$0fff:    invaders.g
- $1000-$17ff:    invaders.f
- $1800-$1fff:    invaders.e

RAM

- $2000-$23ff: work RAM
- $2400-$3fff: video RAM
- $4000-:      RAM mirror


Inputs

Port 0
 bit 0 DIP4 (Seems to be self-test-request read at power up)
 bit 1 Always 1
 bit 2 Always 1
 bit 3 Always 1
 bit 4 Fire
 bit 5 Left
 bit 6 Right
 bit 7 ? tied to demux port 7 ?

Port 1
 bit 0 = CREDIT (1 if deposit)
 bit 1 = 2P start (1 if pressed)
 bit 2 = 1P start (1 if pressed)
 bit 3 = Always 1
 bit 4 = 1P shot (1 if pressed)
 bit 5 = 1P left (1 if pressed)
 bit 6 = 1P right (1 if pressed)
 bit 7 = Not connected

Port 2
 bit 0 = DIP3 00 = 3 ships  10 = 5 ships
 bit 1 = DIP5 01 = 4 ships  11 = 6 ships
 bit 2 = Tilt
 bit 3 = DIP6 0 = extra ship at 1500, 1 = extra ship at 1000
 bit 4 = P2 shot (1 if pressed)
 bit 5 = P2 left (1 if pressed)
 bit 6 = P2 right (1 if pressed)
 bit 7 = DIP7 Coin info displayed in demo screen 0=ON

Port 3
  bit 0-7 Shift register data
