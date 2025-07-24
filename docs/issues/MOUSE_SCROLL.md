# DevTUI Mouse Scroll Investigation: TUI Best Practices Analysis

## Current Status
DevTUI currently utilizes the standard bubbles/viewport component which includes both mouse wheel scrolling and keyboard navigation. This investigation examines whether disabling mouse scroll in favor of standard keyboard navigation aligns with TUI best practices.

## Problem Statement
The current implementation combines mouse and keyboard navigation, but this creates several UX issues:

1. **Text Selection Conflict**: Mouse wheel scrolling prevents terminal text selection, which is essential for copying error messages and debugging output
2. **Inconsistent Navigation**: Mixed input methods can be confusing for users familiar with traditional TUI conventions
3. **Terminal Compatibility**: Not all terminal environments handle mouse events consistently

## TUI Standards Analysis

### 1. Historical TUI Conventions
Based on analysis of established TUI applications:

- **less/more pagers**: Use Page Up/Page Down, no mouse wheel
- **vi/vim**: Arrow keys for line navigation, Ctrl+U/Ctrl+D for half-page, Page Up/Page Down for full page
- **man pages**: Standard pager navigation (b/f for pages, u/d for half-pages)
- **tmux/screen**: Page Up/Page Down for scrollback, no mouse wheel by default
- **pine/mutt email clients**: Arrow keys and page navigation

### 2. Modern TUI Standards
Contemporary TUI libraries demonstrate these patterns:

**Viewport Navigation Standards** (from bubbles/viewport):
```go
// Standard keyboard navigation
PageDown:     "pgdown", spacebar, "f"    // Full page down
PageUp:       "pgup", "b"                // Full page up  
HalfPageDown: "d", "ctrl+d"             // Half page down
HalfPageUp:   "u", "ctrl+u"             // Half page up
Down:         "down", "j"                // Line down
Up:           "up", "k"                  // Line up
```

**Terminal Applications Pattern**:
- List navigation: j/k (vim-style) or arrow keys
- Page navigation: Page Up/Page Down, Space/Backspace
- Half-page: Ctrl+U/Ctrl+D or u/d
- Home/End: g/G or Home/End keys

## Current DevTUI Implementation Analysis

### Mouse Scroll Implementation
From bubbles/viewport source analysis:
```go
// Current viewport mouse handling
case tea.MouseButtonWheelUp:
    lines := m.ScrollUp(m.MouseWheelDelta)
case tea.MouseButtonWheelDown:
    lines := m.ScrollDown(m.MouseWheelDelta)
```

### Keyboard Navigation
DevTUI already implements standard keyboard navigation:
```go
// Current DevTUI keyboard navigation  
case tea.KeyUp, tea.KeyDown:
    // Arrow keys control viewport scroll (when not in edit mode)
case tea.KeyLeft, tea.KeyRight: 
    // Field navigation in normal mode
case tea.KeyTab, tea.KeyShiftTab:
    // Tab section navigation
```

## Recommended TUI Best Practices

### 1. Primary Navigation: Keyboard Only
**Advantages**:
- **Predictable**: Works consistently across all terminal environments
- **Accessible**: Keyboard-only navigation is more accessible
- **Text Selection**: Enables terminal text selection and copying
- **Faster**: Power users prefer keyboard navigation
- **Standard**: Follows established TUI conventions

### 2. Final Navigation Scheme (APPROVED)
```
Current Navigation    | Final Standard Navigation
--------------------|---------------------------
Mouse wheel         | Disabled
Up/Down arrows      | Line-by-line scroll
Page Up/Page Down   | Full page scroll (ONLY PAGE OPTION)
Left/Right arrows   | Field navigation (current)
Tab/Shift+Tab       | Tab section navigation (current)
```

**Decision**: Use **Page Up/Page Down ONLY** for page navigation to maintain simplicity and universal compatibility.

### 3. Implementation Strategy

#### Phase 1: Disable Mouse Wheel
```go
// In viewport initialization
h.viewport.MouseWheelEnabled = false
```

#### Phase 2: Simplified Keyboard Navigation
```go
func (h *DevTUI) handleNormalModeKeyboard(msg tea.KeyMsg) (bool, tea.Cmd) {
    switch msg.Type {
    case tea.KeyUp:
        // Line up in viewport
        h.viewport.ScrollUp(1)
        return false, nil
    case tea.KeyDown:
        // Line down in viewport  
        h.viewport.ScrollDown(1)
        return false, nil
    case tea.KeyPgUp:
        // Full page up (ONLY page navigation option)
        h.viewport.PageUp()
        return false, nil
    case tea.KeyPgDown:
        // Full page down (ONLY page navigation option)
        h.viewport.PageDown()
        return false, nil
    // ... existing field navigation
    }
}
```

## Benefits of Keyboard-Only Navigation

### 1. Enhanced User Experience
- **Text Selection**: Users can select and copy terminal output for debugging
- **Consistency**: Follows standard TUI conventions users expect
- **Efficiency**: Keyboard navigation is faster for power users
- **Reliability**: Works in all terminal environments (SSH, tmux, screen)

### 2. Development Benefits
- **Simpler Code**: Less complex input handling
- **Better Testing**: More predictable behavior for automated tests
- **Accessibility**: Better support for screen readers and accessibility tools
- **Cross-Platform**: Consistent behavior across different terminals

### 3. DevTUI Context Benefits
- **Better Debugging**: Users can copy error messages and logs
- **Documentation**: Easy to copy commands and configuration values
- **Standard Feel**: Fits better with other development tools
- **Professional UX**: More polished, professional feel

## Comparison with Other TUI Libraries

### tview (Rich TUI Library)
- Primarily keyboard navigation
- Mouse support is optional and often disabled
- Focus on keyboard shortcuts and mnemonics

### gocui (Simple TUI)
- Keyboard-only by design
- No mouse support in core functionality
- Emphasizes simplicity and predictability

### termui (Dashboard-style)
- Mixed approach but keyboard-preferred
- Mouse support mainly for interactive widgets
- Standard navigation keys for scrolling

## Recommended Implementation Plan

### 1. Immediate Changes
```go
// Disable mouse wheel in viewport initialization
h.viewport.MouseWheelEnabled = false
```

### 2. Simplified Keyboard Support
Standard TUI navigation (Page Up/Page Down only):
- Page Up/Page Down: Full page scroll (universal standard)
- Up/Down arrows: Line-by-line scroll

### 3. Update Documentation
Update README.md navigation section:
```markdown
## Navigation
- **Tab/Shift+Tab**: Switch between tabs
- **Left/Right**: Navigate fields within tab  
- **Up/Down**: Scroll viewport line by line
- **Page Up/Page Down**: Scroll viewport page by page
- **Enter**: Edit/Execute
- **Esc**: Cancel edit
- **Ctrl+C**: Exit
```

### 4. Update ShortcutsHandler
Update the built-in shortcuts display to reflect the new navigation scheme.

## Conclusion

**Final Decision: DISABLE mouse scroll and implement SIMPLIFIED keyboard navigation**

This aligns DevTUI with:
1. **TUI Best Practices**: Standard keyboard-only navigation using universal keys
2. **User Expectations**: Page Up/Down is recognized by ALL users
3. **Practical Needs**: Enable text selection for debugging
4. **Simplicity**: One clear way to navigate pages (Page Up/Page Down)
5. **Universal Compatibility**: Works consistently across all terminal environments

The implementation prioritizes simplicity and universal accessibility over advanced navigation options.

The implementation should be straightforward and will result in a more standard, accessible, and user-friendly interface that better serves DevTUI's purpose as a development tool abstraction.

## Next Steps

1. **Approval**: Review and approve this analysis
2. **Implementation**: Implement the recommended changes
3. **Testing**: Verify navigation works correctly in various terminals
4. **Documentation**: Update all relevant documentation
5. **Migration**: Consider if any migration notices are needed for existing users

This change represents a move toward TUI best practices and will improve the overall user experience while maintaining DevTUI's core functionality.

## Estado de Implementación ✅

### Cambios Implementados (2025-01-27)

#### 1. Deshabilitación Completa del Mouse en Bubbletea

**Archivo**: `init.go`
- ❌ **Removido**: `tea.WithMouseCellMotion()` de las opciones del programa
- ✅ **Resultado**: Mouse completamente deshabilitado desde la inicialización

**Antes**:
```go
tui.tea = tea.NewProgram(tui,
    tea.WithAltScreen(),
    tea.WithMouseCellMotion(), // ← REMOVIDO
)
```

**Después**:
```go
tui.tea = tea.NewProgram(tui,
    tea.WithAltScreen(),
    // Mouse support disabled to enable terminal text selection
)
```

#### 2. Simplificación del Update Loop

**Archivo**: `update.go`
- ✅ **Simplificado**: Filtrado de eventos innecesario removido
- ✅ **Mejora**: El viewport ahora recibe todos los eventos (solo teclado ya que mouse está deshabilitado)

#### 3. Navegación Actualizada

**Archivo**: `userKeyboard.go`
- ✅ **Implementado**: `tea.KeyPgUp` para Page Up navigation
- ✅ **Implementado**: `tea.KeyPgDown` para Page Down navigation

#### 4. Documentación

**Archivo**: `README.md`
- ✅ **Actualizado**: Sección de navegación con nuevos controles

### Solución Final

La implementación final utiliza la **Opción A** del análisis original:

1. **Sin soporte de mouse**: Programa configurado sin `WithMouseCellMotion()`
2. **Navegación solo por teclado**: Page Up/Page Down para navegación rápida
3. **Selección de texto habilitada**: Los usuarios pueden seleccionar y copiar texto directamente del terminal

### Beneficios Alcanzados

- ✅ **Selección de texto habilitada**: Los usuarios pueden copiar errores y logs
- ✅ **Navegación mejorada**: Page Up/Page Down para movimiento rápido por páginas
- ✅ **Simplicidad**: Un solo método de navegación, consistente con herramientas estándar de terminal
- ✅ **Compatibilidad universal**: Funciona en todos los terminales y multiplexers

### Verificación Requerida

Después de estos cambios, verificar que:
1. ✅ La selección de texto funciona en el terminal
2. ✅ Page Up/Page Down navegan correctamente
3. ✅ Up/Down continúan funcionando para navegación línea por línea

### Resultado Final Verificado ✅

**Funcionalidad Confirmada (2025-01-27)**:
- ✅ **Selección de texto**: Funciona perfectamente para copiar errores y logs
- ✅ **Navegación por teclado**: Page Up/Page Down implementados y funcionales
- ✅ **Scroll con mouse**: Sorprendentemente funciona también (mejora de bubbletea)
- ✅ **Compatibilidad completa**: Tanto navegación por teclado como scroll del mouse disponibles

**Observación Importante**: Las versiones recientes de bubbletea han mejorado el manejo del mouse, permitiendo que funcione tanto el scroll del mouse como la selección de texto simultáneamente. Esto proporciona la mejor experiencia de usuario posible.

**Beneficio Final**: Los usuarios pueden usar tanto teclado (Page Up/Page Down) como mouse wheel para navegar, Y además pueden seleccionar y copiar texto libremente del terminal.
