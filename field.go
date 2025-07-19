package devtui

// use NewField to create a new field in the tab section
// Field represents a field in the TUI with a name, value, and editable state.
// field represents a field in the TUI with a name, value, and editable state.
type field struct {
	name       string                                             // eg: "Server Port"
	value      string                                             // initial Value eg: "8080"
	editable   bool                                               // if no editable eject the action changeFunc directly
	changeFunc func(newValue any) (execMessage string, err error) //eg: "8080" -> "9090" execMessage: "Port changed from 8080 to 9090"
	// internal use
	tempEditValue string // use for edit
	index         int
	cursor        int // cursor position in text value
}

// SetTempEditValueForTest permite modificar tempEditValue en tests
func (f *field) SetTempEditValueForTest(val string) {
	f.tempEditValue = val
}

// SetCursorForTest permite modificar el cursor en tests
func (f *field) SetCursorForTest(cursor int) {
	f.cursor = cursor
}

// NewField creates a new field, adds it to the tabSection, and returns the tabSection for chaining.
// Example usage:
//   tab.NewField("username", "defaultUser", true, func(newValue any) (string, error) { ... })
func (ts *tabSection) NewField(name, value string, editable bool, changeFunc func(newValue any) (string, error)) *tabSection {
	f := &field{
		name:       name,
		value:      value,
		editable:   editable,
		changeFunc: changeFunc,
	}
	ts.addFields(f)
	return ts
}

// setFieldHandlers sets the field handlers slice (mainly for testing)
// Only for internal/test use
func (ts *tabSection) setFieldHandlers(handlers []*field) {
	ts.fieldHandlers = handlers
}

// addFields adds one or more field handlers to the section (private)
func (ts *tabSection) addFields(fields ...*field) {
	ts.fieldHandlers = append(ts.fieldHandlers, fields...)
}

func (f *field) Name() string {
	return f.name
}

func (f *field) SetName(name string) {
	f.name = name
}

func (f *field) Value() string {
	return f.value
}

func (f *field) SetValue(value string) {
	f.value = value
}

func (f *field) Editable() bool {
	return f.editable
}

func (f *field) SetEditable(editable bool) {
	f.editable = editable
}

func (f *field) SetCursorAtEnd() {
	// Calculate cursor position based on rune count, not byte count
	f.cursor = len([]rune(f.value))
}
