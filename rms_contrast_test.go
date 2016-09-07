package main

import (
	"fmt"
	"testing"
	"strconv"
)

func TestF컨트라스트_도우미_JPEG(t *testing.T) {
	const 이미지_파일명_원본 = "test_img_normal.jpg"
	const 이미지_파일명_고대비 = "test_img_sharpened.jpg"

	RMS컨트라스트_원본, 에러 := strconv.ParseInt(f컨트라스트_도우미(이미지_파일명_원본)[1], 10, 64)
	if 에러 != nil {
		fmt.Println(에러.Error())
		t.FailNow()
	}

	RMS컨트라스트_고대비, 에러 := strconv.ParseInt(f컨트라스트_도우미(이미지_파일명_고대비)[1], 10, 64)
	if 에러 != nil {
		fmt.Println(에러.Error())
		t.FailNow()
	}

	fmt.Printf("%v : %v\n", RMS컨트라스트_원본, RMS컨트라스트_고대비)

	if RMS컨트라스트_원본 >= RMS컨트라스트_고대비 {
		t.FailNow()
	}
}

func TestF컨트라스트_도우미_TIFF(t *testing.T) {
	const 이미지_파일명_원본 = "test_img_normal.tif"
	const 이미지_파일명_고대비 = "test_img_sharpened.tif"

	RMS컨트라스트_원본, 에러 := strconv.ParseInt(f컨트라스트_도우미(이미지_파일명_원본)[1], 10, 64)
	if 에러 != nil {
		fmt.Println(에러.Error())
		t.FailNow()
	}

	RMS컨트라스트_고대비, 에러 := strconv.ParseInt(f컨트라스트_도우미(이미지_파일명_고대비)[1], 10, 64)
	if 에러 != nil {
		fmt.Println(에러.Error())
		t.FailNow()
	}

	fmt.Printf("%v : %v\n", RMS컨트라스트_원본, RMS컨트라스트_고대비)

	if RMS컨트라스트_원본 >= RMS컨트라스트_고대비 {
		t.FailNow()
	}
}
