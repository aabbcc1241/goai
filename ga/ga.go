package ga

import (
	"github.com/aabbcc1241/goutils/log"
	"math/rand"
	"sort"
)

type Fitness_i interface {
	Apply(Gene_s) float64
}
type Float64PreCompute struct {
	Value   float64
	IsValid bool
}

func (p *Float64PreCompute) Set(v float64) {
	p.Value = v
	p.IsValid = true
}

type code_t []byte
type Gene_s struct {
	Code    code_t
	Fitness Float64PreCompute
}

//TODO make a constructor for GA_s
type GA_s struct {
	Population  population
	P_CrossOver float64
	/* probability of mutation */
	P_Mutation float64
	/* amount (percentage) of mutation within that gen */
	A_Mutation float64
	Fitness_i  Fitness_i
}

//func (p population) Len() int {
//	return len(p)
//}
func (p population) Swap(i, j int) {
	p.Genes[i], p.Genes[j] = p.Genes[j], p.Genes[i]
}
func (p population) Less(i, j int) bool {
	return p.Genes[i].Fitness.Value < p.Genes[j].Fitness.Value
}

/* updater */
type init_s struct {
	gen_length int
}

func (p init_s) Apply(k int, v Gene_s, r *rand.Rand) Gene_s {
	code := make([]byte, p.gen_length)
	r.Read(code)
	for i := 0; i < p.gen_length; i++ {
		code[i] = code[i] % 2
	}
	return Gene_s{Code: code}
}

func (p *GA_s) Init(n_pop int, gen_length int, nThread int) {
	p.Population = population{Genes: make([]Gene_s, n_pop), NThread: nThread}
	p.Population.Replace(init_s{gen_length}, true)
}

/* equally pick from parent */
func (child *code_t) crossover(parent1, parent2 *code_t) {
	for i := range *child {
		if rand.Intn(2) == 0 {
			(*child)[i] = (*parent1)[i]
		} else {
			(*child)[i] = (*parent2)[i]
		}
	}
}

/* equal-likely single bit flip */
func (p *code_t) mutation(a_mutation float64, r *rand.Rand) {
	for i := range *p {
		if r.Float64() < a_mutation {
			(*p)[i] = 1 - (*p)[i]
		}
	}
}

type measure_and_sort_s struct {
	p *GA_s
}

func (p measure_and_sort_s) Apply(k int, v Gene_s, r *rand.Rand) {
	gene := v
	if !gene.Fitness.IsValid {
		gene.Fitness.Set(p.p.Fitness_i.Apply(gene))
	}
}
func measure_and_sort(p *GA_s) (bestFitness float64) {
	p.Population.For(measure_and_sort_s{p}, false)
	sort.Sort(p.Population)
	return p.Population.Genes[p.Population.Len()-1].Fitness.Value
}

/* updater */
type crossover_s struct {
	n_pop       int
	n_crossover int
	Population  population
}

func (p crossover_s) Apply(k int, v *Gene_s, r *rand.Rand) {
	a := p.n_crossover + r.Intn(p.n_pop-p.n_crossover)
	b := p.n_crossover + r.Intn(p.n_pop-p.n_crossover)
	v.Code.crossover(&(p.Population.Genes[a].Code), &(p.Population.Genes[b].Code))
	v.Fitness.IsValid = false
}

/* REMARK : require pre-sorted population */
func crossover(p *GA_s) {
	n_pop := p.Population.Len()
	n_crossover := (int)(p.P_CrossOver * float64(n_pop))
	p.Population.Inplace_Update(crossover_s{n_pop, n_crossover, p.Population}, true)
}

/* updater */
type mutation_s struct {
	GA_s
}

func (p mutation_s) Apply(k int, v *Gene_s, r *rand.Rand) {
	if r.Float64() < p.P_Mutation {
		v.Code.mutation(p.A_Mutation, r)
	}
}
func mutation(p *GA_s) {
	p.Population.Inplace_Update(mutation_s{}, true)
}
func (p *GA_s) Run(verbose bool) {
	//gen_length := len(p.Population[0].Code)
	/* 1. get fitness */
	bestFitness := measure_and_sort(p)
	/* 2. crossover */
	if verbose {
		log.Info.Println("bestFitness:", bestFitness)
	}
	crossover(p)
	/* 3. mutation */
	mutation(p)
}
func (p *GA_s) RunN(n int, verbose bool) {
	for i := 0; i < n; i++ {
		if verbose {
			log.Info.Printf("%v/%v", i, n)
		}
		p.Run(verbose)
	}
	if verbose {
		log.Info.Printf("%v/%v", n, n)
	}
}
func (p *GA_s) RunUntil(targetBestFitness float64, stepUpperLimit int) (stepCount int, excessUpperLimit bool) {
	currentBestFitness := 0.0
	i := 0
	for {
		i++
		currentBestFitness = measure_and_sort(p)
		//log.Info.Println("step:", i, "bestFitness:", currentBestFitness)
		if currentBestFitness >= targetBestFitness {
			break
		}
		if i > stepUpperLimit {
			return i, true
		}
		crossover(p)
		mutation(p)
	}
	return i, false
}

//func (p *GA_s)RunWhile
