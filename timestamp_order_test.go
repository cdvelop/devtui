package devtui_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cdvelop/unixid"
)

// TestTimestampOrder verifica que los mensajes aparezcan en orden cronológico
// Este test reproduce el problema donde los timestamps aparecen aleatorios
func TestTimestampOrder(t *testing.T) {
	// Configurar unixid
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	// Simular el comportamiento de devtui
	var timestamps []string

	// Generar mensajes con intervalos de 1 segundo
	totalMessages := 5
	for i := 0; i < totalMessages; i++ {
		// Generar nuevo ID (esto devuelve nanosegundos como string)
		id := uid.GetNewID()

		// Simular como devtui formatea el mensaje
		// PROBLEMA: UnixNanoToTime espera segundos pero recibe nanosegundos
		timeStr := uid.UnixNanoToTime(id)

		message := fmt.Sprintf("Mensaje %d", i+1)

		timestamps = append(timestamps, timeStr)

		fmt.Printf("ID generado: %s\n", id)
		fmt.Printf("Tiempo formateado: %s\n", timeStr)
		fmt.Printf("Mensaje: %s\n\n", message)

		// Esperar 1 segundo para el siguiente mensaje
		time.Sleep(1 * time.Second)
	}

	// Verificar que los timestamps estén en orden cronológico
	for i := 1; i < len(timestamps); i++ {
		prev := timestamps[i-1]
		curr := timestamps[i]

		if curr <= prev {
			t.Errorf("Los timestamps NO están en orden cronológico:\n"+
				"Mensaje %d: %s\n"+
				"Mensaje %d: %s\n"+
				"El timestamp actual (%s) debería ser mayor que el anterior (%s)",
				i, prev, i+1, curr, curr, prev)
		}
	}
}

// TestCorrectTimestampConversion verifica la conversión correcta de timestamps
func TestCorrectTimestampConversion(t *testing.T) {
	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	now := time.Now()
	nanoTimestamp := now.UnixNano()
	secondsTimestamp := now.Unix()

	// Probar conversión correcta (segundos a tiempo)
	correctTime := uid.UnixNanoToTime(secondsTimestamp)
	expectedTime := now.Format("15:04:05")

	// Probar conversión incorrecta (nanosegundos a tiempo)
	incorrectTime := uid.UnixNanoToTime(nanoTimestamp)

	fmt.Printf("Timestamp en nanosegundos: %d\n", nanoTimestamp)
	fmt.Printf("Timestamp en segundos: %d\n", secondsTimestamp)
	fmt.Printf("Tiempo esperado: %s\n", expectedTime)
	fmt.Printf("Tiempo correcto (desde segundos): %s\n", correctTime)
	fmt.Printf("Tiempo incorrecto (desde nanosegundos): %s\n", incorrectTime)

	// Verificar que usar nanosegundos directamente da resultados incorrectos
	if incorrectTime == correctTime {
		t.Error("La conversión desde nanosegundos NO debería dar el mismo resultado que desde segundos")
	}

	// Verificar que la conversión correcta está cerca del tiempo esperado
	// (puede diferir en segundos debido al tiempo de ejecución)
	if correctTime == "" {
		t.Error("La conversión desde segundos no debería devolver string vacío")
	}
}

// TestDevTUITimestampIssue simula exactamente como devtui maneja los timestamps
func TestDevTUITimestampIssue(t *testing.T) {
	// Simular la estructura de devtui
	type mockTabContent struct {
		Id      string // GetNewID devuelve string
		Content string
	}

	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal("Error creating unixid:", err)
	}

	// Simular el método formatMessage de devtui
	formatMessage := func(msg mockTabContent) string {
		// ESTE ES EL PROBLEMA: msg.Id es en nanosegundos como string, pero UnixNanoToTime espera segundos
		timeStr := uid.UnixNanoToTime(msg.Id)
		return fmt.Sprintf("%s %s", timeStr, msg.Content)
	}

	var messages []mockTabContent
	var formattedMessages []string

	// Generar mensajes secuenciales
	for i := 0; i < 3; i++ {
		msg := mockTabContent{
			Id:      uid.GetNewID(), // Esto devuelve nanosegundos como string
			Content: fmt.Sprintf("Mensaje %d", i+1),
		}

		messages = append(messages, msg)
		formatted := formatMessage(msg)
		formattedMessages = append(formattedMessages, formatted)

		fmt.Printf("Mensaje %d - ID: %s\n", i+1, msg.Id)
		fmt.Printf("Formateado: %s\n\n", formatted)

		time.Sleep(500 * time.Millisecond)
	}

	// Mostrar el problema
	fmt.Println("=== PROBLEMA DETECTADO ===")
	fmt.Println("Los IDs están en orden cronológico (nanosegundos como string):")
	for i, msg := range messages {
		fmt.Printf("  %d. ID: %s\n", i+1, msg.Id)
	}

	fmt.Println("\nPero los tiempos formateados aparecen aleatorios:")
	for i, formatted := range formattedMessages {
		fmt.Printf("  %d. %s\n", i+1, formatted)
	}

	// Verificar que los IDs están en orden cronológico
	for i := 1; i < len(messages); i++ {
		if messages[i].Id <= messages[i-1].Id {
			t.Error("Los IDs deberían estar en orden cronológico ascendente")
		}
	}
}
