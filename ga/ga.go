package ga

import (
	"math/rand"
	"github.com/aabbcc1241/goutils/log"
	"sort"
)

type Fitness_i interface {
	Apply(Gene_s) float64
}
type Float64PreCompute struct {
	Value   float64
	IsValid bool
}

func (p *Float64PreCompute)Set(v float64) {
	p.Value = v
	p.IsValid = true
}

type Gene_s struct {
	Code    []byte
	Fitness Float64PreCompute
}
type population []Gene_s
type GA_s struct {
	Population  population
	P_CrossOver float64
	P_Mutation  float64
	Fitness_i   Fitness_i
}

func (p population)Len() int {
	return len(p)
}
func (p population)Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p population)Less(i, j int) bool {
	return p[i].Fitness.Value < p[j].Fitness.Value
}

func (p *GA_s)Init(n_pop int, gen_length int) {
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
func (p *GA_s)Run() {
	n_pop := len(p.Population)
	//gen_length := len(p.Population[0].Code)
	/* 1. get fitness */
	for _, v := range (p.Population) {
		v.Fitness.Set(p.Fitness_i.Apply(v))
	}
	/* 2. crossover TODO */
	sort.Sort(p.Population)
	new_pop := make([]Gene_s, n_pop)
	for i_pop := 0; i_pop < n_pop; i_pop++ {
		a := i_pop - i_pop % 2
		//b:=a+1
		code := p.Population[a].Code
		new_pop[i_pop] = Gene_s{
			Code:code,
		}
	}
	//for k, v := range (p.Population) {
	//}
	p.Population = new_pop
	/* 3. mutation TODO */
}
func (p *GA_s)RunN(n int) {
	for i := 0; i < n; i++ {
		log.Info.Printf("%v/%v", i, n)
		p.Run()
	}
	log.Info.Printf("%v/%v", n, n)
}
//func (p *GA_s)RunWhile
