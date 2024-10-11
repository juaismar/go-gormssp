Dialect BigQuery need add on Opt field the result of this example function.

func GormSSPOpt(dataset, tableName, tableAlias string) map[string]interface{} {
	opt := make(map[string]interface{})
	tableInfo := make(map[string]map[string]string)
	tableInfo[tableAlias] = map[string]string{
		"Dataset":   dataset,
		"TableName": tableName,
	}
	opt["TableInfo"] = tableInfo

	return opt
}