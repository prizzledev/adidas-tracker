package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func Info(name string, message string) {
	color := color.New(color.FgBlue).SprintFunc()

	fmt.Println(color(fmt.Sprintf("[%s][%s][%s] %s", time.Now().Format(time.RFC3339Nano), "Info", name, message)))
}

func Error(name string, message string) {
	color := color.New(color.FgRed).SprintFunc()

	fmt.Println(color(fmt.Sprintf("[%s][%s][%s] %s", time.Now().Format(time.RFC3339Nano), "Error", name, message)))
}

func Warning(name string, message string) {
	color := color.New(color.FgYellow).SprintFunc()

	fmt.Println(color(fmt.Sprintf("[%s][%s][%s] %s", time.Now().Format(time.RFC3339Nano), "Warning", name, message)))
}

func Success(name string, message string) {
	color := color.New(color.FgGreen).SprintFunc()

	fmt.Println(color(fmt.Sprintf("[%s][%s][%s] %s", time.Now().Format(time.RFC3339Nano), "Success", name, message)))
}
