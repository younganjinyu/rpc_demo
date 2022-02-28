package common

import (
	"fmt"
	"net/rpc"
	"reflect"
	"sync"
)

type Args struct {
	Level          int
	Param1, Param2 int
}

type Func struct {
	name string
	args []interface{}
	rtn  []interface{}
}

type (
	Notifier interface {
		Register(client *rpc.Client)
		Deregister(client *rpc.Client)
		// notify new events to observers
		Notify()
	}

	Observer interface {
		OnNotify()
	}
)

type ClientsNotifier struct {
	sync.RWMutex
	observers []*rpc.Client
}

func (m *ClientsNotifier) Register(o *rpc.Client) {
	m.Lock()
	m.observers = append(m.observers, o)
	//fmt.Println("here")
	fmt.Println(len(m.observers))
	m.Unlock()
}

func (m *ClientsNotifier) Deregister(o *rpc.Client) {
	m.RLock()
	pos := -1
	for i := range m.observers {
		if reflect.DeepEqual(o, m.observers[i]) {
			pos = i
			break
		}
	}
	m.RUnlock()
	if pos != -1 {
		m.Lock()
		m.observers = append(m.observers[0:pos], m.observers[pos+1:]...)
		m.Unlock()
	}
}

func (m *ClientsNotifier) Notify() {
	//fmt.Println("Notifier notifies all the observers")
	nums := len(m.observers)
	fmt.Println(nums)
}

func NewClientsNotifier() *ClientsNotifier {
	return &ClientsNotifier{
		observers: make([]*rpc.Client, 0),
	}
}

type ObserverClient struct {
}
