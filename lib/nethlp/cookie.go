package nethlp

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"bitbucket.org/pihao/phtool-api/lib/filehlp"
)

type Cookie struct {
	FilePath string
	Data     CookieItems
}
type CookieItems map[string]CookieItem
type CookieItem struct {
	Value string
	Args  string
}

func NewCookie(filePath string) *Cookie {
	c := Cookie{
		FilePath: filePath,
	}
	c.readSaved()
	return &c
}

func (this *Cookie) Value() string {
	arr := []string{}
	for k, v := range this.Data {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v.Value))
	}
	return strings.Join(arr, "; ")
}

func (this *Cookie) Set(rsp *http.Response) {
	new_c := *this.parseResponse(rsp)
	if len(new_c) == 0 {
		fmt.Printf("cookie is null: %v", this.FilePath)
	}
	for k, v := range new_c {
		this.Data[k] = v
	}
	this.save()
}

func (this *Cookie) parseResponse(rsp *http.Response) *CookieItems {
	data := CookieItems{}
	for _, e := range rsp.Header["Set-Cookie"] {
		sp := regexp.MustCompile("; ").Split(e, 2)
		pair := strings.Split(sp[0], "=")
		if len(pair) != 2 {
			fmt.Println("read 'Set-Cookie' failed:", this.FilePath)
			continue
		}
		args := ""
		if len(sp) == 2 {
			args = sp[1]
		}
		data[pair[0]] = CookieItem{pair[1], args}
	}
	return &data
}

func (this *Cookie) readSaved() {
	v := CookieItems{}
	if err := filehlp.ReadJSON(&v, this.FilePath); err != nil && err != io.EOF {
		panic(fmt.Sprintf("read cookie error: %v, %v", this.FilePath, err))
	}
	this.Data = v
}

func (this *Cookie) save() {
	if err := filehlp.WriteJSON(&this.Data, this.FilePath, true); err != nil {
		panic(fmt.Sprintf("write cookie error: %v, %v", this.FilePath, err))
	}
}
