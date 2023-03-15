package tk

import (
	"util"
	"fmt"
//	"log"
)

var FONT_SIZE_24 	= Option{"font","{size 24}"} 

func FontSize(size int) Option {
	return Option{"font",fmt.Sprintf("{size %d}",size)} 
}

type Widget struct {
	id			string
	onClick		func()
	parent		*Widget	
}

var widgets map[string]*Widget

func init() {
	widgets = make(map[string]*Widget)
}

var ROOT = &Widget{"",nil,nil}

func (widget *Widget) CreateElement(parent *Widget) *Widget {
	widget.onClick = nil
	widget.parent  = parent
	widget.id = widget.parent.id + "." + util.NextID()
	widgets[widget.id] = widget
	return widget
}

func GetWidgetByID(id string) *Widget {
	return 	widgets[id]
}

func (widget *Widget) SetOnClick(fn func()) {
	widget.onClick = fn 
}

func (widget *Widget) OnClick() {
	if widget.onClick != nil {
		go widget.onClick()
	} 
}

func (widget *Widget) SetText(text string) {
	Command <- widget.id + " configure -text \"" + text + "\""
}

func (widget *Widget) SetFont(font string) {
	Command <- widget.id + " configure -font \"" + font + "\""
}

func GetFonts() (error, string) {
	Command <- `puts [concat "{\"STRING\":\"" [font families] "\"}"]`
	return WaitForString()
}

/*
 
"{\"EVENT\":\"onClick\",\"ID\":\"`+id+`\"}"

-font Schrift Schrift
-foreground Farbe Textfarbe
-background Farbe
*/