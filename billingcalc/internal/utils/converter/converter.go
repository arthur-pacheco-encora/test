package converter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

func FromStrToDecimal(value string) decimal.Decimal {
	sanitized := sanitization.SanitizeFloatString(value)
	d, _ := decimal.NewFromString(sanitized)
	return d
}

func FromStrToTimeStamp(value string) time.Time {
	if strings.Contains(value, ":") {
		value = strings.Split(value, " ")[0]
	}

	if strings.Contains(value, "T") {
		value = strings.Split(value, "T")[0]
	}

	if strings.Contains(value, "/") {
		value = formatDateFromAmericanToDateOnly(value)
	}

	paddedValue := formatDateWithTwoDigits(value)

	dt, _ := time.Parse("2006-01-02", paddedValue)
	return dt
}

func formatDateFromAmericanToDateOnly(americanDate string) string {
	slices := strings.Split(americanDate, "/")
	return fmt.Sprintf("%s-%s-%s", slices[2], slices[0], slices[1])
}

func formatDateWithTwoDigits(date string) string {
	const (
		appendChar      = "0"
		yearLenght      = 4
		monthDayLenghts = 2
	)

	slices := strings.Split(date, "-")

	slices[0] = strings.Repeat(appendChar, yearLenght-utf8.RuneCountInString(slices[0])) + slices[0]
	slices[1] = strings.Repeat(appendChar, monthDayLenghts-utf8.RuneCountInString(slices[1])) + slices[1]
	slices[2] = strings.Repeat(appendChar, monthDayLenghts-utf8.RuneCountInString(slices[2])) + slices[2]

	return fmt.Sprintf("%s-%s-%s", slices[0], slices[1], slices[2])
}

func FromStrToTimestamppb(date string) *timestamppb.Timestamp {
	t, _ := time.Parse(time.RFC3339, date)
	return timestamppb.New(t)
}

func FromStringBase64ToIoReader(base64String string) io.Reader {
	if len(base64String) == 0 {
		return nil
	}

	buf := bytes.NewBufferString(base64String)
	return base64.NewDecoder(base64.StdEncoding, buf)
}
