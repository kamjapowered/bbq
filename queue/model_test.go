package queue

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

// every implementation is also checked against a slice model driven by a
// deterministic pseudo random script, which reaches the wraparound and
// regrowth states that the hand written cases only sample

func TestModelFifo(t *testing.T) {
	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			for seed := range int64(100) {
				rng := rand.New(rand.NewSource(seed))
				q := newQueue()
				var model []int
				next := 0
				for range 300 {
					switch rng.Intn(10) {
					case 0, 1, 2, 3:
						if len(model) < roomy {
							if !q.enqueue(next) {
								t.Fatalf("seed %d: enqueue rejected at len %d", seed, len(model))
							}
							model = append(model, next)
							next++
						}
					case 4, 5, 6:
						v, ok := q.Pop()
						if len(model) == 0 {
							if ok {
								t.Fatalf("seed %d: dequeue ok on empty", seed)
							}
						} else {
							if !ok || v != model[0] {
								t.Fatalf("seed %d: dequeue %d %v want %d", seed, v, ok, model[0])
							}
							model = model[1:]
						}
					case 7:
						q.Clear()
						model = nil
					case 8:
						n := rng.Intn(5)
						i := 0
						for v := range q.Drain() {
							if v != model[0] {
								t.Fatalf("seed %d: drain %d want %d", seed, v, model[0])
							}
							model = model[1:]
							i++
							if i == n {
								break
							}
						}
					case 9:
						got := slices.Collect(q.Iter())
						if !slices.Equal(got, model) {
							t.Fatalf("seed %d: iter %v want %v", seed, got, model)
						}
					}
					if q.Len() != len(model) {
						t.Fatalf("seed %d: len %d want %d", seed, q.Len(), len(model))
					}
					if !slices.Equal(q.ToSlice(), model) {
						t.Fatalf("seed %d: toslice %v want %v", seed, q.ToSlice(), model)
					}
					pv, pok := q.Peek()
					if len(model) == 0 {
						if pok {
							t.Fatalf("seed %d: peek ok on empty", seed)
						}
					} else if !pok || pv != model[0] {
						t.Fatalf("seed %d: peek %d want %d", seed, pv, model[0])
					}
				}
			}
		})
	}
}

func TestModelDeque(t *testing.T) {
	for seed := range int64(200) {
		rng := rand.New(rand.NewSource(seed))
		d := NewDeque[int]()
		var model []int
		next := 0
		for range 300 {
			switch rng.Intn(8) {
			case 0, 1:
				d.PushBack(next)
				model = append(model, next)
				next++
			case 2, 3:
				d.PushFront(next)
				model = append([]int{next}, model...)
				next++
			case 4:
				v, ok := d.PopFront()
				if len(model) == 0 {
					if ok {
						t.Fatalf("seed %d: popfront ok on empty", seed)
					}
				} else {
					if !ok || v != model[0] {
						t.Fatalf("seed %d: popfront %d %v want %d", seed, v, ok, model[0])
					}
					model = model[1:]
				}
			case 5:
				v, ok := d.PopBack()
				if len(model) == 0 {
					if ok {
						t.Fatalf("seed %d: popback ok on empty", seed)
					}
				} else {
					last := model[len(model)-1]
					if !ok || v != last {
						t.Fatalf("seed %d: popback %d %v want %d", seed, v, ok, last)
					}
					model = model[:len(model)-1]
				}
			case 6:
				d.Clear()
				model = nil
			case 7:
				n := rng.Intn(4)
				i := 0
				for v := range d.Drain() {
					if v != model[0] {
						t.Fatalf("seed %d: drain %d want %d", seed, v, model[0])
					}
					model = model[1:]
					i++
					if i == n {
						break
					}
				}
			}
			if d.Len() != len(model) {
				t.Fatalf("seed %d: len %d want %d", seed, d.Len(), len(model))
			}
			got := slices.Collect(d.Iter())
			if !slices.Equal(got, model) {
				t.Fatalf("seed %d: iter %v want %v", seed, got, model)
			}
		}
	}
}

func TestModelPriorityQueue(t *testing.T) {
	for seed := range int64(200) {
		rng := rand.New(rand.NewSource(seed))
		q := NewPriorityQueue(minInt)
		var model []int
		for range 200 {
			switch rng.Intn(6) {
			case 0, 1, 2:
				v := rng.Intn(50)
				q.Push(v)
				model = append(model, v)
				slices.Sort(model)
			case 3, 4:
				v, ok := q.Pop()
				if len(model) == 0 {
					if ok {
						t.Fatalf("seed %d: pop ok on empty", seed)
					}
				} else {
					if !ok || v != model[0] {
						t.Fatalf("seed %d: pop %d %v want %d", seed, v, ok, model[0])
					}
					model = model[1:]
				}
			case 5:
				n := rng.Intn(4)
				i := 0
				for v := range q.Drain() {
					if v != model[0] {
						t.Fatalf("seed %d: drain %d want %d", seed, v, model[0])
					}
					model = model[1:]
					i++
					if i == n {
						break
					}
				}
			}
			if q.Len() != len(model) {
				t.Fatalf("seed %d: len %d want %d", seed, q.Len(), len(model))
			}
			got := slices.Collect(q.Iter())
			slices.Sort(got)
			if !slices.Equal(got, model) {
				t.Fatalf("seed %d: iter multiset %v want %v", seed, got, model)
			}
		}
	}
}

func TestModelBounded(t *testing.T) {
	for name, newQueue := range boundedFactories {
		for _, capacity := range []int{0, 1, 2, 3, 5} {
			t.Run(fmt.Sprintf("%s/%d", name, capacity), func(t *testing.T) {
				for seed := range int64(100) {
					rng := rand.New(rand.NewSource(seed))
					q := newQueue(capacity)
					var model []int
					next := 0
					for range 200 {
						switch rng.Intn(8) {
						case 0, 1, 2, 3:
							ok := q.Push(next)
							if ok != (len(model) < capacity) {
								t.Fatalf("cap %d seed %d: enqueue ok=%v at len %d", capacity, seed, ok, len(model))
							}
							if ok {
								model = append(model, next)
								next++
							}
						case 4, 5:
							v, ok := q.Pop()
							if len(model) == 0 {
								if ok {
									t.Fatalf("cap %d seed %d: dequeue ok on empty", capacity, seed)
								}
							} else {
								if !ok || v != model[0] {
									t.Fatalf("cap %d seed %d: dequeue %d %v want %d", capacity, seed, v, ok, model[0])
								}
								model = model[1:]
							}
						case 6:
							q.Clear()
							model = nil
						case 7:
							n := rng.Intn(3)
							i := 0
							for v := range q.Drain() {
								if v != model[0] {
									t.Fatalf("cap %d seed %d: drain %d want %d", capacity, seed, v, model[0])
								}
								model = model[1:]
								i++
								if i == n {
									break
								}
							}
						}
						if q.Len() != len(model) {
							t.Fatalf("cap %d seed %d: len %d want %d", capacity, seed, q.Len(), len(model))
						}
						if q.Capacity() != capacity {
							t.Fatalf("cap %d seed %d: capacity drifted to %d", capacity, seed, q.Capacity())
						}
						if q.Full() != (len(model) == capacity) {
							t.Fatalf("cap %d seed %d: full=%v at len %d", capacity, seed, q.Full(), len(model))
						}
						if !slices.Equal(q.ToSlice(), model) {
							t.Fatalf("cap %d seed %d: toslice %v want %v", capacity, seed, q.ToSlice(), model)
						}
						if got := slices.Collect(q.Iter()); !slices.Equal(got, model) {
							t.Fatalf("cap %d seed %d: iter %v want %v", capacity, seed, got, model)
						}
					}
				}
			})
		}
	}
}
