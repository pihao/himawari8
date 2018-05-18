package macos

import (
	"fmt"
	"strconv"
)

func GetDesktopCount() (count int, err error) {
	scpt := `tell application "System Events" to copy count of desktops to stdout`
	out, err := Osascript(scpt)
	if err != nil {
		return 0, err
	}

	c, err := strconv.ParseInt(out, 10, 0)
	if err != nil {
		return 0, err
	}
	return int(c), nil
}

func ApplyDesktop(imagePath string, index int) {
	scpt := fmt.Sprintf(`tell application "System Events"
  tell desktop %v
    set picture to "%v"
  end tell
end tell`, index+1, imagePath)
	Osascript(scpt)
}
