package structs

// MessageDataTable is the response object
type MessageDataTable struct {
	Draw            int           `json:"draw"`
	RecordsTotal    int64         `json:"recordsTotal"`
	RecordsFiltered int64         `json:"recordsFiltered"`
	Data            []interface{} `json:"data,nilasempty"`
}
