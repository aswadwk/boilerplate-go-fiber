package helpers

import (
	"archive/zip"
	"path/filepath"
	"strconv"
	"time"
)

func StringToInt(s string, defaultValue int) int {

	r, err := strconv.Atoi(s)

	if err != nil {
		return defaultValue
	}

	return r
}

func IsImageOrPdf(file *zip.File) bool {
	ext := filepath.Ext(file.Name)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".pdf"
}

func GenerateDate(startDate, endDate string) ([]string, error) {
	if startDate == "" || endDate == "" {
		return nil, nil
	}

	// Parse the input dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// Generate the range of dates
	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	return dates, nil
}
