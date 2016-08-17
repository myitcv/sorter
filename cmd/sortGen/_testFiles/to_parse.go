//go:generate sortGen
package main

import (
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

	sortByAge(people)

	fmt.Printf("Age sorted: %v\n", people)
}

func orderByName(persons []person, i, j int) sorter.Order {
	return persons[i].name < persons[j].name
}

func orderByAge(persons []person, i, j int) sorter.Order {
	return persons[i].age < persons[j].age
}
