// hod-go is a package to dump data in hex or octal format. It is a Go implementation of
// my C program hod (https://github.com/muquit/hod)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var version = "1.0.1"
var offsetInDecimal = false
var dumpInOctal = false

func usage() {
	fmt.Fprintf(os.Stderr, "hod-go v%s\n", version)
	fmt.Fprintf(os.Stderr, "A Go implemention of https://github.com/muquit/hod\n\n")
	fmt.Fprintf(os.Stderr, "Usage: hod-go <file>\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	nargs := len(os.Args)

	flag.Bool("version", false, "Show version infomation")
	flag.BoolVar(&offsetInDecimal, "d", false, "show offsets in decimal. Default is hex")
	flag.BoolVar(&dumpInOctal, "o", false, "dump in octal. Default is hex")

	flag.Parse()

	// remaining non-flag args
	x := len(flag.Args())
	if x != 1 {
		if !isReadingFromStdin() {
			usage()
		}
	}

	var reader io.Reader
	if isReadingFromStdin() {
		reader = bufio.NewReader(os.Stdin)
	} else {
		fileName := os.Args[nargs-1]
		r, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Could not open file %s:%s\n", fileName, err)
		}
		reader = r
	}

	var base int
	if dumpInOctal {
		base = 8
	} else {
		base = 16
	}
	dumpData(reader, base)
}

// This function is almost identical to the C funtion dump_file() in hod.c
func dumpData(reader io.Reader, base int) {

	fmt.Printf("%11d", 0)
	for i := 1; i < base; i++ {
		if base == 16 {
			fmt.Printf("%3x", i)
		} else {
			fmt.Printf("%4o", i)
		}
	}
	if base == 16 {
		fmt.Printf("   ")
	} else {
		fmt.Printf("    ")
	}

	for i := 0; i < base; i++ {
		fmt.Printf("%x", i)
	}
	fmt.Printf("\n")

	lcount := 0
	buf := make([]byte, 0, base)
	for {
		rcount, err := reader.Read(buf[:base])
		if err == io.EOF {
			break
		}
		if rcount <= 0 {
			break
		}
		buf = buf[:rcount]
		if err != nil {
			log.Fatalf("Error reading bytes: %s\n", err)
		}
		// offset
		if offsetInDecimal {
			fmt.Printf("%8d: ", lcount*base)
		} else {
			if base == 16 {
				fmt.Printf("%8x: ", lcount*base)
			} else {
				fmt.Printf("%8o: ", lcount*base)
			}
		}
		for count := 0; count < base; count++ {
			if count < rcount {
				if base == 16 {
					fmt.Printf("%02x ", buf[count])
				} else {
					fmt.Printf("%03o ", buf[count])
				}
			} else {
				if base == 16 {
					fmt.Printf("   ")
				} else {
					fmt.Printf("    ")
				}
			}
		}
		fmt.Printf(" ")
		for i := 0; i < base; i++ {
			if i < rcount {
				if buf[i] >= 32 && buf[i] < 127 {
					fmt.Printf("%c", buf[i])
				} else {
					fmt.Printf(".")
				}
			}
		}
		fmt.Printf("\n")
		lcount++
	}
}

// isReadingFromStdin returns true if data is coming from stdin
func isReadingFromStdin() bool {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	return false
}
