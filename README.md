# Spark gen

Spark is a templating engine for go. It only provides components and `property`/`attribute` passing currently.

## Installation

First, you need to install the generator package:

```bash
go install https://github.com/go-spark/spark-gen
```

### Simple usage
```bash
spark-gen --dir="path/to/spark_templates_dir"
```

### Other flags:

- `--dir` - path to the directory with spark templates
- `--outDir` - path to the output directory. use `@` to alias `--dir` path. default is `@/dist`
- `--ext` - extension for spark templates. default is `.st`
- `--pkg` - package name for the generated code. default is `dist`

## Usage

After generation just import the generated methods.
`hello.st` will be converted to a Hello method under your outDir with the provided package name.

```go
func main() {
	content := dist.Hello().Render()
	fmt.Println(content)
}
```
