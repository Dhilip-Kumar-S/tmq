package web

import (
	"Q"
	"fmt"
	"net/http"
)

func QData (w http.ResponseWriter, r *http.Request) {
	
	fmt.Printf ( "RawQuery %v\n", r.URL.RawQuery)
	nEle, arryEle, qname := Q.ListQ (r.URL.RawQuery)
	
	switch (nEle) {
		case -1:
			fmt.Fprintf(w, "<h1 align=\"center\" color=\"red\"> Sorry, such Queue with id = %s not available... </h1>", qname)
		break
		
		case 0:
			fmt.Fprintf (w, "<h1 align=\"center\">%s</h1> <table border=2> <tr> <th> FRONT </th> <td> ... </td> <td> ... </td> <td> ... </td><th> END </th> </tr> </table>", qname )
		break
		
		default:
			fmt.Fprintf (w, "<h1 align=\"center\">%s</h1> <table border=2> <tr> <th> FRONT </th>", qname)
			for _, val := range (arryEle) {
				fmt.Fprintf (w, "<td> <b>&lt;====</b> </td> <td> %v </td>", val)
			}
			fmt.Fprintf (w, "<th> <b>&lt;====</b> </th> <th> END </th> </tr> </table>")
		break
	
	}
	
	
}

func QStatus(w http.ResponseWriter, r *http.Request) {

	qlist := Q.ListAll()
	
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
