package main

import (
	"errors"
	echo2 "github.com/labstack/echo/v4"
	"net/http"
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

func NewMemDB() MemDB {
	return MemDB{
		data: map[string]Rate{},
	}
}

type GetRequest struct {
	Day string
}

func (s *Service) GetRate(c echo2.Context) error {
	var req GetRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	rate, e := s.access.Get(req.Day)
	if e != nil {
		return e
	}

	return c.JSON(http.StatusOK, rate)
}

func (s *Service) CreateRate(c echo2.Context) error {
	var rate Rate
	if err := c.Bind(&rate); err != nil {
		return err
	}
	r, e := s.access.Create(rate)
	if e != nil {
		return e
	}

	return c.JSON(http.StatusCreated, r)
}

func (s *Service) UpdateRate(c echo2.Context) error {
	var rate Rate
	if err := c.Bind(&rate); err != nil {
		return err
	}
	r, e := s.access.Update(rate)
	if e != nil {
		return e
	}

	return c.JSON(http.StatusOK, r)
}

func (s *Service) DeleteRate(c echo2.Context) error {
	var req GetRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	e := s.access.Delete(req.Day)
	if e != nil {
		return e
	}

	return c.NoContent(http.StatusNoContent)
}

func NewService(access Accessor) *Service {
	return &Service{access}
}

type Service struct {
	access Accessor
}

func main() {
	db := NewMemDB()
	db.Create(Rate{
		Base:  "",
		Date:  "1",
		Rates: Rates{},
	})
	svc := NewService(db)
	e := echo2.New()
	e.GET("/", svc.GetRate)
	e.POST("/", svc.CreateRate)
	e.PUT("/", svc.UpdateRate)
	e.DELETE("/", svc.DeleteRate)
	e.Logger.Fatal(e.Start(":8080"))
}
