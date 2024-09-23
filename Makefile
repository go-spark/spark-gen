build:
	@go build -o bin/gospark

run: build
	@./bin/gospark --dir="./testdata" --outDir="@/dist"

debug: build
	@./bin/gospark --dir="./testdata" --outDir="@/dist" --print=true

generate: build
	@./bin/gospark --dir="./testdata" --outDir="@/dist"
