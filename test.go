package main

import (
	"fmt"
	"github.com/beenotung/goutils/log"
)

func init() {
	log.Init(true, true, true, log.ShortCommFlag)
}
func test(start, end int) {
	N := end - start
	n := 1 // runtime.GOMAXPROCS(0)
	s := N / n
	fmt.Println("N", N, "n", n, "s", s)
	if N < n {
		for ; start < end; start++ {
			fmt.Println(start)
		}
	} else {
		for i := 0; i < n; i++ {
			a := i*s + start
			b := (i+1)*s - 1 + start
			fmt.Println(i, a, b)
		}
		if n*s != N {
			test(n*s+start, end)
		}
	}

}
func main() {
	test(3, 159)
}
