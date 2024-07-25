package types

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var nimiqBase32Encoder = base32.NewEncoding("0123456789ABCDEFGHJKLMNPQRSTUVXY")

var _ fmt.Stringer = (*Address)(nil)

type Address []byte

func FromHex(hexStr string) (Address, error) {
	decoded, err := hex.DecodeString(hexStr)
	return Address(decoded), err
}

func (a Address) String() string {
	const (
		CCODE = "NQ"
		EMPTY = "00"
	)

	base32Encoded := nimiqBase32Encoder.EncodeToString(a)
	ibanNumber := 98 - ibanCheck(base32Encoded+CCODE+EMPTY)
	check := EMPTY + strconv.Itoa(ibanNumber)
	res := CCODE + check[len(check)-2:] + base32Encoded

	chunks := []string{}
	for i := 0; i < len(res); i += 4 {
		chunks = append(chunks, res[i:i+4])
	}

	return strings.Join(chunks, " ")
}

func ibanCheck(str string) int {

	var numberString string
	for _, char := range str {

		code := unicode.ToUpper(char)

		if code >= 48 && code <= 57 {
			numberString = numberString + string(char)
		} else {
			numberString = numberString + strconv.Itoa(int(char-55))
		}
	}

	var ibanNumber int

	var i float64
	for i = 0.0; i < math.Ceil(float64(len(numberString))/6); i++ {
		startIndex := int(i) * 6
		endIndex := min(startIndex+6, len(numberString))
		ibanNumberString := strconv.Itoa(ibanNumber)
		parsedInt, err := strconv.Atoi(ibanNumberString + numberString[startIndex:endIndex])
		if err != nil {
			panic(err)
		}

		ibanNumber = parsedInt % 97
	}

	return ibanNumber
}
