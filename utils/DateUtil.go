package utils

import "time"

/*
	日期工具类
*/

const (
	DATE_PATTERN_DEFAULT = "2006-01-02 15:04:05"
	DATE_PATTERN_PRESS   = "20060102150405"
	DATE_PATTERN_SHORT   = "2006-01-02"
	DATE_PATTERN_SIMPLE  = "20060102"
)

/* 日期转换时间串 */
func ConvertDateToStr(date time.Time, pattern string) string {
	return date.Format(pattern)
}

/* 时间串转换日期 */
func ConvertStrToDate(date, pattern string) (time.Time, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(pattern, date, loc)
}

/* 时间串转换时间戳 */
func ConvertStrToUnix(date, pattern string) (int64, error) {
	t, err := ConvertStrToDate(date, pattern)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

/* 时间戳转日期 */
func ConvertUnixToDate(sec int64, pattern string) time.Time {
	return time.Unix(sec, 0)
}

/* 时间戳转时间串 */
func ConvertUnixToStr(sec int64, pattern string) string {
	return ConvertDateToStr(ConvertUnixToDate(sec, pattern), pattern)
}
