
















package common

import "encoding/hex"





func ToHex(b []byte) string {
	hex := Bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}


func ToHexArray(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = ToHex(b[i])
	}
	return r
}



func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}


func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}


func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}


func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}


func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}


func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}


func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}


func Hex2BytesFixed(str string, flen int) []byte {
	h, _ := hex.DecodeString(str)
	if len(h) == flen {
		return h
	}
	if len(h) > flen {
		return h[len(h)-flen:]
	}
	hh := make([]byte, flen)
	copy(hh[flen-len(h):flen], h)
	return hh
}


func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}


func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}


func TrimLeftZeroes(s []byte) []byte {
	idx := 0
	for ; idx < len(s); idx++ {
		if s[idx] != 0 {
			break
		}
	}
	return s[idx:]
}


func TrimRightZeroes(s []byte) []byte {
	idx := len(s)
	for ; idx > 0; idx-- {
		if s[idx-1] != 0 {
			break
		}
	}
	return s[:idx]
}
