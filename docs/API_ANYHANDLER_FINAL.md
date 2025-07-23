# DevTUI anyHandler - Diseño Final Consolidado

## Resumen Ejecutivo

Reemplazo de la interfaz `fieldHandler` por la estructura privada `anyHandler` para simplificar la implementación de handlers, eliminando 60-85% del código boilerplate manteniendo toda la funcionalidad existente.

## Interfaces Públicas Finales

```go
// Display - Solo lectura con contenido inmediato
type HandlerDisplay interface {
    Name() string    // Identificador para logging: "HelpDisplay", "StatusMonitor"
    Label() string   // Label mostrado en footer derecho
    Content() string // Contenido mostrado inmediatamente al posicionarse
}

// Edit - Campos editables
type HandlerEdit interface {
    Name() string    // Identificador para logging: "ServerPort", "DatabaseURL"
    Label() string   // Label mostrado en footer derecho  
    Value() string   // Valor actual/inicial
    Change(newValue any, progress ...func(string)) error // Sin return string - usar Value() si no hay error
}

// Execution - Botones de acción
type HandlerExecution interface {
    Name() string    // Identificador para logging: "DeployProd", "BuildProject"
    Label() string   // Label mostrado en footer derecho
    Execute(progress ...func(string)) error
}

// Writer básico - Nuevas líneas siempre
type HandlerWriter interface {
    Name() string // Identificador único: "AppLogger", "BuildOutput"
}

// Writer avanzado - Puede actualizar líneas existentes
type HandlerTrackerWriter interface {
    Name() string
    MessageTracker
}

// Tracking opcional para mensajes
type MessageTracker interface {
    GetLastOperationID() string
    SetLastOperationID(id string)
}

// Edit avanzado con tracking
type EditHandlerTracker interface {
    HandlerEdit
    MessageTracker
}

// Execution avanzado con tracking
type ExecutionHandlerTracker interface {
    HandlerExecution
    MessageTracker
}
```

## Estructura anyHandler (Privada)

```go
type handlerType int

const (
    handlerTypeDisplay handlerType = iota
    handlerTypeEdit
    handlerTypeExecution
    handlerTypeWriter
    handlerTypeTrackerWriter
)

// anyHandler - Estructura privada que unifica todos los handlers
type anyHandler struct {
    handlerType handlerType
    timeout     time.Duration // Solo edit/execution
    lastOpID    string       // Tracking interno
    
    // Function pointers - solo los necesarios poblados
    nameFunc      func() string                               // Todos
    labelFunc     func() string                               // Display/Edit/Execution
    valueFunc     func() string                               // Edit/Display
    contentFunc   func() string                               // Display únicamente
    editableFunc  func() bool                                 // Por tipo
    changeFunc    func(any, ...func(string)) error           // Edit/Execution (solo error)
    executeFunc   func(...func(string)) error                // Execution únicamente
    timeoutFunc   func() time.Duration                       // Edit/Execution
    getOpIDFunc   func() string                              // Tracking
    setOpIDFunc   func(string)                               // Tracking
}
```

## Factory Methods Reutilizables

```go
// Base factory con tracking opcional
func newEditHandler(h HandlerEdit, timeout time.Duration, tracker MessageTracker) *anyHandler {
    anyH := &anyHandler{
        handlerType:  handlerTypeEdit,
        timeout:      timeout,
        nameFunc:     h.Name,     // CORREGIDO: Name() no Label()
        labelFunc:    h.Label,    // CORREGIDO: Separado de Name()
        valueFunc:    h.Value,
        editableFunc: func() bool { return true },
        changeFunc:   h.Change,   // CORREGIDO: signature sin string return
        timeoutFunc:  func() time.Duration { return timeout },
    }
    
    // Configurar tracking opcional
    if tracker != nil {
        anyH.getOpIDFunc = tracker.GetLastOperationID
        anyH.setOpIDFunc = tracker.SetLastOperationID
    } else {
        anyH.getOpIDFunc = func() string { return "" }
        anyH.setOpIDFunc = func(string) {}
    }
    
    return anyH
}

func newDisplayHandler(h HandlerDisplay) *anyHandler {
    return &anyHandler{
        handlerType:  handlerTypeDisplay,
        timeout:      0, // Display no requiere timeout
        nameFunc:     h.Name,     // CORREGIDO: Name() para identificación
        labelFunc:    h.Label,    // CORREGIDO: Label() para footer
        valueFunc:    h.Content,  // Content como Value para compatibilidad
        contentFunc:  h.Content,  // Content específico para display
        editableFunc: func() bool { return false },
        getOpIDFunc:  func() string { return "" },
        setOpIDFunc:  func(string) {},
    }
}

func newExecutionHandler(h HandlerExecution, timeout time.Duration) *anyHandler {
    return &anyHandler{
        handlerType:  handlerTypeExecution,
        timeout:      timeout,
        nameFunc:     h.Name,     // CORREGIDO: Name() para identificación
        labelFunc:    h.Label,    // CORREGIDO: Label() para footer
        valueFunc:    h.Label,    // Label como Value
        editableFunc: func() bool { return false },
        executeFunc:  h.Execute,  // NUEVO: Execute específico
        changeFunc: func(val any, progress ...func(string)) error {
            return h.Execute(progress...) // Wrapper para compatibilidad
        },
        timeoutFunc:  func() time.Duration { return timeout },
        getOpIDFunc:  func() string { return "" },
        setOpIDFunc:  func(string) {},
    }
}

// Factory methods específicos
func newBasicEditHandler(h HandlerEdit, timeout time.Duration) *anyHandler {
    return newEditHandler(h, timeout, nil)
}

func newEditHandlerWithTracking(h EditHandlerTracker, timeout time.Duration) *anyHandler {
    return newEditHandler(h, timeout, h)
}

func newWriterHandler(h HandlerWriter) *anyHandler {
    return &anyHandler{
        handlerType: handlerTypeWriter,
        nameFunc:    h.Name,
        getOpIDFunc: func() string { return "" }, // Siempre nuevas líneas
        setOpIDFunc: func(string) {},
    }
}

func newTrackerWriterHandler(h HandlerTrackerWriter) *anyHandler {
    return &anyHandler{
        handlerType: handlerTypeTrackerWriter,
        nameFunc:    h.Name,
        getOpIDFunc: h.GetLastOperationID,
        setOpIDFunc: h.SetLastOperationID,
    }
}
```

## Métodos anyHandler

```go
// Métodos públicos de anyHandler (reemplaza completamente fieldHandler)
func (a *anyHandler) Name() string {
    if a.nameFunc != nil {
        return a.nameFunc()
    }
    return ""
}

func (a *anyHandler) Label() string {
    if a.labelFunc != nil {
        return a.labelFunc()
    }
    return ""
}

func (a *anyHandler) Value() string {
    if a.valueFunc != nil {
        return a.valueFunc()
    }
    return ""
}

func (a *anyHandler) Editable() bool {
    if a.editableFunc != nil {
        return a.editableFunc()
    }
    return false
}

// SIMPLIFICADO: anyHandler.Change() retorna solo error como las interfaces públicas
func (a *anyHandler) Change(newValue any, progress ...func(string)) error {
    if a.changeFunc != nil {
        return a.changeFunc(newValue, progress...)
    }
    return nil
}

func (a *anyHandler) Timeout() time.Duration {
    if a.timeoutFunc != nil {
        return a.timeoutFunc()
    }
    return a.timeout
}

func (a *anyHandler) SetLastOperationID(id string) {
    a.lastOpID = id
    if a.setOpIDFunc != nil {
        a.setOpIDFunc(id)
    }
}

func (a *anyHandler) GetLastOperationID() string {
    if a.getOpIDFunc != nil {
        return a.getOpIDFunc()
    }
    return a.lastOpID
}
```

## Thread-Safety: Slice Implementation

```go
// DECIDIDO: Usar slice para thread-safety
type tabSection struct {
    // ... otros campos
    writingHandlers []*anyHandler // CAMBIO: slice en lugar de map
    mu              sync.RWMutex  // Protección para operaciones concurrentes
}

func (ts *tabSection) RegisterHandlerWriter(handler HandlerWriter) io.Writer {
    ts.mu.Lock()
    defer ts.mu.Unlock()
    
    anyH := newWriterHandler(handler)
    ts.writingHandlers = append(ts.writingHandlers, anyH)
    return &handlerWriter{tabSection: ts, handlerName: anyH.Name()}
}

func (ts *tabSection) getWritingHandler(name string) *anyHandler {
    ts.mu.RLock()
    defer ts.mu.RUnlock()
    
    for _, h := range ts.writingHandlers {
        if h.Name() == name {
            return h
        }
    }
    return nil
}
```

## Comportamiento Display Handlers

**NUEVO REQUISITO**: Los handlers de tipo `Display` deben:

1. **Mostrar contenido inmediatamente** al posicionarse en el campo
2. **En la sección principal**: Mostrar lo que retorne `Content()`
3. **En el footer**: Layout expandido `[Label expandido ________] [Scroll%]`
   - El `Label()` usa todo el espacio disponible restante
   - Scroll % siempre a la derecha
   - Estilo visual normal (mismo que header)

## Comportamiento Execution Handlers

**NUEVO REQUISITO**: Los handlers de tipo `Execution` deben:

1. **En el footer**: Layout normal `[Label] [Value] [Scroll%]` (igual que Edit)
   - Label con ancho fijo (labelWidth) 
   - Value con ancho calculado
   - Scroll % siempre a la derecha
   - **Estilo visual del botón**: Fondo blanco con letras grises (indica que es ejecutable)

## Comportamiento Edit Handlers

**COMPORTAMIENTO ESTÁNDAR**: Los handlers de tipo `Edit` mantienen:

1. **En el footer**: Layout normal `[Label] [Value] [Scroll%]`
   - Label con ancho fijo (labelWidth)
   - Value con ancho calculado
   - Scroll % siempre a la derecha
   - Estilos según estado (edición activa, inactiva, etc.)

```go
// Detección actualizada para display
func (f *field) isDisplayOnly() bool {
    if f.handler == nil {
        return false
    }
    ah, ok := f.handler.(*anyHandler)
    return ok && ah.handlerType == handlerTypeDisplay
}

// NUEVO: Detección para execution con footer expandido
func (f *field) isExecutionHandler() bool {
    if f.handler == nil {
        return false
    }
    ah, ok := f.handler.(*anyHandler)
    return ok && ah.handlerType == handlerTypeExecution
}

// NUEVO: Detección para handlers que usan footer expandido (Display + Execution)
func (f *field) usesExpandedFooter() bool {
    return f.isDisplayOnly() || f.isExecutionHandler()
}

// NUEVO: Método para mostrar contenido en la sección principal
func (f *field) getDisplayContent() string {
    if f.isDisplayOnly() && f.handler != nil {
        ah := f.handler.(*anyHandler)
        if ah.contentFunc != nil {
            return ah.contentFunc() // Content() se muestra en la sección principal
        }
    }
    return ""
}

// NUEVO: Método para footer expandido - Label() usa espacio de label + value
func (f *field) getExpandedFooterLabel() string {
    if f.usesExpandedFooter() && f.handler != nil {
        ah := f.handler.(*anyHandler)
        if ah.labelFunc != nil {
            return ah.labelFunc() // Label() ocupa espacio label+value, ScrollInfo queda intacto
        }
    }
    return ""
}
```

## Registro Type-Safe (Final)

```go
// Builders actualizados - API completa con tracking
func (ts *tabSection) NewEditHandler(handler HandlerEdit) *editHandlerBuilder {
    return &editHandlerBuilder{
        tabSection: ts,
        handler:    handler,
        timeout:    0,
    }
}

func (ts *tabSection) NewEditHandlerWithTracking(handler EditHandlerTracker) *editHandlerBuilder {
    return &editHandlerBuilder{
        tabSection: ts,
        handler:    handler, // EditHandlerTracker extends HandlerEdit
        timeout:    0,
    }
}

func (ts *tabSection) NewExecutionHandler(handler HandlerExecution) *executionHandlerBuilder {
    return &executionHandlerBuilder{
        tabSection: ts,
        handler:    handler,
        timeout:    0,
    }
}

func (ts *tabSection) NewExecutionHandlerTracking(handler ExecutionHandlerTracker) *executionHandlerBuilder {
    return &executionHandlerBuilder{
        tabSection: ts,
        handler:    handler,
        timeout:    0,
    }
}

func (ts *tabSection) NewDisplayHandler(handler HandlerDisplay) *displayHandlerBuilder {
    return &displayHandlerBuilder{
        tabSection: ts,
        handler:    handler,
    }
}

func (ts *tabSection) NewWriterHandler(handler any) *writerHandlerBuilder {
    return &writerHandlerBuilder{
        tabSection: ts,
        handler:    handler,
    }
}

func (ts *tabSection) NewWriterHandlerTracking(handler HandlerTrackerWriter) *writerHandlerBuilder {
    return &writerHandlerBuilder{
        tabSection: ts,
        handler:    handler,
    }
}

// Métodos de registro directo (con auto-detección de tracking)
func (ts *tabSection) RegisterHandlerWriter(handler HandlerWriter) io.Writer {
    // Auto-detecta HandlerTrackerWriter y configura tracking automáticamente
}

// DEPRECATED: Usar RegisterHandlerWriter - detecta tracking automáticamente
func (ts *tabSection) RegisterHandlerTrackerWriter(handler HandlerTrackerWriter) io.Writer {
    return ts.RegisterHandlerWriter(handler)
}
```

## Diferencias Clave Name() vs Label()

| Método | Uso | Ejemplo | Dónde se muestra |
|--------|-----|---------|------------------|
| `Name()` | Identificación única para logging | "ServerPort", "DeployProd" | Logs: `10:55:42 [ServerPort] Port updated` |
| `Label()` | Texto mostrado en UI | "Server Port", "Deploy to Production" | Footer derecho únicamente |

## Ejemplos de Uso Final

```go
// 1. Handler Edit simple (3 métodos)
type PortHandler struct{ port string }
func (h *PortHandler) Name() string  { return "ServerPort" }
func (h *PortHandler) Label() string { return "Server Port" }
func (h *PortHandler) Value() string { return h.port }
func (h *PortHandler) Change(newValue any, progress ...func(string)) error {
    h.port = newValue.(string)
    if len(progress) > 0 {
        progress[0]("Port updated to: " + h.port)
    }
    return nil
}

// 2. Handler Display (3 métodos)
type HelpHandler struct{}
func (h *HelpHandler) Name() string    { return "HelpDisplay" }
func (h *HelpHandler) Label() string   { return "Help" }
func (h *HelpHandler) Content() string { return "Tab: Switch • Enter: Edit • Esc: Cancel" }

// 3. Handler Execution (3 métodos)
type DeployHandler struct{}
func (h *DeployHandler) Name() string  { return "DeployProd" }
func (h *DeployHandler) Label() string { return "Deploy to Production" }
func (h *DeployHandler) Execute(progress ...func(string)) error {
    if len(progress) > 0 {
        progress[0]("Starting deployment...")
        // ... lógica
        progress[0]("Deployment successful")
    }
    return nil
}

// Uso
tab := tui.NewTabSection("Server", "Configuration")
tab.NewEditHandler(&PortHandler{}).WithTimeout(5*time.Second)
tab.NewDisplayHandler(&HelpHandler{})  // Contenido inmediato
tab.NewExecutionHandler(&DeployHandler{}).WithTimeout(30*time.Second)
```

## Plan de Implementación

### Phase 1: Core Structure
1. **ELIMINAR** interfaz `fieldHandler` completamente 
2. Crear `anyHandler` struct en `field.go`
3. Implementar factory methods
4. Actualizar `field.handler` type de `fieldHandler` a `*anyHandler`
5. Actualizar todos los usos de `fieldHandler` a `*anyHandler`

### Phase 2: Thread-Safe Registry  
1. Cambiar `writingHandlers` de `map[string]writingHandler` a `[]*anyHandler`
2. Implementar métodos thread-safe con mutex
3. Actualizar métodos de registro

### Phase 3: Display Behavior
1. Implementar detección automática de contenido display
2. Actualizar lógica de footer para display handlers
3. Asegurar contenido inmediato sin interacción adicional

### Phase 4: Testing Migration
1. Actualizar todos los tests para nueva API
2. Eliminar tests de código obsoleto
3. Validar comportamiento de display handlers

## Análisis de Concurrencia

**Problemática actual**: `map[string]writingHandler` no es thread-safe.

**Solución elegida**: Slice con mutex RWLock
- **Writes** (registro): Lock exclusivo
- **Reads** (búsqueda): Lock compartido  
- **Simplicidad**: Append natural, búsqueda O(n) aceptable para pocos handlers
- **Thread-safety**: Automática con mutex apropiado

## Instrucciones Finales

**IMPORTANTE**: Al completar la implementación:
- ❌ **NO** crear ejemplos de uso
- ❌ **NO** actualizar README.md
- ❌ **NO** entregar informe detallado
- ✅ **SÍ** confirmar únicamente que se terminó

---

*Documento consolidado con decisiones finales aprobadas para implementación*
