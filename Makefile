.DEFAULT_GOAL := dev

clean:
	yarn clean
	rm -rf dist

build: render
	yarn build

%:
	go run ./gen $@

%-docker:
	docker run --rm -it \
		-v ${PWD}:/tkbyexample \
		-p 8000:8000 \
		tkbyexample-dev

%-docker-build:
	docker build -t tkbyexample-dev -f Dockerfile.dev .
