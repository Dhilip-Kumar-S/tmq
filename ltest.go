package main

import (
"container/list"
"fmt"
)

type  p struct{
	name string
	age		int
}

func main () {

	d := p {"dhilip", 31}
	
	l := list.New ()
	
	l.PushBack (1)
	l.PushBack (2)
	l.PushBack ("hi")
	l.PushBack ("how")
	l.PushBack ("Are?")
	l.PushBack (d)
	l.PushBack (3)
	
	fmt.Println (l.Len())
	
	for i:= l.Front (); i != nil;  {
		
		
		pi := i
		i = i.Next ()
		w := l.Remove (pi)
		fmt.Println (w)
	}
	fmt.Println (l.Len())
}