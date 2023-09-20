package ssp

import (
	dialects "github.com/juaismar/go-gormssp/dialects"
	engine "github.com/juaismar/go-gormssp/engine"
)

// Functions exported.
var Simple = engine.Simple
var Complex = engine.Complex
var DataSimple = engine.DataSimple
var DataComplex = engine.DataComplex

// Models exported.
type JoinData = engine.JoinData
type MessageDataTable = engine.MessageDataTable
type Data = dialects.Data
