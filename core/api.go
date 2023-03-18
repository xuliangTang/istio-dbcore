package core

import (
	"dbcore/pbfiles"
	"fmt"
	"gorm.io/gorm"
)

type Select struct {
	Sql string `yaml:"sql"`
}

type API struct {
	Name   string  `yaml:"name"`
	Table  string  `yaml:"table"`
	Sql    string  `yaml:"sql"`
	Select *Select `yaml:"select"`
}

func (this *API) ExecBySql(params *pbfiles.SimpleParam) (int64, map[string]interface{}, error) {
	if this.Sql == "" {
		return 0, nil, fmt.Errorf("error sql")
	}

	if this.Select != nil {
		selectKey := make(map[string]interface{})
		var rowsAffected int64 = 0
		err := GormDB.Transaction(func(tx *gorm.DB) error {
			db := tx.Exec(this.Sql, params.Params.AsMap())
			if db.Error != nil {
				return db.Error
			}
			rowsAffected = db.RowsAffected

			db = tx.Raw(this.Select.Sql, params.Params.AsMap()).Find(&selectKey)
			return db.Error
		})

		if err != nil {
			return 0, nil, err
		}

		return rowsAffected, selectKey, nil

	} else {
		db := GormDB.Exec(this.Sql, params.Params.AsMap())
		return db.RowsAffected, nil, db.Error
	}
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
