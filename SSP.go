package ssp

import (
	engine "github.com/juaismar/go-gormssp/engine"
	"github.com/juaismar/go-gormssp/structs"
)

// Functions exported.
var Simple = engine.Simple
var Complex = engine.Complex
var DataSimple = engine.DataSimple
var DataComplex = engine.DataComplex
var ReaderSimple = engine.ReaderSimple
var ReaderComplex = engine.ReaderComplex

// Models exported.
type JoinData = structs.JoinData
type MessageDataTable = structs.MessageDataTable
type Data = structs.Data
