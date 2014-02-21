package main

import (
"Q"
"fmt"
"strconv"
)


func main () {
	var q1, q2 Q.Q
	var ok bool
	var tele Q.QEle
	Q.Init()
	
	fmt.Println("Create Queues")
	Q.Create ("First", false)
	Q.Create ("Second", false)
	
	fmt.Println("Open Queues")
	q1,ok = Q.Open ("First")
	q2,ok = Q.Open ("Second")
	
	fmt.Println("EnQueue")
	
	for i:=0; i<10; i++ {
		msg  := "RECORD:"+strconv.Itoa(i)
		if (i%2 == 0) {
			q1.EnQ([]byte(msg), int64(len(msg)))
		} else {
			q2.EnQ([]byte(msg), int64(len(msg)))
		}
	}
	
	fmt.Println ("DQ Q1")
	ok = true
	for ok==true {
		tele , ok = q1.DQ()
		fmt.Println(tele.String())
	}
	fmt.Println ("DQ Q2");
	ok = true
	for ok==true {
		tele , ok = q2.DQ()
		fmt.Println(tele.String())
	}
}