# Intel 8080 Space Invaders :space_invader:

## Table of Contents

- [Gameplay](#gameplay)
- [Installation](#installation)
  - [Linux](#linux)
  - [Windows](#windows)
  - [Other platforms](#other-platforms)
- [Input](#input)
- [Testing](#testing)
- [References](#references)

## Gameplay

<img src="https://github.com/gaoliveira21/space-invaders8080/blob/main/gameplay.gif" width="40%">

## Installation

### Linux

You can download it by clicking [here](https://github.com/gaoliveira21/space-invaders8080/releases/download/v1.0.3/linux_amd64_space_invaders), or by running the following commands:

```shell
wget https://github.com/gaoliveira21/space-invaders8080/releases/download/v1.0.3/linux_amd64_space_invaders
chmod +x ./linux_amd64_space_invaders
./linux_amd64_space_invaders # or run ./linux_amd64_space_invaders --sound-off to execute with audio disabled
```

If you have downloaded it manually, remember that you have to give execution permission by running `chmod +x ./linux_amd64_space_invaders`

### Windows

You have to download and unzip [win64_space_invaders.zip available here](https://github.com/gaoliveira21/space-invaders8080/releases/download/v1.0.3/win64_space_invaders.zip), after that just run **space_invaders.exe**

### Other platforms

First of all you have to install `sdl2_dev` and `sdl2_mixer`, you can access the following link to see how to install SDL2 on your OS.
[SDL2 - Installation](https://wiki.libsdl.org/SDL2/Installation)

After having SDL2 installed you have to compile it manually by executing the following commands (Golang is **required** to do that, so if you don't have it installed you can download and install it here: [Golang - Download and install](https://go.dev/doc/install)).

```shell
git clone https://github.com/gaoliveira21/space-invaders8080.git
cd space-invaders8080

go mod download
go build -o ./build/space-invaders ./cmd/invaders/main.go

./build/space-invaders
```

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

- [dustinbowers' intel8080emu ](https://github.com/dustinbowers/intel8080emu/tree/master)
