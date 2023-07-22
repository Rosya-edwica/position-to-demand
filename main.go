package main

import (
	"fmt"
	"os"

	"github.com/Rosya-edwica/position-to-demand/db"
	"github.com/Rosya-edwica/position-to-demand/logger"

	"github.com/joho/godotenv"
)

const InsertLimit = 2000

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
		for i:=0; i<len(statistic); i+=InsertLimit {
			group := statistic[i:]
			if len(group) > InsertLimit {
				group = group[:InsertLimit]
			}
			database.SaveSkillStatistic(group)
		}
		fmt.Println("Done for id:", skill.Id)
		lastId = skill.Id
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
