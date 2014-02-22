package main

import (
	"Q"
	"fmt"
	"strconv"
)

func AvailableQ() {

	fmt.Println("************List of Queues are*******************")
	qlist := Q.ListQ()
	for _, qname := range qlist {
		fmt.Println(qname)
	}
	fmt.Println("*************************************************")

}

func CQ() {
	var name string
	fmt.Printf("Enter Q name:")
	fmt.Scanf("%s", name)
	if Q.Create(name, false) == 0 {
		fmt.Printf("Q %s Created\n", name)
	} else {
		fmt.Printf("Q already exisit\n")
	}
}

func Open() {
	var name string
	fmt.Printf("Enter Q name:")
	fmt.Scanf("%s", name)

	tmpQ, ok := Q.Open(name)

	if ok == true {
		fmt.Printf("Opened Sucessfully")
		SubQOptions(tmpQ)
	} else {
		fmt.Printf("Queue Does not exisist\n")
	}

}

func main() {
	var q1, q2, q3 Q.Q
	var ok bool
	var tele Q.QEle
	Q.Init()

	fmt.Println("Create Queues")
	Q.Create("One", false)
	Q.Create("Two", false)
	Q.Create("Three", false)

	fmt.Println("Open Queues")
	q1, ok = Q.Open("One")
	q2, ok = Q.Open("Two")
	q3, ok = Q.Open("Three")

	fmt.Println("EnQueue")

	for i := 1; i <= 30; i++ {
		msg := "RECORD:" + strconv.Itoa(i)
		if i%3 == 0 {
			q1.EnQ([]byte(msg), int64(len(msg)))
		} else if i%2 == 0 {
			q2.EnQ([]byte(msg), int64(len(msg)))
		} else {
			q3.EnQ([]byte(msg), int64(len(msg)))
		}

	}

	fmt.Println("DQ Q1")
	ok = true
	for ok == true {
		tele, ok = q1.DQ()
		fmt.Println(tele.String())
	}

	fmt.Println("DQ Q2")
	ok = true
	for ok == true {
		tele, ok = q2.DQ()
		fmt.Println(tele.String())
	}

	fmt.Println("DQ Q3")
	ok = true
	for ok == true {
		tele, ok = q3.DQ()
		fmt.Println(tele.String())
	}

}
