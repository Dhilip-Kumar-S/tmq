tmq (Transactional Message Queue)
===

A Simple TCP, Multi-threaded Message Queue written in  golang

***********************************

tmq is Transactional Message Queue.  Once TransStart is seen this buffers all reads / DeQueue locally until TransEnd is seen.  If the Connection fails in the middle or if the client issues a TransAbort intentionally all the buffered Reads or pushed-back into their respective queues. 

'tmq.exe' in windows and 'tmq' in Linux are the final usable binaries.  Have tried my best to follow google go's recommendation on how to organise go source code.  
  To get started you simply have to do below under a desired directory (of-course you need go and git installed in the intended environment)
  
  ```$git pull https://github.com/Dhilip-Kumar-S/tmq.git  
  $cd tmq  
  $GOPATH=${PWD} go build tmq.go  ```


that is all.  Simples :-)

*********************************************************************************************************************************
by default 'tmq' starts the tcp server at port 6161 and web server that displays statistical information at port 6060

  `D:\github\tmq>tmq -v`  
  `2014/03/18 23:02:22 Web server is up http://localhost:6060`  
  `2014/03/18 23:02:22 Q TCP server started at port:6161`  


this can be changed 

  `D:\github\tmq>tmq -port=9191 -web=9090 -v`  
  `2014/03/18 23:04:10 Web server is up http://localhost:9090`  
  `2014/03/18 23:04:10 Q TCP server started at port:9191`  

-v option enables verbose mode to print each and every activity of the message queue server.


Supported Queue Operations:
===
S.NO | Operation
:--:|:---:
  1| CREATE Queue, QName 
  2| OPEN Queue, QName 
  3| EnQueue, mq-id   
  4| DeQueue, mq-id 
  5| CLOSE Queue, mq-id 
  6| DELETE Queu, QName 
  7| SELECT Queue, mq-id  ==> return the number of elements 
  8| Transaction Start 
  9| Transaction End 
  10| Transaction Abort 

Each queue operation is atomic. 



ToDo:
===
  1) Improve the protocol [here] (Available in https://github.com/Dhilip-Kumar-S/tmq/blob/master/doc/TMQ_Protocol_Definition.pdf)  
  2) Make Queues persistent.  
  3) Distributed Message Queue servers.  


Note: There is an C client library under development will be included in the project as soon as its ready.
