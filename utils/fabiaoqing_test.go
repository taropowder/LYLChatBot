package utils

import (
	"fmt"
	"testing"
)

func TestGetFaBiaoQingImagesByParameters(t *testing.T) {
	img, err := GetFaBiaoQingImagesByParameters(0, "哭泣")
	fmt.Println(img, err)
}
