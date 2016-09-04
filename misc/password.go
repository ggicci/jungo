package misc

import (
	"math/rand"
	"time"
)

var base64IndexTableForUrl = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

type ICanGeneratePassword interface {
	Read(p []byte) (n int, err error)
	Seed(seed int64)
}

type randBytesReader struct {
	rand.Source
}

func NewPasswordGenerator(seed ...int64) ICanGeneratePassword {
	if len(seed) > 0 {
		return &randBytesReader{rand.NewSource(seed[0])}
	}
	return &randBytesReader{rand.NewSource(time.Now().UnixNano())}
}

func (r *randBytesReader) Read(p []byte) (n int, err error) {
	offset, remainder := 0, len(p)
	for {
		tmp := r.Int63()
		for i := 0; i < 8; i++ {
			p[offset] = base64IndexTableForUrl[tmp&0x3F]
			remainder--
			if remainder == 0 {
				return len(p), nil
			}
			offset++
			tmp >>= 8
		}
	}
}

func GeneratePassword(length int) string {
	g := NewPasswordGenerator()
	bs := make([]byte, length)
	g.Read(bs)
	return string(bs)
}

func GenerateMultiPasswords(lengths ...int) []string {
	hashs := make([]string, len(lengths))
	if len(hashs) == 0 {
		return hashs
	}
	g := NewPasswordGenerator()
	for i, length := range lengths {
		mbs := make([]byte, length)
		g.Read(mbs)
		hashs[i] = string(mbs)
	}
	return hashs
}
