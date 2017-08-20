package main

import (
	"time"
	"log"
)

func updateDaemon() {
	log.Print("updateDaemon starting")

	update_users := `INSERT INTO users(user_name, deposit)
	select user_name, sum(deposit) from operations group by user_name
		ON CONFLICT (user_name)
		DO
		UPDATE
		SET deposit = EXCLUDED.deposit + users.deposit;`

	update_jackpot := `UPDATE old_jackpot
	SET value = oj.value + t.nj
	FROM old_jackpot as oj
	CROSS JOIN
		(
			SELECT COALESCE(sum(jackpot_part),CAST(0 AS money)) as nj
			FROM operations
		) t;`

	delete_operations := `DELETE FROM operations;`

	ticker := time.NewTicker(time.Second * updaterDeamonTimer)
	all_query := [3]string{update_users, update_jackpot, delete_operations}
	log.Print("updateDaemon start!")
	for {
		select {
		case <-ticker.C:
			db, err := make_db_conn()
			if err != nil {
				log.Panic(err)
			}
			rows, err := db.Query(`SELECT count(*) AS c FROM operations;`)
			if err != nil {
				log.Panic(err)
			}

			var co int
			rows.Next()
			rows.Scan(&co)
			rows.Close()
			if co == 0 {
				log.Print("updateDaemon no operations")
				db.Close()
				continue
			}

			tx, err := db.Begin()
			if err != nil {
				log.Panic(err)
			}
			for _, query := range all_query {
				{
					stmt, err := tx.Prepare(query)
					if err != nil {
						log.Panic(err)
					}
					defer stmt.Close()

					if _, err := stmt.Exec(); err != nil {
						tx.Rollback() // return an error too, we may want to wrap them
						log.Panic(err)
					}
				}
			}
			tx.Commit()
			db.Exec(`VACUUM operations;`)
			db.Close()
			log.Print("updateDaemon tact clear ", co)
		}
	}
}
