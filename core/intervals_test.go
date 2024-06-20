package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type interval struct {
	ID   string
	From time.Time
	To   time.Time
}

func (i *interval) GetID() string {
	return i.ID
}

func (i *interval) GetFrom() time.Time {
	return i.From
}

func (i *interval) GetTo() time.Time {
	return i.To
}

func newInterval(id string, from time.Time, to time.Time) Interval {
	return &interval{
		ID:   id,
		From: from,
		To:   to,
	}
}

func getIds(intervals []Interval) string {
	result := ""
	for _, i := range intervals {
		result += fmt.Sprintf("%s", i.GetID())
	}
	return result
}

func day(day int) time.Time {
	return time.Date(2024, 1, day, 0, 0, 0, 0, time.Local)
}

func TestFindNonOverlappingIntervals(t *testing.T) {

	exp := Set[string]{}

	exp.Add("A", "AF", "AF", "B", "BD", "BDF", "BDF", "BE", "BEF", "BF", "C", "CE", "CE", "CEF", "CF", "D", "DF", "E", "EF", "F")

	var sut = []Interval{
		newInterval("A", day(1), day(7)),
		newInterval("B", day(2), day(4)),
		newInterval("C", day(3), day(5)),
		newInterval("D", day(4), day(6)),
		newInterval("E", day(5), day(7)),
		newInterval("F", day(7), day(8)),
	}

	var result [][]Interval
	FindNonOverlappingIntervals(sut, 0, []Interval{}, &result)

	// Печать результатов

	act := Set[string]{}
	for _, combination := range result {
		act.Add(getIds(combination))
	}

	assert.Equal(t, exp, act)
}
