package tk

import (
	"log"
)

func (widget *Widget) Pack (args... Option) {
	pack := "pack " + widget.id + options(args)
	Command <- pack
	log.Println(pack)
}

var PACK_TOP 	= Option{"side","top"} 
var PACK_LEFT 	= Option{"side","left"} 
var PACK_RIGHT 	= Option{"side","right"} 
var PACK_BOTTOM = Option{"side","bottom"} 

var PACK_EXPAND = Option{"expand","yes"} 
var PACK_X		= Option{"fill","x"}
var PACK_Y		= Option{"fill","y"}
var PACK_BOTH	= Option{"fill","both"}
/*
-side top | left | right | botom
*/