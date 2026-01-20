package csv

type colSchema struct {
	cName   string
	cType   inferredType
	cParser func(b []byte) parsedRes
}

type CSVSchema struct {
	cols []*colSchema
}
