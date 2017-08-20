//Простой генератор данных для тестов

package main

import (
	"github.com/Pallinder/go-randomdata"
	"errors"
	"math/rand"
)

type bet struct {
	name    string
	deposit float32
	jp_part float32
}

type bets []bet

func contains_string(name string, slice []string) bool {
	for _, v := range slice {
		if v == name {
			return true
		}
	}
	return false
}

//Генерация новых ставок
func NewBets(numBets int, numUsers int) (bets, error) {
	resultBets := make(bets, numBets)
	names := make([]string, numUsers)
	for i := range names {
		counter := 1
		for ; counter < 100; {
			counter += 1
			new_name := randomdata.SillyName()
			if !contains_string(new_name, names) {
				names[i] = new_name
				break
			}
		}
		if counter == 100 {
			return resultBets, errors.New("Out of names!")
		}
	}

	for i := range resultBets {
		resultBets[i].name = names[rand.Intn(numUsers)]
		resultBets[i].deposit = rand.Float32() * 100
		resultBets[i].jp_part = rand.Float32() * 10

	}

	return resultBets, nil
}
