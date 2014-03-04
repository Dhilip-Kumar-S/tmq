package web

import (
	"Q"
	"fmt"
	"net/http"
)

func QStatus(w http.ResponseWriter, r *http.Request) {

	qlist := Q.ListQ()
	fmt.Fprintf(w, "<h1>List of Queues</h1>")
	for i, qname := range qlist {
		fmt.Fprintf(w, "<h2>%d. %s</h2>", i+1, qname)
	}
}

func StartHTTP() {

	http.HandleFunc("/index.html", QStatus)
	http.ListenAndServe(":8080", nil)

}
