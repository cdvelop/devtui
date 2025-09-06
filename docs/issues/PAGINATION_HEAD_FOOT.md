# PAGINATION_HEAD_FOOT Implementation Plan

## Overview
Implementation of pagination indicators in header and footer with format `[ current/ total]` (1-based display, e.g. `[1/1]` for a single element) using the same visual style as informational elements (footer style with background color and white text).

## Current State Analysis

### Header Structure
```
[demo/tabName]______________________________________________           
```
- Title: `AppName + "/" + tab.Title()`
- Line: Filled with `─` characters
- Style: `headerTitleStyle` (highlight background, foreground text)

### Footer Structure - Per Handler Type

#### HandlerEdit
```
[Database Co...]  [postgres://localhost:5432/mydb          ][100%]
 ↑                 ↑                                         ↑
Label()           Value()                                  scroll%
```

#### HandlerExecution  
```
[Create Syst...]  [Create System Backup                   ][100%]
 ↑                       ↑                                  ↑
Label()               Value()                            scroll%   
```

#### HandlerDisplay
```
[System Status Information Display                          ][100%]
 ↑                                                           ↑
Name()                                                    scroll% 
```

#### Writers Only
```
____________________________________________________________[100%]
                                                            ↑
                                                         scroll%
```

## Proposed Changes

### Header Pagination
```
[demo/tabName]__________________________________________[ 1/ 1]
```
- Current tab index vs total tabs (1-based for display, developer-friendly internally)
- Exact format `[99/99]` - maximum 2 digits each
- Same style as footer informational elements (same as scroll percentage)

### Footer Pagination
```
[ 1/ 4][help                                                      ][100%]
[ 1/ 4][System Status Information Display                         ][100%]
[ 2/ 4][Database Co...]  [postgres://localhost:5432/mydb          ][100%]
[ 3/ 4][Create Syst...]  [Create System Backup                    ][100%]
[ 4/ 4] __________________________________________________________ [100%]
```
- Current field index vs total field handlers in tab (1-based for display)
- Exact format `[99/99]` - maximum 2 digits each  
- Same informational style as scroll percentage

## Implementation Plan

### Files to Modify

#### 1. `/view.go` - Header Pagination
**Function**: `headerView()`
**Changes**:
- Calculate current tab index (`h.activeTab + 1`) and total tabs (`len(h.tabSections)`)
- Create pagination indicator with format `[ %d/%d]`
- Apply `paginationStyle` for consistency between header and footer
- Adjust line width calculation to accommodate pagination space
- Join: `title` + `line` + `pagination`

**Key Calculations**:
```go
currentTab := h.activeTab  // 0-based internally
totalTabs := len(h.tabSections)

// Check limits and log error if exceeded
if currentTab > 99 || totalTabs > 99 {
    if h.Logger != nil {
        h.Logger("Tab limit exceeded:", currentTab, "/", totalTabs)
    }
}

// Clamp values to display limits
displayCurrent := min(currentTab, 99) + 1 // 1-based for display
displayTotal := min(totalTabs, 99)

pagination := fmt.Sprintf("[%2d/%2d]", displayCurrent, displayTotal)
paginationStyled := h.paginationStyle.Render(pagination)
lineWidth := h.viewport.Width - lipgloss.Width(title) - lipgloss.Width(paginationStyled)
line := h.lineHeadFootStyle.Render(strings.Repeat("─", max(0, lineWidth)))
return lipgloss.JoinHorizontal(lipgloss.Center, title, line, paginationStyled)
```

#### 2. `/footerInput.go` - Footer Pagination 
**Function**: `renderFooterInput()`
**Changes**:
- Calculate current field index (`tabSection.indexActiveEditField`) and total fields (`len(fieldHandlers)`)
- Create field pagination with format `[ %d/%2d]`
- Apply `paginationStyle` for visual consistency
- Modify layout calculations for all handler types

**For HandlerDisplay**:
```go
currentField := tabSection.indexActiveEditField  // 0-based internally
totalFields := len(fieldHandlers)

// Check limits and log error if exceeded
if currentField > 99 || totalFields > 99 {
    if h.Logger != nil {
        h.Logger("Field limit exceeded:", currentField, "/", totalFields)
    }
}

// Clamp values to display limits
displayCurrent := min(currentField, 99) + 1 // 1-based for display
displayTotal := min(totalFields, 99)

fieldPagination := fmt.Sprintf("[%2d/%2d]", displayCurrent, displayTotal)
paginationStyled := h.paginationStyle.Render(fieldPagination)
remainingWidth := h.viewport.Width - lipgloss.Width(info) - lipgloss.Width(paginationStyled) - horizontalPadding*2
// Layout: [Pagination] [Label expanded] [Scroll%]
spacerStyle := lipgloss.NewStyle().Width(horizontalPadding).Render("")
return lipgloss.JoinHorizontal(lipgloss.Left, paginationStyled, spacerStyle, styledLabel, spacerStyle, info)
```

**For HandlerEdit/HandlerExecution**:
```go
currentField := tabSection.indexActiveEditField  // 0-based internally
totalFields := len(fieldHandlers)

// Check limits and log error if exceeded (same as above)
if currentField > 99 || totalFields > 99 {
    if h.Logger != nil {
        h.Logger("Field limit exceeded:", currentField, "/", totalFields)
    }
}

// Clamp values to display limits
displayCurrent := min(currentField, 99) + 1 // 1-based for display
displayTotal := min(totalFields, 99)

fieldPagination := fmt.Sprintf("[%2d/%2d]", displayCurrent, displayTotal)
paginationStyled := h.paginationStyle.Render(fieldPagination)
valueWidth := h.viewport.Width - usedWidth - lipgloss.Width(paginationStyled) - horizontalPadding*2
// Layout: [Pagination] [Label] [Value] [Scroll%]
spacerStyle := lipgloss.NewStyle().Width(horizontalPadding).Render("")
return lipgloss.JoinHorizontal(lipgloss.Left, paginationStyled, spacerStyle, paddedLabel, spacerStyle, styledValue, spacerStyle, info)
```

**Function**: `footerView()` - Writers Only
**Changes**:
- Add field pagination for tabs with only writers (no field handlers)
- When `len(fieldHandlers) == 0`, show pagination as `[ 0/ 0]` (or `[1/1]` if you want to show a single item as selected)
- Apply same limit checking and error logging

#### 3. `/style.go` - New Pagination Style
**Changes**:
- Add `paginationStyle` field to `tuiStyle` struct
- Create new pagination style specifically for pagination indicators
- Background: `Primary` color 
- Foreground: `Foreground` color
- Avoid confusion with existing styles

```go
// Add to tuiStyle struct
type tuiStyle struct {
    *ColorPalette
    // ... existing fields ...
    footerInfoStyle  lipgloss.Style
    paginationStyle  lipgloss.Style // NEW: For pagination indicators
    // ... rest of fields ...
}

// Add in newTuiStyle function
t.paginationStyle = lipgloss.NewStyle().
    Padding(0, 0).
    Background(lipgloss.Color(t.Primary)).
    Foreground(lipgloss.Color(t.Foreground))
```

### Key Implementation Details

#### Tab Navigation Context
- Current tab: `h.activeTab` (0-based, developer-friendly) 
- Total tabs: `len(h.tabSections)`
- Tabs include: SHORTCUTS (auto-created at index 0) + user tabs
- **Limits**: Maximum 99 tabs, log error if exceeded

#### Field Navigation Context
- Current field: `tabSection.indexActiveEditField` (0-based)
- Total fields: `len(tabSection.fieldHandlers)`
- Field types: HandlerDisplay, HandlerEdit, HandlerExecution
- **Limits**: Maximum 99 field handlers per tab, log error if exceeded

#### Layout Calculations
- **Width Priority**: Pagination → Label → Value → Scroll%
- **Space Distribution**: 
  - Pagination: Fixed width `[99/99]` = 6 characters
  - Scroll%: Fixed width `[100%]` = 6 characters  
  - Label: Fixed width `labelWidth` for Edit/Execution
  - Value: Remaining space

#### Edge Cases
- Single tab: Show `[ 1/ 1]` (for clarity, even if only writers are present)
- No fields (writers-only tabs): Should show `[ 1/ 1]` for clarity. If `[ 0/ 0]` is shown, this is confusing and is a bug. See test `TestPaginationWritersOnlyTab`.
- Limit exceeded: Clamp to 99, log error via `Logger`
- Large numbers: Handle with `%2d` format for exact `[99/99]` spacing

### Testing Considerations
- Verify pagination updates on tab navigation (Tab/Shift+Tab)
- Verify field pagination updates on field navigation (Left/Right)
- Test layout with various screen widths
- Test with single tab and multiple tabs
- Test with tabs having different field counts (0, 1, multiple, and only writers)
- Test writers-only tab: pagination should show `[ 1/ 1]`, not `[ 0/ 0]`. The test `TestPaginationWritersOnlyTab` will fail if the code is incorrect.
### Refactor Note

To fix the bug, update the footer pagination logic so that when there are no field handlers (writers-only tab), both current and total are set to 1 for display purposes. This ensures `[ 1/ 1]` is shown, matching the documentation and user expectations.
- **Test limit handling**: 99+ tabs and 99+ field handlers per tab
- **Test error logging**: Verify `Logger` is called when limits exceeded

### Visual Style Consistency
- **Both paginations**: New `paginationStyle` (Primary background + Foreground text)
- **Distinct from scroll%**: Different style to avoid confusion with existing elements
- **Header layout**: Title + Line (no spaces) + Pagination
- **Footer layout**: Pagination + Space + Elements + Space + Scroll%
- **Format**: Exact `[99/99]` spacing for both tabs and fields

### Developer-Friendly Design
- **1-based display indexing**: More intuitive for users, avoids confusion when only one element is present
- **Consistent with internal variables**: Internally 0-based, but display is always 1-based for clarity
- **Practical limits**: 99 tabs/fields is more than reasonable for any application
- **Error handling**: Clear logging when architectural limits are approached

This implementation maintains the existing layout structure while adding pagination context optimized for developer workflow and debugging.
