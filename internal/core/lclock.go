package core

import (
	"sync/atomic"
	"time"
)

// HLCTimestamp представляет Hybrid Logical Clock timestamp
// Состоит из физического времени (в наносекундах) и логического счетчика
type HLCTimestamp struct {
	WallTime uint64 // Физическое время в наносекундах с Unix epoch
	Logical  uint64 // Логический счетчик для событий с одинаковым физическим временем
}

// LClock реализует Hybrid Logical Clock (HLC) алгоритм
// HLC комбинирует физическое время с логическими часами для обеспечения
// причинной упорядоченности и приближенности к реальному времени
type LClock struct {
	wallTime uint64 // Последнее известное физическое время
	logical  uint64 // Логический счетчик
}

func NewLClock() *LClock {
	return &LClock{
		wallTime: uint64(time.Now().UnixNano()),
		logical:  0,
	}
}

// GetPhysicalTime возвращает текущее физическое время в наносекундах
func getPhysicalTime() uint64 {
	return uint64(time.Now().UnixNano())
}

// Get возвращает текущее HLC timestamp как единое значение
// Комбинируем wallTime и logical в одно uint64 значение для обратной совместимости
func (lc *LClock) Get() uint64 {
	pt := getPhysicalTime()
	wallTime := atomic.LoadUint64(&lc.wallTime)
	logical := atomic.LoadUint64(&lc.logical)

	// Возвращаем максимум из физического времени и сохраненного wallTime
	if pt > wallTime {
		return pt
	}
	// Если физическое время не продвинулось, используем wallTime + logical как offset
	return wallTime + logical
}

// GetTimestamp возвращает полное HLC timestamp
func (lc *LClock) GetTimestamp() HLCTimestamp {
	return HLCTimestamp{
		WallTime: atomic.LoadUint64(&lc.wallTime),
		Logical:  atomic.LoadUint64(&lc.logical),
	}
}

// SendOrLocal обновляет HLC при локальном событии или отправке сообщения
func (lc *LClock) SendOrLocal() uint64 {
	for {
		pt := getPhysicalTime()
		currentWallTime := atomic.LoadUint64(&lc.wallTime)
		currentLogical := atomic.LoadUint64(&lc.logical)

		var newWallTime, newLogical uint64

		if pt > currentWallTime {
			// Физическое время продвинулось - используем новое время, сбрасываем logical
			newWallTime = pt
			newLogical = 0
		} else {
			// Физическое время не продвинулось - увеличиваем logical счетчик
			newWallTime = currentWallTime
			newLogical = currentLogical + 1
		}

		// Атомарно обновляем оба значения
		if atomic.CompareAndSwapUint64(&lc.wallTime, currentWallTime, newWallTime) {
			atomic.StoreUint64(&lc.logical, newLogical)
			return newWallTime + newLogical
		}
	}
}

// Recv обновляет HLC при получении сообщения с timestamp отправителя
func (lc *LClock) Recv(msgTimestamp uint64) uint64 {
	for {
		pt := getPhysicalTime()
		currentWallTime := atomic.LoadUint64(&lc.wallTime)
		currentLogical := atomic.LoadUint64(&lc.logical)

		var newWallTime, newLogical uint64

		// Извлекаем wallTime из msgTimestamp (используем как wallTime для простоты)
		msgWallTime := msgTimestamp

		if pt > currentWallTime && pt > msgWallTime {
			// Физическое время больше обоих - используем его
			newWallTime = pt
			newLogical = 0
		} else if currentWallTime > msgWallTime {
			// Наше время больше - увеличиваем logical
			newWallTime = currentWallTime
			newLogical = currentLogical + 1
		} else if msgWallTime > currentWallTime {
			// Время сообщения больше - используем его и сбрасываем logical
			newWallTime = msgWallTime
			newLogical = 0
		} else {
			// Времена равны - увеличиваем logical
			newWallTime = currentWallTime
			newLogical = max(currentLogical, 0) + 1
		}

		// Атомарно обновляем
		if atomic.CompareAndSwapUint64(&lc.wallTime, currentWallTime, newWallTime) {
			atomic.StoreUint64(&lc.logical, newLogical)
			return newWallTime + newLogical
		}
	}
}
