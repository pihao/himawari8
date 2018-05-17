package macos

import (
	"strings"

	"github.com/pihao/himawari8-desktop/lib/oshlp"
)

func Osascript(script string) (stdout string, err error) {
	arr := strings.Split(script, "\n")
	var arg []string
	for _, v := range arr {
		arg = append(arg, "-e", v)
	}
	out, err := oshlp.Cmd("osascript", arg...)
	if err != nil {
		return "", err
	}
	return out, nil
}
