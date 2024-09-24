build:
	@go build -o bin/gospark

run: build
	@./bin/gospark --dir="./testdata"

debug: build
	@./bin/gospark --dir="./testdata" --print=true

generate: build
	@./bin/gospark --dir="./testdata"
