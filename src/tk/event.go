package tk

import (
	"strings"
	"log"
	"encoding/json"
	"errors"
	"time"
)

func onClickEvent(id string) string {
	return ` -command { puts "{\"EVENT\":\"onClick\",\"ID\":\"`+id+`\"}" }`
}

func eventHandler() {
	var ev_map map[string]interface{}
	
	for {
		msg := <- Event
		//log.Println("\nEventHandler",string(msg))
		
		dec := json.NewDecoder(strings.NewReader(string(msg)))
		if err := dec.Decode(&ev_map); err != nil {
			log.Println("DecodePayload: ",err," in ",string(msg))
		}
		//log.Println(ev_map)
		
	    if (ev_map["EVENT"] != nil ) {
		    if ev_map["EVENT"].(string) == "onClick" {
		    	id := ev_map["ID"].(string)
		    	//log.Println("onClickEvent")
		    	widget := GetWidgetByID(id)
		    	if widget != nil {
		    		//log.Println("onClick")
		    		widget.OnClick()
		    	}
		    }
	    }	    
	    
	    if (ev_map["STRING"] != nil ) {
	    	str := ev_map["STRING"].(string)
	    	//log.Println("STRING:\t",str)
	    	Result <- str
	    }
	}
}

func WaitForString() (error, string ) {
	for {
	    select {
		    case str := <- Result:
	        return nil, str
		    case <-time.After(time.Second):
	        return errors.New("TCL Timeout"), ""
	    }
	}
}