// 最新日期信息:
// https://ncthmwrwbtst.cr.chiba-u.ac.jp/img/FULL_24h/latest.json?_=1664026635552
//
// 拼接图片地址示例:
// https://ncthmwrwbtst.cr.chiba-u.ac.jp/img/D531106/2d/550/2022/09/24/131000_0_1.png
// 其中的`4d`表示 4 倍尺寸的图片. 可选 1d, 2d, 4d, 8d, 16d.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

var (
	Version = "?"
	Commit  = "?"
	Ts      = "?"
)

var (
	varDir     = "var"
	httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}

	currentDate = "2022-01-01 00:00:01"
)

const (
	host     = "https://ncthmwrwbtst.cr.chiba-u.ac.jp"
	times    = 4   // 图片尺寸倍数. 裁剪时长宽被等分的倍数的份数. 如: 4 倍时就被裁剪为 4x4=16 张小图
	partSize = 550 // 裁剪图片的边长像素
)

func main() {
	v := flag.Bool("version", false, "show version.")

	// glog default option
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")

	flag.Parse()

	if *v {
		fmt.Printf("version: %s\ncommit: %s\ntimestamp: %s\n", Version, Commit, Ts)
		return
	}

	if err := os.MkdirAll(varDir, 0755); err != nil {
		glog.Fatal(err)
	}

	for {
		check()
		time.Sleep(time.Minute * 10)
	}
}

func check() {
	latestDate, err := getLatestDate()
	if err != nil {
		return
	}
	if latestDate == currentDate {
		return
	}
	currentDate = latestDate

	w := NewWallpaper(latestDate)
	if err := w.downloadLatestParts(); err != nil {
		return
	}
	if err := w.concat(); err != nil {
		return
	}

	count, err := GetMacOSDesktopCount()
	if err != nil {
		return
	}
	for i := 0; i < count; i++ {
		ApplyMacOSWallpaper(w.path(), i)
	}
	glog.Infof("apply wallpaper complete.")
}

type Wallpaper struct {
	Times     int
	Date      string
	PartsPath [][]string
}

func NewWallpaper(date string) *Wallpaper {
	return &Wallpaper{
		Times: times,
		Date:  date, // 2018-05-16 23:50:00
	}
}

func (w *Wallpaper) path() string {
	hms := strings.ReplaceAll(w.Date[11:], ":", "") // hour + minute + second
	return fmt.Sprintf("%s/%s.png", varDir, hms)
}

func (w *Wallpaper) downloadLatestParts() error {
	paths := make([][]string, w.Times)
	for r := 0; r < w.Times; r++ {
		paths[r] = make([]string, w.Times)
		for c := 0; c < w.Times; c++ {
			url := fmtPartURL(w.Times, w.Date, c, r)
			path := fmtPartPath(c, r)
			if err := downloadFile(url, path); err != nil {
				return HitErrorf("failed to download: %v", url)
			}
			glog.Infof("downloaded: %v", path)
			paths[r][c] = path
		}
	}
	w.PartsPath = paths
	return nil
}

// https://ncthmwrwbtst.cr.chiba-u.ac.jp/img/D531106/2d/550/2022/09/24/131000_0_1.png
func fmtPartPath(colum, row int) string {
	return fmt.Sprintf("%s/part_%d_%d.png", varDir, colum, row)
}

// https://ncthmwrwbtst.cr.chiba-u.ac.jp/img/D531106/2d/550/2022/09/24/131000_0_1.png
func fmtPartURL(times int, date string, colum, row int) string {
	dateurl := dateToURL(date)
	return fmt.Sprintf(
		"%s/img/D531106/%dd/%d/%s_%d_%d.png",
		host, times, partSize, dateurl, colum, row)
}

// 2018-05-16 23:50:00 => 2018/05/16/235000
func dateToURL(s string) string {
	s = strings.ReplaceAll(s, " ", "/")
	s = strings.ReplaceAll(s, "-", "/")
	s = strings.ReplaceAll(s, ":", "")
	return s
}

func (w *Wallpaper) concat() error {
	img := image.NewRGBA(image.Rect(0, 0,
		partSize*w.Times,
		partSize*w.Times,
	))
	for r := 0; r < w.Times; r++ {
		for c := 0; c < w.Times; c++ {
			if err := drawImage(img, w.PartsPath[r][c], c, r); err != nil {
				return HitErrorf("failed to draw image: %v", w.PartsPath[r][c])
			}
		}
	}

	concatPath := w.path()
	file, err := os.Create(concatPath)
	if err != nil {
		return HitErrorf("failed to create file: f=%s, err=%w", concatPath, err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return HitErrorf("failed to encode image: f=%s, err=%w", concatPath, err)
	}
	glog.Infof("concated: %v", concatPath)
	return nil
}

func drawImage(img *image.RGBA, partImgFile string, colum, row int) error {
	f, err := os.Open(partImgFile)
	if err != nil {
		return HitErrorf("failed to open file: f=%v, err=%w", f, err)
	}
	defer f.Close()
	part, err := png.Decode(f)
	if err != nil {
		return HitErrorf("failed to decode image file: f=%v, err=%w", f, err)
	}
	draw.Draw(
		img,
		part.Bounds().Add(image.Pt(partSize*colum, partSize*row)),
		part,
		part.Bounds().Min,
		draw.Over,
	)
	return nil
}

// return example: 2018-05-16 23:50:00
func getLatestDate() (date string, err error) {
	var payload struct {
		Date string `json:"date"` // 2018-05-16 23:50:00
	}
	if err = request(host+"/img/FULL_24h/latest.json", &payload); err != nil {
		return "", err
	}
	return payload.Date, nil
}

func downloadFile(url string, path string) error {
	rsp, err := httpClient.Get(url)
	if err != nil {
		return HitErrorf("failed to get himawari8 image: url=%v, err=%w", url, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return HitErrorf("failed to create file: path=%v, err=%w", path, err)
	}

	defer rsp.Body.Close()
	defer f.Close()

	_, err = io.Copy(f, rsp.Body)
	if err != nil {
		return HitErrorf("failed to write file: path=%v, err=%w", path, err)
	}
	return nil
}

func request(url string, dst any) error {
	rsp, err := httpClient.Get(url)
	if err != nil {
		return HitErrorf("failed to get himawari8 info json: url=%v, err=%w", url, err)
	}
	if rsp.StatusCode != http.StatusOK {
		return HitErrorf("failed to request himawari8: url=%v, code=%v", url, rsp.StatusCode)
	}

	defer rsp.Body.Close()
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return HitErrorf("failed to read himawari8 response body: url=%v, err=%v", url, err)
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return HitErrorf("failed to decode himawari8 response json: body=%s, err=%v", b, err)
	}
	return nil
}

func GetMacOSDesktopCount() (count int, err error) {
	scpt := `tell application "System Events" to copy count of desktops to stdout`
	output, err := Osascript(scpt)
	if err != nil {
		return 0, HitErrorf("failed to get desktop count: err=%w", err)
	}
	c, err := strconv.Atoi(output)
	if err != nil {
		return 0, HitErrorf("invalid desktop count: output=%v, err=%w", output, err)
	}
	return c, nil
}

func ApplyMacOSWallpaper(imagePath string, index int) error {
	abs, err := filepath.Abs(imagePath)
	if err != nil {
		return HitErrorf("failed to get image absolute path: rel=%v, err=%w", imagePath, err)
	}

	script := fmt.Sprintf(`
tell application "System Events"
  tell desktop %v
    set picture to "%v"
  end tell
end tell`, index+1, abs)

	if output, err := Osascript(script); err != nil {
		return HitErrorf("failed to apply wallpaper: output=%v, err=%w", output, err)
	}
	return nil
}

func Osascript(script string) (output string, err error) {
	var args []string
	for _, v := range strings.Split(script, "\n") {
		args = append(args, "-e", v)
	}
	return Cmd("osascript", args...)
}

func Cmd(name string, args ...string) (output string, err error) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = cmd.Stdout
	out, err := cmd.Output()
	if err != nil {
		return "", HitErrorf("execute command failed: cmd=%v, err=%w", cmd, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func HitErrorf(format string, a ...any) error {
	return HitErrorfDepth(2, format, a...)
}

func HitErrorfDepth(depth int, format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	glog.ErrorDepth(depth, err)
	return err
}
