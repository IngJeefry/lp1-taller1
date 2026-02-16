package main

import (
	"fmt"
	"sync"
	"time"
)

// Objetivo: Implementar una versión del problema de los Filósofos Comensales.
// Hay 5 filósofos y 5 tenedores (recursos). Cada filósofo necesita 2 tenedores para comer.
// Estrategia segura: imponer un **orden global** al tomar los tenedores (primero el menor ID, luego el mayor)
// para evitar deadlock. También puedes limitar concurrencia (ej. mayordomo).
// completa la lógica de toma/soltado de tenedores y bucle de pensar/comer.

type tenedor struct{ mu sync.Mutex }

// Versión 1: Orden Global
func filosofo(id int, izq, der *tenedor, wg *sync.WaitGroup) {
	// desarrolla el código para el filósofo
	defer wg.Done()

	for i := 0; i < 3; i++ {
		pensar(id)

		if id < (id+1) %5 {
			// Orden: primero el de menor ID
            fmt.Printf("[filósofo %d] esperando tenedor %d...\n", id, id)
            izq.mu.Lock()
            fmt.Printf("[filósofo %d] tomó tenedor %d\n", id, id)
            
            fmt.Printf("[filósofo %d] esperando tenedor %d...\n", id, (id+1)%5)
            der.mu.Lock()
            fmt.Printf("[filósofo %d] tomó tenedor %d\n", id, (id+1)%5)
        } else {
            // Orden inverso: primero el de menor ID (que sería el derecho)
            fmt.Printf("[filósofo %d] esperando tenedor %d...\n", id, (id+1)%5)
            der.mu.Lock()
            fmt.Printf("[filósofo %d] tomó tenedor %d\n", id, (id+1)%5)
            
            fmt.Printf("[filósofo %d] esperando tenedor %d...\n", id, id)
            izq.mu.Lock()
            fmt.Printf("[filósofo %d] tomó tenedor %d\n", id, id)
        }
        
        comer(id)
        
        // Liberar tenedores (orden inverso al que los tomamos)
        der.mu.Unlock()
        izq.mu.Unlock()
        fmt.Printf("[filósofo %d] soltó tenedores\n", id)
	}
	fmt.Printf("[filósofo %d] satisfecho\n", id)
}

func pensar(id int) {
	fmt.Printf("[filósofo %d] pensando...\n", id)
	// simular tiempo de pensar
	time.Sleep(time.Duration(500+id*100) * time.Millisecond)
}

func comer(id int) {
	fmt.Printf("[filósofo %d] COMIENDO\n", id)
	// simular tiempo de pensar
	time.Sleep(time.Duration(300+id*50) * time.Millisecond)
}

func main() {
	const n = 5
	var wg sync.WaitGroup
	wg.Add(n)

	// crear tenedores
	forks := make([]*tenedor, n)
	for i := 0; i < n; i++ {
		// inicializar cada tenedor i
		forks[i] = &tenedor{}
	}

	// lanzar filósofos
	for i := 0; i < n; i++ {
		izq := forks[i]
		der := forks[(i+1)%n]
		// lanzar goroutine para el filósofo i
		go filosofo(i, izq, der, &wg)
	}

	wg.Wait()
	fmt.Println("Todos los filósofos han comido sin deadlock.")
}

// Versión 2: Con Mayordomo
func filosofoConMayordomo(id int, izq, der *tenedor, mayordomo chan struct{}, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for i := 0; i < 3; i++ {
        pensar(id)
        
        // Pedir permiso al mayordomo
        mayordomo <- struct{}{} // Si el canal está lleno, espera
        
        // Tomar tenedores (cualquier orden, porque el mayordomo limita)
        izq.mu.Lock()
        der.mu.Lock()
        
        comer(id)
        
        // Soltar tenedores
        der.mu.Unlock()
        izq.mu.Unlock()
        
        // Devolver permiso al mayordomo
        <-mayordomo
    }
    fmt.Printf("[filósofo %d] SATISFECHO\n", id)
}

func mainConMayordomo() {
    const n = 5
    var wg sync.WaitGroup
    wg.Add(n)
    
    // Mayordomo: permite máximo 4 filósofos simultáneamente
    mayordomo := make(chan struct{}, 4)

    forks := make([]*tenedor, n)
    for i := 0; i < n; i++ {
        forks[i] = &tenedor{}
    }

    for i := 0; i < n; i++ {
        izq := forks[i]
        der := forks[(i+1)%n]
        go filosofoConMayordomo(i, izq, der, mayordomo, &wg)
    }

    wg.Wait()
    fmt.Println("✅ Todos los filósofos han comido sin deadlock (con mayordomo).")
}