/*
nes-header-decoder Decodes NES files.

Usage:

	nes-header-decoder [flags] [file]

The flags are:

	-h
		Show this help message.
	-v
		Show version.
	-d
		Show debug messages.

The file is:

	An NES file to decode.

Examples:

	nes-header-decoder ./zelda.nes
	nes-header-decoder -v
	nes-header-decoder -d -v
	nes-header-decoder -d -v -h
	nes-header-decoder -d -v -h ./zelda.nes
*/
package main

// Imports
import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/fatih/color"
)

// Constants
const (
	KB = 1024 // 1 KB = 1024 bytes
)

// Global Variables
var print_info = color.New(color.FgCyan)
var print_error = color.New(color.FgRed)
var print_debug = color.New(color.FgYellow)
var print_good = color.New(color.FgGreen)
var isset_debug = false
var isset_filename = ""
var version = "0.0.1"

// NesHeader is the header of an NES file.
type NesHeader struct {
	Magic      [4]byte
	PRGROMSize uint8
	CHRROMSize uint8
	Flags6     uint8
	Flags7     uint8
	PRGRAMSize uint8
	Flags9     uint8
	Flags10    uint8
	Zero       [5]byte
}

// check_error checks for errors and exits if there is one.
// It takes an error and a message to print if there is an error.
func check_error(err error, msg string) {
	if err != nil {
		print_error.Printf(msg)
		print_error.Printf("%v\n", err)
		log.Fatalln(err)
	}
}

// print_help_short prints the short help message.
func print_help_short() {
	fmt.Printf("Usage:\tnes-header-decoder [flags] [file]\n\n")
}

// print_help prints the help message.
func print_help() {
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tnes-header-decoder [flags] [file]\n\n")
	fmt.Printf("The flags are:\n\n")
	fmt.Printf("\t-h\tShow this help message.\n")
	fmt.Printf("\t-v\tShow version.\n")
	fmt.Printf("\t-d\tShow debug messages.\n\n")
	fmt.Printf("The file is:\n\n")
	fmt.Printf("\tAn NES file to decode.\n\n")
	fmt.Printf("Examples:\n\n")
	fmt.Printf("\tnes-header-decoder ./zelda.nes\n")
	fmt.Printf("\tnes-header-decoder -v\n")
	fmt.Printf("\tnes-header-decoder -d -v\n")
	fmt.Printf("\tnes-header-decoder -d -v -h\n")
	fmt.Printf("\tnes-header-decoder -d -v -h ./zelda.nes\n\n")
}

// pretty prints a struct of nesheader
func pretty(nesHeader NesHeader) {
	pretty, err := json.MarshalIndent(nesHeader, "*", "    ")
	if err != nil {
		fmt.Println("Failed to generate json", err)
	}
	fmt.Printf("*%s\n", pretty)
}

func readNumBytes(f *os.File, numBytes int) []byte {
	bytes := make([]byte, numBytes)
	_, err := f.Read(bytes)
	check_error(err, "** ERROR: Reading File.\n\n")
	return bytes
}

// init is the first function to run.
func init() {
	// Check for args
	if len(os.Args) < 2 {
		print_error.Printf("** ERROR: No Args.\n")
		print_help_short()
		os.Exit(1)
	}

	// Print Header
	print_good.Printf("=== NES Header Decoder ===\n")

	// Check for flags
	show_version := false
	show_help := false
	for n, args := range os.Args {
		if n == 0 {
			continue
		}
		switch args {
		case "-h":
			if len(os.Args) == 2 {
				print_help()
				os.Exit(0)
			}
			show_help = true
		case "-v":
			show_version = true
		case "-d":
			isset_debug = true
		default:
			isset_filename = args
		}
	}

	// Print Args if Debug is on
	if isset_debug {
		for n, args := range os.Args {
			print_debug.Printf("#%d - Args: %s\n", n, args)
		}
	}

	// Check for filename
	if isset_filename == "" {
		print_error.Printf("** ERROR: No Filename.\n")
		print_help_short()
		os.Exit(2)
	}

	// Show Version if Enabled
	if show_version {
		print_info.Printf("Info: Version = %v\n", version)
	}

	// Show Help if Enabled
	if show_help {
		print_help()
	}

	// Show Debug
	if isset_debug {
		print_debug.Printf("Debug: Enabled.\n")
	}
}

// main is the main function.
func main() {
	// Announce
	print_good.Printf("==========================\n")
	print_info.Printf("\nInfo: Opening File.\n")
	if isset_debug {
		print_debug.Printf("Debug: Filename = %s\n", isset_filename)
	}

	// Check File Exists
	if _, err := os.Stat(isset_filename); os.IsNotExist(err) {
		print_error.Printf("** ERROR: File Does Not Exist.\n")
		os.Exit(3)
	}
	print_good.Printf("Info: File Exists.\n")

	// Open and File
	f, err := os.Open(isset_filename)
	check_error(err, "** ERROR: Opening File.\n\n")
	defer f.Close()

	// Get File Size
	fi, err := f.Stat()
	check_error(err, "** ERROR: Getting File Size.\n\n")
	if isset_debug {
		print_debug.Printf("Debug: File Size = %d KB\n", fi.Size()/KB)
	}

	header := NesHeader{}
	data := readNumBytes(f, int(unsafe.Sizeof(header)))
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.LittleEndian, &header)
	check_error(err, "** ERROR: Decoding Header.\n\n")

	if isset_debug {
		pretty(header)
	}

	print_good.Printf("\nInfo: Decoded Header.\n")
	fmt.Printf("Magic:    %c%c%c x%02x\n", header.Magic[0], header.Magic[1], header.Magic[2], header.Magic[3])
	fmt.Printf("PRG ROM:  %d KB\n", header.PRGROMSize*16)
	fmt.Printf("CHR ROM:  %d KB\n", header.CHRROMSize*8)
	fmt.Printf("Flags 6:  %08b\n", header.Flags6)
	fmt.Printf("Flags 7:  %08b - (Mapper)\n", header.Flags7)
	fmt.Printf("Flags 8:  %d KB - (PRG RAM Size)\n", header.PRGRAMSize*8)
	fmt.Printf("Flags 9:  %08b\n", header.Flags9)
	fmt.Printf("Flags 10: %08b\n", header.Flags10)
}
