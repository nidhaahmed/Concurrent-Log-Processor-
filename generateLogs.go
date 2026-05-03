package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Create("logs.txt")
	defer file.Close()

	for i := 0; i < 50000; i++ {
		switch i % 3 {
		case 0:
			fmt.Fprintf(file, "2026-04-27 INFO User logged in %d\n", i)
		case 1:
			fmt.Fprintf(file, "2026-04-27 WARNING Disk usage high %d\n", i)
		case 2:
			fmt.Fprintf(file, "[2026-04-27 10:20:%02d] ERROR: Database failure %d\n", i%60, i)
		}
	}
}