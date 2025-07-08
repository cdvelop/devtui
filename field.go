package devtui

// use NewField to create a new field in the tab section
// Field represents a field in the TUI with a name, value, and editable state.
type Field struct {
	name       string                                             // eg: "Server Port"
	value      string                                             // initial Value eg: "8080"
	editable   bool                                               // if no editable eject the action changeFunc directly
	changeFunc func(newValue any) (execMessage string, err error) //eg: "8080" -> "9090" execMessage: "Port changed from 8080 to 9090"
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
//   It receives the new value as any type and returns an execution message and an error, if any.
//
// Returns:
// - A pointer to the newly created Field instance.
//
// Example usage:
//   // String field
//   field := NewField("username", "defaultUser", true, func(newValue any) (string, error) {
//       strValue := newValue.(string)
//       if len(strValue) == 0 {
//           return "", errors.New("value cannot be empty")
//       }
//       return "Username updated to " + strValue, nil
//   })
//
//   // Numeric field
//   portField := NewField("port", "8080", true, func(newValue any) (string, error) {
//       switch v := newValue.(type) {
//       case string:
//           if port, err := strconv.Atoi(v); err == nil && port > 0 && port < 65536 {
//               return fmt.Sprintf("Port changed to %d", port), nil
//           }
//           return "", errors.New("invalid port number")
//       case int:
//           if v > 0 && v < 65536 {
//               return fmt.Sprintf("Port changed to %d", v), nil
//           }
//           return "", errors.New("port out of range")
//       default:
//           return "", errors.New("unsupported type")
//       }
//   })
//
//   // Boolean field
//   enabledField := NewField("enabled", "false", true, func(newValue any) (string, error) {
//       switch v := newValue.(type) {
//       case bool:
//           return fmt.Sprintf("Enabled set to %t", v), nil
//       case string:
//           if b, err := strconv.ParseBool(v); err == nil {
//               return fmt.Sprintf("Enabled set to %t", b), nil
//           }
//           return "", errors.New("invalid boolean value")
//       default:
//           return "", errors.New("unsupported type")
//       }
//   })
func NewField(name, value string, editable bool, changeFunc func(newValue any) (string, error)) *Field {
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
