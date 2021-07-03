package utils

import (
	"errors"
	"strconv"
)

const (
	WeiDecimalPlaces = 1e18
)

func ConvertStringToBNB(amount string) (float64, error) {
	if amount == "" {
		return 0, errors.New("cannot convert empty string")
	}

	converted, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	return converted / WeiDecimalPlaces, nil
}

func ComputeFee(gasUsed string, gasPrice string) (float64, error) {
	if gasUsed == "" || gasPrice == "" {
		return 0, errors.New("cannot convert empty strings")
	}

	gasUsedConverted, err := strconv.ParseFloat(gasUsed, 64)
	if err != nil {
		return 0, err
	}

	gasPriceConverted, err := strconv.ParseFloat(gasPrice, 64)
	if err != nil {
		return 0, err
	}

	return gasPriceConverted * gasUsedConverted / WeiDecimalPlaces, nil
}
