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
}

type PlaybackLog struct {
	ID         int
	TrackID    int
	PlayedAt   time.Time
	AmountPaid float64
}

type TrackRepository interface {
	GetTrackByID(id int) (*Track, error)
	UpdateTrackPrice(id int, newPrice float64) error
}

type PlaybackLogRepository interface {
	CreateLog(log PlaybackLog)
	GetAllLogs() []PlaybackLog
}
