package main

import (
	"errors"
	"fmt"
)

type Accessor interface {
	Get(day string) (Rate, error)
	Update(Rate) (Rate, error)
	Delete(day string) error
	Create(Rate) (Rate, error)
}

type Rate struct {
	Base  string `json:"base"`
	Date  string `json:"date"`
	Rates Rates  `json:"rates"`
}

type Rates struct {
	USD float64 `json:"USD"`
	GBP float64 `json:"GBP"`
	EUR float64 `json:"EUR"`
	CZK float64 `json:"CZK"`
}

var (
	ErrNotFound = errors.New("not found")
)

type MemDB struct {
	data map[string]Rate
}

func (m MemDB) Get(day string) (Rate, error) {
	rate, ok := m.data[day]
	if !ok {
		return Rate{}, ErrNotFound
	}
	return rate, nil
}

func (m MemDB) Update(rate Rate) (Rate, error) {
	r, ok := m.data[rate.Date]
	if !ok {
		return r, ErrNotFound
	}
	m.data[rate.Date] = rate
	return rate, nil
}

func (m MemDB) Delete(day string) error {
	_, ok := m.data[day]
	if !ok {
		return ErrNotFound
	}
	delete(m.data, day)
	return nil
}

func (m MemDB) Create(rate Rate) (Rate, error) {
	m.data[rate.Date] = rate
	return rate, nil
}

func main() {
	memDb := MemDB{data:make(map[string]Rate)}
	rate1 := Rate{"","",Rates{}}
	memDb.Create(rate1)

	fmt.Println(memDb.Get(rate1.Date))

}
