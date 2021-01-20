package idevice

/*
#cgo pkg-config: libimobiledevice-1.0
#include <libimobiledevice/libimobiledevice.h>
#include <stdio.h>

extern void eventCallback();

static void event_callback(const idevice_event_t *event, void *user_data)
{
	eventCallback(event);
}

static idevice_error_t event_subscribe()
{
	return idevice_event_subscribe(event_callback, NULL);
}
*/
import "C"

import "sync"

// Event provides information about the occurred event.
type Event struct {
	Type           uint32
	UDID           string
	ConnectionType uint32
}

// The event type for device add or removal.
const (
	EventAdd uint32 = iota
	EventRemove
	EventPaired
)

var (
	eventChannel     chan Event
	eventChannelLock sync.Mutex
)

//export eventCallback
func eventCallback(event *C.idevice_event_t) {
	eventChannelLock.Lock()
	defer eventChannelLock.Unlock()

	eventChannel <- Event{
		Type:           event.event,
		UDID:           C.GoString(event.udid),
		ConnectionType: event.conn_type,
	}
}

// EventSubscribe returns a channel that receives IDevice events. If a channel
// has already been allocated, the old one is closed.
func EventSubscribe() (ch <-chan Event, err error) {
	eventChannelLock.Lock()
	defer eventChannelLock.Unlock()

	if eventChannel != nil {
		if e := EventUnsubscribe(); e != nil {
			return nil, e
		}
	}

	eventChannel = make(chan Event)
	defer func() {
		if err != nil {
			eventChannel = nil
		}
	}()

	if e := iDeviceError(C.event_subscribe()); e != nil {
		return nil, e
	}

	return eventChannel, nil
}

// EventUnsubscribe closes the channel returned by SubscribeEvent and cleans up
// resources.
func EventUnsubscribe() error {
	eventChannelLock.Lock()
	defer eventChannelLock.Unlock()

	if eventChannel != nil {
		if err := iDeviceError(C.idevice_event_unsubscribe()); err != nil {
			return err
		}

		close(eventChannel)
		eventChannel = nil
	}
	return nil
}
