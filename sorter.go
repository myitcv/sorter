package sorter

type Order bool

// Inspired by https://github.com/mattn/sorter
type (
	FLen  func() int
	FLess func(i, j int) bool
	FSwap func(i, j int)
)

type Wrapper struct {
	LenFunc  FLen
	LessFunc FLess
	SwapFunc FSwap
}

func (w *Wrapper) Len() int {
	return w.LenFunc()
}

func (w *Wrapper) Less(i, j int) bool {
	return w.LessFunc(i, j)
}

func (w *Wrapper) Swap(i, j int) {
	w.SwapFunc(i, j)
}
