package test

import (
	"github.com/beenotung/goai/ga"
	"github.com/beenotung/goutils/log"
	"testing"
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
 *       0.000912     | 642	(suggested)
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
 * parallel tune (n_pop:1000, gen_len:10000, n_step:50)
 *   all parallel
 *     1 thread  : 20.535 s
 *     8 threads : 54.526 s
 *   exclude init
 *     1 thread  : 20.790 s
 *     8 threads : 56.872 s
 *   exclude measure
 *     1 thread  : 20.470 s
 *     8 threads : 55.046 s
 *   exclude crossover		(selected)
 *     1 thread  : 20.840 s
 *     8 threads : 20.167 s
 *   exclude crossover and mutation
 *     1 thread  : 23.105 s
 *     8 threads : 22.665 s
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
func TestGa(t *testing.T) {
	ga_s := ga.GA_s{
		P_CrossOver: 0.8,
		P_Mutation:  0.2,
		A_Mutation:  0.000912,
		Fitness_i:   Fitness_i{},
	}
	nThread := 1
	ga_s.Init(1000, 10000, nThread)
	ga_s.RunN(50, false)
	//stepCount, excessLimit := ga_s.RunUntil(1000, 10000)
	//log.Info.Println("stepCount", stepCount, "earlyTerm", excessLimit)
}
