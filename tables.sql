


--Таблица пользователей и суммарных счетов.
CREATE TABLE users(
  user_name CHAR(50) PRIMARY KEY ,
  deposit money
);

CREATE unique INDEX user_index
ON users (user_name);

--Не актуальный джекпот полученный во время последнего обновления
CREATE TABLE old_jackpot(
  value money
);
insert into old_jackpot values (0);

-- Таблица свежих ставок.
-- Раз в 10 секунд (время стоит подобрать в зависимоти от нагрузки), запускается updateDaemon записи очищаются.
-- Деньги пользователей идут на их счета в users, джекпот прибавляется к old_jackpot
-- autovacuum нужен редко поскольку updateDaemon делает обычный vacuum

CREATE TABLE operations (
  id SERIAL primary key,
  user_name CHAR(50),
  deposit money,
  jackpot_part money
) with (FILLFACTOR = 90);

-- Вьюшка для быстрого получения актуального джектпота.
create view real_jackpot as
    select sum(t.jp) from
        (
             select jackpot_part as jp from operations
             union
             select value as jp from old_jackpot
        ) t;


-------------------------------------------------

INSERT INTO users(user_name, deposit)
select user_name, sum(deposit) from operations group by user_name
ON CONFLICT (user_name)
DO
    UPDATE
        SET deposit = EXCLUDED.deposit + users.deposit;

update jackpot set jackpot.value = jackpot.value + t.sum from ()


UPDATE old_jackpot
SET value = oj.value + t.nj
FROM old_jackpot as oj
CROSS JOIN
    (
        SELECT sum(jackpot_part) as nj
        FROM operations
    ) t