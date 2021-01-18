package distro

import (
	"bufio"
	"os"
	"strings"
)

// Release returns the contents of /etc/os-release, if it exists. It is guaranteed
// to exist, if the system uses systemd.
func Release() (map[string]string, error) {
	ret := make(map[string]string)

	f, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		pair := strings.Split(scan.Text(), "=")
		if len(pair) >= 2 {
			ret[pair[0]] = strings.Trim(pair[1], "\"")
		} else if len(pair) == 1 {
			ret[pair[0]] = ""
		}
	}
	return ret, nil
}
