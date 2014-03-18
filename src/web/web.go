package web

import (
	"Q"
	"fmt"
	"net/http"
)

func QData (w http.ResponseWriter, r *http.Request) {

	//Q.Trace (fmt.Printf, "request %v\n", http.Request)
	fmt.Printf ( "request %v\n", r.URL.Path)
	fmt.Printf ( "RawQuery %v\n", r.URL.RawQuery)
	fmt.Printf ( "Fragment %v\n", r.URL.Fragment)
	
}

func QStatus(w http.ResponseWriter, r *http.Request) {

	qlist := Q.ListQ()
	
	fmt.Fprintf(w, "<h1 align=\"center\">Queues</h1>")
	fmt.Fprintf (w, "<table align=\"center\" border=2> <tr align=\"center\"> <th> Name </th> <th> Len </th> <th> Opened </th></tr>\n")
	for _, tque := range (qlist) {
		fmt.Fprintf(w, "<tr align=\"center\"> <td><a href=\"Q.html?%s\">%s</a></td> <td>%d</td> <td>%d</td>\n", tque.Id(), tque.Name(), tque.Len(), tque.Ref())
	}
	
	fmt.Fprintf (w, "</table>")
	
}

func StartHTTP(wport string) {

	http.HandleFunc ("/index.html", QStatus)
	http.HandleFunc ("/Q.html", QData)
	http.ListenAndServe(":"+wport, nil)

}
