package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Objetivo: Simular "futuros" en Go usando canales. Una función lanza trabajo asíncrono
// y retorna un canal de solo lectura con el resultado futuro.
// completa las funciones y experimenta con varios futuros a la vez.

func asyncCuadrado(x int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		// simular trabajo
		tiempo := time.Duration(100+rand.Intn(400)) * time.Millisecond
		fmt.Printf("[async] calculando cuadrado de %d (tardará %v)\n", x, tiempo)
		time.Sleep(tiempo)
		ch <- x * x
		fmt.Printf("[async] %d² = %d listo\n", x, x*x)
	}()
	return ch
}

func fanIn(canales ...<-chan int) <-chan int {
    salida := make(chan int)
    var wg sync.WaitGroup

    // Por cada canal de entrada, lanzamos una goroutine
    for i, ch := range canales {
        wg.Add(1)
        go func(id int, c <-chan int) {
            defer wg.Done()
            // Leemos todos los valores del canal hasta que se cierre
            for v := range c {
                fmt.Printf("[fan-in] reenviando valor %d del canal %d\n", v, id+1)
                salida <- v
            }
            fmt.Printf("[fan-in] canal %d cerrado\n", id+1)
        }(i, ch)
    }

    // Cerramos el canal de salida cuando todas las goroutines terminen
    go func() {
        wg.Wait()
        fmt.Println("[fan-in] todos los canales cerrados, cerrando salida")
        close(salida)
    }()

    return salida
}

func main() {
	// crea varios futuros y recolecta sus resultados: f1, f2, f3
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== FUTUROS EN GO (Opción 1: Secuencial) ===")

	f1 := asyncCuadrado(5)
	f2 := asyncCuadrado(8)
	f3 := asyncCuadrado(12)

	// Opción 1: esperar cada futuro secuencialmente
	fmt.Println("\n[main] Esperando resultados secuancialmente. . .")

	r1 := <-f1
	fmt.Printf("[main] Resultado 1: %d\n", r1)

	r2 := <-f2
	fmt.Printf("[main Resultado 2: %d\n]", r2)

	r3 := <-f3
	fmt.Printf("[main Resultado 3: %d\n]", r3)
	
	// Opción 2: fan-in (combinar múltiples canales)
	fmt.Println("=== FUTUROS EN GO (Opción 2: Combinar multiples canales) ===")
	// Pista: crea una función fanIn que recibe múltiples <-chan int y retorna un único <-chan int
	f4 := asyncCuadrado(7)
    f5 := asyncCuadrado(9)
    f6 := asyncCuadrado(15)
    f7 := asyncCuadrado(4)
    f8 := asyncCuadrado(11)

	combinado := fanIn(f4, f5, f6, f7, f8)
	// que emita todos los valores. Requiere goroutines y cerrar el canal de salida cuando todas terminen.
	fmt.Println("\n[main] Recolectando resultados en orden de llegada:")
    contador := 0
    for r := range combinado {
        contador++
        fmt.Printf("[main] Resultado %d: %d\n", contador, r)
    }
    
    fmt.Println("\n¡Todos los futuros procesados!")
}
