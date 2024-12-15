build-linux:
	CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 go build -tags static -ldflags "-s -w" -o $(O_DIR)/linux_amd64_space_invaders ./cmd/invaders/main.go

build-windows:
	env CGO_ENABLED="1" CC="/usr/bin/x86_64-w64-mingw32-gcc" GOOS="windows" CGO_LDFLAGS="-lmingw32 -lSDL2" CGO_CFLAGS="-D_REENTRANT" go build -x -ldflags "-H windowsgui" -o $(O_DIR)/win64_space_invaders/space_invaders.exe cmd/invaders/main.go
