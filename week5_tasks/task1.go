package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Task 1: Notifications

type Notifier interface {
	Notify(msg string) error
}

type Email struct {
	Address string
}

func (e Email) Notify(msg string) error {
	fmt.Printf("email -> %s: %s\n", e.Address, msg)
	return nil
}

type SMS struct {
	Phone string
}

func (s SMS) Notify(msg string) error {
	if len(msg) > 20 {
		return fmt.Errorf("sms to %s: message too long (%d > 20 chars)", s.Phone, len(msg))
	}
	fmt.Printf("sms -> %s: %s\n", s.Phone, msg)
	return nil
}

// Telegram counts how many messages it has sent, so it needs a pointer
// receiver: a value receiver would increment a copy and the caller would
// never see the updated count.
type Telegram struct {
	Chat  string
	Count int
}

func (t *Telegram) Notify(msg string) error {
	t.Count++
	fmt.Printf("telegram -> %s: %s (sent #%d)\n", t.Chat, msg, t.Count)
	return nil
}

// Compile-time checks that each type implements Notifier.
var _ Notifier = Email{}
var _ Notifier = SMS{}
var _ Notifier = &Telegram{}

func Broadcast(ns []Notifier, msg string) []error {
	var errs []error
	for _, n := range ns {
		if err := n.Notify(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Participation Task (Part 1): a function that accepts an io.Writer and
// prints a table of grades. The function doesn't care whether it is
// writing to the terminal or into a string - that's the whole point of
// a small interface.

type Grade struct {
	Name  string
	Score int
}

func printGrades(w io.Writer, grades []Grade) {
	fmt.Fprintf(w, "%-12s %s\n", "Name", "Score")
	for _, g := range grades {
		fmt.Fprintf(w, "%-12s %5d\n", g.Name, g.Score)
	}
}

func main() {
	tg := &Telegram{Chat: "golang-fall-2026"}

	notifiers := []Notifier{
		Email{Address: "student@example.com"},
		SMS{Phone: "+77001234567"},
		tg,
	}

	fmt.Println("=== Task 1: Notifications ===")

	fmt.Println("--- short message ---")
	errs := Broadcast(notifiers, "Hi there!")
	fmt.Println("errors:", errs)

	fmt.Println("--- long message (SMS should fail) ---")
	errs = Broadcast(notifiers, "This message is definitely longer than twenty characters")
	fmt.Println("errors:", errs)

	fmt.Println("--- second short message (Telegram counter should increase) ---")
	errs = Broadcast(notifiers, "Second ping")
	fmt.Println("errors:", errs)

	fmt.Println("telegram total sent:", tg.Count)

	fmt.Println("\n=== Participation Task: grades table via io.Writer ===")

	grades := []Grade{
		{"Aruzhan", 95},
		{"Dias", 80},
		{"Yerlan", 67},
		{"Saniya", 88},
	}

	fmt.Println("--- writing to os.Stdout ---")
	printGrades(os.Stdout, grades)

	var sb strings.Builder
	printGrades(&sb, grades)

	fmt.Println("--- writing to *strings.Builder, then printed ---")
	fmt.Print(sb.String())
}
