package main

import (
	"fmt"
	"github.com/aabbcc1241/goai/ga"
)

func main() {
	fmt.Println("start")

	ga_profile := ga.Init_Basic(100, 1000, 0.5, 0.01)
	c1 := 0
	c2 := 0
	for _, v := range (ga_profile.Population) {
		for _, v := range (v.Code) {
			if v == 0 {
				c1++
			} else if v == 1 {
				c2++
			} else {
				fmt.Println("wrong value", v, v % 2)
			}
		}
	}
	fmt.Println("population", ga_profile.Population)
	fmt.Println("0 count:", c1, "1 count:", c2)
	fmt.Println("ratio", float64(c1) / float64(c1 + c2))

	fmt.Println("end")
}
