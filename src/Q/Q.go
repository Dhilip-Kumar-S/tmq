/*
 * This package provides all the related functions and wrappers for the basic Q operation
 * Here is how this will work.
 * 1) Each Q will be reffered by a handle.
 * 2) Each created handle will be stored in a map
 * 3) The hash map will be
 * 4) A Handle is a 32-byte Sha256 message digest created using the Q's name
 */

package Q

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	//"fmt"
)

/*
 ******************************************************************************
 ********************** Type Declaration **************************************
 ******************************************************************************
 */

/*
 * This is the main structure which holds all the Q
 * When we open a Q we basically lookup this map and returns its corresponding Q structure.
 *
 */
type QTree struct {
	rwlock sync.RWMutex // Read/Write
	nodes  map[string]*Q // Map of Q handle with key as Hash.hash
}

/*
 * A Global variable of all the Q
 * TODO: We need to make the init function read from a configuration file
 * and load a list of all active persistant Qs
 */

var root QTree
var IsVerbose bool

type Q struct {
	id    string
	mutex sync.Mutex // Lock the Q
	name  string     // Name of the queue
	ref   uint64     // Number of users currently opend the queu
	Store bool       // Persistant Queu store in the disk
	l     *list.List // Actual Queue values
}

type QEle struct {
	msg []byte
	len int64
}

/*
 ******************************************************************************
 ***********************Function Definitions **********************************
 ******************************************************************************
 */

/* This function takes a function as the first argumetn adn use it to print the format and the string */
func Trace (fn  func (string, ...interface{}), fmt string, v ...interface{}) {
	if IsVerbose {
		fn (fmt , v...)
	}
}


func ListQ (id string) (int, []QEle, string) {
	
	var len int
	var rlist []QEle
	var qname string
	
	tQ := GetQ (id)
	if tQ != nil {
		len = tQ.l.Len()
		rlist = make([]QEle, len)
		i := 0
		for e := tQ.l.Front(); e != nil; e = e.Next() {
			rlist[i] = e.Value.(QEle)
			i++
		}
		return i, rlist, tQ.name
		
	} else {
	
		return -1, rlist, qname
		
	}
	
}
 
/* Will loop throught the root's map of Qs and returns a array of Queues */
func ListAll() []*Q {

	var i int

	rlist := make([]*Q, len(root.nodes))

	for _, val := range root.nodes {
		rlist[i] = val
		i++
	}

	return rlist
}

/* A string Converter of a Q message */
func (q QEle) String() string {
	return string(q.msg)
}

func (q *Q) Name () string {
	return q.name
}

func (q *Q) Len () int {
	return q.l.Len()
}

func (q *Q) Id () string {
	return q.id
}

func (q *Q) Ref () int {
	return int(q.ref)
}

/* Code to perform Enqueu Operation */
func (q *Q) EnQ(msg []byte, size int64) int {

	tmsg := QEle{msg, size}
	q.mutex.Lock()
	e := q.l.PushBack(tmsg)
	q.mutex.Unlock()

	if e != nil {
		return 1
	}
	return 0
}

/*
 * DQ Message returns QElement, bool
 */
func (q *Q) DQ() (QEle, bool) {

	var tmsg QEle
	rc := false

	q.mutex.Lock()
	if q.l.Len() > 0 {
		e := q.l.Front()
		if e != nil {
			tmsg = q.l.Remove(e).(QEle)
			rc = true
		}
	}
	q.mutex.Unlock()

	return tmsg, rc
}

func Init() {
	root.rwlock.Lock()
	root.nodes = make(map[string]*Q)
	root.rwlock.Unlock()
}

func MakeQID(name string) string {

	h := sha256.New()
	h.Write([]byte(name))
	id := hex.EncodeToString(h.Sum(nil))
	return id
}

func GetQ (id string) *Q {
	root.rwlock.RLock()
	tQ, ok := root.nodes[id]
	root.rwlock.RUnlock()
	if ok {
		return tQ
	} 
	return nil
}

func Create(name string, store bool) byte {

	var tmpQ *Q
	var rc byte

	tmpQ = new(Q)
	tmpQ.id = MakeQID(name)
	tmpQ.name = name
	tmpQ.Store = store
	tmpQ.ref = 0
	tmpQ.l = list.New()

	root.rwlock.Lock()
	_, ok := root.nodes[tmpQ.id]
	if ok == false {
		root.nodes[tmpQ.id] = tmpQ
		rc = 0x00
	} else {
		rc = 0x01
	}
	root.rwlock.Unlock()

	return rc
}

func Open(name string) (*Q, bool) {

	id := MakeQID(name)
	root.rwlock.RLock()
	tQ, ok := root.nodes[id]
	if ok == true {
	
		tQ.mutex.Lock()
		tQ.ref++
		tQ.mutex.Unlock()
		//root.nodes[id] = tQ		
	
	}
	root.rwlock.RUnlock()

	return tQ, ok
}

func (q *Q) Close() byte{

	q.mutex.Lock()
	cnt := q.ref
	if q.ref > 0 {
		q.ref--
	}
	q.mutex.Unlock()

	if cnt <= 0 {
		return 0x01
	} else {
		return 0x00
	}

}
