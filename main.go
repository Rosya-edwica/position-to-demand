package main

import (
	"github.com/Rosya-edwica/position-to-demand/db"
	"github.com/Rosya-edwica/position-to-demand/logger"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	checkErr(err)
}

func main() {
	database := db.Database{
		Name:     os.Getenv("MYSQL_DATABASE"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	}
	database.NewConnection()
	defer database.CloseConnection()

	lastId := 0
	for {
		skill := database.GetNextSkill(lastId, false)
		if skill.Name == "" {
			logger.Log.Printf("Навыки закончились на ID: %d", lastId)
			break
		}

		logger.Log.Printf("Навык: '%s'. Его ID: %d", skill.Name, skill.Id)
		statistic := database.GetVacanciesContainSkill(skill)
		logger.Log.Printf("Навык упоминается %d раз (Для каждой профессии в каждом городе)", len(statistic))		
		database.SaveSkillStatistic(statistic)
		lastId = skill.Id
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
