package core

import (
	"time"
)

// Interval интерфейс для работы с интервалами
type Interval interface {
	GetID() string
	GetFrom() time.Time
	GetTo() time.Time
}

// Функция для проверки пересечения двух интервалов
func intersect(a, b Interval) bool {
	return a.GetFrom().Before(b.GetTo()) && b.GetFrom().Before(a.GetTo())
}

// FindNonOverlappingIntervals Рекурсивная функция для поиска всех непересекающихся комбинаций
func FindNonOverlappingIntervals(intervals []Interval, index int, current []Interval, result *[][]Interval) {
	if index == len(intervals) {
		if len(current) > 0 {
			*result = append(*result, append([]Interval{}, current...))
		}
		return
	}
	// Пропустить текущий интервал и продолжить
	FindNonOverlappingIntervals(intervals, index+1, current, result)

	// Добавить текущий интервал, если он не пересекается с уже добавленными
	overlaps := false
	for _, i := range current {
		if intersect(intervals[index], i) {
			overlaps = true
			break
		}
	}
	if !overlaps {
		current = append(current, intervals[index])
		FindNonOverlappingIntervals(intervals, index+1, current, result)
		current = current[:len(current)-1] // Вернуться назад
	}
}
