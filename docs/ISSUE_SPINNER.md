# ISSUE: Spinner y Feedback Async en DevTUI

## Problema Actual

Cuando se ejecuta una acción async en DevTUI:
1. **No hay feedback visual inmediato** - la interfaz parece congelada
2. **El spinner no se muestra** durante la operación
3. **Los mensajes de progreso son hardcodeados** en lugar de ser responsabilidad del handler
4. **No hay indicación de tiempo transcurrido** o progreso

## Análisis de la Situación Actual

### Ejemplo de Código Problemático
```go
f.sendProgressMessage("Operation started...")  // ❌ HARDCODEADO
```

### Flujo Actual
```
Usuario presiona Enter → executeAsyncChange() → sendProgressMessage("Operation started...") → Handler.Change()
```

## Casos de Uso Reales

### Ejemplos de Diferentes Handlers
1. **ChatHandler**: "Enviando mensaje...", "Conectando al servidor...", "Mensaje enviado"
2. **PortConfigHandler**: "Validando puerto...", "Verificando disponibilidad...", "Puerto configurado"
3. **DockerBuildHandler**: "Iniciando build...", "Descargando imagen base...", "Construyendo capa 1/5...", "Build completado"
4. **BrowserHandler**: "Reiniciando navegador...", "Cerrando procesos...", "Navegador reiniciado"

## Decisiones Tomadas ✅

### 1. Control de Mensajes de Progreso
**✅ DECIDIDO: Handler controla completamente los mensajes**
- DevTUI NO proporciona mensajes, solo el formato/estética 
- DevTUI maneja spinner/animación y tiempo transcurrido
- Handler proporciona el contenido específico del mensaje

### 2. Nivel de Complejidad
**✅ DECIDIDO: Implementación SIMPLE**
- Evitar channels complejos o contexts avanzados
- API mínima e intuitiva
- Solución transparente para el handler

### 3. Compatibilidad hacia Atrás
**✅ DECIDIDO: Mantener método Change() actual**
- Necesario opciones de API mínimas/intuitivas/opcionales
- Refactorización gradual sin romper código existente

### 4. Información en el Spinner
**✅ DECIDIDO: Información completa controlada por handler**
- Solo animación + mensaje del handler
- También tiempo transcurrido ("⟳ Conectando... (5s)")
- También porcentaje si está disponible ("⟳ Descargando... 45% (12s)")
- **IMPORTANTE**: Mensajes son responsabilidad del handler, NO de DevTUI

### 5. Comunicación de Progreso
**✅ DECIDIDO: Método/canal/callback transparente e intuitivo**
- Preferible estándar (pub/sub, channel, etc.)
- DevTUI debe poder consultar o recibir información de manera transparente

## Preguntas Pendientes por Resolver

### A. ¿Cómo refactorizar la firma de Change() manteniendo compatibilidad?

#### Opción A1: Interface Opcional (RECOMENDADA)
```go
// Interface opcional para handlers que quieren progreso
type ProgressAware interface {
    ChangeWithProgress(newValue any, progressFunc func(message string, percent ...float64)) (string, error)
}

// Implementación en DevTUI
func (f *field) executeAsyncChange(valueToSave any) {
    if progressHandler, ok := f.handler.(ProgressAware); ok {
        // Handler soporta progreso
        result, err := progressHandler.ChangeWithProgress(valueToSave, f.sendProgressCallback)
    } else {
        // Handler usa método tradicional
        result, err := f.handler.Change(valueToSave)
    }
}

// Ejemplo de handler con progreso
type TinyGoHandler struct {
    // ... campos existentes
}

func (h *TinyGoHandler) Change(newValue any) (string, error) {
    // Comportamiento tradicional para compatibilidad
    return "TinyGo compilation completed", nil
}

func (h *TinyGoHandler) ChangeWithProgress(newValue any, progress func(string, ...float64)) (string, error) {
    progress("Iniciando compilación TinyGo...")
    time.Sleep(1 * time.Second)
    
    progress("Verificando dependencias...", 25)
    time.Sleep(2 * time.Second)
    
    progress("Compilando módulos WASM...", 50)
    time.Sleep(3 * time.Second)
    
    progress("Optimizando binario...", 75)
    time.Sleep(1 * time.Second)
    
    progress("Finalizando...", 100)
    return "TinyGo compilation completed successfully", nil
}
```

#### Opción A2: Context con Channel
```go
// Extender método Change para recibir context opcional
type ContextAwareHandler interface {
    ChangeWithContext(ctx context.Context, newValue any) (string, error)
}

// Implementación
func (f *field) executeAsyncChange(valueToSave any) {
    if ctxHandler, ok := f.handler.(ContextAwareHandler); ok {
        progressChan := make(chan ProgressUpdate, 10)
        ctx := context.WithValue(context.Background(), "progress", progressChan)
        
        go f.listenProgress(progressChan)
        result, err := ctxHandler.ChangeWithContext(ctx, valueToSave)
    }
}

type ProgressUpdate struct {
    Message string
    Percent *float64
}
```

#### Opción A3: Callback en Constructor (MÁS SIMPLE)
```go
// Handler recibe callback al crearse
type TinyGoHandler struct {
    progressCallback func(string, ...float64)
}

func NewTinyGoHandler(wasmHandler *tinywasm.Handler) *TinyGoHandler {
    return &TinyGoHandler{
        // ... inicialización normal
    }
}

// DevTUI inyecta callback después de crear el handler
func (ts *tabSection) NewField(handler FieldHandler) *tabSection {
    f := &field{handler: handler, parentTab: ts, /* ... */}
    
    // Inyectar callback si el handler lo soporta
    if progressAware, ok := handler.(interface{ SetProgressCallback(func(string, ...float64)) }); ok {
        progressAware.SetProgressCallback(f.sendProgressCallback)
    }
    
    // ... resto igual
}
```

### B. ¿Cómo debe comunicarse el progreso de forma estándar?

#### Opción B1: Callback Function (SIMPLE Y ESTÁNDAR)
```go
// Firma estándar para progreso
type ProgressCallback func(message string, percent ...float64)

// Uso en handler
func (h *DockerHandler) ChangeWithProgress(value any, progress ProgressCallback) (string, error) {
    progress("Iniciando build Docker...")           // Sin porcentaje
    progress("Descargando imagen base...", 20.0)    // Con porcentaje  
    progress("Construyendo capa 1/3...", 40.0)     // Progreso específico
    progress("Docker build completed", 100.0)       // Completado
}
```

#### Opción B2: Channel con Struct (MÁS FLEXIBLE)
```go
type ProgressMessage struct {
    Text    string
    Percent *float64
    Type    string // "info", "warning", "error"
}

// Handler envía por channel
progressChan <- ProgressMessage{Text: "Compilando...", Percent: &percent}
```

#### Opción B3: Observer Pattern (ESTÁNDAR ENTERPRISE)
```go
type ProgressObserver interface {
    OnProgress(message string, percent ...float64)
    OnError(error)
    OnComplete(result string)
}

// Handler recibe observer
func (h *Handler) ChangeWithObserver(value any, observer ProgressObserver) (string, error) {
    observer.OnProgress("Starting...")
    // ... trabajo
    observer.OnComplete("Done")
}
```

## Casos de Uso Reales de Godev

Basado en `section-build.go`, los handlers reales necesitarán:

### TinyGo Compilation
```go
func (h *TinyGoHandler) ChangeWithProgress(value any, progress ProgressCallback) (string, error) {
    progress("Verificando TinyGo installation...")
    progress("Compilando archivo WASM...", 30)
    progress("Optimizando binario...", 70)
    progress("Generando JavaScript inicializador...", 90)
    return "WASM compilation completed", nil
}
```

### Server Handler  
```go
func (h *ServerHandler) ChangeWithProgress(value any, progress ProgressCallback) (string, error) {
    progress("Deteniendo servidor actual...")
    progress("Compilando servidor Go...", 40)
    progress("Iniciando nuevo servidor...", 80)
    progress(fmt.Sprintf("Servidor iniciado en puerto %s", port), 100)
    return "Server restarted successfully", nil
}
```

### Browser Handler
```go
func (h *BrowserHandler) ChangeWithProgress(value any, progress ProgressCallback) (string, error) {
    progress("Cerrando procesos del navegador...")
    progress("Reiniciando navegador...", 50)
    progress("Navegador abierto en http://localhost:8080", 100)
    return "Browser reloaded", nil
}
```

## Recomendación Final

### Solución Recomendada: **Opción A1 + B1** 

**Interface Opcional + Callback Function**

**Justificación:**
- ✅ **Simple**: Solo una interface opcional y callback function
- ✅ **Intuitivo**: `ChangeWithProgress(value, progressCallback)`  
- ✅ **Compatible**: Mantiene `Change()` existente
- ✅ **Estándar**: Callback es patrón conocido
- ✅ **Flexible**: Soporta mensaje + porcentaje opcional
- ✅ **Transparente**: DevTUI detecta automáticamente capacidad de progreso

**Implementación mínima:**
```go
type ProgressAware interface {
    ChangeWithProgress(newValue any, progress func(string, ...float64)) (string, error)
}
```

### Alternativa si Callback no convence: **Opción A3**

Si prefieres evitar interfaces opcionales, usar **callback en constructor** es más explícito.

## Análisis de Bubbletea Examples

### progress-download.md
```go
type progressWriter struct {
    onProgress func(float64)  // Callback para progreso
}

// Reporta progreso con porcentaje
pw.onProgress(float64(pw.downloaded) / float64(pw.total))
```

**Lecciones:**
- Usa callback function para reportar progreso
- Proporciona información cuantitativa (porcentaje)
- El componente que hace el trabajo controla el mensaje

## Propuestas de Solución

### Propuesta A: Interface Opcional con Callback
```go
// Nueva interface opcional
type ProgressReporter interface {
    ReportProgress(message string)
    ReportProgressWithPercent(message string, percent float64)
}

// Handler que quiere reportar progreso
type DockerBuildHandler struct {
    progressCallback func(string)
}

func (h *DockerBuildHandler) Change(newValue any) (string, error) {
    if h.progressCallback != nil {
        h.progressCallback("Iniciando build...")
        // trabajo...
        h.progressCallback("Construyendo capa 1/5...")
        // más trabajo...
    }
    return "Build completado", nil
}
```

### Propuesta B: Context con Channel
```go
func (f *field) executeAsyncChange(valueToSave any) {
    progressChan := make(chan string, 10)
    ctx := context.WithValue(context.Background(), "progressChan", progressChan)
    
    go func() {
        for msg := range progressChan {
            f.sendProgressMessage(msg)
        }
    }()
    
    result, err := f.handler.ChangeWithContext(ctx, valueToSave)
}
```

### Propuesta C: Callback Simple (RECOMENDADA)
```go
// Modificar la interface FieldHandler
type FieldHandler interface {
    // ... métodos existentes ...
    ChangeWithProgress(newValue any, progressCallback func(string)) (string, error)
}

// Handlers que no necesitan progreso pueden usar implementación vacía
func (h *SimpleHandler) ChangeWithProgress(newValue any, progressCallback func(string)) (string, error) {
    return h.Change(newValue) // Delegar al método simple
}

// Handlers que SÍ necesitan progreso
func (h *DockerBuildHandler) ChangeWithProgress(newValue any, progressCallback func(string)) (string, error) {
    progressCallback("Iniciando build de imagen Docker...")
    
    // Simular trabajo
    time.Sleep(1 * time.Second)
    progressCallback("Descargando imagen base...")
    
    time.Sleep(2 * time.Second)  
    progressCallback("Construyendo capa 1/3...")
    
    // ... más trabajo ...
    
    return "Imagen Docker construida exitosamente", nil
}
```

## Recomendación

### Solución Propuesta: **Propuesta C - Callback Simple**

**Ventajas:**
- ✅ Simple de implementar
- ✅ No rompe compatibilidad (método opcional)
- ✅ Handler controla completamente sus mensajes
- ✅ Flexible para diferentes tipos de progreso
- ✅ No requiere channels o contexts complejos

**Implementación:**
1. Agregar método opcional `ChangeWithProgress` a la interface
2. DevTUI detecta si el handler implementa el método extendido
3. Si lo implementa, usa callback; si no, comportamiento actual
4. Spinner se muestra automáticamente cuando hay operación async
5. Tiempo transcurrido se puede agregar opcionalmente

## 🎯 SOLUCIÓN DEFINITIVA: Modificar Change() con Parámetro Variádico

### Propuesta Elegante: Change con Progress Opcional

```go
// Cambiar la firma de Change() en FieldHandler
type FieldHandler interface {
    Label() string
    Value() string  
    Editable() bool
    Change(newValue any, progress ...func(string, ...float64)) (string, error) // ✅ SOLO ESTE CAMBIO
    Timeout() time.Duration
    WritingHandler
}
```

### Ventajas de esta Aproximación:

1. ✅ **NO necesitas nuevas interfaces**
2. ✅ **Compatibilidad hacia atrás completa** - handlers existentes siguen funcionando  
3. ✅ **Simple y elegante** - solo un parámetro opcional
4. ✅ **Intuitivo** - si quieres progreso, usas el callback; si no, lo ignoras

### Ejemplos de Uso:

#### Handler SIN progreso (comportamiento actual)
```go
func (h *SimpleHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
    // Ignora completamente el parámetro progress
    time.Sleep(2 * time.Second) 
    return "Operation completed", nil
}
```

#### Handler CON progreso  
```go
func (h *BuildActionHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
    if len(progress) > 0 {
        progressCallback := progress[0]
        
        progressCallback("Iniciando build...")
        time.Sleep(1 * time.Second)
        
        progressCallback("Descargando dependencias...", 25.0)
        time.Sleep(2 * time.Second)
        
        progressCallback("Compilando código...", 75.0)  
        time.Sleep(1 * time.Second)
        
        progressCallback("Finalizando...", 100.0)
    } else {
        // Fallback si no hay callback (no debería pasar)
        time.Sleep(4 * time.Second)
    }
    
    return "Build completed successfully", nil
}
```

### Implementación en DevTUI:

```go
// En field.go - executeAsyncChange()
func (f *field) executeAsyncChange(valueToSave any) {
    f.asyncState.isRunning = true
    f.asyncState.startTime = time.Now()
    
    // Crear callback para progreso
    progressCallback := func(message string, percent ...float64) {
        f.sendProgressMessage(message) // Envía mensaje con spinner automático
    }
    
    // Ejecutar handler con callback
    go func() {
        result, err := f.handler.Change(valueToSave, progressCallback)
        
        // ... resto de la lógica actual igual
        if err != nil {
            f.sendErrorMessage(err.Error())
        } else {
            f.sendSuccessMessage(result)  
        }
        f.asyncState.isRunning = false
    }()
}
```

### ¿Cómo afecta handlers existentes?

```go
// Handler existente en cmd/main.go - NO CAMBIAR NADA
func (h *HostConfigHandler) Change(newValue any) (string, error) {
    // ❌ Esta firma YA NO COMPILARÁ
}

// Necesario cambiar a:
func (h *HostConfigHandler) Change(newValue any, progress ...func(string, ...float64)) (string, error) {
    // ✅ Funciona igual que antes si ignoras progress
    host := strings.TrimSpace(newValue.(string))
    if host == "" {
        return "", fmt.Errorf("host cannot be empty")
    }
    time.Sleep(1 * time.Second)
    h.currentHost = host
    return fmt.Sprintf("Host configured: %s", host), nil
}
```

## ⚠️ **Consideración Importante:**

Esta solución **ROMPE** la compatibilidad hacia atrás porque cambia la firma del método `Change()`. Todos los handlers existentes necesitarán actualizar su firma.

## 🤔 **¿Es esto aceptable?**

**Pros:**
- ✅ Solución más elegante y simple
- ✅ No necesitas múltiples interfaces  
- ✅ API unificada y clara

**Contras:**
- ❌ Rompe compatibilidad - necesitas actualizar TODOS los handlers
- ❌ En [`cmd/main.go`](cmd/main.go ) hay ~6 handlers que actualizar
- ❌ En godev probablemente hay más handlers

## Alternativa Híbrida - Mantener Compatibilidad:

Si quieres evitar romper compatibilidad, podemos implementar ambas firmas temporalmente:

```go
// Detectar si handler soporta la nueva firma
type ProgressAwareHandler interface {
    ChangeWithProgress(newValue any, progress func(string, ...float64)) (string, error)
}

func (f *field) executeAsyncChange(valueToSave any) {
    if progressHandler, ok := f.handler.(ProgressAwareHandler); ok {
        // Nueva API con progreso
        result, err := progressHandler.ChangeWithProgress(valueToSave, f.createProgressCallback())
    } else {
        // API actual - mostrar mensaje genérico
        f.sendProgressMessage(f.handler.Label() + " en progreso...")
        result, err := f.handler.Change(valueToSave)
    }
    // ... resto igual
}
```

## Recomendación Final:

**¿Prefieres?**

1. **Cambiar Change()** directamente - Más elegante pero rompe compatibilidad
2. **Interface adicional** - Mantiene compatibilidad pero más complejo
3. **Híbrido** - Soporta ambas durante transición
