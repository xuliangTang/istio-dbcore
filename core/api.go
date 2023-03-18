package core

type API struct {
	Name  string `yaml:"name"`
	Table string `yaml:"table"`
	Sql   string `yaml:"sql"`
}

func (this *API) QueryByTableName() ([]map[string]interface{}, error) {
	dbResult := make([]map[string]interface{}, 0)
	db := GormDB.Table(this.Table).Find(&dbResult)
	return dbResult, db.Error
}
