package util

import "time"

func IsSameMonth(time1 int64, time2 int64) bool {
	t1 := time.Unix(time1, 0)
	t2 := time.Unix(time2, 0)
	if t1.Year() != t2.Year() {
		return false
	}

	if t1.Month() != t2.Month() {
		return false
	}
	return true
}

func GetLastTimeByMonth(year int, month int, currentLocation *time.Location) int64 {
	lastOfMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, currentLocation)
	return lastOfMonth.Unix() - 1
}

func GetLastYearAndMonth(year, month int) (lastYear int, lastMonth int) {
	lastMonth = month - 1
	lastYear = year
	if lastMonth <= 0 {
		lastMonth = 12
		lastYear--
	}
	return lastYear, lastMonth
}

func GetLastTimeOfLastMonthByNow(timeNow int64) int64 {
	t := time.Unix(timeNow, 0)
	lastYear, lastMonth := GetLastYearAndMonth(t.Year(), int(t.Month()))
	currentLocation := t.Location()
	return GetLastTimeByMonth(lastYear, lastMonth, currentLocation)
}

func GetFirstTimeByMonth(year int, month int, currentLocation *time.Location) int64 {
	lastOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, currentLocation)
	return lastOfMonth.Unix()
}

func GetNextYearAndMonth(year, month int) (nextYear int, nextMonth int) {
	nextMonth = month + 1
	nextYear = year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	return nextYear, nextMonth
}

func GetNextYearAndMonthByNow() (nextYear int, nextMonth int) {
	year, month := GetCurrentYearAndMonth()
	return GetNextYearAndMonth(year, month)
}

func GetFirstTimeOfNextMonthByNow(timeNow int64) int64 {
	t := time.Unix(timeNow, 0)
	nextYear, nextMonth := GetNextYearAndMonth(t.Year(), int(t.Month()))
	currentLocation := t.Location()
	return GetFirstTimeByMonth(nextYear, nextMonth, currentLocation)
}

func GetCurrentYearAndMonth() (year int, month int) {
	timeNow := time.Now().Unix()
	t := time.Unix(timeNow, 0)
	return t.Year(), int(t.Month())
}

func GetCurrentTimeOfCurMonth() int64 {
	year, month := GetCurrentYearAndMonth()

	now := time.Now()
	return GetLastTimeByMonth(year, month, now.Location())
}

func GetMiddleTimeOfCurDay() int64 {
	t := time.Now()
	middle_tm := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location()).Unix()
	return middle_tm
}

//判断是否为闰年
func IsLeapYear(year int) bool { //y == 2000, 2004
	//判断是否为闰年
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}

	return false
}

func GetDaysOfMonth(year, month int) int {
	if month == 2 {
		if IsLeapYear(year) {
			return 29
		} else {
			return 28
		}
	}

	if month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12 {
		return 31
	}

	return 30
}

func GetDaysOfCurMonth() int {
	year, month := GetCurrentYearAndMonth()
	return GetDaysOfMonth(year, month)
}

func GetZeroUnixTime() int64 {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.Unix()
	return tm2
}
