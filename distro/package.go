package distro

import (
	"fmt"
)

// PackageManager represents common commands associated with a particular package
// manager.
type PackageManager interface {
	// Command returns the name of the package manager executable.
	Command() string

	// Install forms a command that will install the given packages.
	Install(...string) (cmd string, args []string)

	// Sync forms a command that will rebuild the package cache. For some package
	// managers, this command is unnecessary or discouraged (e.g. `pacman`). In
	// those cases, the returned args are nil.
	Sync() (cmd string, args []string)
}

type (
	apk      struct{}
	apt      struct{}
	pacman   struct{}
	emerge   struct{}
	slackpkg struct{}
	urpmi    struct{}
	yum      struct{}
	zypper   struct{}

	brew       struct{}
	chocolatey struct{}
	pkg        struct{}
)

var (
	Apk      = &apk{}
	Apt      = &apt{}
	Pacman   = &pacman{}
	Emerge   = &emerge{}
	SlackPkg = &slackpkg{}
	Urpmi    = &urpmi{}
	Yum      = &yum{}
	Zypper   = &zypper{}

	Brew       = &brew{}
	Chocolatey = &chocolatey{}
	Pkg        = &pkg{}
)

var packageManagers = map[string]PackageManager{
	Amazon:     Yum,
	Arch:       Pacman,
	CentOS:     Yum,
	CloudLinux: Yum,
	Debian:     Apt,
	Fedora:     Yum,
	Gentoo:     Emerge,
	Kali:       Apt,
	Mageia:     Urpmi,
	Mandriva:   Urpmi,
	Mint:       Apt,
	OpenSUSE:   Zypper,
	Oracle:     Yum,
	Raspbian:   Apt,
	RHEL:       Yum,
	Slackware:  SlackPkg,
	SLES:       Zypper,
	Ubuntu:     Apt,

	Darwin:  Brew,
	FreeBSD: Pkg,
	Windows: Chocolatey,
}

func GetPackageManager() (PackageManager, error) {
	distro := Detect()
	if packman, ok := packageManagers[distro]; ok {
		return packman, nil
	}
	return nil, fmt.Errorf("no known package manager for distro/os: %s", distro)
}

func (p *yum) Command() string {
	return "yum"
}

func (p *yum) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"install"}, packages...)
}

func (p *yum) Sync() (string, []string) {
	return p.Command(), []string{"makecache"}
}

func (p *apt) Command() string {
	return "apt-get"
}

func (p *apt) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"install"}, packages...)
}

func (p *apt) Sync() (string, []string) {
	return p.Command(), []string{"update"}
}

func (p *pacman) Command() string {
	return "pacman"
}

func (p *pacman) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"-S"}, packages...)
}

func (p *pacman) Sync() (string, []string) {
	return p.Command(), nil
}

func (p *apk) Command() string {
	return "apk"
}

func (p *apk) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"add"}, packages...)
}

func (p *apk) Sync() (string, []string) {
	return "", nil
}

// Command return the command for installing packages.
func (p *urpmi) Command() string {
	return "urpmi"
}

func (p *urpmi) Install(packages ...string) (string, []string) {
	return "urpmi", packages
}

func (p *urpmi) Sync() (string, []string) {
	return "urpmi.update", []string{"-a"}
}

func (p *zypper) Command() string {
	return "zypper"
}

func (p *zypper) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"install"}, packages...)
}

func (p *zypper) Sync() (string, []string) {
	return p.Command(), []string{"refresh"}
}

func (p *slackpkg) Command() string {
	return "slackpkg"
}

func (p *slackpkg) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"install"}, packages...)
}

func (p *slackpkg) Sync() (string, []string) {
	return p.Command(), []string{"update"}
}

func (p *brew) Command() string {
	return "brew"
}

func (p *brew) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"install"}, packages...)
}

func (p *brew) Sync() (string, []string) {
	return "", nil
}

func (p *chocolatey) Command() string {
	return "choco"
}

func (p *chocolatey) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"install"}, packages...)
}

func (p *chocolatey) Sync() (string, []string) {
	return "", nil
}

func (p *emerge) Command() string {
	return "emerge"
}

func (p *emerge) Install(packages ...string) (string, []string) {
	return p.Command(), packages
}

func (p *emerge) Sync() (string, []string) {
	return "", nil
}

func (p *pkg) Command() string {
	return "pkg"
}

func (p *pkg) Install(packages ...string) (string, []string) {
	return p.Command(), append([]string{"install"}, packages...)
}

func (p *pkg) Sync() (string, []string) {
	return p.Command(), []string{"update"}
}
