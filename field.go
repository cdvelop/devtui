package devtui

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
