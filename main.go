// main.go позволяет протестировать нагрузку на базу данных в условиях постоянных ставок
// Перед запуском необходимо настроить доступ к базе в db.go и инициализировать объекты из tables.sql
package main

import (
	"fmt"
	"time"
	"database/sql"
	"log"
)

const (
	updateDaemonTimer = 3 * time.Second       //seconds             Частота обновления джекпотов/счетов
	insertPackTimer   = 100 * time.Millisecond //        Частота обнволения операций
	insertPackSize    = 2000                  // Память под пул операций в одном обновлении
	numBets_start     = 500000                 // Число ставок для теста
	numUsers_start    = 1000                   // Число юзеров которые будут делать ставки (поставить могут не все)
	betDelay          = 10 * time.Microsecond  //microseconds
)

var (
	db *sql.DB
	err error
	bets_ch   = make(chan bet)
	finish_ch = make(chan int)
)

func main() {
	db, err = make_db_conn()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	bets, _ := NewBets(numBets_start, numUsers_start)

	go updateDaemon()
	go betLoop()

	s_time := time.Now()
	defer func() { fmt.Println(time.Since(s_time)) }()
	go func() {
		for _, bet := range bets {
			time.Sleep(betDelay)
			addBet(bet.name, bet.deposit, bet.jp_part)
		}
	}()

	fmt.Println("Waiting...")
	<-finish_ch
	fmt.Println("All bets in. Goodbye	")
}
