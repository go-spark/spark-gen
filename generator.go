package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"
)

type Generator struct {
	data     map[string]*Component
	outPath  string
	out      string
	outParts string
	ext      string
}

func NewGenerator(out, pkg, ext string, data map[string]*Component) *Generator {
	g := &Generator{data: data, outPath: out, ext: ext + "-"}
	g.out += "// This file is auto generated. Please do not edit.\n// To generate a file, run: spark-gen -dir <your templates dir> -out <your out dir> -pkg <go package name>"
	g.out += fmt.Sprintf("\npackage %s\n\nimport \"github.com/go-spark/spark\"\n\n", pkg)
	return g
}

func (g *Generator) Make() error {
	for name, component := range g.data {
		err := g.createComponent(name, component)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) createComponent(name string, component *Component) error {
	if len(component.Elements) == 0 {
		Warning(fmt.Sprintf("component %s has no elements, skipping...", name))
		return nil
	}

	start := time.Now()
	defer func(start time.Time) {
		Info(fmt.Sprintf("⌚ Generated component: %s in %s", name, time.Since(start)))
	}(start)

	g.out += fmt.Sprintf("// %s %s %s\nfunc %s() spark.Component {\n", name, "is a spark.Component ->", component.File, name)
	g.out += fmt.Sprintf("\t_component := spark.NewV1Component([]string{%s})\n", strings.Join(component.Props, ","))
	for i, el := range component.Elements {
		elName := g.elementName(name, el.Name, i)
		g.out += fmt.Sprintf("\t_component.Push(%s(_component))\n", elName)
		g.createElement(el, elName, component.Props)
	}
	g.out += "\treturn _component\n}\n\n"

	return nil
}

func (g *Generator) createElement(el *Element, name string, props []string) {
	g.outParts += fmt.Sprintf("func %s(base spark.Component) spark.Element {\n", name)
	if strings.HasPrefix(el.Name, g.ext) {
		g.createReferenceElement(el, name, props, componentCase(strings.TrimPrefix(el.Name, g.ext)))
	} else {
		g.createSimpleElement(el, name, props)
	}
}

func (g *Generator) createReferenceElement(el *Element, _ string, props []string, component string) {
	g.outParts += fmt.Sprintf("\t_ref := spark.Ref(%s(), base)\n", component)

	ref, ok := g.data[component]
	if !ok {
		g.outParts += "\treturn _ref\n}\n\n"
		return
	}

	for name := range el.Attributes {
		if slices.Contains(ref.Props, `"`+name+`"`) {
			g.outParts += fmt.Sprintf("\t_ref.Component().Bind(\"%s\", func() string {\n\t\treturn %s\n\t})\n", name, g.getAttribute(name, el.Attributes[name], props, "_ref.GetProp"))
		} else {
			g.outParts += fmt.Sprintf("\t_ref.FirstChild().SetAttribute(\"%s\", func() string {\n\t\treturn %s\n\t})\n", camelToKebab(name), g.getAttribute(name, el.Attributes[name], props, "_ref.GetProp"))
		}
	}

	g.outParts += "\treturn _ref\n}\n\n"
}

func (g *Generator) createSimpleElement(el *Element, name string, props []string) {
	g.outParts += fmt.Sprintf("\t_el := spark.NewV1Element(\"%s\", %t, base)\n", el.Name, el.SelfClosing)
	for key, val := range el.Attributes {
		g.outParts += fmt.Sprintf("\t_el.SetAttribute(\"%s\", func() string {\n\t\treturn %s\n\t})\n", key, g.getAttribute(key, val, props, "_el.GetProp"))
	}

	for i, child := range el.Children {
		elName := g.elementName(name, child.Name, i)
		g.outParts += fmt.Sprintf("\t_el.AddChild(%s(base))\n", elName)
	}

	g.outParts += fmt.Sprintf("\t_el.Content(func() string {\n\t\treturn %s\n\t})\n", g.propContent(el.Content, props))

	g.outParts += "\treturn _el\n}\n\n"

	for i, child := range el.Children {
		elName := g.elementName(name, child.Name, i)
		g.createElement(child, elName, props)
	}
}

func (g *Generator) Save() error {
	return os.WriteFile(g.outPath, []byte(g.out+"// Helpers \n\n"+g.outParts), os.ModePerm)
}

func (g *Generator) getAttribute(name string, attr Attribute, props []string, propGetter string) string {
	if !attr.Go {
		return fmt.Sprintf(`"%s"`, attr.Value)
	}

	if attr.Value == "" && attr.Go {
		return fmt.Sprintf("%s(\"%s\")", propGetter, name)
	}

	return g.propString(attr.Value, props)
}

func (g *Generator) propString(val string, props []string) string {
	for _, prop := range props {
		prop = strings.Trim(prop, `"`)
		pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(prop))
		if matched, _ := regexp.MatchString(pattern, val); matched {
			val = strings.ReplaceAll(val, prop, fmt.Sprintf(`_el.GetProp("%s")`, prop))
		}
	}

	return strings.TrimSpace(val)
}

func (g *Generator) elementName(name string, element string, i int) string {
	return fmt.Sprintf("%s_%d_%s", camelCase(name), i, pascalCase(componentCase(element)))
}

func (g *Generator) propContent(val string, props []string) string {
	val = strings.TrimSpace(val)

	pattern := regexp.MustCompile(`\{\{\s*(.*?)\s*}}`)

	val = pattern.ReplaceAllStringFunc(val, func(match string) string {
		content := pattern.FindStringSubmatch(match)[1]
		return "` + " + g.propString(strings.TrimSpace(content), props) + " + `"
	})

	return "`" + val + "`"
}
