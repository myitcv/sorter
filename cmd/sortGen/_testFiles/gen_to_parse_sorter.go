package main

import "sort"
import "github.com/myitcv/sorter"

import "bytes"

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

func sortBufferByContents(vs []bytes.Buffer) {
	sort.Sort(&sorter.Wrapper{
		LenFunc: func() int {
			return len(vs)
		},
		LessFunc: func(i, j int) bool {
			return bool(orderBufferByContents(vs, i, j))
		},
		SwapFunc: func(i, j int) {
			vs[i], vs[j] = vs[j], vs[i]
		},
	})
}

func sortMap(vs []map[string]bool) {
	sort.Sort(&sorter.Wrapper{
		LenFunc: func() int {
			return len(vs)
		},
		LessFunc: func(i, j int) bool {
			return bool(orderMap(vs, i, j))
		},
		SwapFunc: func(i, j int) {
			vs[i], vs[j] = vs[j], vs[i]
		},
	})
}

func (e *example) sortBanana(vs []string) {
	sort.Sort(&sorter.Wrapper{
		LenFunc: func() int {
			return len(vs)
		},
		LessFunc: func(i, j int) bool {
			return bool(e.orderBanana(vs, i, j))
		},
		SwapFunc: func(i, j int) {
			vs[i], vs[j] = vs[j], vs[i]
		},
	})
}
