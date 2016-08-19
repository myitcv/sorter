package main

import "sort"
import "github.com/myitcv/sorter"

func sortByName(vs []person) {
	sort.Sort(&sorter.Wrapper{
		LenFunc: func() int {
			return len(vs)
		},
		LessFunc: func(i, j int) bool {
			return bool(orderByName(vs, i, j))
		},
		SwapFunc: func(i, j int) {
			vs[i], vs[j] = vs[j], vs[i]
		},
	})
}

func sortByAge(vs []person) {
	sort.Sort(&sorter.Wrapper{
		LenFunc: func() int {
			return len(vs)
		},
		LessFunc: func(i, j int) bool {
			return bool(orderByAge(vs, i, j))
		},
		SwapFunc: func(i, j int) {
			vs[i], vs[j] = vs[j], vs[i]
		},
	})
}
