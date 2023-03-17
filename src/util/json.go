package util

import (
	"fmt"
	"log"
	"strings"
	"encoding/json"
//	conv "github.com/cstockton/go-conv"
)

// -----------------------------------------------------------------------------------------------------------

func DecodePayload(data []byte) (map[string]interface{}) {
	dec := json.NewDecoder(strings.NewReader(string(data)))	
	for {
		var m map[string]interface{}
		if err := dec.Decode(&m); err != nil {
			log.Println("Sheet:\tDecodePayload: ",err," in ",string(data))
            return nil
        }
        return m
	}
	
	return nil;
}

func DumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			DumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}

func GetJsonValue(jPath string, jMap map[string]interface{} ) interface{} {
	
	items := strings.Split(jPath,"/")
	
	for i:=0; i<len(jMap); i++ {
		item := jMap[items[i]]
		//fmt.Println("Loop ",i,items[i],item)
		switch item.(type) {
		case map[string]interface{}:
			jMap = item.(map[string]interface{})
			break
		default:
			return item
		}
	}
	
	return nil
}
