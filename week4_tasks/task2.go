package main

import "fmt"

type Person struct {
	Name  string
	Email string
}

func (p Person) Contact() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)

}

type Student struct {
	Person
	GPA float64
}

type Teacher struct {
	Person
	Department string
}

func (t Teacher) Contact() string {
	return fmt.Sprintf("%s<%s>, %s Dept.", t.Name, t.Email, t.Department)
}

func main() {
	s := Student{
		Person: Person{Name: "Alice", Email: "alice@uni.edu"},
		GPA:    3.8,
	}

	t := Teacher{
		Person:     Person{Name: "Dr. Smith", Email: "smith@uni.edu"},
		Department: "Computer Science",
	}

	fmt.Println(s.Contact()) // Student uses promoted Person.Contact()
	fmt.Println(t.Contact()) // Teacher's own Contact() (shadows Person's)
	fmt.Println(t.Person.Contact())
}
