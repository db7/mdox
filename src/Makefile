TARGET=$(HOME)/.local

all: install

mdox:
	go build -mod=vendor -o mdox .

install:mdox 
	cp mdox $(TARGET)/bin

clean:
	rm -rfmdox 

.PHONY: all clean install mdox
