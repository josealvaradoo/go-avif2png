package main

import (
	"avif2png/internal/cli"
	"fmt"
	"os"
)

func main() {
	config, err := cli.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Println("🚀 Starting AVIF to PNG conversion...")
	}

	if err := cli.Run(config); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Println("🎉 Conversion completed successfully!")
	} else {
		fmt.Println("✅ Done")
	}
}
