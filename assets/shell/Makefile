EXE = shell
SRC = .
LDFLAGS = -ldflags="-s -w"

windows:
	GOOS=windows go build -o $(EXE).exe $(LDFLAGS) $(SRC)

macos:
	GOOS=darwin go build -o $(EXE).macho $(LDFLAGS) $(SRC)

linux:
	GOOS=linux go build -o $(EXE).elf $(LDFLAGS) $(SRC)

all: windows macos linux
	echo "done."

clean:
	rm -f *.exe *.elf *.macho