package main

import (
	"fmt"

	"github.com/aabbcc1241/goai/ga"
	"time"
)

func main() {
	fmt.Println("start")
	fmt.Println(time.Now())

	ga_s:=ga.GA_s{}
	ga_s.Init(100,1000)

	fmt.Println("end")
	fmt.Println(time.Now())
}
