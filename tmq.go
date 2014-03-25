package main 

import (
"log"
"Q"
"web"
"os"
"runtime"
"flag"
"runtime/pprof"
"os/signal"
)



func main () {
	
	if len (os.Args) < 2 { 
		log.Fatal("Not enough argument")
		return
	}
	
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	var tport, wport, pfile string
	
	Q.Init()
	flag.BoolVar (&(Q.IsVerbose), "v", false, "User This for Trace message")
	flag.StringVar ( &tport, "port", "6161", "Port number to start the tcp server")
	flag.StringVar ( &wport, "web", "6060", "Port number of the web server")
	flag.StringVar ( &pfile, "proffile", "", "Profiling filename")
	flag.Parse ()

	qsig := make (chan os.Signal, 1)
	signal.Notify (qsig, os.Interrupt, os.Kill)
	
	if pfile != "" {
		/* This means we want to profile the server */
		pfd, err := os.Create(pfile)
		if err != nil {
			log.Fatal ("unaable to create the profile file:", err)
		}
		pprof.StartCPUProfile (pfd)
		defer pprof.StopCPUProfile ()
	}
	
	go web.StartHTTP (wport)
	Q.Trace (log.Printf, "Web server is up http://localhost:%s\n", wport)
	Q.Trace (log.Printf, "Q TCP server started at port:%s\n",tport)
	go Q.StartTCP (tport)
	
	gsig := <-qsig
	
	log.Println ("Q finished signal:",gsig)
		
}
