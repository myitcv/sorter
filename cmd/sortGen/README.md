## `sortGen`

A [`go generate`](https://blog.golang.org/generate)-or that makes sorting arbitrary slices easier via order functions, removing
the need to implement [`sort.Interface`](https://godoc.org/sort#Interface) etc.

Simply define an order function on the slice type of interest:

```go
// main.go
//go:generate sortGen

func orderByName(persons []person, i, j int) sorter.Ordered {
	return persons[i].name < persons[j].name
}
```

then run `go generate` and a corresponding sort function will have been generated for you:

```go
// gen_main_sorter.go

func sortByName(vs []person) {
	...
}
```

See the example and rules below for more details.

### Install

```
go get -u github.com/myitcv/sorter/cmd/sortGen
```

### Example

Taking the example from [`example/main.go`](https://github.com/myitcv/sorter/blob/master/example/main.go):

```go
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

func orderByName(persons []person, i, j int) sorter.Ordered {
	return persons[i].name < persons[j].name
}

func orderByAge(persons []person, i, j int) sorter.Ordered {
	return persons[i].age < persons[j].age
}
```

Then:

```
$ go generate
$ go build
$ ./example
Before: [{Sarah 60} {Jill 34} {Paul 25}]
Name sorted: [{Jill 34} {Paul 25} {Sarah 60}]
Age sorted: [{Paul 25} {Jill 34} {Sarah 60}]
```

Examine the contents of `gen_main_sorter.go` to see the generated functions.

### Rules

`sortGen` generates sort functions/methods according to the following simple rules:

1. The file, e.g. `my_file.go`, containing the order function/method must include the directive `//go:generate sortGen`
2. The order function/method name must be of the form `"order*"` or `"Order*"` (more strictly `^[oO]rder[[:word:]]+` in a [regex](https://godoc.org/regexp)
   [pattern](https://github.com/google/re2/wiki/Syntax))
3. The parameters of the order function/method must be a slice type, followed by two `int`'s
4. The return type must be `github.com/myitcv/sorter.Ordered`

The sort functions/methods generated will be of the form `"sort*"` or `"Sort*"` (following the capitalisation
of the order function). They will be written to a file with a name corresponding to the input file,
`gen_my_file_sorter.go` in the case of the file mentioned in point 1.

Notice the use of the term "function/method"; if the order function defines a receiver (i.e. it is a method)
then the generated sort function will use the same receiver, and hence be a method.

### Implementation

The current implementation of the generator simply wraps a call to `sort.Sort`; this of course can be improved...

### Bugs

This is only an initial proof-of-concept, probably lots of bugs and edge cases missed. Please raise issues...
