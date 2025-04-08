package devtui

type Field struct {
	Name             string                                                // eg: "Server Port"
	Value            string                                                //initial Value eg: "8080"
	Editable         bool                                                  // if no editable eject the action FieldValueChange directly
	FieldValueChange func(newValue string) (execMessage string, err error) //eg: "8080" -> "9090" execMessage: "Port changed from 8080 to 9090"
	//internal use
	tempEditValue string // use for edit
	index         int
	cursor        int // cursor position in text value
}

func (f *Field) SetCursorAtEnd() {
	// Calculate cursor position based on rune count, not byte count
	f.cursor = len([]rune(f.Value))
}
