package main

import (
	"fmt"
	"github.com/aabbcc1241/goai/ga"
	"github.com/aabbcc1241/goutils/log"
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
func main() {
	fmt.Println("start")
	fmt.Println(time.Now())

	ga_s := ga.GA_s{
		P_CrossOver: 0.8,
		P_Mutation:  0.2,
		A_Mutation:  0.000912,
		Fitness_i:   Fitness_i{},
	}
	ga_s.Init(100, 1000)
	//ga_s.RunN(10000)
	stepCount, excessLimit := ga_s.RunUntil(1000, 10000)
	log.Info.Println("stepCount", stepCount, "earlyTerm", excessLimit)

	fmt.Println("end")
	fmt.Println(time.Now())
}
