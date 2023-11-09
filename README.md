# gormssp

_Using Datatables pagination with golang_

[![Build Status](https://travis-ci.org/juaismar/go-gormssp.svg?branch=master)](https://travis-ci.org/juaismar/go-gormssp)
[![Go Report Card](https://goreportcard.com/badge/github.com/juaismar/go-gormssp)](https://goreportcard.com/report/github.com/juaismar/go-gormssp)
[![codecov](https://codecov.io/gh/juaismar/go-gormssp/branch/master/graph/badge.svg)](https://codecov.io/gh/juaismar/go-gormssp)
[![MIT licensed](https://img.shields.io/github/license/juaismar/go-gormssp)](https://raw.githubusercontent.com/juaismar/go-gormssp/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-go_gormssp-blue.svg)](https://godoc.org/github.com/juaismar/go-gormssp)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/juaismar/go-gormssp)](https://pkg.go.dev/github.com/juaismar/go-gormssp)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/juaismar/go-gormssp)

### Pre-requisites 📋

This package use the new version of gorm in https://github.com/go-gorm/gorm, for previous version, go to https://github.com/juaismar/gormssp

Database compatible: Postgres (stable), SQLite (without REGEXP)

* Obviously use it in a golang project
* Gorm package (https://gorm.io/) (https://github.com/go-gorm/gorm)
* Beego package (https://beego.me/) (https://github.com/astaxie/beego)

### Installation 🔧

_Install with the next command:_

```
go get github.com/juaismar/go-gormssp
```

_and import the package with:_

```
import ("github.com/juaismar/go-gormssp")
```
## Working example 🚀

A working example on https://github.com/juaismar/go-GormSSP_Example

-This is a simple code that sends data to the Datatables JS client.
```
import ("github.com/juaismar/go-gormssp")

func (c *User) Pagination() {

  // Array of database columns which should be read and sent back to DataTables.
  // The `db` parameter represents the column name in the database, while the `dt`
  // parameter represents the DataTables column identifier. In this case simple
  // indexes but can be a string
  // Formatter is a function to customize the value of field , can be nil.
  columns := []SSP.Data{
    SSP.Data{Db: "name", Dt: 0, Formatter: nil},
    SSP.Data{Db: "role", Dt: 1, Formatter: nil},
    SSP.Data{Db: "email", Dt: 2, Formatter: nil},
  }

  // Send the data to the client
  // "users" is the name of the table
  c.Data["json"], _ = SSP.Simple(c, model.ORM, "users", columns)
  c.ServeJSON()
}
```

-This is an example of data formatting.
```
SSP.Data{Db: "registered", Dt: 3, Formatter: func(
  data interface{}, row map[string]interface{}) (interface{}, error) {
  //data is the value id column, row is a map whit the values of all columns
  if data != nil {
    return data.(time.Time).Format("2006-01-02 15:04:05"), nil
  }
  return "", nil
}}
```

-This is a complex example.
```
import ("github.com/juaismar/go-gormssp")

func (c *User) Pagination() {
    columns := []SSP.Data{
      SSP.Data{Db: "id", Dt: "id", Formatter: nil},
    }
	
    //whereResult is a WHERE condition to apply to the result set
    //whereAll is a WHERE condition to apply to all queries
    var whereResult []string
    var whereAll []string
    var whereJoin = make([]SSP.JoinData, 0)
    whereAll = append(whereAll, "deleted_at IS NULL")

    c.Data["json"], _ = SSP.Complex(c, model.ORM, "events", columns, whereResult, whereAll, whereJoin)
    c.ServeJSON()
}
```

-Fields with spaces in the name or start whit numbers must be escaped.
```
SSP.Data{Db: "\"1My special table\"", Dt: 3, Formatter: nil}
```

-Can implement functions in search, example:accent insensitive search 
```
SELECT * FROM users WHERE f_unaccent(name) LIKE '%maria%';.
```
can be implemented as:
```
SSP.Data{Db: "name", Dt: 2, Formatter: nil, Sf: "f_unaccent(name)"}
```

-This project is based in the PHP version of datatables pagination in https://datatables.net/examples/data_sources/server_side
-Original file can be found in https://github.com/DataTables/DataTables/blob/master/examples/server_side/scripts/ssp.class.php

## Important
If you are using Bigquery must set env var LW_DATASET to your dataset

## Author ✒️

* **Juan Iscar** - (https://github.com/juaismar)

## Thanks 🎁
* All my friends at work.
* Sergio(https://github.com/serveba) and Mario (https://github.com/mapreal19) who taught me how to program golang and showed me the wonderful world of good practices.
* Juan, Juan and Joaquin.


_Readme.md based in https://gist.github.com/Villanuevand/6386899f70346d4580c723232524d35a_
