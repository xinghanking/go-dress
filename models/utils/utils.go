package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/disintegration/imaging"
	"image"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"
)

func JsonDecode(jsonStr string) any {
	var jsonData any
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return nil
	}
	return jsonData
}

func Rand(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

func ReqPost(urls *[]string, data map[string]any, invalidUrls *[]string) (map[string]any, error) {
	n := len(*urls)
	if n == 0 {
		return nil, errors.New("urls is empty")
	}
	if n == 1 {
		return PostJson((*urls)[0], data)
	}
	i := Rand(n)
	url := (*urls)[i]
	resp, err := PostJson(url, data)
	if err != nil {
		*urls = append((*urls)[:i], (*urls)[i+1:]...)
		*invalidUrls = append(*invalidUrls, url)
		return ReqPost(urls, data, invalidUrls)
	}
	return resp, nil
}

func PostJson(url string, data map[string]any) (map[string]any, error) {
	var d []byte
	var err error
	if data != nil {
		d, err = json.Marshal(data)
		if err != nil {
			panic(err)
		}
	}
	var resp *http.Response
	resp, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 500 {
		return nil, errors.New(resp.Status)
	}
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if string(body) == "" {
		return nil, nil
	}
	var response map[string]any
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ConvertMap(input any) any {
	switch d := input.(type) {
	case map[any]any:
		output := make(map[string]any)
		for k, v := range d {
			output[k.(string)] = ConvertMap(v)
		}
		return output
	case []any:
		return d[:]
	default:
		return input
	}
}

func CheckType(i any) string {
	switch i.(type) {
	case string:
		return "string"
	case int:
		return "int"
	case int8:
		return "int8"
	case int16:
		return "int16"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case uint:
		return "uint"
	case uint8:
		return "uint8"
	case uint16:
		return "uint16"
	case uint32:
		return "uint32"
	case uint64:
		return "uint64"
	case float64:
		return "float64"
	case float32:
		return "float32"
	case bool:
		return "bool"
	case complex64:
		return "complex64"
	case complex128:
		return "complex128"
	case uintptr:
		return "uintptr"
	case []byte:
		return "[]byte"
	case chan<- string:
		return "chan<- string"
	case chan<- int:
		return "chan<- int"
	case chan<- int8:
		return "chan<- int8"
	case chan<- int16:
		return "chan<- int16"
	case chan int32:
		return "chan<- int32"
	case chan int64:
		return "chan<- int64"
	case chan string:
		return "chan<- string"
	case chan<- float32:
		return "chan<- float32"
	case chan<- float64:
		return "chan<- float64"
	case chan<- complex64:
		return "chan<- complex64"
	case chan<- complex128:
		return "chan<- complex128"
	case map[string]interface{}:
		return "map"
	case []interface{}:
		return "array"
	case struct{}:
		return "struct"
	default:
		return reflect.TypeOf(i).Kind().String()
	}
}

func CovertToPng(image_path string, output_path string) error {
	imageData, err := imaging.Open(image_path)
	if err != nil {
		return err
	}
	out, err := os.Create(output_path)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		panic(err)
	}(out)
	err = imaging.Encode(out, imageData, imaging.PNG)
	if err != nil {
		return err
	}
	return nil
}

func ImageCrop(imageData image.Image, rect image.Rectangle, outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}(out)
	rgb := imaging.Crop(imageData, rect)
	err = imaging.Encode(out, rgb, imaging.PNG)
	return err
}
func MaxInt(num ...int) int {
	n := num[0]
	for _, i := range num {
		if i > n {
			n = i
		}
	}
	return n
}
func MinInt(num ...int) int {
	n := num[0]
	for _, i := range num {
		if i < n {
			n = i
		}
	}
	return n
}
func MinFloat(num ...float64) float64 {
	n := num[0]
	for _, i := range num {
		if i < n {
			n = i
		}
	}
	return n
}

func ChangeCronRang(rect map[string]float64, ratio float64) map[string]int {
	if ratio < 1 {
		rect["x"] += rect["w"] * (1 - ratio) / 2
		rect["w"] = rect["w"] * ratio
		rect["y"] += rect["h"] * (1 - ratio) / 2
		rect["h"] = rect["h"] * ratio
	}
	if ratio > 1 {
		widthMaxRatio := (rect["w"] + 2*rect["x"]) / rect["w"]
		heightMaxRatio := (2*rect["y"] + rect["h"]) / rect["h"]
		ratio := MinFloat(ratio, widthMaxRatio, heightMaxRatio)
		rect["x"] -= rect["w"] * (ratio - 1) / 2
		rect["y"] -= rect["h"] * (ratio - 1) / 2
		rect["w"] = rect["w"] * ratio
		rect["h"] = rect["h"] * ratio
	}
	return map[string]int{
		"x": MaxInt(int(rect["x"]), 0),
		"y": MaxInt(int(rect["y"]), 0),
		"w": int(MinFloat(rect["w"], rect["width"]-rect["x"])),
		"h": int(MinFloat(rect["h"], rect["height"]-rect["y"])),
	}
}
