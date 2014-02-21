/*
 * This package provides all the related functions and wrappers for the basic Q operation
 * Here is how this will work.
 * 1) Each Q will be reffered by a handle.
 * 2) Each created handle will be stored in a map
 * 3) The hash map will be
 */

package Q

import (
	"container/list"
	"crypto/sha256"
	"sync"
	"encoding/hex"
)

type QTree struct {
	rwlock  sync.RWMutex   // On the hash map
	nodes map[string]Q // Map of Q handle with key as Hash.hash
}

var root QTree

type Q struct {
	id        string
	mutex     sync.Mutex //
	name      string     // Name of the queue
	ref	  uint64     // Number of users currently opend the queu
	Store 	bool       // Persistant Queu store in the disk
	len       uint64     // Total number of elements in the Queue
	l         *list.List      // Actual Queue values
}

type QEle struct {
	msg []byte
	len int64
}

func (q QEle) String() string {
	return string(q.msg)
}


func (q Q) EnQ(msg []byte, size int64) int {

	tmsg := QEle {msg, size}
	q.mutex.Lock()
	e := q.l.PushBack(tmsg)
	q.len++
	q.mutex.Unlock()

	if e != nil {
		return 1
	}
	return 0
}

func (q Q) DQ() (QEle, bool) {

	var tmsg QEle
	rc := false
	
	q.mutex.Lock()
	if q.l.Len() > 0 {
		e := q.l.Back()
		if e != nil {
			tmsg = q.l.Remove(e).(QEle)
			q.len--
			rc = true
		}
	}
	q.mutex.Unlock()
	
	return tmsg, rc
}

func Init () {
	root.rwlock.Lock()
	root.nodes = make (map[string]Q)
	root.rwlock.Unlock()
}

func MakeQID (name string) string {

	h := sha256.New()
	h.Write ([]byte(name))
	id := hex.EncodeToString(h.Sum(nil))
	return id
}

func Create(name string, store bool) int {
	
	var tmpQ Q
	var rc int
	
	tmpQ.id = MakeQID (name)
	tmpQ.name = name
	tmpQ.Store = store
	tmpQ.ref = 0
	tmpQ.len = 0
	tmpQ.l = list.New()
	
	root.rwlock.Lock()
	_ , ok :=  root.nodes[tmpQ.id]
	if  ok == false {
		root.nodes[tmpQ.id] =  tmpQ
		rc = 0
	} else {
		rc  = 1
	}
	root.rwlock.Unlock()	
	
	return rc
}

func Open (name string) (Q , bool) {
		
	root.rwlock.RLock ()
	tQ , ok := root.nodes[MakeQID(name)]
	root.rwlock.RUnlock ()
	
	return tQ , ok
	
}
