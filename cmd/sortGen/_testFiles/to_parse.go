//go:generate sortGen
package main

import (
	"bytes"
	"fmt"

	"github.com/myitcv/sorter"
)

type person struct {
	name string
	age  int
}

func main() {
	people := []person{
		person{"Sarah", 60},
		person{"Jill", 34},
		person{"Paul", 25},
	}

	fmt.Printf("Before: %v\n", people)

	sortByName(people)

	fmt.Printf("Name sorted: %v\n", people)

	SortByAge(people)

	fmt.Printf("Age sorted: %v\n", people)
}

// MATCH
func orderByName(persons []person, i, j int) sorter.Ordered {
	return persons[i].name < persons[j].name
}

// fail
func order(persons []person, i, j int) sorter.Ordered {
	return persons[i].name < persons[j].name
}

// MATCH
func orderPointerByName(persons []*person, i, j int) sorter.Ordered {
	return persons[i].name < persons[j].name
}

// MATCH
func orderBufferByContents(buffers []bytes.Buffer, i int, j int) sorter.Ordered {
	return buffers[i].String() < buffers[j].String()
}

func orderMap(buffers []map[string]bool, i int, j int) sorter.Ordered {
	return true
}
