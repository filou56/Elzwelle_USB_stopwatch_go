package main 

import (
	"tk"
	"fmt"
	"log"
	"time"
	"os"
	"flag"
	"github.com/mkch/gpio"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"mqttpipe"
	"github.com/goburrow/serial"
	"strconv"
	"notify"
	"util"
//	"sysinfo"
	"strings"
)

var (
	mqttHost 			string 	= "localhost"
	debounceStart  		uint
	debounceFinish  	uint
	enableButtons		bool
	usbStopwatch		bool
	usbDevice			string
	usbPollIntervall	int
 	usbStartMillis		int64
	usbFinishMillis		int64
	usbSyncEvent		int64
	masterMillis		int64	= 0
	publishInterval		int
	
	baseTime			time.Time
)

func init() {	
	baseTime = time.Now()
}

//------------------------- MQTT ------------------------
func mqttConnect(clientId string, uri *url.URL) mqtt.Client {
	opts := mqttClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func mqttClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
//  ---------- Options -------------	
//	opts.SetUsername(uri.User.Username())
//	password, _ := uri.User.Password()
//	opts.SetPassword(password)
//	opts.SetClientID(clientId)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true)
//	opts.SetMaxReconnectInterval(10 * time.Second)
//  --------------------------------

	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("MQTT connection lost error: %s\n" + err.Error())
	})
	
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		log.Println("MQTT reconnecting")
	})
	
	opts.SetDefaultPublishHandler(mqttReceive)
	
	opts.SetOnConnectHandler(func(c mqtt.Client) {
	    log.Printf("MQTT Client connected\n")
        //Subscribe here, otherwise after connection lost, you may not receive any message
        if token := c.Subscribe(fmt.Sprintf("stopwatch/cmd/+"), 0, nil); token.Wait() && token.Error() != nil {
            log.Println(token.Error())
            // TODO handle Error
        }
    })
	return opts
}

func mqttReceive(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Sheet:\tMQTT received [%s] %s\n", msg.Topic(), string(msg.Payload()))
	if strings.Contains(msg.Topic(),"stopwatch/cmd/sync") {
		var items map[string]interface{}		
		items = util.DecodePayload(msg.Payload())
		if (items["MASTER"] != nil) { 
			switch items["MASTER"].(type) {
			case float64 :
				masterOffset := int64(items["MASTER"].(float64)*1000)
				log.Println("Master Clock: ",masterOffset)
				masterMillis = masterMillis + masterOffset				
			}		
		}
	}
}

func publish(channel int, now time.Time, stamp int64) {
	stamp = stamp + masterMillis
	payload :=fmt.Sprintf(`{"Channel":%d,"Time":"%s","Stamp":%4.2f}`,channel,now.Format("15:04:05"),float64(stamp)/1000.0)
	mqttpipe.Send <- mqttpipe.Message{"stopwatch/data",payload}
}

func timestamp(now time.Time, stamp int64) string {
	stamp = stamp + masterMillis
	return fmt.Sprintf("(%s) %4.2fs",now.Format("15:04:05"),float64(stamp)/1000.0)
}

/*

https://stackoverflow.com/questions/71689285/how-to-get-monotonic-part-of-time-time-in-go

The monotonic clock is just used for differences between times. 
The absolute value of the monotonic clock is undefined and you 
should not try to get it. I think what you really want for your 
timestamp is the duration from a base time.

func init() {
    baseTime = time.Now()
}

// NowTimestamp returns really just the duration from the base time
func NowTimestamp() time.Duration {
    return time.Now().Sub(baseTime)
}

*/

func millis( clock time.Duration ) int64 {
	return int64(clock)/1000000
}

//func uptime() time.Duration {
//	return sysinfo.Get().Uptime	
//}

func uptime() time.Duration {
    return time.Now().Sub(baseTime)
}

func main() {	
	var (
		startTime 		time.Time
		finishTime 		time.Time
	)
	
	flag.StringVar(&mqttHost,		"mqtt", 	"//localhost:1883/","MQTT Host")
	flag.StringVar(&usbDevice,		"device", 	"/dev/ttyUSB0", 	"USB Device")
	flag.UintVar(&debounceStart,	"debsta", 	30, 				"Debounce Start")
	flag.UintVar(&debounceFinish,	"debfin", 	30, 				"Debounce Finish")
	flag.BoolVar(&enableButtons,	"buttons",  true, 				"Manaul Buttons")
	flag.BoolVar(&usbStopwatch,		"usb", 		true, 	    		"USB Adapter")
	flag.IntVar(&usbPollIntervall,	"poll", 	100, 				"USB Poll Intervall")
	flag.IntVar(&publishInterval,	"lock", 	300, 				"Publish Lock")
	
	flag.Parse()
	
	startLocked 	:= util.TimedLock()
	finishLocked 	:= util.TimedLock()
	
	// -------------- MQTT --------------
	uri, err := url.Parse(mqttHost)
	if err != nil {
		log.Fatal("MQTT: ",err)
	}
		
	mqttClient := mqttConnect("STOPWATCH_XX", uri)
	defer mqttClient.Disconnect(0)
	
	mqttpipe.Send = make(chan mqttpipe.Message, 100)
	defer close(mqttpipe.Send)
	
	go mqttpipe.Sender(mqttClient)
	
	tk.Title("Elzwelle USB Stopuhr")
	
	hostname,_ := os.Hostname()
	if hostname == "pady" {
		tk.Zoomed()
		//tk.Fullscreen()
		//tk.Geometry(0,0,800,400)
	} else {
		//tk.Zoomed()
		tk.Geometry(0,0,800,400)	
	}
	
	//startupTime		:= time.Now()	
	eventStartup 	:= time.Unix(0,0)
		
	//la := tk.Label(tk.ROOT,timestamp(startupTime,millis(time.Now().Sub(startupTime))),tk.FontSize(18))
	la := tk.Label(tk.ROOT,timestamp(time.Now(),millis(uptime())),tk.FontSize(18))
	la.Pack(tk.PACK_X)
	
	mainFrame := tk.Frame(tk.ROOT)
	mainFrame.Pack(tk.PACK_EXPAND,tk.PACK_BOTH)
	
	startFrame := tk.Frame(mainFrame)
	startFrame.Pack(tk.PACK_LEFT,tk.PACK_EXPAND,tk.PACK_BOTH)
		
	startList := tk.List(startFrame,tk.FontSize(12))
	startList.Pack(tk.PACK_BOTTOM,tk.PACK_EXPAND,tk.PACK_BOTH)
		
	startButton := tk.Button(startFrame,"Start",tk.FontSize(16))
	startButton.Pack(tk.PACK_TOP,tk.PACK_X)
	
	if enableButtons {	
		startButton.SetOnClick( func() {
				now := time.Now()
				//startList.ListAppend(timestamp(now, millis(now.Sub(startupTime))))	
				//publish(1,now,millis(now.Sub(startupTime)))
				startList.ListAppend(timestamp(now, millis(uptime())))	
				publish(1,now,millis(uptime()))
			},
		)
	}
	
	finshFrame := tk.Frame(mainFrame)
	finshFrame.Pack(tk.PACK_RIGHT,tk.PACK_EXPAND,tk.PACK_BOTH)
		
	finishList := tk.List(finshFrame,tk.FontSize(12))
	finishList.Pack(tk.PACK_BOTTOM,tk.PACK_EXPAND,tk.PACK_BOTH)
		
	finishButton := tk.Button(finshFrame,"Finish",tk.FontSize(16))
	finishButton.Pack(tk.PACK_TOP,tk.PACK_X)
	
	if enableButtons {		
		finishButton.SetOnClick( func() {
				now := time.Now()
				//finishList.ListAppend(timestamp(now, millis(now.Sub(startupTime))))		
				//publish(2,now,millis(now.Sub(startupTime)))
				finishList.ListAppend(timestamp(now, millis(uptime())))		
				publish(2,now,millis(uptime()))
			},
		)
	}
	
//	syncEntry := tk.Entry(tk.ROOT,"TEST",tk.FontSize(16))
//	syncEntry.Pack(tk.PACK_X)
//	
//	syncButton := tk.Button(tk.ROOT,"Sync",tk.FontSize(16))
//	syncButton.Pack(tk.PACK_X)
//	syncButton.SetOnClick( func() {
//			e,v := syncEntry.GetVariable()
//			if e == nil {
//				log.Println("Sync:",v)
//			}
//		},
//	)
	
	go notify.Listen()	
	
	//---------------------------- PI GPIO using Linux Kernel Driver----------------------------
	
	// Init GPIO, check for pi otherwise use only buttons or USB
	if _,err := os.Stat("/boot/config.txt"); err == nil {
		chip, err := gpio.OpenChip("gpiochip0")
		if err != nil {
			return
		}
		defer chip.Close()
		
		// GPIO 26 START Trigger
		line26, err := chip.OpenLineWithEvents(26, gpio.Input, gpio.FallingEdge, "gpio")
		if err != nil {
			return
		}
		defer line26.Close()	
		
		// GPIO 19 FINISH Trigger
		line19, err := chip.OpenLineWithEvents(19, gpio.Input, gpio.FallingEdge, "gpio")
		if err != nil {
			return
		}
		defer line19.Close()
		
		tick := make(chan int)
		
		go func() {
			for {
				time.Sleep(time.Duration(1000 * time.Millisecond))
				tick <- 0
			}
		}()
		//---------------------------- GPIO Inerrupt Worker Loop -----------------------------
		
		for {
			select {
				case event := <- line26.Events():
				now := time.Now()
				if event.Time.Sub(startTime) > time.Millisecond * time.Duration(debounceStart) {
					startTime = event.Time	
					if ! startLocked(publishInterval) {
						startList.ListAppend(timestamp(now, millis(event.Time.Sub(eventStartup))))					
						publish(1,now,millis(now.Sub(eventStartup)))	
					}
				}
				case event := <- line19.Events():
				now := time.Now()
				if event.Time.Sub(finishTime) > time.Millisecond * time.Duration(debounceFinish) {		
					finishTime = event.Time
					if ! finishLocked(publishInterval) {
						finishList.ListAppend(timestamp(now, millis(event.Time.Sub(eventStartup))))			
						publish(2,now,millis(now.Sub(eventStartup)))
					}		
				}
				case <- tick:
				now := time.Now()
				la.SetText(timestamp(now,millis(uptime())))
//				la.SetText(timestamp(now, millis(now.Sub(startupTime))))
//				stamp := timestamp(now, millis(now.Sub(startupTime)))
//				la.SetText(stamp[0:11]+"00)")
			}
		}
	} else {
		if usbStopwatch {
			usbNotify := notify.Notification()
			
			buf 	:= make([]byte, 1)
			loop 	:= true
			msg 	:= make([]byte, 20)
			idx		:= 0
			
			port, err := serial.Open(&serial.Config{
				Address:	usbDevice,
				BaudRate:	115200,
				DataBits:	8,
				Parity:		"N",
				StopBits:	1,
				Timeout:	100 * time.Millisecond,
			})
			if err != nil {
				usbNotify("ERROR Opend USB Serial Port!\n\nCheck USB Adapter.\nRestart App.",8000)
				tk.Done()
				log.Fatal(err)
			}
			defer port.Close()
			
			go func(port serial.Port) {
				for {
					time.Sleep(time.Duration(usbPollIntervall) * time.Millisecond)
					port.Write([]byte("$"))
				}
			}(port)
			
			for loop {
				n,err := port.Read(buf)

				if err != nil {
					//log.Printf("ERROR %v\n",err)
				} else {
					if n == 0 {
						usbNotify("ERROR read: EOF!\n\nCheck USB Adapter.\nRestart App.",8000)
					}
				}

				if n > 0 {
					switch buf[0] {
						case 0x0A:
						case 0x0D:
							if msg[0] == 'F' {
								now := time.Now()
								t,_ := strconv.Atoi(string(msg[1:idx]))
								usbFinishMillis = int64(t)
								if ! finishLocked(publishInterval) {
									finishList.ListAppend(timestamp(now, usbFinishMillis))
									publish(2,now,usbFinishMillis)
								}
							} else if msg[0] == 'S' {
								now := time.Now()
								t,_ := strconv.Atoi(string(msg[1:idx]))
								usbStartMillis = int64(t)
								if ! startLocked(publishInterval) {
									startList.ListAppend(timestamp(now, usbStartMillis))
									publish(1,now,usbStartMillis)	
								}
							} else if msg[0] == '@' {
								now := time.Now()
								t,_ := strconv.Atoi(string(msg[1:idx]))
								usbSyncEvent = (int64(t)/100)*100
								//log.Println(usbSyncEvent)
								la.SetText(timestamp(now, usbSyncEvent))
							}
							//log.Println(string(msg[0:idx]))
							//log.Printf("%v\n",msg[0])
							idx=0
						default:
							if idx < 20 {
								msg[idx] = buf[0]
								idx++	
							}				
					}
				}
			}
		}
	}
	
	tk.Done()
}

