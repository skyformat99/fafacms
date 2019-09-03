package util

import (
	"encoding/hex"
	"fmt"
	"github.com/hunterhug/go_image"
	"github.com/hunterhug/parrot/util"
	"testing"
)

func TestListFile(t *testing.T) {
	err := go_image.ScaleF2F("./timg.jpeg", "./timg_x.jpeg", 100)
	fmt.Printf("%#v", err)
}

func TestMd5(t *testing.T) {
	raw, err := util.ReadfromFile("./timg_x.jpeg")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(Md5(raw))

	raw, err = util.ReadfromFile("./timg.jpeg")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(Md5(raw))

	raw, err = util.ReadfromFile("./timg.jpeg")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	d, _ := Sha256(raw)
	fmt.Println(len(d))

	fmt.Println(hex.EncodeToString([]byte("123456789")))
}
