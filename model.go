package main

import "time"

/**
Track (Довідник пісень):
id: Primary Key.
title: Назва композиції.
artist: Виконавець.
price: Поточна вартість (Decimal).



PlaybackLog (Історія програвань):
id: Primary Key.
track: Зв'язок Many-to-One з сутністю Track.
played_at: Дата та час події.
amount_paid: Сума, яка була фактично сплачена в момент замовлення.
*/

type Track struct {
	ID     int
	Title  string
	Artist string
	Price  float64
	Plays  []*PlaybackLog
}

type PlaybackLog struct {
	ID         int
	Track      *Track
	PlayedAt   time.Time
	AmountPaid float64
}
