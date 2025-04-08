package devtui

// use NewField to create a new field in the tab section
// Field represents a field in the TUI with a name, value, and editable state.
type Field struct {
	name       string                                                // eg: "Server Port"
	value      string                                                // initial Value eg: "8080"
	editable   bool                                                  // if no editable eject the action changeFunc directly
	changeFunc func(newValue string) (execMessage string, err error) //eg: "8080" -> "9090" execMessage: "Port changed from 8080 to 9090"
	//internal use
	tempEditValue string // use for edit
	index         int
	cursor        int // cursor position in text value
}

// NewField creates and returns a new Field instance.
//
// Parameters:
// - name: The name of the field.
// - value: The initial value of the field.
// - editable: A boolean indicating whether the field is editable.
// - changeFunc: A callback function that is invoked when the field's value changes.
//   It receives the new value as a parameter and returns the updated value and an error, if any.
//
// Returns:
// - A pointer to the newly created Field instance.
//
// Example usage:
//   field := NewField("username", "defaultUser", true, func(newValue string) (string, error) {
//       if len(newValue) == 0 {
//           return "", errors.New("value cannot be empty")
//       }
//       return newValue, nil
//   })
func NewField(name, value string, editable bool, changeFunc func(newValue string) (string, error)) *Field {
	return &Field{
		name:       name,
		value:      value,
		editable:   editable,
		changeFunc: changeFunc,
	}
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) SetName(name string) {
	f.name = name
}

func (f *Field) Value() string {
	return f.value
}

func (f *Field) SetValue(value string) {
	f.value = value
}

func (f *Field) Editable() bool {
	return f.editable
}

func (f *Field) SetEditable(editable bool) {
	f.editable = editable
}

func (f *Field) SetCursorAtEnd() {
	// Calculate cursor position based on rune count, not byte count
	f.cursor = len([]rune(f.value))
}
