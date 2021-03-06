//|------------------------------------------------------------------
//|        __
//|     __/  \
//|  __/  \__/_
//| /  \__/    \
//|/\__/CellGo /_
//|\/_/NetFW__/  \
//|  /\__ _/  \__/
//|  \/_/  \__/_/
//|    /\__/_/
//|    \/_/
//|------------------------------------------------------------------
//| Cellgo Framework session/encrypt file
//| All rights reserved: By cellgo.cn CopyRight
//| You are free to use the source code, but in the use of the process,
//| please keep the author information. Respect for the work of others
//| is respect for their own
//|-------------------------------------------------------------------
// Author:Tommy.Jin Dtime:2016-7-29

package session

import (
	"bytes"
	"encoding/gob"

	"fmt"
	"strconv"

	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

var (
	HashKey string = "9597f4KpYTsJ5tD6"
)

// serialize database
func Serialize(obj map[interface{}]interface{}) ([]byte, error) {
	for _, v := range obj {
		gob.Register(v)
	}
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(obj)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

// unserialize database
func Unserialize(encoded []byte) (map[interface{}]interface{}, error) {
	buf := bytes.NewBuffer(encoded)
	dec := gob.NewDecoder(buf)
	var out map[interface{}]interface{}
	err := dec.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Authcode(valueFormer []byte, operation string, hashkey string) (string, error) {
	var value []byte
	if operation == "ENCODE" {
		temp := encode(valueFormer)
		value = temp
	} else {
		value = valueFormer
	}
	auth_key := If(hashkey != "", hashkey, HashKey).(string)
	h := md5.New()
	h.Write([]byte(auth_key)) // md5加密
	cipherStr := h.Sum(nil)
	key := hex.EncodeToString(cipherStr)
	key_length := len(key)
	var valueStr string
	if operation == "DECODE" {
		temp, err := decode(value)
		if err != nil {
			return "", err
		}
		valueStr = fmt.Sprintf("%s", temp)

	} else {
		h.Write([]byte(fmt.Sprintf("%s", value) + key)) // md5加密
		cipherStr = h.Sum(nil)
		valueStr = Substr(hex.EncodeToString(cipherStr), 0, 8) + fmt.Sprintf("%s", value)

	}
	var (
		rndkey [256]int
		box    [256]int
		result string
		keys   []rune
	)

	keys = []rune(key)
	for i := 0; i <= 255; i++ {
		tempInt, err := strconv.ParseInt(fmt.Sprintf("%x", keys[i%key_length]), 16, 10)
		if err != nil {

			return "", err
		}
		rndkey[i] = int(tempInt)
		box[i] = i
	}

	for k, j := 0, 0; j < 256; j++ {
		k = (k + box[j] + rndkey[j]) % 256
		box[j], box[k] = box[k], box[j]
	}

	temprune := []rune(valueStr)
	valueStr_length := len(temprune)
	for x, y, z := 0, 0, 0; z < valueStr_length; z++ {
		x = (x + 1) % 256
		y = (y + box[x]) % 256
		box[x], box[y] = box[y], box[x]
		tempInt, err := strconv.ParseInt(fmt.Sprintf("%x", temprune[z]), 16, 64)
		if err != nil {
			return "", err
		}
		str := fmt.Sprintf("%c", int(tempInt)^(box[(box[x]+box[y])%256]))
		result += str
	}

	if operation == "DECODE" {

		h.Write([]byte(Substr(result, 8, 0) + key)) // md5加密
		cipherStr = h.Sum(nil)
		if Substr(result, 0, 8) == Substr(hex.EncodeToString(cipherStr), 0, 8) {
			temp, err := decode([]byte(Substr(result, 8, 0)))
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s", temp), nil
		} else {
			return "", nil
		}
	} else {
		return fmt.Sprintf("%s", encode([]byte(result))), nil

	}
}

//Intercept string function
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	if length == 0 {
		return string(rs[start:])
	}
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//Simulated three element operation
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// encode encodes a value using base64.
func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

// decode decodes a cookie using base64.
func decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}
