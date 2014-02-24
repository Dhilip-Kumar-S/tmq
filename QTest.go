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
	fmt.Scanf("%s", &name)
	fmt.Println ("Creating Q with name:" + name)
	if Q.Create(name, false) == 0 {
		fmt.Printf("Q %s Created\n", name)
	} else {
		fmt.Printf("Q already exisit\n")
	}
}

func Open() {
	var name string
	fmt.Printf("Enter Q name:")
	fmt.Scanf("%s", &name)

	tmpQ, ok := Q.Open(name)

	if ok == true {
		fmt.Printf("Opened Sucessfully")
		SubQOptions(tmpQ)
	} else {
		fmt.Printf("Queue Does not exisist\n")
	}
}

func SubQOptions(tmpQ Q.Q) {

	var opt int
	var msg string
	fmt.Printf("1.EnQ\n2.DeQ\nChoice:")
	fmt.Scanf("%d", &opt)

	if opt == 1 {
		fmt.Printf("Enter message\n")
		fmt.Scanf("%d", &msg)
		tmpQ.EnQ([]byte(msg), int64(len(msg)))
		fmt.Println ("Enqued")
	} else {
		fmt.Printf("Message in the Queue is:")
		val, ok := tmpQ.DQ()
		if ok {
			fmt.Println(val.String())
		} else {
			fmt.Println ("Q Empty")
		}

	}

}

func main() {
	var q1, q2, q3 Q.Q
	var ok bool
	var tele Q.QEle
	Q.Init()
	
	isloop := true
	

	for isloop == true {
		var opt int
		AvailableQ()
		fmt.Printf ("1.CREATE\n2.OPEN\n3.EXIT\nEnter Option:")
		fmt.Scanf ("%d", &opt)
		switch opt {
			case 1:
				CQ()
				break;
			case 2:
				Open()
				break;
			case 3:
				isloop = false
				break
		}
	}
	
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
