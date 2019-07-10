package main

import (
	"errors"
	echo2 "github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"time"
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
	Day  string
	From string
	To   string
}

func (s *Service) ensureRateDay(day string) (Rate, error) {
	s.access.Get(day)
	createdRate, err := s.ratesApi.Get(day)
	if err != nil {
		return Rate{}, errors.New("Server error")
	}
	return s.access.Create(createdRate)
}

func (s *Service) GetRate(c echo2.Context) error {
	var req GetRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.Day != "" {
		rate, _ := s.ensureRateDay(req.Day)
		return c.JSON(http.StatusOK, rate)
	} else {
		rates, _ := s.findRange(req.From, req.To)
		return c.JSON(http.StatusOK, rates)
	}

}

const (
	DayLoyout = "2006-01-01"
)

func (service *Service) findRange(from, to string) ([]Rate, error) {
	var rates []Rate
	fromDay, err := time.Parse(DayLoyout, from)
	if err != nil {
		return rates, err
	}
	toDay, err := time.Parse(DayLoyout, to)
	if err != nil {
		return rates, err
	}

	for day := fromDay; day.Before(toDay); day = day.AddDate(0, 0, 1) {
		rate, err := service.ensureRateDay(day.String())
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
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

func NewService(access Accessor, rateService *RatesApiService) *Service {
	return &Service{access, rateService}
}

type Service struct {
	access   Accessor
	ratesApi *RatesApiService
}

func main() {
	viper.SetDefault("port", "8080")
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.json")

	db := NewMemDB()
	rateSvc := &RatesApiService{}

	svc := NewService(db, rateSvc)
	e := echo2.New()
	e.GET("/", svc.GetRate)
	e.POST("/", svc.CreateRate)
	e.PUT("/", svc.UpdateRate)
	e.DELETE("/", svc.DeleteRate)
	e.Logger.Fatal(e.Start(":8080"))
}
