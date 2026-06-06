package sidebar

import (
	CommonWidgets "Zbolt/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type popup_skeleton struct {
	gui.DefaultWidget

	heading      CommonWidgets.Description
	input_widget CommonWidgets.TextInputWithContextMenu
	button       widget.Button
}
