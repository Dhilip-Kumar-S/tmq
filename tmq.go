package main 

import (
"log"
"Q"
"web"
"os"
)



func main () {
	
	if len (os.Args) < 2 { 
		log.Fatal("Not enough argument")
		return
	}
	Q.Init()
	portstr := os.Args[1]
	go web.StartHTTP ()
	log.Println ("Web server is up http://localhost:6060")
	log.Println ("Q TCP server started at:",portstr)
	Q.StartTCP (portstr)
	log.Println ("Q finished")
		
}