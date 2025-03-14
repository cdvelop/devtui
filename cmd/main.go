package main

import (
	"fmt"
	"sync"

	"github.com/cdvelop/devtui"
)

func main() {

	// Inicializar la UI
	tui := devtui.DefaultTUIForTest(func(messageErr any) {
		fmt.Println("Error: ", messageErr)
	})

	// Usar un WaitGroup para esperar a que la UI termine
	var wg sync.WaitGroup
	wg.Add(1)

	// Iniciar la UI con el WaitGroup para control de sincronización
	go tui.InitTUI(&wg)

	// Esperar hasta que la UI termine
	wg.Wait()

	fmt.Println("Aplicación finalizada correctamente")
}
