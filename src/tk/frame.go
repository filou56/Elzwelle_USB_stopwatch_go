package tk

import (
//	"log"
)

func Frame(parent *Widget, args... Option) *Widget {
	var frame Widget
	
	frame.CreateElement(parent)
	Command <- "frame " + frame.id + options(args)
	
	return &frame
}

