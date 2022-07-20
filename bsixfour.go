package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// set up command-line flag parsing
	appendNewline := true
	const appendNewlineHelpMessage = "By default this tool appends a final newline character. Set to false to disable this behavior."

	flagSet := flag.NewFlagSet("encode", flag.ExitOnError)
	flagSet.BoolVar(&appendNewline, "append-newline", true, appendNewlineHelpMessage)

	// default Usage message didn't seem to work with my subcommands,
	// so I'm just setting it manually
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s encode|decode [-append-newline=true|false]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  -append-newline")
		fmt.Fprintf(os.Stderr, "\t%s\n", appendNewlineHelpMessage)
	}

	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(1)
	}

	flagSet.Parse(os.Args[2:])

	switch os.Args[1] {
	case "encode":
		base64Channel := make(chan string)
		go encode(os.Stdin, base64Channel)
		for {
			val, isOpen := <-base64Channel
			if !isOpen {
				break
			}
			fmt.Print(val)
		}
	case "decode":
		originalDataChannel := make(chan []byte)
		go decode(os.Stdin, originalDataChannel)
		for {
			val, isOpen := <-originalDataChannel
			if !isOpen {
				break
			}
			fmt.Print(string(val))
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	if appendNewline {
		fmt.Println()
	}
}
