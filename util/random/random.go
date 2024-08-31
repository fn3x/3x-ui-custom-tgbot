package random

import (
	"crypto/rand"
	"fmt"
	mathRand "math/rand"
	"strings"
)

var (
	numSeq      [10]rune
	lowerSeq    [26]rune
	upperSeq    [26]rune
	numLowerSeq [36]rune
	numUpperSeq [36]rune
	allSeq      [62]rune
)

var seq = strings.Split("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")

func Init() {
	for i := 0; i < 10; i++ {
		numSeq[i] = rune('0' + i)
	}
	for i := 0; i < 26; i++ {
		lowerSeq[i] = rune('a' + i)
		upperSeq[i] = rune('A' + i)
	}

	copy(numLowerSeq[:], numSeq[:])
	copy(numLowerSeq[len(numSeq):], lowerSeq[:])

	copy(numUpperSeq[:], numSeq[:])
	copy(numUpperSeq[len(numSeq):], upperSeq[:])

	copy(allSeq[:], numSeq[:])
	copy(allSeq[len(numSeq):], lowerSeq[:])
	copy(allSeq[len(numSeq)+len(lowerSeq):], upperSeq[:])
}

func Seq(n int) string {
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
		runes[i] = allSeq[mathRand.Intn(len(allSeq))]
	}
	return string(runes)
}

func Num(n int) int {
	return mathRand.Intn(n)
}

func RandomLowerAndNum(n int) string {
	str := ""
	for i := 0; i < n; i++ {
		str += seq[mathRand.Intn(len(seq))]
	}

	return str
}

func RandomUUID() string {
	uuid := make([]byte, 16)
	_, _ = rand.Read(uuid)

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set the version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x8

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16])
}
