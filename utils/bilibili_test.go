package utils

import (
	"fmt"
	"testing"
)

func TestCard(t *testing.T) {

	initConfig()
	contnt := DealWithBliBli("BV1aF4m1L7XL")

	fmt.Println(contnt)
	//DealWithBliBli
	//fmt.Println(regexp.MustCompile(`@.*?\s|\s`).ReplaceAllString("@三胖子 你好", ""))

}
