package main

import (
	"fmt"
	"github.com/aabbcc1241/goai/ga"
	"github.com/aabbcc1241/goutils/log"
	"runtime"
	"time"
)

/* demo application of the ga
 *   maximizing number of 1 in gen code
 *     a_mutation     | stepCount for all code to be 1
 *       0.0001       | 7507
 *       0.0005833333 | 3865
 *       0.00065      | 3499
 *       0.000825     | 1323
 *       0.00086875   | 1199
 *       0.0009       | 10001 (excess limit)
 *       0.00091      | 830
 *       0.000911     | 642
 *       0.000912     | 642
 *       0.0009125    | 642
 *       0.0009128    | 642
 *       0.000913     | 926
 *       0.000914     | 1281
 *       0.000915     | 1351
 *       0.00092      | 10001 (excess limit)
 *       0.000934375  | 1821
 *       0.001        | 2153
 * the parameter is for user application initial guess reference
 *
 * parallel support
 *   when run for 1000 steps (100 pop, 1000 gen_len)
 *     single thread             : 4 seconds
 *     8 thread on 8 core system : 7 seconds
 *   when run for 10000 steps (100 pop, 1000 gen_len)
 *     single thread             : 39 seconds
 *     8 thread on 8 core system : 75 seconds
 *   when run for 1000 steps (1000 pop, 1000 gen_len)
 *     single thread             : 39 seconds
 *     8 thread on 8 core system : 64 seconds
 *   when run for 100 steps (10000 pop, 1000 gen_len)
 *     single thread             : 39 seconds
 *     8 thread on 8 core system : 67 seconds
 *   when run for 1000 steps (32 pop, 10000 gen_len)
 *     single thread             : 12 seconds
 *     8 thread on 8 core system : 40 seconds
 *   when run for 1000 steps (16 pop, 100000 gen_len)
 *     single thread             : 57 seconds
 *     8 thread on 8 core system : 26 seconds
 *       using ParallelArray     : xx seconds
 */
func init() {
	log.Init(true, true, true, log.ShortCommFlag)
}

type Fitness_i struct {
}

func (Fitness_i) Apply(gen ga.Gene_s) float64 {
	i := float64(0)
	for _, v := range gen.Code {
		i += float64(v)
	}
	//log.Debug.Println("fitness:",i)
	return i
}
func test(start, end int) {
	N := end - start
	n := runtime.GOMAXPROCS(0)
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
	//test(56, 159)
	//return
	fmt.Println("start")
	//fmt.Println("cpu",runtime.NumCPU())
	//fmt.Println("goroutine",runtime.NumGoroutine())
	//fmt.Println("max go routine",runtime.GOMAXPROCS(0))
	fmt.Println(time.Now())

	ga_s := ga.GA_s{
		P_CrossOver:   0.8,
		P_Mutation:    0.2,
		A_Mutation:    0.000912,
		Fitness_i:     Fitness_i{},
		IsMultiThread: false,
	}
	ga_s.Init(16, 100000)
	ga_s.RunN(1000, false)
	//stepCount, excessLimit := ga_s.RunUntil(1000, 10000)
	//log.Info.Println("stepCount", stepCount, "earlyTerm", excessLimit)

	fmt.Println("end")
	fmt.Println(time.Now())
}
