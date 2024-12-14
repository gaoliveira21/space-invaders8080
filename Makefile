build:
	CGO_ENABLED=1 CC=gcc GOOS=$(OS) GOARCH=$(ARCH) go build -tags static -ldflags "-s -w" -o $(O_DIR)/$(OS)_$(ARCH)_space_invaders ./cmd/invaders/main.go
