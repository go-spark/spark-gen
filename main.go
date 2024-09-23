package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	FileExt = "st" // spark template file extension
)

var (
	start      = time.Now()
	rootDir    = flag.String("dir", "", "Directory to parse")
	printDebug = flag.Bool("print", false, "Print the parsed data")
	outDir     = flag.String("outDir", "@/dist", "Output directory (start with @ to use the root directory)")
	ext        = flag.String("ext", FileExt, "File extension. (default: st)")
	pkg        = flag.String("pkg", "dist", "Package name for dist. (default: dist)")
)

func main() {
	flag.Parse()
	generate()
	fmt.Println("done in", time.Since(start))
}

func generate() {
	oDir := *outDir
	oDir = strings.Replace(oDir, "@", *rootDir, 1)
	outDir = &oDir

	if *rootDir == "" {
		log.Fatal("dir name cannot be empty.")
	}

	var data = make(map[string]*Component)

	files, err := walkDir(*rootDir, *ext) // get all st (spark template) files.
	if err != nil {
		log.Fatal("failed to read files from dir:", err)
	}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal("failed to open file:", err)
		}

		parser := NewParser(file, filePath)

		els, err := parser.Parse()
		if err != nil {
			log.Fatal("failed to parse file:", err)
		}

		filePath, err = getNormalizedPath(filePath, *rootDir)
		if err != nil {
			log.Fatal("failed to get normalized path:", err)
		}

		data[componentCase(filePath)] = els
	}

	if *printDebug {
		for filePath, els := range data {
			fmt.Println("file:", filePath)
			for _, el := range els.Elements {
				printElement(el, 0)
			}
		}
	}

	b, _ := json.MarshalIndent(data, "", "  ")
	os.MkdirAll(*outDir, os.ModePerm)

	os.WriteFile(filepath.Join(*outDir, "spark_gen.json"), b, os.ModePerm)

	g := NewGenerator(filepath.Join(*outDir, "out.go"), *pkg, *ext, data)

	err = g.Make()
	if err != nil {
		log.Fatal("failed to generate go files:", err)
	}

	err = g.Save()
	if err != nil {
		log.Fatal("failed to save go files:", err)
	}
}
