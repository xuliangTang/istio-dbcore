package core

import (
	"dbcore/pbfiles"
)

type API struct {
	Name  string `yaml:"name"`
	Table string `yaml:"table"`
	Sql   string `yaml:"sql"`
}

func (this *API) QueryByTableName(params *pbfiles.SimpleParam) ([]map[string]interface{}, error) {
	dbResult := make([]map[string]interface{}, 0)
	db := GormDB.Table(this.Table)

	if params != nil {
		for k, v := range params.Params.AsMap() {
			db = db.Where(k, v)
		}
	}

	db = db.Find(&dbResult)

	return dbResult, db.Error
}
