# DevTUI - Complete Purpose and Functionality Description

## What is DevTUI?

DevTUI is a **reusable generic abstraction** for terminal user interfaces (TUI) built on top of [bubbletea](https://github.com/charmbracelet/bubbletea) and [bubbles](https://github.com/charmbracelet/bubbles). It provides a pre-configured, fixed, and minimalist interface where you can **inject different handlers** to display their messages in the terminal in an organized way.

## Problem it Solves

### The Original Problem
During development of Go-to-WASM compilation tools with:
- File change detection
- Browser hot reload
- CSS/JS file minification
- Complex configuration for full-stack Go applications

**The problems were:**
- Too many scattered logs everywhere
- Complex and confusing visual information
- Difficult to understand what was really happening
- The view layer grew too much and became unmanageable

### The Solution: DevTUI
A TUI that **separates the view layer** from business logic, enabling:
- **Organized logs in the same space** (they don't accumulate infinitely)
- **Automatic reordering** always showing what happened last
- **Arrow navigation** maintaining focus without filling the UI with unnecessary elements
- **Injection of specialized handlers** according to need

## Ideal Use Cases

### Minimalist Development Tools
- **Real-time compilers** (Go-to-WASM, etc.)
- **File monitors** with automatic actions
- **Development dashboards** with live metrics
- **Configuration interfaces** for complex pipelines
- **Build tools** with multiple steps
- **Asset minifiers** with visual progress

### Integration with Development Environments
- **VS Code**: Integrated terminal with limited space
- **Text editors**: Auxiliary panel for monitoring
- **CI/CD pipelines**: Real-time monitoring interface
- **Docker workflows**: Container configuration and status

## Architecture and Key Concepts

### What are "Handlers"?
Handlers are **injectable components** that encapsulate:
- Specific business logic (compilation, minification, etc.)
- Information presentation (logs, status, configuration)
- User interaction (input fields, action buttons)

**They are not traditional widgets**, but **functionality abstractions** that are automatically rendered according to their type.

### Tab System and Organization
- **Thematic tabs**: Group related handlers (Config, Build, Logs, etc.)
- **One active element**: Only one handler is shown at a time (maintains focus)
- **Arrow navigation**: Left/Right to switch between handlers, Tab/Shift+Tab to switch between tabs
- **Informative footer**: Context of the active handler
- **Automatic ShortcutsHandler**: "SHORTCUTS" tab automatically loaded at position 0 with navigation help

### MessageTracker: The Key Differentiator
**Traditional problem**: Logs accumulate infinitely creating visual noise.

**MessageTracker solution**: 
- Handlers can **update existing messages** instead of creating new ones
- **Reuse the same space** on screen through `operationID`
- **Progressive logs** showing current state of long operations with `progress` callbacks
- **Clean history** without information saturation
- **Thread-safe**: Protection with `sync.RWMutex` for concurrent operations

## Comparison with Other TUI Libraries

### vs bubbletea + bubbles (base)
- **DevTUI**: Pre-configured abstraction with specific patterns and integrated viewport
- **bubbletea/bubbles**: General framework, requires implementing all UI logic

### vs tview, termui, gocui
- **DevTUI**: Focus on injectable handlers for development
- **Others**: General widgets for complete applications

### Unique Advantage: Functional Minimalism
- **1-4 methods per handler** vs complex implementations
- **Specialized interfaces** by specific purpose
- **Minimal configuration** for common development use cases
- **Method chaining**: `.WithTimeout()` and `.Register()` for fluid configuration
- **Auto-detection**: Automatically recognizes capabilities like MessageTracker

## Handler Types and Their Purposes

### 1. HandlerDisplay (2 methods)
**Purpose**: Read-only information that displays immediately
**Cases**: System status, metrics, help (like ShortcutsHandler), current configuration

### 2. HandlerEdit (4 methods)
**Purpose**: Interactive input fields with validation
**Cases**: Port configuration, URLs, file paths, compilation parameters

### 3. HandlerExecution (3 methods)
**Purpose**: Action buttons with optional progress callbacks
**Cases**: Compile, deploy, clear cache, restart services, backups

### 4. HandlerWriter (1 method)
**Purpose**: Basic logging (always new lines)
**Cases**: Application logs, system events, command output

### 5. HandlerWriterTracker (3 methods)
**Purpose**: Advanced logging (can update existing lines)
**Cases**: Compilation progress, deployment status, continuous monitoring

### Handlers with Tracking (*Tracker interfaces)
All handlers can implement `MessageTracker` for advanced capabilities:
- **HandlerEditTracker**: Edit + MessageTracker
- **HandlerExecutionTracker**: Execution + MessageTracker  
- **HandlerWriterTracker**: Writer + MessageTracker (built-in)

## When NOT to use DevTUI

- **Complex end-user applications** (use tview, bubbletea directly)
- **Multi-window GUIs** (DevTUI is single-window)
- **Highly customized interfaces** (DevTUI prioritizes consistency)
- **Web or desktop applications** (DevTUI is terminal-specific)

## Technical Benefits

### For the Developer
- **Fast implementation**: Minimal and clear interfaces
- **Separation of concerns**: View separated from business logic
- **Reusability**: Portable handlers between projects
- **Maintenance**: UI changes don't affect business logic
- **Testing**: Integrated test mode (`testMode`) for automated verification
- **Focus management**: Automatic focus return system (`ReturnFocus()`)
- **Error handling**: Integration with `messagetype` for automatic classification

### For the End User
- **Consistent experience**: Standard navigation across all tools
- **Organized information**: No visual saturation, viewport with automatic scroll
- **Real-time feedback**: Progress callbacks and message tracking
- **Efficient space**: Maximum terminal utilization
- **Integrated help**: Automatic ShortcutsHandler with navigation commands

## Relationship with GoDEV App

DevTUI is the **main interface** of GoDEV App, a Go development tool that includes:
- Go-to-WASM compilation
- Browser hot reload
- Asset minification
- Dependency management
- File monitoring

GoDEV demonstrated the need to separate the view layer, giving rise to DevTUI as an independent project.
