package tk

import (
	"fmt"
//	"log"
)

func Entry(parent *Widget, text string, args... Option) *Widget {
	var entry Widget
	
	entry.CreateElement(parent)
	//Command <- "entry " + entry.id + options(args)
	Command <- fmt.Sprintf("entry %s -textvariable %s_entry %s",entry.id,entry.id[1:],options(args))
	
	return &entry
}

func (widget *Widget) GetVariable() (error,string) {
	//log.Printf( `puts [ concat "{\"STRING\":\"" $`+widget.id[1:]+ `_entry "\"}" ]`)
	Command <- `puts [ concat "{\"STRING\":\"" $`+widget.id[1:]+ `_entry "\"}" ]`
	return WaitForString() 
}
