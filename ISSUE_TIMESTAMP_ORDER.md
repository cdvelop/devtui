# Issue: Timestamps aparecen en orden aleatorio en la TUI

## 🐛 Descripción del Problema

Los timestamps mostrados en la interfaz TUI aparecen en orden aleatorio/incorrecto, a pesar de que los mensajes se generan secuencialmente. Esto confunde al usuario ya que los mensajes no aparecen en orden cronológico.

## 🔍 Análisis del Problema

### Causa Raíz
El problema está en la conversión entre nanosegundos y segundos:

1. **`GetNewID()`** devuelve un timestamp en **nanosegundos** (como string)
2. **`UnixSecondsToTime()`** espera un timestamp en **segundos** 
3. **`formatMessage()`** pasa nanosegundos a una función que espera segundos

### Flujo Problemático
```go
// En devtui/print.go - newContent()
Id: h.id.GetNewID(), // ← Devuelve nanosegundos como string (ej: "1748914249659368000")

// En devtui/print.go - formatMessage()  
timeStr := t.id.UnixSecondsToTime(msg.Id) // ← Trata nanosegundos como segundos
```

### Resultado
- **ID 1**: `1748914249659368000` (nanosegundos) → `01:53:20` (incorrecto)
- **ID 2**: `1748914250160692200` (nanosegundos) → `10:36:40` (incorrecto)  
- **ID 3**: `1748914250661485400` (nanosegundos) → `15:50:00` (incorrecto)

Los timestamps deberían ser consecutivos (diferencia de ~500ms), pero aparecen con diferencias de horas.

## 🧪 Reproducción

### Test que Demuestra el Problema
El archivo `timestamp_order_test.go` contiene tres tests que demuestran el issue:

```bash
cd devtui
go test -v -run TestDevTUITimestampIssue
```

**Salida esperada:**
```
=== PROBLEMA DETECTADO ===
Los IDs están en orden cronológico (nanosegundos como string):
  1. ID: 1748914249659368000
  2. ID: 1748914250160692200  
  3. ID: 1748914250661485400

Pero los tiempos formateados aparecen aleatorios:
  1. 01:53:20 Mensaje 1
  2. 10:36:40 Mensaje 2
  3. 15:50:00 Mensaje 3
```

### Reproducción Manual
1. Ejecutar `godev` en cualquier proyecto
2. Observar los timestamps en la TUI
3. Verificar que no siguen orden cronológico secuencial

## 🔧 Soluciones Propuestas

### Opción 1: Convertir Nanosegundos a Segundos en formatMessage()
```go
func (t *DevTUI) formatMessage(msg tabContent) string {
    // Convertir nanosegundos a segundos antes de formatear
    nanoSeconds, err := strconv.ParseInt(msg.Id, 10, 64)
    if err != nil {
        // Handle error
    }
    seconds := nanoSeconds / 1e9
    timeStr := t.timeStyle.Render(t.id.UnixSecondsToTime(seconds))
    // ...resto del código
}
```

### Opción 2: Agregar método UnixNanoToTime() a unixid
```go
// En unixid
func (u UnixID) UnixNanoToTime(input any) string {
    var unixNano int64
    // ...conversión similar a UnixSecondsToTime
    unixSeconds := unixNano / 1e9
    t := time.Unix(unixSeconds, 0)
    return t.Format("15:04:05")
}

// En devtui/print.go
timeStr := t.timeStyle.Render(t.id.UnixNanoToTime(msg.Id))
```

### Opción 3: Cambiar GetNewID() para devolver segundos
```go
// Modificar GetNewID para que devuelva segundos en lugar de nanosegundos
func (id *UnixID) GetNewID() string {
    return strconv.FormatInt(time.Now().Unix(), 10)
}
```

## ✅ Plan de Implementación

### Paso 1: Agregar UnixNanoToTime() a unixid
- [ ] Crear método `UnixNanoToTime()` en `unixid_back.go`
- [ ] Agregar tests para el nuevo método
- [ ] Documentar el método

### Paso 2: Actualizar devtui
- [ ] Cambiar `formatMessage()` para usar `UnixNanoToTime()` 
- [ ] Validar que los timestamps aparezcan en orden
- [ ] Ejecutar tests existentes

### Paso 3: Testing
- [ ] Ejecutar `timestamp_order_test.go` para verificar la corrección
- [ ] Test manual con `godev` 
- [ ] Verificar que no hay regresiones

## 📋 Checklist de Testing

- [ ] `go test -v -run TestTimestampOrder` pasa
- [ ] `go test -v -run TestDevTUITimestampIssue` muestra timestamps ordenados
- [ ] Test manual con `godev` muestra tiempos cronológicos
- [ ] No hay regresiones en otros tests de `devtui`
- [ ] No hay regresiones en otros tests de `unixid`

## 🎯 Criterios de Aceptación

1. **Timestamps cronológicos**: Los mensajes aparecen en orden temporal correcto
2. **Diferencias realistas**: Diferencias de tiempo reflejan intervalos reales entre mensajes
3. **Sin regresiones**: Funcionalidad existente sigue funcionando
4. **Tests pasan**: Todos los tests nuevos y existentes pasan

## 🔗 Archivos Afectados

- `devtui/print.go` - Método `formatMessage()`
- `unixid/unixid_back.go` - Nuevo método `UnixNanoToTime()`
- `devtui/timestamp_order_test.go` - Tests de verificación

## 📸 Evidencia Visual

### Antes (Problema)
```
06:13:20 Info: TinyGo installation verified
04:40:00 Server Start ...
06:46:40 Watch path added...
```
*Timestamps aparecen en orden aleatorio*

### Después (Esperado)  
```
14:32:01 Info: TinyGo installation verified
14:32:02 Server Start ...
14:32:03 Watch path added...
```
*Timestamps en orden cronológico secuencial*
