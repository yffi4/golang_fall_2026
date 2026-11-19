package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Task 4: Parse all the lines

var ErrBadFormat = errors.New("bad format")

type Record struct {
	Name  string
	Grade int
}

func ParseAll(input string) ([]Record, error) {
	var records []Record
	var errs []error
	for i, line := range strings.Split(input, "\n") {
		lineNo := i + 1
		name, num, ok := strings.Cut(line, ",")
		if !ok {
			errs = append(errs, fmt.Errorf("line %d: %w: %q", lineNo, ErrBadFormat, line))
			continue
		}
		grade, err := strconv.Atoi(strings.TrimSpace(num))
		if err != nil {
			errs = append(errs, fmt.Errorf("line %d: %w", lineNo, err))
			continue
		}
		records = append(records, Record{Name: strings.TrimSpace(name), Grade: grade})
	}
	// errors.Join returns nil when errs is empty, so a fully correct input
	// gives a nil error.
	return records, errors.Join(errs...)
}

func main() {
	input := `Aruzhan, 95
Dias, 8O
Timur
Bota, 99999999999999999999
Ali, 70`

	records, err := ParseAll(input)
	fmt.Println("records:", records)
	fmt.Println("errors:")
	fmt.Println(err)

	fmt.Println("Is ErrBadFormat:  ", errors.Is(err, ErrBadFormat))
	fmt.Println("Is ErrSyntax:     ", errors.Is(err, strconv.ErrSyntax))
	fmt.Println("Is ErrRange:      ", errors.Is(err, strconv.ErrRange))

	var ne *strconv.NumError
	if errors.As(err, &ne) {
		fmt.Printf("first NumError: func=%s num=%q err=%v\n", ne.Func, ne.Num, ne.Err)
	}

	_, err = ParseAll("Aruzhan, 95\nAli, 70")
	fmt.Println("all correct -> err == nil:", err == nil) // true
}
