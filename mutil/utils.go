package mutil

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
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
	bys := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bys[r.Intn(len(bys))])
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

// return uuid
func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", 32)
}

// return uuid and no split line
func UUIDNoSplit() string {
	return uuid.New().String()
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