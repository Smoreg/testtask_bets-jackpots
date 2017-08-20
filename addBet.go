package main

import (
	"log"
	"fmt"
	"time"
	"sync"
)

func addBet(name string, deposit, jp_part float32) {
	bets_ch <- bet{name, deposit, jp_part}
}

func betLoop() {
	log.Print("betLoop starting")
	bets_counter := 0
	betsPool := make([]bet, 0, insertPackSize)
	ticker := time.NewTicker(time.Millisecond * insertPackTimer)
	var wg sync.WaitGroup
	for {
		select {
		case <-ticker.C:
			log.Print(
				fmt.Sprintf("betLoop start inserting Pool %v All %v",
					len(betsPool),
					bets_counter))
			if len(betsPool) == 0 {
				continue
			}

			wg.Add(1)
			go func(pool []bet) {

				tmp_db, err := make_db_conn()
				if err != nil {
					log.Panic(err)
				}
				defer tmp_db.Close()
				defer wg.Done()

				tx, err := tmp_db.Begin()
				if err != nil {
					log.Panic(err)
				}
				stmt, err := tx.Prepare(
					"INSERT INTO operations (user_name, deposit, jackpot_part) VALUES ($1, $2, $3)")
				if err != nil {
					log.Panic(err)
				}
				for _, iBet := range pool {
					if _, err := stmt.Exec(iBet.name, iBet.deposit, iBet.jp_part); err != nil {
						tx.Rollback() // return an error too, we may want to wrap them
						log.Panic(err)
					}
				}
				stmt.Close()
				if err = tx.Commit(); err != nil {
					log.Panic(err)
				}

			}(betsPool)

			if bets_counter == numBets_start {
				wg.Wait()
				finish_ch <- 1
				log.Print("betLoop says goodbye")
				return
			}
			if bets_counter > numBets_start {
				log.Panic("Too many bets!")
			}
			betsPool = make([]bet, 0, insertPackSize)

		case newBet := <-bets_ch:
			bets_counter++
			betsPool = append(betsPool, newBet)

		}
	}
}
