// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

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

	m := &myStruct{}
	m.sortByAge(people)

	fmt.Printf("Age sorted: %v\n", people)
}

func orderByName(persons []person, i, j int) sorter.Ordered {
	return persons[i].name < persons[j].name
}

type myStruct struct {
	// some fields
}

func (m *myStruct) orderByAge(persons []person, i, j int) sorter.Ordered {
	return persons[i].age < persons[j].age
}
