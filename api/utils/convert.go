package utils

import (
	"errors"
	"strconv"
)

const (
	WeiDecimals int8 = 18
)

func ConvertWeiToMain(amount string) (float64, error) {
	return ConvertStringToAmount(amount, WeiDecimals)
}

func getDivider(decimals int8) float64 {
	if decimals == 0 {
		return 1
	}
	divider := float64(10)
	for i := int8(0); i < decimals; i++ {
		divider *= 10
	}
	return divider
}

func ConvertStringToAmount(amount string, decimals int8) (float64, error) {
	if amount == "" {
		return 0, errors.New("cannot convert empty string")
	}

	converted, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	return converted / getDivider(decimals), nil
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

	fee := gasPriceConverted * gasUsedConverted / getDivider(WeiDecimals)

	return fee, nil
}
