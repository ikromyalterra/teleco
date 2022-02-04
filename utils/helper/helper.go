package helper

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GenerateUniqueID func to generate reference number using timestamp and random number
func GenerateUniqueID() string {
	uniqueTS := fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
	// Get last 12 digit of millisecond timestamp + 4 random number
	refNo := uniqueTS[len(uniqueTS)-12:] + RandomNumber(4)
	return refNo
}

// RandomNumber func to get random number in x digits
func RandomNumber(digit int) string {
	digit--
	pattern := "%0" + strconv.Itoa(digit) + "d"

	low := "1" + fmt.Sprintf(pattern, 0)
	lowInt, _ := strconv.Atoi(low)

	hi := low + "0"
	hiInt, _ := strconv.Atoi(hi)

	randNumber := lowInt + rand.Intn(hiInt-lowInt)
	return strconv.Itoa(randNumber)
}
