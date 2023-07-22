package db

import (
	"fmt"

	"github.com/Rosya-edwica/position-to-demand/logger"
	"github.com/Rosya-edwica/position-to-demand/models"
)

func (d *Database) GetVacanciesContainSkill(skill models.Skill) (skillStatistic []models.Statistic) {
	query := fmt.Sprintf(`SELECT position_id, city_id, count(*) as vacancy_count, vacancy_date, parsing_date 
	FROM h_vacancy 
	WHERE LOWER(key_skills) LIKE '%%%s%%' AND position_id != 0 AND city_id != 0 
	GROUP by position_id, city_id
	ORDER BY count(*) DESC`, skill.Name)
	rows, err := d.Connection.Query(query)
	checkErr(err)

	for rows.Next() {
		var position_id, city_id, vacancies_count int
		var vacancy_date, parsing_date string

		err = rows.Scan(&position_id, &city_id,  &vacancies_count, &vacancy_date, &parsing_date)
		checkErr(err)
		skillStatistic = append(skillStatistic, models.Statistic{
			SkillID: skill.Id,
			PositionID: position_id,
			CityID: city_id,
			LastDate: vacancy_date,
			VacanciesCount: vacancies_count,
			ParsingDate: parsing_date,
		})
	}

	return
}

func (d *Database) SaveSkillStatistic(skillStatistic []models.Statistic) {
	if len(skillStatistic) == 0 {
		return
	}
	vals := []interface{}{}
	query := `INSERT INTO n_position_to_demand(
		demand_id, position_id, city_id, count_in_vac, is_custom, last_listed
		) VALUES `

	for _, stat := range skillStatistic {
		query += "(?, ?, ?, ?, ?, ?),"
		var date string
		if stat.LastDate == "" {
			date = stat.ParsingDate
		} else {
			date = stat.LastDate
		}
		vals = append(vals, stat.SkillID, stat.PositionID, stat.CityID, stat.VacanciesCount, false, date)
	}
	query = query[0:len(query) - 1]
	stmt, err := d.Connection.Prepare(query)
	checkErr(err)

	_, err = stmt.Exec(vals...)
	checkErr(err)

	logger.Log.Printf("Успешно сохранили %d строк", len(skillStatistic))
}
