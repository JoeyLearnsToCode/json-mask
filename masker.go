package jsonmask

import (
	"strings"
	"unicode/utf8"
)

type (
	MaskHandler interface {
		Mask(raw string) (masked string)
	}
	MaskHandlerFunc func(string) string
)

func (f MaskHandlerFunc) Mask(raw string) (masked string) {
	return f(raw)
}

// ToAsterisks 变成全星号
func ToAsterisks(data string) string {
	cnt := utf8.RuneCountInString(data)
	sb := strings.Builder{}
	for i := 0; i < cnt; i++ {
		sb.WriteByte('*')
	}
	return sb.String()
}

func ToEmpty(data string) string {
	return ""
}

func Name(data string) string {
	runes := []rune(data)
	if len(runes) == 2 {
		return string(runes[0]) + "*"
	} else {
		return string(runes[0]) + "*" + string(runes[len(runes)-1])
	}
}

// Telephone 对11位手机号，用星号掩码中间4位
func Telephone(data string) string {
	return data[:3] + "****" + data[7:]
}
