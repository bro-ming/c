/*
	aes加密
*/
package mutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

/* 调用示例
// aes加密key
var ic = util.Iaes{
	Key:[]byte("12345678901234567890123456789012"), //可以是16位,24位,32位
	Iv:[]byte("0000000000000000"), //必须为16位
}

//Aes加密
encode,_ := ic.Encrypt([]byte(`{"code":1,"data":"msgBody","msg":"Success"}`))
fmt.Println(string(encode))

//Aes解密
decode,_ := ic.Decrypt(string(encode))
fmt.Println(string(decode))
*/

type Aes struct {
	Key []byte //可以是16位,24位,32位
	Iv  []byte //必须为16位
}

// Encrypt 加密数据,返回一个base64类型的string
func (ia *Aes) Encrypt(origData []byte) (aesStr []byte, err error) {
	block, err := aes.NewCipher(ia.Key)
	if err != nil {
		return aesStr, err
	}
	blockSize := block.BlockSize()
	origData = pkcs7padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, ia.Iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	b64Str := base64.StdEncoding.EncodeToString(crypted)

	return []byte(b64Str), nil
}

// Decrypt 解密数据
func (ia *Aes) Decrypt(base64string string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(base64string)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(ia.Key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, ia.Iv)
	origData := make([]byte, len(b))
	blockMode.CryptBlocks(origData, b)
	origData = pkcs7unpadding(origData)
	return origData, nil
}

// pkcs7padding 增加填充
func pkcs7padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pkcs7padding 删除填充
func pkcs7unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
