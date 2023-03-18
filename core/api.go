package core

import (
	"dbcore/pbfiles"
	"fmt"
)

type API struct {
	Name  string `yaml:"name"`
	Table string `yaml:"table"`
	Sql   string `yaml:"sql"`
}

func (this *API) Exec(params *pbfiles.SimpleParam) (int64, error) {
	if this.Sql == "" {
		return 0, fmt.Errorf("error sql")
	}

	db := GormDB.Exec(this.Sql, params.Params.AsMap())

	return db.RowsAffected, db.Error
}

func (this *API) Query(params *pbfiles.SimpleParam) ([]map[string]interface{}, error) {
	if this.Table == "" || this.Sql == "" {
		return nil, fmt.Errorf("error sql or table")
	}

	if this.Sql != "" {
		return this.QueryBySql(params)
	}

	return this.QueryByTableName(params)
}

// deprecated
func (this *API) QueryByTableName(params *pbfiles.SimpleParam) ([]map[string]interface{}, error) {
	dbResult := make([]map[string]interface{}, 0)
	db := GormDB.Table(this.Table)

	if params != nil && params.Params != nil {
		for k, v := range params.Params.AsMap() {
			db = db.Where(k, v)
		}
	}

	db = db.Find(&dbResult)

	return dbResult, db.Error
}

func (this *API) QueryBySql(params *pbfiles.SimpleParam) ([]map[string]interface{}, error) {
	dbResult := make([]map[string]interface{}, 0)

	values := map[string]interface{}{}
	if params != nil && params.Params != nil {
		values = params.Params.AsMap()
	}

	db := GormDB.Raw(this.Sql, values).Find(&dbResult)

	return dbResult, db.Error
}
