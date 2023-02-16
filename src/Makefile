TARGET=$(HOME)/.local

all: install

markdox:
	go build -mod=vendor -o markdox .

install: markdox
	cp markdox $(TARGET)/bin

clean:
	rm -rf markdox

.PHONY: all clean install