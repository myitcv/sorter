package main

import "testing"

func TestSort(t *testing.T) {
	sarah := person{"Sarah", 60}
	jill := person{"Jill", 34}
	paul := person{"Paul", 25}
	people := []person{
		sarah,
		jill,
		paul,
	}

	byName := []person{
		jill,
		paul,
		sarah,
	}

	byAge := []person{
		paul,
		jill,
		sarah,
	}

	checkEqual := func(ref []person) {
		for i := range people {
			if people[i] != ref[i] {
				t.Fatalf("%v should have been %v\n", people[i], ref[i])
			}
		}
	}

	sortByName(people)

	checkEqual(byName)

	SortByAge(people)

	checkEqual(byAge)
}
