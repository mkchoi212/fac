package main

import "fmt"

func Red(s string) string {
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", 1, 1, s)
}

func Green(s string) string {
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", 2, 1, s)
}
