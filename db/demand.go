package db

import (
	"database/sql"
	"github.com/Rosya-edwica/position-to-demand/models"
	"fmt"
)

// FIXME: ParentId

// Если хотите взять первый скилл, то передайте в качестве lastId значение 0
func (d *Database) GetNextSkill(lastID int, equalID bool) (skill models.Skill) {
	var query  string
	if equalID {
		query = fmt.Sprintf("SELECT id, name, parent_id FROM demand WHERE id = %d ORDER BY id ASC LIMIT 1;", lastID)
	} else {
		query = fmt.Sprintf("SELECT id, name, parent_id FROM demand WHERE id > %d ORDER BY id ASC LIMIT 1;", lastID)
	}
	// Временное значение query. Так как непонятно как быть с теми навыками, у которых есть родитель. И как быть с дубликатами (навык, профессия, город)
	query = fmt.Sprintf("SELECT id, name, parent_id FROM demand WHERE parent_id IS NULL AND id > %d ORDER BY id ASC LIMIT 1;", lastID)
	rows, err := d.Connection.Query(query)
	checkErr(err)
	for rows.Next() {
		var id int
		var parentId sql.NullInt64
		var name string
		err = rows.Scan(&id, &name, &parentId)
		checkErr(err)
		skill = models.Skill{
			Id: id,
			Name: name,
			ParentId: parentId,
		}
	}
	// Если parent_id != NULL, то возвращаем родителя навыка
	if skill.ParentId.Valid {
		return d.GetNextSkill(int(skill.ParentId.Int64), true)
	} else {
		return skill
	}

}