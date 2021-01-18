// Package distro provides utilities for common operations that differ based on
// what operating system is currently running.
package distro

import (
	"runtime"
)

const (
	Generic = "linux"
	Linux   = "linux"

	Alpine     = "alpine"
	Amazon     = "amazon"
	Arch       = "arch"
	CentOS     = "centos"
	CloudLinux = "cloudlinux"
	Debian     = "debian"
	Fedora     = "fedora"
	Gentoo     = "gentoo"
	Kali       = "kali"
	Mageia     = "mageia"
	Mandriva   = "mandriva"
	Mint       = "linuxmint"
	OpenSUSE   = "opensuse"
	Oracle     = "oracle"
	Raspbian   = "raspbian"
	RHEL       = "rhel"
	Slackware  = "slackware"
	SLES       = "sles"
	Ubuntu     = "ubuntu"

	// Regular GOOS values.
	AIX       = "aix"
	Android   = "android"
	Darwin    = "darwin"
	Dragonfly = "dragonfly"
	FreeBSD   = "freebsd"
	Hurd      = "hurd"
	Illumos   = "illumos"
	IOS       = "ios"
	JS        = "js"
	NaCl      = "nacl"
	NetBSD    = "netbsd"
	OpenBSD   = "openbsd"
	Plan9     = "plan9"
	Solaris   = "solaris"
	Windows   = "windows"
	ZOS       = "zos"
)

// Detect returns a string representing the current running Linux distribution.
// If the system is not Linux, runtime.GOOS is returned. If the distribution
// cannot be determined, then "linux" is returned.
func Detect() string {
	if runtime.GOOS != Linux {
		return runtime.GOOS
	}

	release, err := Release()
	if err == nil {
		if os, ok := release["ID"]; ok {
			return os
		}
	}
	return "linux"
}
