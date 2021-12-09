package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var sessionPath = path.Join(os.Getenv("HOME"), ".advent-of-code-session")
var cacheDir = path.Join(os.Getenv("HOME"), ".cache/advent-of-code")

// Input returns input for the specified date, e.g. "2020/13" (leading zero optional).
// Input is downloaded and cached. If a "-" command-line flag was provided, data is
// read from stdin instead.
func Input(date string) string {
	var year, day int
	Extract(date, `^(\d{4})/(\d{1,2})$`, &year, &day)

	// Read from stdin if a "-" arg was supplied.
	if len(os.Args) == 2 && os.Args[1] == "-" {
		var b bytes.Buffer
		_, err := io.Copy(&b, os.Stdin)
		if err != nil {
			panic(fmt.Sprint("Failed to read input: ", err))
		}
		return b.String()
	}

	// Try to read cached input.
	p := filepath.Join(cacheDir, strconv.Itoa(year), strconv.Itoa(day))
	if _, err := os.Stat(p); err == nil {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			panic(fmt.Sprint("Failed reading cached input: ", err))
		}
		return string(b)
	}

	// Download input.
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	log.Print("Downloading ", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed creating request to %v: %v", url, err))
	}

	session, err := ioutil.ReadFile(sessionPath)
	if err != nil {
		panic(fmt.Sprint("Failed reading session ID: ", err))
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: strings.TrimSpace(string(session))})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Failed requesting %v: %v", url, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Getting %v failed: %v", url, resp.Status))
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed reading %v: %v", url, err))
	}

	// Save input.
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		panic(fmt.Sprint("Failed creating cache dir: ", err))
	}
	if err := ioutil.WriteFile(p, b, 0644); err != nil {
		panic(fmt.Sprint("Failed writing input: ", err))
	}

	return string(b)
}

// InputInts extracts integers from the day's input. See ExtractInts.
func InputInts(date string) []int {
	return ExtractInts(Input(date))
}

// InputInt64s extracts 64-bit integers from the day's input. See ExtractInt64s.
func InputInt64s(date string) []int64 {
	return ExtractInt64s(Input(date))
}

// InputLines returns newline-separated lines of input for the specified day.
func InputLines(date string) []string {
	var lines []string
	for _, ln := range strings.Split(Input(date), "\n") {
		if len(ln) > 0 {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 {
		panic("No lines found")
	}
	return lines
}

// InputLinesBytes returns newline-separated lines of input for the specified day.
// If valid is non-empty, panics if any unlisted bytes are encountered.
func InputLinesBytes(date string, valid ...byte) [][]byte {
	return NewByteLines(Input(date), valid...)
}

// InputByteGrid returns a ByteGrid holding the input for the specified day.
// If valid is non-empty, panics if any unlisted bytes are encountered.
func InputByteGrid(date string, valid ...byte) ByteGrid {
	return NewByteGridString(Input(date), valid...)
}

var newlinesRegexp = regexp.MustCompile(`\n\n+`)

// InputParagraphs returns the day's input split into paragraphs on multiple
// newlines. Each paragraph is split further into individual lines.
func InputParagraphs(date string) [][]string {
	var pgs [][]string
	all := strings.Trim(Input(date), "\n")
	for _, pg := range newlinesRegexp.Split(all, -1) {
		if len(pg) > 0 {
			pgs = append(pgs, strings.Split(pg, "\n"))
		}
	}
	if len(pgs) == 0 {
		panic("No paragraphs found")
	}
	return pgs
}
