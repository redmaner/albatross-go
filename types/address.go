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

const (
	COUNTRY_CODE = "NQ"
	EMPTY_CODE   = "00"
)

var (
	ErrInvalidAddressCountryCode = InvalidAddressError{code: 1, message: "invalid country code"}
	ErrInvalidAddressLength      = InvalidAddressError{code: 2, message: "invalid address lenght"}
	ErrInvalidAddressChecksum    = InvalidAddressError{code: 3, message: "invalid address checksum"}
)

type InvalidAddressError struct {
	code    int
	message string
}

func (i InvalidAddressError) Error() string { return "invalid address: " + i.message }

type Address [20]byte

func NewAddressFromHex(hexStr string) (Address, error) {
	decoded, err := hex.DecodeString(hexStr)
	return Address(decoded), err
}

func NewAddressFromFriendly(friendlyAddress string) (Address, error) {
	friendlyAddress = strings.ReplaceAll(friendlyAddress, " ", "")

	if friendlyAddress[:2] != COUNTRY_CODE {
		return Address{}, ErrInvalidAddressCountryCode
	}

	if ibanCheck(friendlyAddress[4:]+friendlyAddress[:4]) != 1 {
		return Address{}, ErrInvalidAddressChecksum
	}

	decoded, err := nimiqBase32Encoder.DecodeString(friendlyAddress[4:])
	if err != nil {
		return Address{}, err
	}

	return Address(decoded), nil
}

func (a Address) String() string {
	base32Encoded := nimiqBase32Encoder.EncodeToString(a[:])
	ibanNumber := 98 - ibanCheck(base32Encoded+COUNTRY_CODE+EMPTY_CODE)
	check := EMPTY_CODE + strconv.Itoa(ibanNumber)
	res := COUNTRY_CODE + check[len(check)-2:] + base32Encoded

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
			panic(err) // should be unreachable
		}

		ibanNumber = parsedInt % 97
	}

	return ibanNumber
}
