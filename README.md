# Intel 8080 Space Invaders :space_invader:

## Input

| Key                | Description     |
|--------------------|-----------------|
| C                  | Inser coin      |
| 1                  | 1P Start        |
| 2                  | 2P Start        |
| W                  | 1P Shot         |
| A/D                | 1P Left/Right   |
| [Arrow Up]         | 2P Shot         |
| [Arrow Left/Right] | 2P Left/Right   |
| T                  | Tilt(Game over) |

## Testing

```shell
go run ./cmd/cpudiag/main.go
Running a test ROM - roms/tests/TST8080.COM
1536 bytes loaded
MICROCOSM ASSOCIATES 8080/8085 CPU DIAGNOSTIC
 VERSION 1.0  (C) 1980

 CPU IS OPERATIONAL
 ```

## References

- [Emulator 101](http://www.emulator101.com/welcome.html)

- [Intel 8080 Opcode](http://www.emulator101.com/reference/8080-by-opcode.html)

- [Space Invaders - Hardware](http://computerarcheology.com/Arcade/SpaceInvaders/Hardware.html)

- [Emutalk - Space Invaders](https://www.emutalk.net/threads/space-invaders.38177/)
