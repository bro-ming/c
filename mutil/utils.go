package mutil

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"math/rand"
	"reflect"
	"time"
)

// gzip
func GZipDecompress(input []byte) ([]byte, error) {
	buf := bytes.NewBuffer(input)
	reader, gzipErr := gzip.NewReader(buf)
	if gzipErr != nil {
		return nil, gzipErr
	}
	defer reader.Close()


	result, readErr := ioutil.ReadAll(reader)
	if readErr != nil {
		return nil, readErr
	}
	return result, nil
}


func GZipCompress(input string) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write([]byte(input))
	if err != nil {
		return nil, err
	}

	err = gz.Flush()
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetRandomString  
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// StructToMap struct to map
func StructToMap(obj interface{}) map[string]interface{} {
	v := reflect.ValueOf(obj).Elem()
	vType := v.Type()
	var result = map[string]interface{}{}
	for i := 0; i < vType.NumField(); i++ {
		result[vType.Field(i).Name] = v.Field(i).String()
	}
	return result
}

// Verification attributes args is verification field
func Verification(s interface{}, args ...string) error {
	attributes := StructToMap(s)
	if len(attributes) == 0 {
		return errors.New("验证失败！")
	}

	for _, v := range args {
		for k, cv := range attributes {
			if v == k && cv == "" {
				return errors.New(k + "不能为空！")
			}
		}
	}

	return nil
}

// TimeMinAgo 返回前一分钟整秒
func TimeMinAgo(ts int64) int64 {
	Ts := time.Unix(ts-60, ts)

	//小时Ts
	hour := Ts.Hour()

	// 分钟Ts
	minute := Ts.Minute()
	mu := time.Date(Ts.Year(), Ts.Month(), Ts.Day(), hour, minute, 0, 0, Ts.Location())
	return mu.Unix()
}

// TimeMinZero 返回时间戳的分钟整秒
func TimeMinZero(ts int64) int64 {
	timer := time.Unix(ts, 0)
	z := time.Date(timer.Year(), timer.Month(), timer.Day(), timer.Hour(), timer.Minute(), 0, 0, timer.Location())
	return z.Unix()
}

// TimeZeroForToday 返回今日零点时间戳
func TimeZeroForToday() (zeroTs int64) {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return newTime.Unix()
}
