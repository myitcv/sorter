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

func orderByName(persons []person, i, j int) sorter.Order {
	return persons[i].name < persons[j].name
}

// this will not get matched
func order(persons []person, i, j int) sorter.Order {
	return persons[i].name < persons[j].name
}

// TODO
func orderBufferByContents(buffers []bytes.Buffer, i int, j int) sorter.Order {
	return buffers[i].String() < buffers[j].String()
}
