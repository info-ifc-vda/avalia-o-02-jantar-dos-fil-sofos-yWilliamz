package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	numFilosofos        = 5
	tempoMaxComendo     = 50 * time.Millisecond
	tempoMaxPensando    = 30 * time.Millisecond
	tempoMaxEsperando   = 20 * time.Millisecond
	tempoMaxInatividade = 80 * time.Millisecond
)

type EstadoFilosofo int

const (
	Pensando EstadoFilosofo = iota
	Faminto
	Comendo
)

type Garfo struct {
	sync.Mutex
	utilizadoPor int
}

type Filosofo struct {
	id               int
	garfoEsquerdo    *Garfo
	garfoDireito     *Garfo
	estado           EstadoFilosofo
	ultimaRefeicao   time.Time
	quantRefeicoes   int
	prioridade       bool
	temGarfoEsquerdo bool
	temGarfoDireito  bool
}

var (
	filosofos []*Filosofo
	mesa      sync.Mutex
)

func (f *Filosofo) rotinaFilosofo(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		f.ficarPensando()
		f.ficarFaminto()

		tempoInicioEspera := time.Now()

		for !f.tentarComer() {
			if time.Since(tempoInicioEspera) > tempoMaxEsperando {
				if !f.prioridade {
					f.soltarGarfos()
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
		}

		if f.temGarfoEsquerdo && f.temGarfoDireito {
			f.comer()
		}
	}
}

func (f *Filosofo) ficarPensando() {
	f.estado = Pensando
	imprimirStatus()
	time.Sleep(time.Duration(rand.Intn(int(tempoMaxPensando))))
}

func (f *Filosofo) ficarFaminto() {
	f.estado = Faminto
	imprimirStatus()
}

func (f *Filosofo) tentarComer() bool {
	if f.estado != Faminto {
		return false
	}

	if !f.temGarfoEsquerdo && f.garfoEsquerdo.utilizadoPor == -1 {
		f.garfoEsquerdo.Lock()
		f.garfoEsquerdo.utilizadoPor = f.id
		f.temGarfoEsquerdo = true

		return f.temGarfoEsquerdo && f.temGarfoDireito
	}

	if !f.temGarfoDireito && f.garfoDireito.utilizadoPor == -1 {
		f.garfoDireito.Lock()
		f.garfoDireito.utilizadoPor = f.id
		f.temGarfoDireito = true

		return f.temGarfoEsquerdo && f.temGarfoDireito
	}

	return f.temGarfoEsquerdo && f.temGarfoDireito
}

func (f *Filosofo) comer() {
	f.estado = Comendo

	imprimirStatus()

	tempoComendo := time.Duration(rand.Intn(int(tempoMaxComendo)))
	time.Sleep(tempoComendo)

	f.quantRefeicoes++
	f.ultimaRefeicao = time.Now()
	f.prioridade = false

	imprimirStatus()

	f.soltarGarfos()
	f.ficarPensando()
}

func (f *Filosofo) soltarGarfos() {
	if f.temGarfoEsquerdo {
		f.garfoEsquerdo.Unlock()
		f.garfoEsquerdo.utilizadoPor = -1
		f.temGarfoEsquerdo = false
	}

	if f.temGarfoDireito {
		f.garfoDireito.Unlock()
		f.garfoDireito.utilizadoPor = -1
		f.temGarfoDireito = false
	}
}

func imprimirStatus() {
	mesa.Lock()
	defer mesa.Unlock()

	estados := []string{"P", "F", "C"}
	statusFilosofos := make([]string, numFilosofos)
	garfos := make([]string, numFilosofos)
	prioridades := make([]string, numFilosofos)
	refeicoes := make([]string, numFilosofos)

	for i, f := range filosofos {
		statusFilosofos[i] = estados[f.estado]

		if f.temGarfoEsquerdo {
			garfos[i] = fmt.Sprintf("%d", f.id+1)
		} else if f.garfoEsquerdo.utilizadoPor != -1 {
			garfos[i] = fmt.Sprintf("%d", f.garfoEsquerdo.utilizadoPor+1)
		} else {
			garfos[i] = "-"
		}

		if f.prioridade {
			prioridades[i] = "!"
		} else {
			prioridades[i] = "-"
		}

		if f.quantRefeicoes < 10 {
			refeicoes[i] = fmt.Sprintf("%d ", f.quantRefeicoes)
		} else {
			refeicoes[i] = fmt.Sprintf("%d", f.quantRefeicoes)
		}
	}

	fmt.Printf("Filosofos:     %s\n", strings.Join(statusFilosofos, "  | "))
	fmt.Printf("Garfos:        %s\n", strings.Join(garfos, "  | "))
	fmt.Printf("Prioridades:   %s\n", strings.Join(prioridades, "  | "))
	fmt.Printf("Refeicoes:     %s\n", strings.Join(refeicoes, " | "))
	fmt.Print("\n\n")
}

func monitorarInatividade() {
	for {
		time.Sleep(100 * time.Millisecond)

		contadorPrioridade := 0
		for _, f := range filosofos {
			if f.prioridade {
				contadorPrioridade++
			}
		}

		if contadorPrioridade >= 4 {
			continue
		}

		for _, f := range filosofos {
			if !f.prioridade && time.Since(f.ultimaRefeicao) > tempoMaxInatividade && f.estado != Comendo {
				f.prioridade = true
				f.estado = Faminto
				contadorPrioridade++
				
				if contadorPrioridade >= 4 {
					break
				}
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	garfos := make([]*Garfo, numFilosofos)
	filosofos = make([]*Filosofo, numFilosofos)

	for i := 0; i < numFilosofos; i++ {
		garfos[i] = &Garfo{utilizadoPor: -1}
	}

	for i := 0; i < numFilosofos; i++ {
		filosofos[i] = &Filosofo{
			id:             i,
			garfoEsquerdo:  garfos[i],
			garfoDireito:   garfos[(i+1)%numFilosofos],
			ultimaRefeicao: time.Now(),
		}

		wg.Add(1)

		go filosofos[i].rotinaFilosofo(&wg)
	}

	go monitorarInatividade()

	wg.Wait()
}
