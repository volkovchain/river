package entity

import (
	"time"
)

type Salary struct {
	ID       int64
	Status   string
	CreateAt *time.Time
}

type Employee struct {
	ID           int64
	Name         string
	SalaryAmount int
	Addr         string
}

type Payment struct {
	ID         int64
	EmployeeID int64
	SalaryID   int64
	Amount     int64
	Addr       string
	Status     string
	Error      string
	CreateAt   *time.Time
}
