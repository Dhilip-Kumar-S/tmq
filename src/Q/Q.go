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
	"fmt"
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
	nodes  map[string]Q // Map of Q handle with key as Hash.hash
}

/*
 * A Global variable of all the Q
 * TODO: We need to make the init function read from a configuration file
 * and load a list of all active persistant Qs
 */

var root QTree

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

/* Will loop throught the root's map of Qs and returns a string array */
func ListQ() []string {

	var i int

	list := make([]string, len(root.nodes))

	for _, val := range root.nodes {
		list[i] = fmt.Sprintf("%s has %d messages, opened by %d clients", val.name,  val.l.Len(), val.ref)
		i++
	}

	return list
}

/* A string Converter of a Q message */
func (q QEle) String() string {
	return string(q.msg)
}

/* Code to perform Enqueu Operation */
func (q Q) EnQ(msg []byte, size int64) int {

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
func (q Q) DQ() (QEle, bool) {

	var tmsg QEle
	rc := false

	q.mutex.Lock()
	if q.l.Len() > 0 {
		e := q.l.Back()
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
	root.nodes = make(map[string]Q)
	root.rwlock.Unlock()
}

func MakeQID(name string) string {

	h := sha256.New()
	h.Write([]byte(name))
	id := hex.EncodeToString(h.Sum(nil))
	return id
}

func Create(name string, store bool) byte {

	var tmpQ Q
	var rc byte

	tmpQ.id = MakeQID(name)
	tmpQ.name = name
	tmpQ.Store = store
	tmpQ.ref = 0
	tmpQ.l = list.New()

	root.rwlock.Lock()
	_, ok := root.nodes[tmpQ.id]
	if ok == false {
		root.nodes[tmpQ.id] = tmpQ
		rc = 0x01
	} else {
		rc = 0x00
	}
	root.rwlock.Unlock()

	return rc
}

func Open(name string) (Q, bool) {

	id := MakeQID(name)
	root.rwlock.RLock()
	tQ, ok := root.nodes[id]
	if ok == true {
	
		root.nodes[id].mutex.Lock()
		tQ.ref++
		tQ.mutex.Unlock()
		//root.nodes[id] = tQ		
	
	}
	root.rwlock.RUnlock()

	return tQ, ok
}

func (q Q) Close() {
	q.mutex.Lock()
	q.ref--
	q.mutex.Unlock()
}
