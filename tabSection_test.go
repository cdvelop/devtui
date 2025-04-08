package devtui

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cdvelop/unixid"
)

func TestTabSectionWriter(t *testing.T) {
	// Configuración mínima
	exitChan := make(chan bool)
	id, _ := unixid.NewUnixID(sync.Mutex{})

	tui := &DevTUI{
		TuiConfig: &TuiConfig{
			ExitChan: exitChan,
		},
		id:              id,
		tabContentsChan: make(chan tabContent, 1),
	}

	// Crear tab section de prueba
	tab := tui.NewTabSection("TEST", "")

	// Testear el Writer
	testMsg := "Mensaje de prueba"
	n, err := fmt.Fprintln(tab, testMsg)
	if err != nil {
		t.Fatalf("Error escribiendo en el Writer: %v", err)
	}
	if n != len(testMsg)+1 { // +1 por el newline
		t.Errorf("Bytes escritos incorrectos: esperado %d, obtenido %d", len(testMsg)+1, n)
	}

	// Verificar que el mensaje llegó al canal
	select {
	case msg := <-tui.tabContentsChan:
		if msg.Content != testMsg {
			t.Errorf("Contenido incorrecto: esperado '%s', obtenido '%s'", testMsg, msg.Content)
		}
		if msg.Type != 0 { // 0 es el tipo por defecto para mensajes normales
			t.Errorf("Tipo de mensaje incorrecto: esperado 0, obtenido %v", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout: el mensaje no llegó al canal")
	}
}
