package services

import (
	"math"
)

func RemoveEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func GetPrice(s []float64) float64 {
	var totalPrice float64
	i := 0
	for _, value := range s {
		if value != 0.0 {
			totalPrice += value
			i++
		}
	}
	if i > 2 && i < 4 {
		totalPrice = totalPrice - (totalPrice * 0.1)
		return math.Round(totalPrice*100) / 100
	}
	if i > 4 {
		totalPrice = totalPrice - (totalPrice * 0.15)
		return math.Round(totalPrice*100) / 100
	}
	return totalPrice
}

func CheckArray(arr []string) []string {
	s := len(arr)
	str := ""
	if s < 5 {
		for i := s; i < 5; i++ {
			arr = append(arr, str)
		}
	}
	return arr
}
