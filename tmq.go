package main 

import (
"log"
"Q"
"web"
"os"
"runtime"
"flag"
)



func main () {
	
	if len (os.Args) < 2 { 
		log.Fatal("Not enough argument")
		return
	}
	
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	var tport, wport string
	
	Q.Init()
	flag.BoolVar (&(Q.IsVerbose), "v", false, "User This for Trace message")
	flag.StringVar ( &tport, "port", "6161", "Port number to start the tcp server")
	flag.StringVar ( &wport, "web", "6060", "Port number of the web server")
	flag.Parse ()
	
	
	go web.StartHTTP (wport)
	Q.Trace (log.Printf, "Web server is up http://localhost:%s\n", wport)
	Q.Trace (log.Printf, "Q TCP server started at port:%s\n",tport)
	Q.StartTCP (tport)
	log.Println ("Q finished")
		
}