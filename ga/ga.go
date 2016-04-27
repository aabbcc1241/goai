package ga

import (
	"github.com/aabbcc1241/goutils/log"
	"math/rand"
	"sort"
	"sync"
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
type population []Gene_s

//TODO make a constructor for GA_s
type GA_s struct {
	Population  population
	P_CrossOver float64
	/* probability of mutation */
	P_Mutation float64
	/* amount (percentage) of mutation within that gen */
	A_Mutation float64
	Fitness_i  Fitness_i
	/* at least 1 */
	NumberOfThread int
}

func (p population) Len() int {
	return len(p)
}
func (p population) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p population) Less(i, j int) bool {
	return p[i].Fitness.Value < p[j].Fitness.Value
}

func (p *GA_s) Init(n_pop int, gen_length int) {
	if p.NumberOfThread == 0 {
		p.NumberOfThread = 1
	}
	p.Population = make([]Gene_s, n_pop)
	for i_pop := 0; i_pop < n_pop; i_pop++ {
		code := make([]byte, gen_length)
		rand.Read(code)
		for i := 0; i < gen_length; i++ {
			code[i] = code[i] % 2
		}
		p.Population[i_pop] = Gene_s{
			Code: code,
		}
	}
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
func (p *code_t) mutation(a_mutation float64) {
	for i := range *p {
		if rand.Float64() < a_mutation {
			(*p)[i] = 1 - (*p)[i]
		}
	}
}
func measure_and_sort(p *GA_s) (bestFitness float64) {
	for i := range p.Population {
		if !p.Population[i].Fitness.IsValid {
			p.Population[i].Fitness.Set(p.Fitness_i.Apply(p.Population[i]))
		}
	}
	sort.Sort(p.Population)
	return p.Population[len(p.Population)-1].Fitness.Value
}

/* REMARK : require pre-sorted population */
func crossover(p *GA_s) {
	n_pop := len(p.Population)
	n_crossover := (int)(p.P_CrossOver * float64(n_pop))
	if p.NumberOfThread == 1 {
		/* replace first p_mutation percent for bad gens */
		for i_pop := 0; i_pop < n_crossover; i_pop++ {
			/* select from last (1-p_mutation) parent pool */
			a := n_crossover + rand.Intn(n_pop-n_crossover)
			b := n_crossover + rand.Intn(n_pop-n_crossover)
			p.Population[i_pop].Code.crossover(&p.Population[a].Code, &p.Population[b].Code)
			p.Population[i_pop].Fitness.IsValid = false
		}
	} else {
		/* parallel version */
		wg := sync.WaitGroup{}
		wg.Add(p.NumberOfThread)
		for i_thread := 0; i_thread < p.NumberOfThread; i_thread++ {
			//log.Debug.Println("crossover:fork")
			go func(i_thread int) {
				defer wg.Done()
				for i_pop := 0; i_pop < n_crossover; i_pop++ {
					if i_pop%p.NumberOfThread != i_thread {
						continue
					}
					/* select from last (1-p_mutation) parent pool */
					a := n_crossover + rand.Intn(n_pop-n_crossover)
					b := n_crossover + rand.Intn(n_pop-n_crossover)
					p.Population[i_pop].Code.crossover(&p.Population[a].Code, &p.Population[b].Code)
					p.Population[i_pop].Fitness.IsValid = false
				}
			}(i_thread)
		}
		//log.Debug.Println("crossover:wait")
		wg.Wait()
		//log.Debug.Println("crossover:joined")
	}
}
func mutation(p *GA_s) {
	if p.NumberOfThread == 1 {
		for i := range p.Population {
			if rand.Float64() < p.P_Mutation {
				p.Population[i].Code.mutation(p.A_Mutation)
			}
		}
	} else {
		wg := sync.WaitGroup{}
		wg.Add(p.NumberOfThread)
		for i_thread := 0; i_thread < p.NumberOfThread; i_thread++ {
			//log.Debug.Println("mutation:fork")
			go func(i_thread int) {
				defer wg.Done()
				for i := range p.Population {
					if i%p.NumberOfThread != i_thread {
						continue
					}
					if rand.Float64() < p.P_Mutation {
						p.Population[i].Code.mutation(p.A_Mutation)
					}
				}
			}(i_thread)
		}
		//log.Debug.Println("mutation:wait")
		wg.Wait()
		//log.Debug.Println("mutation:joined")
	}
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
