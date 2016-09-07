package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // 명시적으로 사용되지는 않지만 초기화를 위해서 import 되었음.
	_ "golang.org/x/image/tiff" // 명시적으로 사용되지는 않지만 초기화를 위해서 import 되었음.
	"math"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"runtime"
)

func main() {
	이미지_파일_경로_모음 := F이미지_파일_목록()
	fmt.Printf("%v개 파일을 찾았습니다.\n", len(이미지_파일_경로_모음))

	결과_맵 := make(map[string]([]string))
	ch초기화 := make(chan bool)
	in := make(chan string, len(이미지_파일_경로_모음))
	out := make(chan []string)
	ch종료 := make(chan bool)

	// 병렬처리를 위해서 CPU코어 수량만큼 컨트라스트 계산 루틴 실행.
	for i:=0 ; i<runtime.NumCPU() ; i++ {
		go F컨트라스트(ch초기화, in, out, ch종료)
		<-ch초기화
	}

	// 컨트라스트 계산 루틴에 입력값 제공.
	for _, 이미지_파일_경로 := range 이미지_파일_경로_모음 {
		in <- 이미지_파일_경로
	}

	// 컨트라스트 계산 루틴으로부터 결과값 받음.
	// 병렬처리에서는 순서가 뒤바뀔 수 있으므로, 순서에 무관한 맵에 임시 보관.
	for _, _ = range 이미지_파일_경로_모음 {
		결과 := <-out
		결과_맵[결과[0]] = 결과
		fmt.Printf("%v : %v\n", 결과[0], 결과[1])
	}

	// 맵에 저장된 결과를 파일 경로명에 따라서 정렬.
	결과_모음 := make([]([]string), len(이미지_파일_경로_모음))
	for i, 이미지_파일_경로 := range 이미지_파일_경로_모음 {
		결과_모음[i] = 결과_맵[이미지_파일_경로]
	}

	csv파일명 := "대비값_" + time.Now().Format("060102-150405") + ".csv"
	if 에러 := F_CSV쓰기(결과_모음, csv파일명); 에러 != nil {
		fmt.Println(에러.Error())
	}

	fmt.Printf("기록 완료 : %v\n\n", csv파일명)

	close(ch종료)
	fmt.Printf("연산을 마쳤습니다.\n Ctrl-C를 눌러서 프로그램을 종료하거나, 프로그램 창을 닫아주세요.\n")

	time.Sleep(24*time.Hour)
}

func F컨트라스트(ch초기화 chan bool, in chan string, out chan []string, ch종료 chan bool) {
	ch초기화 <- true

	for {
		select {
		case <-ch종료:
			return
		case 이미지_파일_경로 := <-in:
			out <- f컨트라스트_도우미(이미지_파일_경로)
		}
	}
}

func f컨트라스트_도우미(이미지_파일_경로 string) []string {
	const 정수64_최대값 int64 = 9223372036854775807
	const 무부호정수16_최대값 int64 = 65535

	파일, 에러 := os.Open(이미지_파일_경로)
	F에러_패닉(에러)

	원본_이미지, _, 에러 := image.Decode(파일)
	F에러_패닉(에러)

	경계 := 원본_이미지.Bounds()
	가로폭, 세로길이 := 경계.Max.X, 경계.Max.Y
	흑백_이미지 := image.NewGray16(경계)
	var 원본_색상 color.Color
	var 흑백_색상값 uint16

	최종_합계 := big.NewInt(0) // 크기에 제한이 없는 big.Int형식으로 지정.
	중간_합계 := int64(0)      // big.Int는 연산 속도가 느리므로, int64형식의 중간_합계를 버퍼로 사용
	const 중간_합계_리셋_기준점 = 정수64_최대값 - 무부호정수16_최대값 - 1

	// 흑백 변환 및 평균값 계산.
	// 이미지의 경계는 (0, 0)이 아닐 수도 있으므로, 반복문은 '경계.Min.Y', '경계.Min.X'에서 시작한다.
	// X보다는 Y에 대해서 먼저 반복하는 것이 (Y보다 X에 대해서 반복하는 것보다) 메모리 사용 효율이 좋은 경향이 있다.
	for y := 경계.Min.Y; y < 세로길이; y++ {
		for x := 경계.Min.X; x < 가로폭; x++ {
			// 흑백으로 변환
			원본_색상 = 원본_이미지.At(x, y)
			흑백_이미지.Set(x, y, 원본_색상)
			흑백_색상값 = 흑백_이미지.Gray16At(x, y).Y
			중간_합계 += int64(흑백_색상값)

			if 중간_합계 >= 중간_합계_리셋_기준점 {
				최종_합계.Add(최종_합계, big.NewInt(중간_합계))
				중간_합계 = 0
			}
		}
	}

	최종_합계.Add(최종_합계, big.NewInt(중간_합계))
	픽셀_수량 := int64(가로폭-경계.Min.X) * int64(세로길이-경계.Min.Y)

	최종_합계_Rat := new(big.Rat).SetInt(최종_합계)
	픽셀_수량_Rat := new(big.Rat).SetInt64(픽셀_수량)
	평균값, _ := new(big.Rat).Quo(최종_합계_Rat, 픽셀_수량_Rat).Float64()

	// RMS 컨트라스트 계산.
	var 제곱근_내_합계 float64

	for y := 경계.Min.Y; y < 세로길이; y++ {
		for x := 경계.Min.X; x < 가로폭; x++ {
			흑백_색상값 = 흑백_이미지.Gray16At(x, y).Y
			차이 := float64(흑백_색상값) - 평균값
			차이_제곱 := 차이 * 차이
			제곱근_내_합계 += 차이_제곱
		}
	}

	컨트라스트_제곱 := 제곱근_내_합계 / float64(픽셀_수량)
	RMS_컨트라스트 := math.Sqrt(컨트라스트_제곱)

	결과 := make([]string, 2)
	결과[0] = 이미지_파일_경로
	결과[1] = strconv.FormatInt(int64(RMS_컨트라스트 + 0.5), 10)   // 양수 반올림.

	return 결과
}

func F이미지_파일_목록() []string {
	드라이브_모음 := []string{"D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R",
		"S", "T", "U", "V", "W", "X", "Y", "Z"}
	이미지_파일_목록 := make([]string, 0)

	for _, 드라이브명 := range 드라이브_모음 {
		디렉토리명 := 드라이브명 + ":/DCIM/"
		if 존재함, 에러 := F존재함(디렉토리명); !존재함 || 에러 != nil {
			continue
		}

		filepath.Walk(디렉토리명, func(파일경로 string, 파일정보 os.FileInfo, 에러 error) error {
			switch {
			case 에러 != nil:
				if strings.Contains(에러.Error(), "Access is denied.") {
					return nil
				}

				fmt.Printf("예상하지 못한 에러 발생 : %v\n%v", 파일정보.Name(), 에러.Error())
				return 에러
			case 파일정보.IsDir():
				return nil
			case strings.HasSuffix(strings.ToLower(파일경로), ".jpg"):
				이미지_파일_목록 = append(이미지_파일_목록, 파일경로)
			case strings.HasSuffix(strings.ToLower(파일경로), ".tif"):
				이미지_파일_목록 = append(이미지_파일_목록, 파일경로)
			}

			return nil
		})
	}

	sort.Strings(이미지_파일_목록)

	return 이미지_파일_목록
}

func F존재함(경로 string) (bool, error) {
	_, 에러 := os.Stat(경로)

	switch {
	case 에러 == nil:
		return true, nil
	case os.IsNotExist(에러):
		return false, nil
	}

	return true, 에러
}

func F_CSV쓰기(레코드_모음 [][]string, 파일명 string) (에러 error) {
	defer func() {
		if r := recover(); r != nil {
			switch 값 := r.(type) {
			case error:
				에러 = 값
			default:
				에러 = fmt.Errorf("%v", 값)
			}
		}
	}()

	파일, 에러 := os.Create(파일명)
	F에러_패닉(에러)
	defer 파일.Close()

	csv기록기 := csv.NewWriter(파일)

	for _, 레코드 := range 레코드_모음 {
		F에러_패닉(csv기록기.Write(레코드))
	}

	csv기록기.Flush()
	F에러_패닉(csv기록기.Error())

	return nil
}

func F에러_패닉(에러 error) {
	if 에러 != nil {
		panic(에러)
	}
}
