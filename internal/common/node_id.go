package common

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cast"

	"sharefood/pkg/util"
)

// GetHostname  get machine hostname
func GetHostname() string {
	nodeName, err := os.Hostname()
	if err != nil {
		return ""
	}

	return nodeName
}

func GetNodeID(nodeName string) int {
	nodeDefault := fmt.Sprint(time.Now().UnixNano() / (1 << 22))
	x := strings.Split(nodeName, ".")
	rgx, err := regexp.Compile(`[^0-9]+`)

	if err != nil || len(x) == 0 {
		return cast.ToInt(util.SubString(nodeDefault, (len(nodeDefault) - 3), 3))
	}

	n := rgx.ReplaceAllString(x[0], "")

	if n == "" {
		n = fmt.Sprint(time.Now().UnixNano() / (1 << 22))
	}

	if len(n) <= 3 {
		return cast.ToInt(n)
	}

	return cast.ToInt(util.SubString(n, (len(n) - 3), 3))
}
