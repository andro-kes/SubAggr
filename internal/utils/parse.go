package utils

import (
	"time"
)

// string MM-YYYY -> time.Time
func ParseTime(startDate, endDate string) (*time.Time, *time.Time, error) {
	start, err := time.Parse("01-2006", startDate)
	if !Ok(err == nil, "Не удалось распарсить дату старта при фильтрации") {
		return nil, nil, err
	}

	if endDate == "" {
		return &start, nil, nil
	}
	
	end, err := time.Parse("01-2006", endDate)
	if !Ok(err == nil, "Не удалось распарсить дату окончания при фильтрации") {
		return nil, nil, err
	}

	return &start, &end, nil
}