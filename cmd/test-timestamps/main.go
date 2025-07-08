package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cdvelop/devtui"
	"github.com/cdvelop/unixid"
)

// TestApp simula una aplicación que genera mensajes secuenciales
// para verificar el orden de timestamps en la TUI
func main() {
	// Configuración del test
	totalMessages := 10
	intervalSeconds := 1 // Configurable: tiempo entre mensajes

	fmt.Printf("=== TEST TIMESTAMP ORDER ===\n")
	fmt.Printf("Generando %d mensajes con intervalos de %d segundo(s)\n\n", totalMessages, intervalSeconds)

	// Configurar devtui (simulado)
	uid, err := unixid.NewUnixID()
	if err != nil {
		log.Fatal("Error creating unixid:", err)
	}

	// Simular DevTUI
	tui, err := devtui.New("TEST-TIMESTAMPS")
	if err != nil {
		log.Fatal("Error creating devtui:", err)
	}

	// Generar mensajes secuenciales
	for i := 1; i <= totalMessages; i++ {
		message := fmt.Sprintf("Mensaje secuencial #%d - timestamp debería incrementar", i)

		// Enviar mensaje a TUI (esto debería mostrar timestamps ordenados)
		tui.Print(message)

		// También mostrar en consola para comparar
		id := uid.GetNewID()
		timeFormatted := uid.UnixSecondsToTime(id) // ← Aquí está el problema
		fmt.Printf("[CONSOLA] %s %s (ID: %s)\n", timeFormatted, message, id)

		// Esperar intervalo configurable
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
	}

	fmt.Printf("\n=== FIN DEL TEST ===\n")
	fmt.Printf("Si los timestamps en TUI no están en orden cronológico, hay un bug\n")

	// Mantener TUI abierta para observar
	time.Sleep(5 * time.Second)
}
