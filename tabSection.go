package devtui

// NewTabSection creates a new tab section
func (t *DevTUI) NewTabSection(title string, fhAdapters ...fieldHandlerAdapter) *tabSection {
	// create a new tab section with the given title and field handlers
	newSection := &tabSection{
		title:                title,
		fieldHandlers:        make([]fieldHandler, 0, len(fhAdapters)),
		tuiMessages:          make([]tuiMessage, 0),
		indexActiveEditField: 0, //  selected field index
		tui:                  t,
	}

	// build fieldHandler from fieldHandlers
	for i, fha := range fhAdapters {
		stdHandler := fieldHandler{
			fieldHandlerAdapter: fha,
			tempEditValue:       "",
			index:               i,
			cursor:              len(fha.Value()), // position cursor at end of text by default
		}
		newSection.fieldHandlers = append(newSection.fieldHandlers, stdHandler)
	}

	// Set the index of the tab section (usually managed by DevTUI internally)
	newSection.index = len(t.tabSections)

	// Add to the DevTUI sections
	t.tabSections = append(t.tabSections, *newSection)

	return newSection
}
