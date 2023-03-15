package notify

import (
	"log"
	"github.com/godbus/dbus"
	"time"
)

var (
	Receiver chan string
)

func Listen() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
	Receiver = make(chan string,10)	
	defer close(Receiver)
	
	log.Println("Notify:\tListen!")
	
	loop := true
	
	for loop {
		select {
			case incoming := <- Receiver:
			log.Printf("Notify:\tMessage: %s\n", incoming)	
			
			obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
			call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
				"", "Stopwatch Notify", incoming, []string{},
				map[string]dbus.Variant{}, int32(5000))
			if call.Err != nil {
				log.Println("Notify:\tDBus Send:",err)	
				loop = false
			}	
		}
	}
}

func millisTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func Notification() func(msg string, interval int64) {
    var stamp int64 = 0
    
    return func(msg string, interval int64) {
    	var now int64 = millisTimestamp()
    	
        if (stamp + interval) < now {
        	Receiver <- msg
        	stamp = now
        } 
        return
    }
}