package model

type ValueFactory func() MapValueModel

type baseObjectMap struct {
	baseMap
	valueFactory ValueFactory
}
