package main

import (
	"fmt"
	"sync"

	"github.com/cdvelop/devtui"
)

func main() {

	// Inicializar la UI usando la nueva API encapsulada
	tui := devtui.NewTUI(&devtui.TuiConfig{
		AppName:       "Ejemplo DevTUI",
		TabIndexStart: 0,
		ExitChan:      make(chan bool),
		Color: &devtui.ColorStyle{
			Foreground: "#F4F4F4",
			Background: "#000000",
			Highlight:  "#FF6600",
			Lowlight:   "#666666",
		},
		LogToFile: func(messages ...any) {
			fmt.Println(append([]any{"Error:"}, messages...)...)
		},
	})

	// Configurar la sección y los campos usando la API encadenada
	tui.NewTabSection("Datos personales", "Información básica").
		NewField("Nombre", "", true, nil).
		NewField("Edad", "", true, nil).
		NewField("Email", "", true, nil)

	// Usar un WaitGroup para esperar a que la UI termine
	var wg sync.WaitGroup
	wg.Add(1)

	// Iniciar la UI con el WaitGroup para control de sincronización
	go tui.InitTUI(&wg)

	// Esperar hasta que la UI termine
	wg.Wait()

	fmt.Println("Aplicación finalizada correctamente")
}
