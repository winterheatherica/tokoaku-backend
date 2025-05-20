package fetcher

import (
	"fmt"
	"strconv"
)

func parseUintFromString(val interface{}) (uint, error) {
	str := fmt.Sprintf("%v", val)
	n, err := strconv.ParseUint(str, 10, 32)
	return uint(n), err
}
