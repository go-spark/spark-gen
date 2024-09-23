package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func printElement(element *Element, indent int) {
	fmt.Printf("%s<Element: %s, Content: %s, SelfClosing: %t, Attributes: %v>\n", strings.Repeat("  ", indent), element.Name, element.Content, element.SelfClosing, element.Attributes)
	for _, child := range element.Children {
		printElement(child, indent+1)
	}
}

func walkDir(root, ext string) ([]string, error) {
	var a []string
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == "."+ext {
			a = append(a, s)
		}
		return nil
	})
	return a, err
}

func pascalCase(input string) string {
	var result strings.Builder

	input = strings.ReplaceAll(input, "_", " ")
	input = strings.ReplaceAll(input, "-", " ")

	capitalizeNext := true
	for i, ch := range input {
		if i > 0 && unicode.IsUpper(ch) && unicode.IsLower(rune(input[i-1])) {
			capitalizeNext = true
		}

		if capitalizeNext {
			result.WriteRune(unicode.ToUpper(ch))
			capitalizeNext = false
		} else {
			result.WriteRune(unicode.ToLower(ch))
		}

		if unicode.IsSpace(ch) {
			capitalizeNext = true
		}
	}

	return strings.ReplaceAll(result.String(), " ", "")
}

func camelCase(s string) string {
	s = pascalCase(s)
	return strings.ToLower(s[:1]) + s[1:]
}

func componentCase(s string) string {
	s = strings.ReplaceAll(s, ".", "/")
	var fileName []string

	for _, p := range strings.Split(s, "/") {
		fileName = append(fileName, pascalCase(p))
	}

	return strings.Join(fileName, "_")
}

func camelToKebab(input string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	kebab := re.ReplaceAllString(input, "${1}-${2}")
	return strings.ToLower(kebab)
}

func getNormalizedPath(path, root string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(absPath, absRoot) {
		trimmedPath := strings.TrimPrefix(absPath, absRoot)
		path := strings.TrimPrefix(trimmedPath, string(filepath.Separator))
		return strings.TrimSuffix(path, filepath.Ext(path)), nil
	}

	return strings.TrimSuffix(absPath, filepath.Ext(absPath)), nil
}

func keys[K any, T map[string]K](m T) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
