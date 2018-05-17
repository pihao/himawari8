package oshlp

import (
	"fmt"
	"os/exec"
	"strings"
)

func Cmd(name string, arg ...string) (stdout string, err error) {
	c := exec.Command(name, arg...)
	out, err := c.Output()
	if err != nil {
		fmt.Println("ERROR::Cmd:", name, arg)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
