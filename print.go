package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

func Error(v any, exit ...bool) {
	c := color.New(color.FgRed)
	fmt.Println(c.Sprint("🚨 Error: ", v))
	osExit(exit)
}

func Log(v any) {
	c := color.New(color.FgBlue)
	fmt.Println(c.Sprint("🗒️ Log: ", v))
}

func Warning(v any) {
	c := color.New(color.FgYellow)
	fmt.Println(c.Sprint("⚠️ Warning: ", v))
}

func Info(v any) {
	c := color.New(color.FgMagenta)
	fmt.Println(c.Sprint(v))
}

func Done(t time.Duration) {
	c := color.New(color.FgGreen)
	fmt.Println(c.Sprint("✅ Done in: ", t))
}

func osExit(exit []bool) {
	var real bool

	if len(exit) > 0 {
		real = exit[0]
	}

	if real {
		os.Exit(1)
	}
}
