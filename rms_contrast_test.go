package main

import (
	"fmt"
	"testing"
)

func TestRMS컨트라스트_JPEG(t *testing.T) {
	const 이미지_파일명_원본 = "test_img_normal.jpg"
	const 이미지_파일명_고대비 = "test_img_sharpened.jpg"

	RMS컨트라스트_원본 := F_RMS_컨트라스트(이미지_파일명_원본)
	RMS컨트라스트_고대비 := F_RMS_컨트라스트(이미지_파일명_고대비)

	fmt.Printf("%g : %g", RMS컨트라스트_원본, RMS컨트라스트_고대비)

	if RMS컨트라스트_원본 >= RMS컨트라스트_고대비 {
		t.FailNow()
	}
}

func TestRMS컨트라스트_TIFF(t *testing.T) {
	const 이미지_파일명_원본 = "test_img_normal.tif"
	const 이미지_파일명_고대비 = "test_img_sharpened.tif"

	RMS컨트라스트_원본 := F_RMS_컨트라스트(이미지_파일명_원본)
	RMS컨트라스트_고대비 := F_RMS_컨트라스트(이미지_파일명_고대비)

	fmt.Printf("%g : %g", RMS컨트라스트_원본, RMS컨트라스트_고대비)

	if RMS컨트라스트_원본 >= RMS컨트라스트_고대비 {
		t.FailNow()
	}
}
