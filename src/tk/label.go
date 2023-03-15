package tk

import (
//	"log"
)

func Label(parent *Widget, text string, args... Option) *Widget {
	var label Widget
	
	label.CreateElement(parent)
	Command <- "label " + label.id + " -text \"" + text + "\"" + options(args)
	//log.Println(button)
	
	return &label
}

