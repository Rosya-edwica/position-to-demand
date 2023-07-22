package models

import "database/sql"

type Statistic struct {
	SkillID        int
	PositionID     int
	CityID         int
	VacanciesCount int
	LastDate       string
}

type Skill struct {
	Id       int
	ParentId sql.NullInt64
	Name     string
}
