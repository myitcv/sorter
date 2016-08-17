//go:generate sortGen
package main

import mysorter "github.com/myitcv/sorter"

func orderByAge(persons []person, i, j int) mysorter.Order {
	return persons[i].age < persons[j].age
}
