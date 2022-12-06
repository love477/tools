package demo

import (
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type EventType int

const (
	Update EventType = iota
	Create
	Delete
)

func (e EventType) String() string {
	return [...]string{"Update", "Create", "Delete"}[e]
}

type Event struct {
	Key       string
	EventType EventType
	Obj       interface{}
}

func HandleAndQueueEvent(eventType EventType, obj interface{}, queue workqueue.RateLimitingInterface) *Event {
	// var err error
	queueQbj := &Event{
		EventType: eventType,
		Obj:       obj,
	}
	queueQbj.Key = strconv.Itoa(obj.(int))
	queue.Add(queueQbj)
	// queueQbj.Key, err = cache.MetaNamespaceKeyFunc(obj)
	// if err == nil {
	// 	queue.Add(queueQbj)
	// } else {
	// 	klog.Errorf("can not get object key for %v", err.Error())
	// }

	return queueQbj
}

type Controller struct {
	ResourceType string
	Queue        workqueue.RateLimitingInterface
	EventHandler func(event *Event) error
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// let the workers stop when we are done
	defer c.Queue.ShutDown()
	klog.Infof("Starting %s controller", c.ResourceType)

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Infof("stop controller for %s", c.ResourceType)
}

func (c *Controller) runWorker() {
	// processNextWorkItem will automatically wait until there's work available
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.Queue.Get()
	if quit {
		return false
	}

	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed parallel
	defer c.Queue.Done(key)

	err := c.EventHandler(key.(*Event))

	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of an outdated
		// error history
		c.Queue.Forget(key)
	} else if c.Queue.NumRequeues(key) < 5 {
		// not reach the maximum retry count, should retry
		klog.Infof("Error sync %s %s, retry", c.ResourceType, key)
		c.Queue.AddRateLimited(key)
	} else {
		c.Queue.Forget(key)
		klog.Infof("Drop %s %s because maximum retry has reached", c.ResourceType, key)
	}

	return true
}
