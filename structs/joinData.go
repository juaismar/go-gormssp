package structs

type JoinData struct {
	Table string //name of column
	Alias string //id of column in client (int or string)
	Query string //case sensitive - optional default false
}
