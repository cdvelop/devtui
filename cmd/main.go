package main

import (
	"fmt"
	"sync"

	"github.com/cdvelop/devtui"
)

func main() {

	// Inicializar la UI usando la configuración por defecto
	tui := devtui.DefaultTUIForTest(func(messages ...any) {
		fmt.Println(append([]any{"Error:"}, messages...)...)
	})

	// Usar un WaitGroup para esperar a que la UI termine
	var wg sync.WaitGroup
	wg.Add(1)

	// Iniciar la UI con el WaitGroup para control de sincronización
	go tui.Start(&wg)

	// Esperar hasta que la UI termine
	wg.Wait()

	fmt.Println("Aplicación finalizada correctamente")
}
