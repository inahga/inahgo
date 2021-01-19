// Package imobiledevice provides CGo bindings for libimobiledevice.
package imobiledevice

/*
#cgo pkg-config: libimobiledevice-1.0
#include <stdlib.h>
#include <stdint.h>
#include <libimobiledevice/libimobiledevice.h>

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

import (
	"errors"
	"sync"
	"unsafe"
)

var (
	ErrIDeviceInvalidArgs   = errors.New("invalid argument")
	ErrIDeviceUnknown       = errors.New("unknown error")
	ErrIDeviceNoDevice      = errors.New("no device")
	ErrIDeviceNotEnoughData = errors.New("not enough data")
	ErrIDeviceSSLError      = errors.New("SSL error")
	ErrIDeviceTimeout       = errors.New("timeout")
	ErrIDeviceUnspecified   = errors.New("unspecified error")
)

func iDeviceError(ierr C.idevice_error_t) error {
	switch ierr {
	case 0:
		return nil
	case -1:
		return ErrIDeviceInvalidArgs
	case -2:
		return ErrIDeviceUnknown
	case -3:
		return ErrIDeviceNoDevice
	case -4:
		return ErrIDeviceNotEnoughData
	case -6:
		return ErrIDeviceSSLError
	case -7:
		return ErrIDeviceTimeout
	default:
		return ErrIDeviceUnspecified
	}
}

const (
	IDeviceLookupUsbmux uint32 = iota << 1
	IDeviceLookupNetwork
	IDeviceLookupPreferNetwork
)

const (
	ConnectionUsbmuxd uint32 = iota
	ConnectionNetwork
)

type IDevice struct {
	idevice C.idevice_t
}

// GetIDeviceList gets a list of UDIDs of currently available devices (USBMUX
// devices only).
func GetIDeviceList() ([]string, error) {
	var (
		count C.int
		clist **C.char
		ret   []string
	)
	if err := iDeviceError(C.idevice_get_device_list(&clist, &count)); err != nil {
		return nil, err
	}
	defer C.idevice_device_list_free(clist)

	golist := (*[1 << 30]*C.char)(unsafe.Pointer(clist))[:count:count]
	for _, ptr := range golist {
		ret = append(ret, C.GoString(ptr))
	}
	return ret, nil
}

// NewIDevice creates an IDevice structure for the device specified by UDID, if
// the device is available (USBMUX devices only). If you need to connect to a
// device available via network, use NewIDeviceWithOptions() and include
// IDeviceLookupNetwork in options.
func NewIDevice(udid string) (*IDevice, error) {
	return NewIDeviceWithOptions(udid, 0)
}

// NewIDeviceWithOptions creates an IDevice structure for the device specified
// by UDID, if the device is available, with the given lookup options.
//
// Options specifies what connection types should be considered when looking up
// devices. Accepts bitwise or'ed values of idevice_options. If 0 (no option) is
// specified it will default to IDeviceLookupUsbmux. To lookup both USB and
// network-connected devices, pass IDeviceLookupUsbmux | IDeviceLookupNetwork.
// If a device is available both via USBMUX and network, it will select the USB
// connection. This behavior can be changed by adding IDeviceLookupPreferNetwork
// to the options in which case it will select the network connection.
//
// To select the first available device, pass an empty string to udid.
func NewIDeviceWithOptions(udid string, options uint32) (*IDevice, error) {
	var (
		ret   IDevice
		cudid *C.char
	)

	if udid == "" {
		cudid = nil
	} else {
		cudid := C.CString(udid)
		defer C.free(unsafe.Pointer(cudid))
	}

	if err := iDeviceError(C.idevice_new_with_options(
		&ret.idevice, cudid, options)); err != nil {
		return nil, err
	}
	return &ret, nil
}

// Close cleans up the structure, then frees the structure itself.
func (i *IDevice) Close() error {
	if err := iDeviceError(C.idevice_free(i.idevice)); err != nil {
		return err
	}
	return nil
}

// Handle gets the usbmux device id of the device.
func (i *IDevice) Handle() (uint32, error) {
	var ret C.uint32_t
	if err := iDeviceError(C.idevice_get_handle(i.idevice, &ret)); err != nil {
		return 0, err
	}
	return uint32(ret), nil
}

// UDID gets the unique ID for the device.
func (i *IDevice) UDID() (string, error) {
	var ptr *C.char
	if err := iDeviceError(C.idevice_get_udid(i.idevice, &ptr)); err != nil {
		return "", err
	}
	defer C.free(unsafe.Pointer(ptr))
	return C.GoString(ptr), nil
}

// IDeviceEvent provides information about the occurred event.
type IDeviceEvent struct {
	Type           uint32
	UDID           string
	ConnectionType uint32
}

var (
	eventChannel     chan IDeviceEvent
	eventChannelLock sync.Mutex
)

const (
	IDeviceEventAdd uint32 = iota
	IDeviceEventRemove
	IDeviceEventPaired
)

//export eventCallback
func eventCallback(event *C.idevice_event_t) {
	eventChannelLock.Lock()
	defer eventChannelLock.Unlock()

	eventChannel <- IDeviceEvent{
		Type:           event.event,
		UDID:           C.GoString(event.udid),
		ConnectionType: event.conn_type,
	}
}

// SubscribeIDeviceEvent returns a channel that receives IDevice events.
func SubscribeIDeviceEvent() (ch <-chan IDeviceEvent, err error) {
	eventChannelLock.Lock()
	defer eventChannelLock.Unlock()

	if eventChannel != nil {
		if e := UnsubscribeIDeviceEvent(); e != nil {
			return nil, e
		}
	}

	eventChannel = make(chan IDeviceEvent)
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

// UnsubscribeIDeviceEvent closes the channel returned by SubscribeIDeviceEvent,
// and cleans up resources.
func UnsubscribeIDeviceEvent() error {
	eventChannelLock.Lock()
	defer eventChannelLock.Unlock()

	if eventChannel != nil {
		close(eventChannel)
		eventChannel = nil

		if err := iDeviceError(C.idevice_event_unsubscribe()); err != nil {
			return err
		}
	}
	return nil
}
