















package trie



















func hexToCompact(hex []byte) []byte {
	terminator := byte(0)
	if hasTerm(hex) {
		terminator = 1
		hex = hex[:len(hex)-1]
	}
	buf := make([]byte, len(hex)/2+1)
	buf[0] = terminator << 5 
	if len(hex)&1 == 1 {
		buf[0] |= 1 << 4 
		buf[0] |= hex[0] 
		hex = hex[1:]
	}
	decodeNibbles(hex, buf[1:])
	return buf
}



func hexToCompactInPlace(hex []byte) int {
	var (
		hexLen    = len(hex) 
		firstByte = byte(0)
	)
	
	if hexLen > 0 && hex[hexLen-1] == 16 {
		firstByte = 1 << 5
		hexLen-- 
	}
	var (
		binLen = hexLen/2 + 1
		ni     = 0 
		bi     = 1 
	)
	if hexLen&1 == 1 {
		firstByte |= 1 << 4 
		firstByte |= hex[0] 
		ni++
	}
	for ; ni < hexLen; bi, ni = bi+1, ni+2 {
		hex[bi] = hex[ni]<<4 | hex[ni+1]
	}
	hex[0] = firstByte
	return binLen
}

func compactToHex(compact []byte) []byte {
	if len(compact) == 0 {
		return compact
	}
	base := keybytesToHex(compact)
	
	if base[0] < 2 {
		base = base[:len(base)-1]
	}
	
	chop := 2 - base[0]&1
	return base[chop:]
}

func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	return nibbles
}



func hexToKeybytes(hex []byte) []byte {
	if hasTerm(hex) {
		hex = hex[:len(hex)-1]
	}
	if len(hex)&1 != 0 {
		panic("can't convert hex key of odd length")
	}
	key := make([]byte, len(hex)/2)
	decodeNibbles(hex, key)
	return key
}

func decodeNibbles(nibbles []byte, bytes []byte) {
	for bi, ni := 0, 0; ni < len(nibbles); bi, ni = bi+1, ni+2 {
		bytes[bi] = nibbles[ni]<<4 | nibbles[ni+1]
	}
}


func prefixLen(a, b []byte) int {
	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}


func hasTerm(s []byte) bool {
	return len(s) > 0 && s[len(s)-1] == 16
}
