package tk

import (
	"fmt"
)

func basicList(id string, args []Option)  string {
	return fmt.Sprintf(`
		ttk::frame %s
		listbox %s.listbox -width 15 -height 15 -yscrollcommand {%s.scrollbar_y set} -listvariable %s_list %s
		ttk::scrollbar %s.scrollbar_y -command {%s.listbox yview}
		pack %s.scrollbar_y -side right -fill y
		pack %s.listbox -side top -expand yes -fill both
	`,id,id,id,id,options(args),id,id,id,id)
}

func List(parent *Widget, args... Option) *Widget {
	var list Widget
	
	list.CreateElement(parent)
	Command <- basicList(list.id,args)	
	return &list
}

func (widget *Widget) ListAppend(item string) {
	Command <- "lappend " + widget.id + "_list \"" + item + "\""
	// show last entry in list 	
	Command <- widget.id + ".listbox see [" + widget.id + ".listbox size]"	
}
