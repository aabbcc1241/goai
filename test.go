package main

import (
	"fmt"

	"github.com/aabbcc1241/goai/ga"
	"time"
	"github.com/aabbcc1241/goutils/log"
)

func init() {
	log.Init(true, true, true, log.ShortCommFlag)
}

type Fitness_i struct {

}

func (Fitness_i)Apply(gen ga.Gene_s) float64 {
	i := float64(0)
	for _, v := range (gen.Code) {
		i += float64(v)
	}
	return i
}
func main() {
	fmt.Println("start")
	fmt.Println(time.Now())

	ga_s := ga.GA_s{
		Fitness_i:Fitness_i{},
	}
	ga_s.Init(100, 1000)
	ga_s.RunN(100)

	fmt.Println("end")
	fmt.Println(time.Now())
}
