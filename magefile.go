// +build mage

package main

import (
	"fmt"
	"strings"

	"github.com/inahga/inahgo/pkg/distro"
)

// RunDep provides a list of runtime dependencies, based on the running OS/distribution.
func RunDep() error {
	packman, err := distro.GetPackageManager()
	if err != nil {
		return err
	}

	var pkgs []string
	switch packman {
	case distro.Yum:
		pkgs = []string{
			"libimobiledevice",
		}
	case distro.Apt:
		pkgs = []string{
			"libimobiledevice6",
		}
	default:
		return fmt.Errorf("unknown or unsupported platform for runDep: %s/%s",
			distro.Detect(), packman.Command())
	}
	fmt.Println(strings.Join(pkgs, " "))
	return nil
}

// BuildDep provides a list of compile-time dependencies, based on the running
// OS/distribution.
func BuildDep() error {
	packman, err := distro.GetPackageManager()
	if err != nil {
		return err
	}

	var pkgs []string
	switch packman {
	case distro.Yum:
		pkgs = []string{
			"libimobiledevice-devel",
		}
	case distro.Apt:
		pkgs = []string{
			"libimobiledevice-dev",
		}
	default:
		return fmt.Errorf("unknown or unsupported platform for runDep: %s/%s",
			distro.Detect(), packman.Command())
	}
	fmt.Println(strings.Join(pkgs, " "))
	return nil
}
