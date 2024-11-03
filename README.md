# Intel 8080 Space Invaders :space_invader:

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

# References

- [Emulator 101](http://www.emulator101.com/welcome.html)

- [Intel 8080 Opcode](http://www.emulator101.com/reference/8080-by-opcode.html)

- [Space Invaders - Hardware](http://computerarcheology.com/Arcade/SpaceInvaders/Hardware.html)

- [Emutalk - Space Invaders](https://www.emutalk.net/threads/space-invaders.38177/)
