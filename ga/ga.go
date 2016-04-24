package ga

import (
	"math/rand"
)

type Fitness_i interface {
	Fitness() float64
}
type Float64PreCompute struct {
	Value   float64
	IsValid bool
}
type Gene_s struct {
	Code    []byte
	Fitness Float64PreCompute
}
type GA_s struct {
	P_Crossover float64
	P_Mutation  float64
	Population  []Gene_s
}

func Init_Basic(n_pop int, gen_length int, p_crossover float64, p_mutation float64) GA_s {
	population := make([]Gene_s, n_pop)
	for i := 0; i < n_pop; i++ {
		code := make([]byte, gen_length)
		rand.Read(code)
		for i := 0; i < gen_length; i++ {
			code[i] = code[i] % 2
		}
		population[i] = Gene_s{
			Code:code,
		}
	}
	return GA_s{
		P_Crossover:p_crossover,
		P_Mutation:p_mutation,
		Population:population,
	}
}
