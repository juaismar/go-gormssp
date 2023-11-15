package structs

// Data is a line in map that link the database field with datatable field
type Data struct {
	Db        string                                                                  //name of column
	Dt        interface{}                                                             //id of column in client (can be int or string)
	Cs        bool                                                                    //case sensitive - optional default false
	Sf        string                                                                  //Search Function - for custom functions declared in your ddbb
	Formatter func(data interface{}, row map[string]interface{}) (interface{}, error) //can run code in this function to edit results - optional, can be nil
}

type DataParsed struct {
	Data

	ParsedDT string
}
