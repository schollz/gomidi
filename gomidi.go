package gomidi

import (
	"fmt"
	"strings"
)

func filterName(name string) (foundName string, foundNum int, err error) {
	names := Devices()
	for i, n := range names {
		if strings.Contains(strings.ToLower(n), strings.ToLower(name)) {
			foundName = n
			foundNum = i
			break
		}
	}
	if foundNum == -1 {
		err = fmt.Errorf("could not find device with name %s", name)
	}
	return
}
