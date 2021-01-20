// Package idevice provides device/connection handling and communication. It
// corresponds to the libimobiledevice/libimobiledevice.h header.
package idevice

/*
#cgo pkg-config: libimobiledevice-1.0
#include <stdlib.h>
#include <stdint.h>
#include <libimobiledevice/libimobiledevice.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

var (
	ErrInvalidArgs   = errors.New("invalid argument")
	ErrUnknown       = errors.New("unknown error")
	ErrNoDevice      = errors.New("no device")
	ErrNotEnoughData = errors.New("not enough data")
	ErrSSLError      = errors.New("SSL error")
	ErrTimeout       = errors.New("timeout")
	ErrUnspecified   = errors.New("unspecified error")
)

func iDeviceError(ierr C.idevice_error_t) error {
	switch ierr {
	case 0:
		return nil
	case -1:
		return ErrInvalidArgs
	case -2:
		return ErrUnknown
	case -3:
		return ErrNoDevice
	case -4:
		return ErrNotEnoughData
	case -6:
		return ErrSSLError
	case -7:
		return ErrTimeout
	default:
		return ErrUnspecified
	}
}

// Options for NewWithOptions.
const (
	// LookupUsbmux indicates to include usbmux devices during lookup.
	LookupUsbmux uint32 = iota << 1
	// LookupNetwork indicates to include network devices during lookup.
	LookupNetwork
	// LookupPreferNetwork indicates to prefer network connection if device is
	// available over network and usbmux.
	LookupPreferNetwork
)

// Type of connection a device is available on.
const (
	ConnectionUsbmuxd uint32 = iota
	ConnectionNetwork
)

// IDevice is the device handle.
type IDevice struct {
	idevice C.idevice_t
}

// List gets a list of UDIDs of currently available devices (USBMUX devices only).
func List() ([]string, error) {
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

// Debug sets the debug level for libimobiledevice. Set to 1 to enable, or 0
// to disable. This only works if libimobiledevice has been compiled with
// debugging enabled.
func Debug(level int) {
	C.idevice_set_debug_level(C.int(level))
}

// New creates an IDevice structure for the device specified by UDID, if the
// device is available (USBMUX devices only). If you need to connect to a device
// available via network, use NewWithOptions() and include LookupNetwork in
// options.
//
// To select the first available device, pass an empty string to udid.
func New(udid string) (*IDevice, error) {
	return NewWithOptions(udid, 0)
}

// NewWithOptions creates an IDevice structure for the device specified
// by UDID, if the device is available, with the given lookup options.
//
// To select the first available device, pass an empty string to udid.
func NewWithOptions(udid string, options uint32) (*IDevice, error) {
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
