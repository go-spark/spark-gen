package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	FileExt = "st"          // spark template file extension
	Version = "0.1.1-alpha" // version of spark-gen
)

var (
	start   = time.Now()
	rootDir = flag.String("dir", "", "Directory to parse")
	out     = flag.String("outDir", "@/dist/out.go", "Output file path (start with @ to use the root directory)")
	ext     = flag.String("ext", FileExt, "File extension. (default: st)")
	pkg     = flag.String("pkg", "dist", "Package name for dist. (default: dist)")
)

func main() {
	flag.Parse()
	generate()
	Done(time.Since(start))
}

func generate() {
	outPath := *out
	outPath = strings.Replace(outPath, "@", *rootDir, 1)
	outPath = filepath.Clean(outPath)
	out = &outPath

	if *rootDir == "" {
		Error("flag --dir cannot be empty", true)
	}

	if *ext != FileExt {
		Warning("the --ext flag is not the default, it may cause possible problems")
	}

	fmt.Println("Spark-Gen v" + Version)

	Log("Tokenizing spark templates...")

	var data = make(map[string]*Component)

	files, err := walkDir(*rootDir, *ext) // get all st (spark template) files.
	if err != nil {
		Error(fmt.Sprint("failed to read files from dir:", err), true)
	}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			Error(fmt.Sprint("failed to open file:", err), true)
		}

		parser := NewParser(file, filePath)

		els, err := parser.Parse()
		if err != nil {
			Error(fmt.Sprint("failed to parse file:", err), true)
		}

		filePath, err = getNormalizedPath(filePath, *rootDir)
		if err != nil {
			Error(fmt.Sprint("failed to get normalized path:", err), true)
		}

		data[componentCase(filePath)] = els
	}

	os.MkdirAll(filepath.Dir(*out), os.ModePerm)
	g := NewGenerator(*out, *pkg, *ext, data)

	err = g.Make()
	if err != nil {
		Error(fmt.Sprint("failed to generate go files:", err), true)
	}

	err = g.Save()
	if err != nil {
		Error(fmt.Sprint("failed to save go files:", err), true)
	}

	fmt.Printf("The created file location is: %s\n", *out)
}
