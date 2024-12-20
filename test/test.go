package test

import (
	"fmt"
	"time"

	ssp "github.com/juaismar/go-gormssp"
	engine "github.com/juaismar/go-gormssp/engine"
	"github.com/juaismar/go-gormssp/structs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

const layoutISO = "2006-01-02"

// ControllerEmulated emulate the beego controller
type ControllerEmulated struct {
	Params map[string]string
}

// GetString emulate the beego controoller method
func (c *ControllerEmulated) GetString(key string, def ...string) string {
	return c.Params[key]
}

// FunctionsTest internal function test
func FunctionsTest() {
	Describe("flated", func() {
		It("returns Empty", func() {

			var whereArray []string

			result := engine.Flated(whereArray)

			Expect(result).To(Equal(""))
		})
		It("returns one query", func() {

			var whereArray []string
			whereArray = append(whereArray, "number = 1")

			result := engine.Flated(whereArray)

			Expect(result).To(Equal("number = 1"))
		})
		It("returns two query", func() {

			var whereArray []string
			whereArray = append(whereArray, "number = 1")
			whereArray = append(whereArray, "name = 'John'")

			result := engine.Flated(whereArray)

			Expect(result).To(Equal("number = 1 AND name = 'John'"))
		})
	})
	Describe("search", func() {
		It("returns -1", func() {

			columns := []structs.DataParsed{
				{Data: structs.Data{Db: "name", Dt: 0, Formatter: nil}, ParsedDT: "0"},
				{Data: structs.Data{Db: "role", Dt: "role", Formatter: nil}, ParsedDT: "role"},
				{Data: structs.Data{Db: "email", Dt: 2, Formatter: nil}, ParsedDT: "2"},
			}
			result := engine.Search(columns, "")

			Expect(result).To(Equal(-1))
		})
		It("returns -1", func() {

			columns := []structs.DataParsed{
				{Data: structs.Data{Db: "name", Dt: 0, Formatter: nil}, ParsedDT: "0"},
				{Data: structs.Data{Db: "role", Dt: "role", Formatter: nil}, ParsedDT: "role"},
				{Data: structs.Data{Db: "email", Dt: 2, Formatter: nil}, ParsedDT: "2"},
			}
			result := engine.Search(columns, "instrument")

			Expect(result).To(Equal(-1))
		})
		It("returns 1", func() {

			columns := []structs.DataParsed{
				{Data: structs.Data{Db: "name", Dt: 0, Formatter: nil}, ParsedDT: "0"},
				{Data: structs.Data{Db: "role", Dt: "role", Formatter: nil}, ParsedDT: "role"},
				{Data: structs.Data{Db: "email", Dt: 2, Formatter: nil}, ParsedDT: "2"},
			}
			result := engine.Search(columns, "role")

			Expect(result).To(Equal(1))
		})
		It("returns 0", func() {

			columns := []structs.DataParsed{
				{Data: structs.Data{Db: "name", Dt: 0, Formatter: nil}, ParsedDT: "0"},
				{Data: structs.Data{Db: "role", Dt: "role", Formatter: nil}, ParsedDT: "role"},
				{Data: structs.Data{Db: "email", Dt: 2, Formatter: nil}, ParsedDT: "2"},
			}
			result := engine.Search(columns, "0")

			Expect(result).To(Equal(0))
		})

	})
}

// ComplexFunctionTest test for Complex method
func ComplexFunctionTest(db *gorm.DB) {
	Describe("Complex", func() {
		//filter whereall (where in all queries)
		It("returns fun only Juan Joaquin Laura", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "62"
			mapa["start"] = "0"
			mapa["length"] = "4"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
			}
			whereResult := make([]string, 0)
			whereJoin := make([]structs.JoinData, 0)

			whereAll := make([]string, 0)
			whereAll = append(whereAll, "fun = '1'")

			result, err := ssp.Complex(&c, db, "users", columns, whereResult, whereAll, whereJoin, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(62))
			Expect(result.RecordsTotal).To(Equal(int64(3)))
			Expect(result.RecordsFiltered).To(Equal(int64(3)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Joaquin"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Laura"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		//filter whereResult (where in only result sended)
		It("returns fun only Juan Joaquin Laura", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "62"
			mapa["start"] = "0"
			mapa["length"] = "5"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
			}
			whereResult := make([]string, 0)
			whereResult = append(whereResult, "fun = '1'")

			whereJoin := make([]structs.JoinData, 0)
			whereAll := make([]string, 0)

			result, err := ssp.Complex(&c, db, "users", columns, whereResult, whereAll, whereJoin, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(62))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(3)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Joaquin"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Laura"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		//check join compatibility

		It("Join test", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "62"
			mapa["start"] = "0"
			mapa["length"] = "3"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "users.name", Dt: 0, Formatter: nil},
				{Db: "pets.name", Dt: 1, Formatter: nil},
				{Db: "name", Dt: 2, Formatter: nil},
			}
			whereResult := make([]string, 0)

			whereJoin := make([]structs.JoinData, 0)

			whereJoin = append(whereJoin, structs.JoinData{
				Table: "pets",
				Alias: "",
				Query: "left join pets on pets.master_id = users.uuid",
			})

			whereAll := make([]string, 0)

			result, err := ssp.Complex(&c, db, "users", columns, whereResult, whereAll, whereJoin, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(62))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(6)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			row["1"] = "Cerverus"
			row["2"] = "Juan"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "JuAn"
			row["1"] = "Mikey"
			row["2"] = "JuAn"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Joaquin"
			row["1"] = "Epona"
			row["2"] = "Joaquin"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})

		It("Join alias", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "62"
			mapa["start"] = "0"
			mapa["length"] = "3"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "users.name", Dt: 0, Formatter: nil},
				{Db: "animal.name", Dt: 1, Formatter: nil},
				{Db: "name", Dt: 2, Formatter: nil},
				{Db: "beast.name", Dt: 3, Formatter: nil},
			}
			whereResult := make([]string, 0)

			whereJoin := make([]structs.JoinData, 0)

			whereJoin = append(whereJoin, structs.JoinData{
				Table: "pets",
				Alias: "animal",
				Query: "left join pets AS animal on animal.master_id = users.uuid",
			})

			whereJoin = append(whereJoin, structs.JoinData{
				Table: "pets",
				Alias: "beast",
				Query: "left join pets AS beast on beast.master_id = users.uuid",
			})

			whereAll := make([]string, 0)

			result, err := ssp.Complex(&c, db, "users", columns, whereResult, whereAll, whereJoin, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(62))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(6)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			row["1"] = "Cerverus"
			row["2"] = "Juan"
			row["3"] = "Cerverus"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "JuAn"
			row["1"] = "Mikey"
			row["2"] = "JuAn"
			row["3"] = "Mikey"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Joaquin"
			row["1"] = "Epona"
			row["2"] = "Joaquin"
			row["3"] = "Epona"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})

		It("Join search test", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "62"
			mapa["start"] = "0"
			mapa["length"] = "3"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "1"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][search][value]"] = "Cerverus"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "users.name", Dt: 0, Formatter: nil},
				{Db: "pets.name", Dt: 1, Formatter: nil},
				{Db: "name", Dt: 2, Formatter: nil},
			}
			whereResult := make([]string, 0)

			whereJoin := make([]structs.JoinData, 0)

			whereJoin = append(whereJoin, structs.JoinData{
				Table: "pets",
				Alias: "",
				Query: "left join pets on pets.master_id = users.uuid",
			})

			whereAll := make([]string, 0)

			result, err := ssp.Complex(&c, db, "users", columns, whereResult, whereAll, whereJoin, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(62))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(1)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			row["1"] = "Cerverus"
			row["2"] = "Juan"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
	})
}

// RegExpTest test for regular expression
func RegExpTest(db *gorm.DB) {
	Describe("RegExp", func() {
		It("Global search regex", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "1"
			mapa["order[0][dir]"] = "desc"

			mapa["search[value]"] = "^Eze"
			mapa["search[regex]"] = "true"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
				{Db: "instrument", Dt: 1, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(1)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Ezequiel"
			row["1"] = "Trompeta"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		It("returns names whit 5 chars (regex)", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][orderable]"] = "true"
			mapa["columns[0][search][value]"] = "^.{5}$"
			mapa["columns[0][search][regex]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(2)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Laura"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Marta"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		It("returns names 2 names", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][orderable]"] = "true"
			mapa["columns[0][search][value]"] = "Marta|Laura"
			mapa["columns[0][search][regex]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(2)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Laura"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Marta"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		It("returns 2 ages int", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][search][value]"] = "13|18"
			mapa["columns[0][search][regex]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "age", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(2)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = int64(18)
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = int64(13)
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		It("returns 2 money float", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][orderable]"] = "true"
			mapa["columns[0][search][value]"] = "22.11|0.1"
			mapa["columns[0][search][regex]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "money", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(2)))

			f1 := result.Data[0].(map[string]interface{})["0"].(float64)
			f2 := result.Data[1].(map[string]interface{})["0"].(float64)
			Expect(0.09 < f1 && f1 < 0.11).To(BeTrue())
			Expect(22.109 < f2 && f2 < 22.111).To(BeTrue())
			Expect(result.Data).To(HaveLen(2))
		})
		It("returns 2 money float", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][orderable]"] = "true"
			mapa["columns[0][search][value]"] = "22,11|0,1"
			mapa["columns[0][search][regex]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "money", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(2)))

			f1 := result.Data[0].(map[string]interface{})["0"].(float64)
			f2 := result.Data[1].(map[string]interface{})["0"].(float64)
			Expect(0.09 < f1 && f1 < 0.11).To(BeTrue())
			Expect(22.109 < f2 && f2 < 22.111).To(BeTrue())
			Expect(result.Data).To(HaveLen(2))
		})
	})
}

// Types test for types
func Types(db *gorm.DB) {
	Describe("Types", func() {
		Describe("uint", func() {
			It("returns 2 Age 15", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "15"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "age", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(2)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "JuAn"
				row["1"] = int64(15)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				row["1"] = int64(15)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("int", func() {
			It("returns 1 Candies 10", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "10"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "candies", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(1)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Joaquin"
				row["1"] = int64(10)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("int 8", func() {
			It("returns 2 users", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "1"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "toys", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(2)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "JuAn"
				row["1"] = int64(1)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				row["1"] = int64(1)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("bool", func() {
			It("returns fun only Juan Joaquin Laura", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "true"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "fun", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(3)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = true
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Joaquin"
				row["1"] = true
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Laura"
				row["1"] = true
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("float32", func() {
			It("returns money only Juan Marta", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "2.0"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "money", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(2)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = float64(2.0)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				row["1"] = float64(2.0)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
			It("returns all with decimals", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "money", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = float64(2.0)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "JuAn"
				row["1"] = float64(3.0999999046325684)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Joaquin"
				row["1"] = float64(3.4000000953674316)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Ezequiel"
				row["1"] = float64(22.110000610351562)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				row["1"] = float64(2.0)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Laura"
				row["1"] = float64(0.10000000149011612)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("float64", func() {
			It("returns bitcoins only Juan Marta", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "3.0"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "bitcoins", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(2)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = float64(3.0)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				row["1"] = float64(3.0)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
			It("returns all with decimals", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "bitcoins", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = float64(3.0)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "JuAn"
				row["1"] = float64(4.3)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Joaquin"
				row["1"] = float64(7.18)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Ezequiel"
				row["1"] = float64(82.14)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				row["1"] = float64(3.0)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Laura"
				row["1"] = float64(22.71)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("time.TIME", func() {
			It("returns a time and formatter", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "62"
				mapa["start"] = "0"
				mapa["length"] = "1"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "birth_date", Dt: 0, Formatter: func(
						data interface{}, row map[string]interface{}) (interface{}, error) {
						time := data.(time.Time)
						var err error
						Expect(time.Format(layoutISO)).To(Equal("2011-12-11"))
						return time, err
					}},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(62))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))
			})
		})
		Describe("UUID", func() {
			It("returns Juan", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "bfe44cb2-c65c-4f37-9672-8437b6718d70"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "uuid", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(1)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = "bfe44cb2-c65c-4f37-9672-8437b6718d70"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
	})
}

// SimpleFunctionTest test for ssp.Simplex method
func SimpleFunctionTest(db *gorm.DB) {
	Describe("Simple and basic features", func() {
		It("returns from 0 to 4", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "62"
			mapa["start"] = "0"
			mapa["length"] = "4"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(62))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(6)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "JuAn"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Joaquin"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Ezequiel"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		Describe("Length is negative", func() {
			It("returns from 10 elements", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "62"
				mapa["start"] = "0"
				mapa["length"] = "-1"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(62))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "JuAn"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Joaquin"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Ezequiel"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Laura"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Start is negative", func() {
			It("returns from 0 to 4", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "62"
				mapa["start"] = "-1"
				mapa["length"] = "4"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(62))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "JuAn"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Joaquin"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Ezequiel"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Paginate", func() {
			It("returns from 2 to 6", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "63"
				mapa["start"] = "2"
				mapa["length"] = "4"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(63))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Joaquin"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Ezequiel"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Laura"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Global search", func() {
			It("returns 2 Juan", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["search[value]"] = "uAn"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = ""

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = ""

				mapa["columns[2][data]"] = "2"
				mapa["columns[2][searchable]"] = "true"
				mapa["columns[2][search][value]"] = ""

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "instrument", Dt: 1, Formatter: nil},
					{Db: "age", Dt: 2, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(2)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = "Tambor"
				row["2"] = int64(10)
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "JuAn"
				row["1"] = "Trompeta"
				row["2"] = int64(15)
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Multiple individual search", func() {
			It("returns 1 Juan", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = "Juan"

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = "Tambor"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "instrument", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(1)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = "Tambor"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		It("global search and individual search together", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["search[value]"] = "uAn"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][search][value]"] = ""

			mapa["columns[1][data]"] = "1"
			mapa["columns[1][searchable]"] = "true"
			mapa["columns[1][search][value]"] = "Tambor"

			mapa["columns[2][data]"] = "2"
			mapa["columns[2][searchable]"] = "true"
			mapa["columns[2][search][value]"] = ""

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
				{Db: "instrument", Dt: 1, Formatter: nil},
				{Db: "age", Dt: 2, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(1)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			row["1"] = "Tambor"
			row["2"] = int64(10)
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
		Describe("Naming a row", func() {
			It("returns all", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "3"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[supername][data]"] = "0"
				mapa["columns[supername][searchable]"] = "true"
				mapa["columns[supername][search][value]"] = ""

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: "supername", Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["supername"] = "Juan"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["supername"] = "JuAn"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["supername"] = "Joaquin"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Search LIKE string case insensitive", func() {
			It("returns 2 Juan", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = "uAn"

				mapa["columns[1][data]"] = "1"
				mapa["columns[1][searchable]"] = "true"
				mapa["columns[1][search][value]"] = ""

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "instrument", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(2)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = "Tambor"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "JuAn"
				row["1"] = "Trompeta"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Search on varchar LIKE string case insensitive", func() {
			It("returns 2 Tambor", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "1"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = "ambor"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: nil},
					{Db: "instrument", Dt: 1, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(2)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "Juan"
				row["1"] = "Tambor"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "Marta"
				row["1"] = "Tambor"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Search LIKE string case sensitive", func() {
			It("returns 2 Juan", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "64"
				mapa["start"] = "0"
				mapa["length"] = "10"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				mapa["columns[0][data]"] = "0"
				mapa["columns[0][searchable]"] = "true"
				mapa["columns[0][search][value]"] = "uAn"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Cs: true, Formatter: nil},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(64))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(1)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "JuAn"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		Describe("Format", func() {
			It("return name whit prefix and age", func() {

				mapa := make(map[string]string)
				mapa["draw"] = "62"
				mapa["start"] = "0"
				mapa["length"] = "4"
				mapa["order[0][column]"] = "0"
				mapa["order[0][dir]"] = "asc"

				c := ControllerEmulated{Params: mapa}

				columns := []structs.Data{
					{Db: "name", Dt: 0, Formatter: func(
						data interface{}, row map[string]interface{}) (interface{}, error) {
						return fmt.Sprintf("PREFIX_%v_%v", data, row["age"]), nil
					}},
				}
				result, err := ssp.Simple(&c, db, "users", columns, nil)

				Expect(err).To(BeNil())
				Expect(result.Draw).To(Equal(62))
				Expect(result.RecordsTotal).To(Equal(int64(6)))
				Expect(result.RecordsFiltered).To(Equal(int64(6)))

				testData := make([]interface{}, 0)
				row := make(map[string]interface{})
				row["0"] = "PREFIX_Juan_10"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "PREFIX_JuAn_15"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "PREFIX_Joaquin_18"
				testData = append(testData, row)
				row = make(map[string]interface{})
				row["0"] = "PREFIX_Ezequiel_13"
				testData = append(testData, row)

				Expect(result.Data).To(Equal(testData))
			})
		})
		It("Ordering by instrument asc", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "1"
			mapa["order[0][dir]"] = "asc"

			mapa["search[value]"] = "uAn"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][orderable]"] = "true"
			mapa["columns[0][search][value]"] = ""

			mapa["columns[1][data]"] = "0"
			mapa["columns[1][searchable]"] = "true"
			mapa["columns[1][orderable]"] = "true"
			mapa["columns[1][search][value]"] = ""

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
				{Db: "instrument", Dt: 1, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(2)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Juan"
			row["1"] = "Tambor"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "JuAn"
			row["1"] = "Trompeta"
			testData = append(testData, row)
			//
			Expect(result.Data).To(Equal(testData))
		})
		It("Ordering by instrument desc", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "1"
			mapa["order[0][dir]"] = "desc"

			mapa["search[value]"] = "uAn"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][orderable]"] = "true"
			mapa["columns[0][search][value]"] = ""

			mapa["columns[1][data]"] = "0"
			mapa["columns[1][searchable]"] = "true"
			mapa["columns[1][orderable]"] = "true"
			mapa["columns[1][search][value]"] = ""

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
				{Db: "instrument", Dt: 1, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(2)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "JuAn"
			row["1"] = "Trompeta"
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = "Juan"
			row["1"] = "Tambor"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})

		It("Order non string fields (for sqlserver)", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "desc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][orderable]"] = "true"
			mapa["columns[0][search][value]"] = ""

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "toys", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(6)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = int64(3)
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = int64(2)
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = int64(2)
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = int64(1)
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = int64(1)
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = int64(0)
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
	})
	It("Order dates (for sqlserver)", func() {

		mapa := make(map[string]string)
		mapa["draw"] = "64"
		mapa["start"] = "0"
		mapa["length"] = "10"
		mapa["order[0][column]"] = "0"
		mapa["order[0][dir]"] = "desc"

		mapa["columns[0][data]"] = "0"
		mapa["columns[0][searchable]"] = "true"
		mapa["columns[0][orderable]"] = "true"
		mapa["columns[0][search][value]"] = ""

		c := ControllerEmulated{Params: mapa}

		columns := []structs.Data{
			{Db: "birth_date", Dt: 0, Formatter: nil},
		}
		result, err := ssp.Simple(&c, db, "users", columns, nil)

		Expect(err).To(BeNil())
		Expect(result.Draw).To(Equal(64))
		Expect(result.RecordsTotal).To(Equal(int64(6)))
		Expect(result.RecordsFiltered).To(Equal(int64(6)))

		/*
			TODO Test date order
			date, _ := time.Parse(layoutISO, "2011-11-11")
			date2, _ := time.Parse(layoutISO, "2011-12-11")
			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = date2
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = date
			testData = append(testData, row)
			testData = append(testData, row)
			testData = append(testData, row)
			testData = append(testData, row)
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		*/
	})
	Describe("Field with space", func() {
		It("return favorite song ", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][search][value]"] = ""

			mapa["columns[1][data]"] = "1"
			mapa["columns[1][searchable]"] = "true"
			mapa["columns[1][search][value]"] = "Español"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
				{Db: "\"Favorite song\"", Dt: 1, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(1)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "JuAn"
			row["1"] = "Himno Español"
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
	})
}

// Errors test some errors
func Errors(db *gorm.DB) {
	Describe("Column not found", func() {
		It("Return error", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "2"
			mapa["order[0][column]"] = "1"
			mapa["order[0][dir]"] = "desc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "bike", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = nil
			testData = append(testData, row)
			row = make(map[string]interface{})
			row["0"] = nil
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
	})
	Describe("Dt is nil", func() {
		It("Return error", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "2"
			mapa["order[0][column]"] = "1"
			mapa["order[0][dir]"] = "desc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "bike", Dt: nil, Formatter: nil},
			}
			_, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(fmt.Sprintf("%v", err)).To(Equal("Dt cannot be nil in column[0]"))
		})
	})
	Describe("Format error", func() {
		It("return name whit prefix and age", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "62"
			mapa["start"] = "0"
			mapa["length"] = "4"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: func(
					data interface{}, row map[string]interface{}) (interface{}, error) {
					layout := "2006-01-02T15:04:05.000Z"
					//try convert name to date
					return time.Parse(layout, data.(string))
				}},
			}

			_, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).ToNot(BeNil())
		})
	})
	Describe("Column with reserved word", func() {
		It("returns 2 Age 15", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][search][value]"] = ""

			mapa["columns[1][data]"] = "1"
			mapa["columns[1][searchable]"] = "true"
			mapa["columns[1][search][value]"] = "2"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
				{Db: "end", Dt: 1, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(1)))

			testData := make([]interface{}, 0)
			row := make(map[string]interface{})
			row["0"] = "Joaquin"
			row["1"] = int64(2)
			testData = append(testData, row)

			Expect(result.Data).To(Equal(testData))
		})
	})
	Describe("Prevent SQL injection", func() {
		It("no return error", func() {

			mapa := make(map[string]string)
			mapa["draw"] = "64"
			mapa["start"] = "0"
			mapa["length"] = "10"
			mapa["order[0][column]"] = "0"
			mapa["order[0][dir]"] = "asc"

			mapa["columns[0][data]"] = "0"
			mapa["columns[0][searchable]"] = "true"
			mapa["columns[0][search][value]"] = "Juan`'"

			c := ControllerEmulated{Params: mapa}

			columns := []structs.Data{
				{Db: "name", Dt: 0, Formatter: nil},
			}
			result, err := ssp.Simple(&c, db, "users", columns, nil)

			Expect(err).To(BeNil())
			Expect(result.Draw).To(Equal(64))
			Expect(result.RecordsTotal).To(Equal(int64(6)))
			Expect(result.RecordsFiltered).To(Equal(int64(0)))

			testData := make([]interface{}, 0)

			Expect(result.Data).To(Equal(testData))
		})
	})
}
