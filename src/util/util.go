package util

import (
	"fmt"
	"time"
	"github.com/sony/sonyflake"
)

var unique *sonyflake.Sonyflake

func init() {
	unique =  sonyflake.NewSonyflake(sonyflake.Settings{time.Unix(0,0),func()(uint16, error){return 0x0c29,nil},nil})
	if unique == nil {
		panic("sonyflake not created")
	}
}

func NextID() string {
	id,err := unique.NextID()
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	return fmt.Sprintf("%X",id)
}

func millisTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func TimedLock() func(interval int) bool {
    var stamp int64 = 0
    
    return func(interval int) bool {
    	var now int64 = millisTimestamp()
    	
        if (stamp + int64(interval)) < now {
        	stamp = now
        	return false
        } 
        return true
    }
}