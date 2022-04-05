package idevice

/*
#cgo pkg-config: libimobiledevice-1.0
#include <libimobiledevice/libimobiledevice.h>
*/
import "C"

import "time"

// Connection is the connection handle.
type Connection struct {
	connection C.idevice_connection_t
}

// Connect sets up a connection to the device on the given destination port.
func (i *IDevice) Connect(port uint16) (*Connection, error) {
	return nil, nil
}

// Close disconnects from the device and cleans up the connection strucure.
func (c *Connection) Close() error {
	return nil
}

// Send data to the device via the connection. Returns the number of bytes sent.
func (c *Connection) Send(buf []byte) (uint32, error) {
	return 0, nil
}

// Receive data from the connection. Returns the number of bytes sent.
func (c *Connection) Receive(buf []byte) (uint32, error) {
	return 0, nil
}

// ReceiveTimeout receives data from the connection, terminating if the given
// timeout is elapsed. Returns the number of bytes sent.
func (c *Connection) ReceiveTimeout(buf []byte, timeout time.Duration) (uint32, error) {
	return 0, nil
}

// Fd returns the underlying file descriptor for the connection.
func (c *Connection) Fd() uintptr {
	return 0
}

// EnableSSL enables SSL on the connection. Returns ErrSSLError if SSL initialization,
// setup, or handshake fails.
func (c *Connection) EnableSSL() error {
	return nil
}

// DisableSSL disables SSL on the connection. Returns success even if SSL is not
// currently enabled.
func (c *Connection) DisableSSL() {

}
