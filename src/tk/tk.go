package tk

import (
	"log"
	"fmt"
	"bufio"
	"os"
    "os/exec"
    "io"
)

type Option struct {
    Name 	string
    Value	interface{}
}

var Command chan string
var Event   chan string
var Result	chan string
 
//var piWindow = string(`
//	source /usr/share/tcltk/ttkthemes/themes/pkgIndex.tcl
//	source /usr/share/tcltk/ttkthemes/png/pkgIndex.tcl
//	
//	namespace /usr/share/tcltk/import ttk::*
//	ttk::setTheme breeze
//	
//	wm attributes . -fullscreen 1
//`)

var frameWindow = string(`
	
#	source "azure.tcl"
#	set_theme dark
#   ------- ttkthemes -------
#	source /usr/share/tcltk/ttkthemes/themes/pkgIndex.tcl
#	source /usr/share/tcltk/ttkthemes/png/pkgIndex.tcl
#	package require ttkthemes 1.0
	
	namespace import ttk::*
	
#	option add *Font {Helvetica 18} widgetDefault
#	ttk::setTheme breeze

#	wm attributes . -fullscreen 1
#	wm geometry . 800x800+0+0

#	wm title . "Test"
`)

func options(args []Option) (string) {
	res := ""
	
	for _,item := range args {
		res = res + " -" + item.Name + " " + item.Value.(string)
	}
	return res
}

var wish *exec.Cmd

func init() {
	// Execute wish tcl/tk window shell
    wish = exec.Command("wish")

	Command = make(chan string, 100)
	Event   = make(chan string, 100)
	Result  = make(chan string)
	
	// Get stdout pipe
    stdout, err := wish.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }

	// Get stderr pipe
	stderr, err := wish.StderrPipe()
    if err != nil {
        log.Fatal(err)
    }
    
    // Get stdin pipe
	stdin, err := wish.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }

	// Catch errors from wish 
	go func() {
        defer stderr.Close()
        _, err := io.Copy(os.Stderr, stderr)

	    if err != nil {
	        log.Fatal(err)
	    }
    }()
	
	// Start wish
    if err := wish.Start(); err != nil {
        log.Fatal(err)
    }
	
	go eventHandler()
	
//	// Write tcl commands to wish stdin pipe, check for pi using fullscreen
//	if _,err := os.Stat("/boot/config.txt"); err == nil {
//		io.WriteString(stdin, piWindow)
//	} else {
//		io.WriteString(stdin, frameWindow)
//	}
	io.WriteString(stdin, frameWindow)
	// Output pipe event loop
	go func() {
		defer stdout.Close()
		// New buffered io reader on wish stdout
		buf := bufio.NewReader(stdout) 		
	    for {
	        line, _, err := buf.ReadLine()
	        //line, err := buf.ReadString('\n')	        
	        if err != nil {
		        log.Fatal(err)
		    }
		    //log.Println("ReadLine:\t",string(line))
		    Event <- string(line)
	    }	    
	}()
	
	// Input pipe command loop
	go func() {
		defer stdin.Close()
		for {
			command := <- Command
			//log.Println("Command:\t",string(command))
			io.WriteString(stdin, command + "\n")
		}		
	}()
	
}

func Done() {
	// Wait until done
    if err := wish.Wait(); err != nil {
        log.Fatal(err)
    }
}

//	wm attributes . -fullscreen 1
func Fullscreen() {
	Command <- "wm attributes . -fullscreen 1"
}

//	wm geometry . 800x800+0+0
func Geometry(x int, y int, width int, height int) {
	Command <- fmt.Sprintf("wm geometry . %dx%d+%d+%d",width,height,x,y)
}

// wm attributes . -zoomed 1
func Zoomed() {
	Command <-"wm attributes . -zoomed 1"
}

//	wm title . "Test"
func Title(title string) {
	Command <- fmt.Sprintf("wm title . \"%s\"",title)
}

