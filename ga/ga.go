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
	Population  []Gene_s
	P_CrossOver float64
	P_Mutation  float64
}

func (p*GA_s)Init(n_pop int, gen_length int) {
	p.Population = make([]Gene_s, n_pop)
	for i_pop := 0; i_pop < n_pop; i_pop++ {
		code := make([]byte, gen_length)
		rand.Read(code)
		for i := 0; i < gen_length; i++ {
			code[i] = code[i] % 2
		}
		p.Population[i_pop] = Gene_s{
			Code:code,
		}
	}
}
