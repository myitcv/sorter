package main

import "sort"
import "github.com/myitcv/sorter"

func sortPointerByName(vs []*person) {
	sort.Sort(&sorter.Wrapper{
		LenFunc: func() int {
			return len(vs)
		},
		LessFunc: func(i, j int) bool {
			return bool(orderPointerByName(vs, i, j))
		},
		SwapFunc: func(i, j int) {
			vs[i], vs[j] = vs[j], vs[i]
		},
	})
}
