package main

import (
	"fmt"
	"sync"

	"github.com/cdvelop/devtui"
)

func main() {
	// Crear una configuración por defecto
	config := &devtui.TuiConfig{
		TabIndexStart: 0,               // Iniciar en la primera pestaña
		ExitChan:      make(chan bool), // Canal para señalar salida
		Color:         nil,             // Usar colores por defecto
		LogToFile: func(messageErr any) {
			fmt.Printf("Error log: %v\n", messageErr)
		},
	}

	// Inicializar la UI
	tui := devtui.NewTUI(config)

	// Usar un WaitGroup para esperar a que la UI termine
	var wg sync.WaitGroup
	wg.Add(1)

	// Iniciar la UI con el WaitGroup para control de sincronización
	go tui.InitTUI(&wg)

	// Esperar hasta que la UI termine
	wg.Wait()

	fmt.Println("Aplicación finalizada correctamente")
}
