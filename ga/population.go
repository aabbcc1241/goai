/* ParallelArray of Gene_s */

package ga

import (
	. "github.com/beenotung/goutils/lang"
	"math/rand"
	"runtime"
	"sync"
)

type population struct {
	Genes   []Gene_s
	Lock    sync.Mutex
	NThread int
}
type consumer interface {
	Apply(k int, v Gene_s, r *rand.Rand)
}
type producer interface {
	Apply(k int, v Gene_s, r *rand.Rand) Gene_s
}
type inplace_updater interface {
	Apply(k int, v *Gene_s, r *rand.Rand)
}

func (p population) Len() int {
	return len(p.Genes)
}

/* [start,end) : end is excluded
 * REMARK : this function does not handle lock
 */
func _for(p population, f consumer, withRandom bool, start, end int) {
	N := end - start
	n := runtime.GOMAXPROCS(0)
	if p.NThread > 0 {
		n = p.NThread
	}
	if N < n {
		sem := make(Semaphore, N)
		sem.P(N)
		for ; start < end; start++ {
			go func(i int) {
				var r *rand.Rand
				if withRandom {
					r = rand.New(rand.NewSource(int64(i)))
				}
				f.Apply(i, p.Genes[i], r)
				sem.Signal()
			}(start)
		}
		sem.Wait(N)
	} else {
		s := N / n
		sem := make(Semaphore, n)
		sem.P(n)
		for i := 0; i < n; i++ {
			go func(i int) {
				var r *rand.Rand
				if withRandom {
					r = rand.New(rand.NewSource(int64(i)))
				}
				for j := i*s + start; j < (i+1)*s+start; j++ {
					f.Apply(j, p.Genes[j], r)
				}
				sem.Signal()
			}(i)
		}
		sem.Wait(n)
		if n*s != N {
			_for(p, f, withRandom, n*s+start, end)
		}
	}
}
func _replace(p *population, f producer, withRandom bool, start, end int) {
	N := end - start
	n := runtime.GOMAXPROCS(0)
	if p.NThread > 0 {
		n = p.NThread
	}
	if N < n {
		sem := make(Semaphore, N)
		sem.P(N)
		for ; start < end; start++ {
			go func(i int) {
				var r *rand.Rand
				if withRandom {
					r = rand.New(rand.NewSource(int64(i)))
				}
				p.Genes[i] = f.Apply(i, p.Genes[i], r)
				sem.Signal()
			}(start)
		}
		sem.Wait(N)
	} else {
		s := N / n
		sem := make(Semaphore, n)
		sem.P(n)
		for i := 0; i < n; i++ {
			go func(i int) {
				var r *rand.Rand
				if withRandom {
					r = rand.New(rand.NewSource(int64(i)))
				}
				for j := i*s + start; j < (i+1)*s+start; j++ {
					p.Genes[j] = f.Apply(j, p.Genes[j], r)
				}
				sem.Signal()
			}(i)
		}
		sem.Wait(n)
		if n*s != N {
			_replace(p, f, withRandom, n*s+start, end)
		}
	}
}
func _inplace_update(p *population, f inplace_updater, withRandom bool, start, end int) {
	N := end - start
	n := runtime.GOMAXPROCS(0)
	if p.NThread > 0 {
		n = p.NThread
	}
	if N < n {
		sem := make(Semaphore, N)
		sem.P(N)
		for ; start < end; start++ {
			go func(i int) {
				var r *rand.Rand
				if withRandom {
					r = rand.New(rand.NewSource(int64(i)))
				}
				f.Apply(i, &p.Genes[i], r)
				sem.Signal()
			}(start)
		}
		sem.Wait(N)
	} else {
		s := N / n
		sem := make(Semaphore, n)
		sem.P(n)
		for i := 0; i < n; i++ {
			go func(i int) {
				var r *rand.Rand
				if withRandom {
					r = rand.New(rand.NewSource(int64(i)))
				}
				for j := i*s + start; j < (i+1)*s+start; j++ {
					f.Apply(j, &p.Genes[j], r)
				}
				sem.Signal()
			}(i)
		}
		sem.Wait(n)
		if n*s != N {
			_inplace_update(p, f, withRandom, n*s+start, end)
		}
	}
}
func (p population) For(f consumer, withRandom bool) {
	p.Lock.Lock()
	_for(p, f, withRandom, 0, len(p.Genes))
	p.Lock.Unlock()
}

func (p *population) Replace(f producer, withRandom bool) {
	p.Lock.Lock()
	_replace(p, f, withRandom, 0, len(p.Genes))
	p.Lock.Unlock()
}

func (p *population) Inplace_Update(f inplace_updater, withRandom bool) {
	p.Lock.Lock()
	_inplace_update(p, f, withRandom, 0, len(p.Genes))
	p.Lock.Unlock()
}

//TODO map, reduce, fold
