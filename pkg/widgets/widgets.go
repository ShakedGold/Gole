package widgets

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ShakedGold/Gole/pkg/widgets/grid"
)

type GoleTheme struct{}

func Grid(cols int) grid.Grid {
	list := new(widget.List)
	list.Axis = layout.Vertical
	return grid.Grid{
		Columns:    cols,
		List:       list,
		ItemWidth:  200,
		ItemHeight: 200,
	}
}
