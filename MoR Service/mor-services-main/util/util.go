package util

import (
	"github.com/oldjon/gutil"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	com "gitlab.com/morbackend/mor_services/common"
)

func CheckEmailAddr(emailAddr string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`)
	return regex.MatchString(emailAddr)
}

func GenerateRandomCode(length int) string {
	codeInt := rand.Intn(int(math.Pow(float64(10), float64(length))))
	code := strconv.Itoa(codeInt)
	for len(code) < com.VCodeLen {
		code = "0" + code
	}
	return code
}

func ReadUint32Slice(str string, separator string) []uint32 {
	if str == "" {
		return nil
	}
	out := make([]uint32, 0, 1)
	for _, v := range strings.Split(str, separator) {
		out = append(out, gutil.Str2Uint32(v))
	}
	return out
}
