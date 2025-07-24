# FOOTER_VIEW: Diseño del Footer para HandlerDisplay después de Refactorización API

## Introducción

Este documento describe el comportamiento y diseño visual del footer para handlers del tipo `HandlerDisplay` después de la refactorización de la API que simplifica la interface a solo 2 métodos: `Name()` y `Content()`.

## Interface Refactorizada

```go
type HandlerDisplay interface {
    Name() string    // Texto completo para mostrar en footer (handler responsable del contenido)
    Content() string // Contenido mostrado automáticamente en la sección principal
}
```

## Comportamiento del Footer

### Layout Visual

El footer para `HandlerDisplay` utiliza un layout expandido que aprovecha todo el ancho disponible:

```
┌─ Content Area ───────────────────────────────────────────────────┐
│ Status: Running                                                  │
│ PID: 12345                                                       │
│ Uptime: 2h 30m                                                   │
│ Memory: 45MB                                                     │
│ CPU: 12%                                                         │
│ (contenido mostrado automáticamente via Content())              │
└──────────────────────────────────────────────────────────────────┘
[System Status Information Display                              ][100%]
 ↑                                                                 ↑
Name() completo                                               Scroll%
```

### Responsabilidades

#### DevTUI (Capa Visual)
- **Renderiza directamente** el texto proporcionado por `Name()`
- **NO modifica** ni agrega información al contenido del handler
- **Gestiona el layout** y espacio disponible
- **Aplica estilos** consistentes con el resto de la UI

#### Handler (Lógica de Negocio)  
- **Proporciona el texto completo** para el footer via `Name()`
- **Controla totalmente** qué se muestra
- **Formatea apropiadamente** el contenido
- **Decide el contexto** y información a mostrar

## Implementación Técnica

### Código del Footer

```go
// En footerInput.go
func (f *field) getExpandedFooterLabel() string {
    if f.isDisplayOnly() && f.handler != nil {
        if f.handler.nameFunc != nil {
            return f.handler.nameFunc() // Renderizar directamente Name()
        }
    }
    return ""
}
```

### Renderizado

```go
// En footerInput.go - renderFooterInput()
if field.isDisplayOnly() {
    // Layout for Display: [Name() expandido] [Scroll%]
    remainingWidth := h.viewport.Width - lipgloss.Width(info) - horizontalPadding
    labelText := tinystring.Convert(field.getExpandedFooterLabel()).Truncate(remainingWidth-1, 0).String()
    
    displayStyle := lipgloss.NewStyle().
        Width(remainingWidth).
        Padding(0, horizontalPadding)
    styledLabel := displayStyle.Render(labelText)
    
    spacerStyle := lipgloss.NewStyle().Width(horizontalPadding).Render("")
    return lipgloss.JoinHorizontal(lipgloss.Left, styledLabel, spacerStyle, info)
}
```

## Ejemplos de Uso

### Ejemplo 1: Status Handler
```go
type StatusHandler struct{}

func (h *StatusHandler) Name() string { 
    return "System Status Information Display"
}

func (h *StatusHandler) Content() string {
    return "Status: Running\nPID: 12345\nUptime: 2h 30m\nMemory: 45MB\nCPU: 12%"
}
```

**Footer resultante:**
```
[System Status Information Display                              ][100%]
```

### Ejemplo 2: Help Handler
```go
type HelpHandler struct{}

func (h *HelpHandler) Name() string { 
    return "DevTUI Help & Navigation Guide"
}

func (h *HelpHandler) Content() string {
    return "Tab/Shift+Tab: Switch tabs\nLeft/Right: Navigate fields\nEnter: Edit/Execute\nEsc: Cancel"
}
```

**Footer resultante:**
```
[DevTUI Help & Navigation Guide                                ][100%]
```

### Ejemplo 3: Server Info Handler
```go
type ServerInfoHandler struct{}

func (h *ServerInfoHandler) Name() string { 
    return "Production Server Monitoring Dashboard"
}

func (h *ServerInfoHandler) Content() string {
    return "Server: prod-web-01\nLoad: 0.75\nConnections: 124\nUptime: 15d 6h"
}
```

**Footer resultante:**
```
[Production Server Monitoring Dashboard                         ][100%]
```

## Características del Diseño

### Ventajas Visuales

1. **Espacio Optimizado**: Todo el ancho del footer se utiliza efectivamente
2. **Claridad**: El handler controla completamente el mensaje mostrado
3. **Flexibilidad**: Puede mostrar desde identificadores simples hasta descripciones completas
4. **Consistencia**: Layout uniforme para todos los handlers Display

### Comportamiento Responsivo

- **Truncado automático**: Texto largo se trunca con `...` si excede el ancho
- **Espacio mínimo**: Scroll% siempre visible en la esquina derecha
- **Padding consistente**: Espaciado uniforme con el resto de la UI

## Comparación con Otros Handlers

### HandlerEdit
```
[Database Co...]  [postgres://localhost:5432/mydb          ][100%]
 ↑                 ↑
Label()           Value()
```

### HandlerExecution  
```
[Create Syst...]  [Create System Backup                   ][100%]
 ↑                 ↑
Label()           Value()
```

### HandlerDisplay (Nuevo)
```
[System Status Information Display                          ][100%]
 ↑
Name() completo
```

## Estilos Visuales

### Estilo Base
- **Fuente**: Misma que el resto del footer
- **Padding**: Consistente con otros elementos
- **Alineación**: Izquierda para el texto, derecha para scroll%

### Estados Visuales
- **Normal**: Estilo estándar del footer
- **Sin interacción**: No responde a clicks o teclas
- **Solo navegación**: Permite navegación entre fields pero sin edición

## Flujo de Datos

```
Handler.Name() → anyHandler.nameFunc() → getExpandedFooterLabel() → renderFooterInput() → Footer Visual
     ↑                                                                                           ↓
Lógica de Negocio                                                                      Capa Visual
```

## Testing del Footer

### Casos de Prueba

```go
func TestDisplayHandlerFooter(t *testing.T) {
    handler := &StatusHandler{}
    
    // Verificar que Name() no esté vacío
    if handler.Name() == "" {
        t.Error("Name() no debe estar vacío")
    }
    
    // Verificar que el footer se renderiza correctamente
    // (test de integración con TUI)
}

func TestFooterTruncation(t *testing.T) {
    handler := &LongNameHandler{}
    // Verificar truncado cuando Name() es muy largo
}

func TestFooterLayout(t *testing.T) {
    // Verificar que scroll% siempre esté visible
    // Verificar que el espaciado sea correcto
}
```

## Conclusión

El diseño del footer para `HandlerDisplay` después de la refactorización API logra:

- **Simplicidad**: Una sola fuente de información (`Name()`)
- **Control total**: El handler decide qué mostrar
- **Separación clara**: DevTUI renderiza, Handler controla
- **Diseño limpio**: Sin redundancia ni espacio desperdiciado
- **Flexibilidad**: Adaptable a cualquier tipo de información

Este enfoque respeta el principio de que DevTUI es solo la capa visual, mientras que los handlers son responsables de la lógica y contenido de la aplicación.