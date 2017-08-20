package main

import (
	"fmt"
	"time"
)

const (
	updaterDeamonTimer = 10    //seconds
	insertPackTimer    = 500   //milliseconds
	insertPackSize     = 50000 //Inserts per commit
	numBets_start      = 500000
	numUsers_start     = 1000
	betDelay           = 30
)

var (
	bets_ch   = make(chan bet)
	finish_ch = make(chan int)
)

func main() {
	bets, _ := NewBets(numBets_start, numUsers_start)

	go updateDaemon()
	go betLoop()

	s_time := time.Now()
	defer func() { fmt.Println(time.Since(s_time)) }()
	go func() {
		for _, bet := range bets {
			time.Sleep(time.Microsecond * betDelay)
			addBet(bet.name, bet.deposit, bet.jp_part)
		}
	}()

	fmt.Println("Waiting...")
	<-finish_ch
	fmt.Println("All bets in. Goodbye	")
}
