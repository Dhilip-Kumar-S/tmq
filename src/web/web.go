package web

import (
	"Q"
	"fmt"
	"net/http"
)

func QStatus(w http.ResponseWriter, r *http.Request) {

	qlist := Q.ListQ()
	fmt.Fprintf(w, "<h1 align=\"center\">Queues</h1>")
	for i, qname := range qlist {
		fmt.Fprintf(w, "<h2 align=\"center\">%d. %s</h2>", i+1, qname)
	}
	
}

func StartHTTP(wport string) {

	http.HandleFunc("/index.html", QStatus)
	http.ListenAndServe(":"+wport, nil)

}
