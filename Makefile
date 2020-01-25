.DEFAULT_GOAL := dev

clean:
	yarn clean
	rm -rf dist

build: render
	yarn build

%:
	go run ./gen $@
