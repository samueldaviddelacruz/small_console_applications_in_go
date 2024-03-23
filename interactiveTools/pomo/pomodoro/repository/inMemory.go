//go:build inmemory || containers
// +build inmemory containers

package repository

import (
	"fmt"
	"small_console_applications_go/interactiveTools/pomo/pomodoro"
	"strings"
	"sync"
	"time"
)

type inMemoryRepo struct {
	sync.RWMutex
	intervals []pomodoro.Interval
}

func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		intervals: []pomodoro.Interval{},
	}
}
func (r *inMemoryRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()
	i.ID = int64(len(r.intervals) + 1)
	r.intervals = append(r.intervals, i)
	return i.ID, nil
}
func (r *inMemoryRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()
	if i.ID == 0 {
		return fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, i.ID)
	}
	r.intervals[i.ID-1] = i
	return nil
}
func (r *inMemoryRepo) ByID(id int64) (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()
	i := pomodoro.Interval{}
	if id == 0 {
		return i, fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, id)
	}
	i = r.intervals[id-1]
	return i, nil
}
func (r *inMemoryRepo) Last() (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()
	i := pomodoro.Interval{}
	if len(r.intervals) == 0 {
		return i, pomodoro.ErrNoIntervals
	}
	return r.intervals[len(r.intervals)-1], nil
}
func (r *inMemoryRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()
	data := []pomodoro.Interval{}
	for k := len(r.intervals) - 1; k >= 0; k-- {
		if r.intervals[k].Category == pomodoro.CategoryPomodoro {
			continue
		}
		data = append(data, r.intervals[k])
		if len(data) == n {
			break
		}
	}
	return data, nil
}

func (r *inMemoryRepo) CategorySummary(day time.Time, filter string) (time.Duration, error) {
	r.RLock()
	defer r.RUnlock()
	var d time.Duration
	filter = strings.Trim(filter, "%")
	for _, i := range r.intervals {
		if i.StartTime.Year() == day.Year() && i.StartTime.YearDay() == day.YearDay() {
			if strings.Contains(i.Category, filter) {
				d += i.ActualDuration
			}
		}
	}
	return d, nil
}
