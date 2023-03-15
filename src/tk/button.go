package tk

import (
//	"log"
)

func Button(parent *Widget, text string, args... Option) *Widget {
	var button Widget
	
	button.CreateElement(parent)
	Command <- "button " + button.id + " -text \"" + text + "\"" + onClickEvent(button.id) + options(args)
	//log.Println(button)
	
	return &button
}

