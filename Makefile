EXE = gshell
SRC = .
LDFLAGS = -ldflags="-s -w"

darwin:
	GOOS=darwin packr build -o $(EXE).macho $(LDFLAGS) $(SRC)

windows:
	GOOS=windows packr build -o $(EXE).exe $(LDFLAGS) $(SRC)

linux:
	GOOS=linux packr build -o $(EXE).elf $(LDFLAGS) $(SRC)

all: darwin windows linux

clean:
	rm -f $(EXE).macho $(EXE).exe $(EXE).elf
