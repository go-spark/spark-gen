# Spark gen

Spark is a templating engine for go. It only provides components and `property`/`attribute` passing currently.

## Installation

First, you need to install the generator package:

```bash
go install github.com/go-spark/spark-gen@latest
```

### Simple usage
```bash
spark-gen --dir="path/to/spark_templates_dir"
```

### Other flags:

- `--dir` - path to the directory with spark templates
- `--out` - path to the output file. use `@` to alias `--dir` path. default is `@/dist/out.go`
- `--ext` - extension for spark templates. default is `st`
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

### Component code

A component, for example `hello.st` looks like this:

```html
!! label, id
<div :id>
  <p>{{ label }}</p>
  <st-input :placeholder="label" />
</div>
```

And `input.st` is imported with `st-input`, so let's create that too:

```html
<input name="username" />
```

The first line of the component can contain two `!` sign to define props.
The `label` and `id` are passed to the component. The `:id` attribute automatically referred to the prop id in `hello.st`.

### Import component from other directory

You can import code from other directory like this:

Component `/utils/button.st`
```html
!! title
<button>
  {{ title }}
</button>
```

To import this component from any other Spark Template (`st`) you need to provide the full path to the component.
You have to change the `/` to `.` and also add the `st-` prefix to help the compiler identify that you want to import a component.

```html
<div>
  <st-utils.button title="Hello" />
</div>
```
## Support us 🩵

If you’d like to support us financially, you can do so by donating through our [Ko-fi](https://ko-fi.com/bndrmrtn) page.
