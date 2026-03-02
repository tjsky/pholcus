package data

import (
	"sync"
)

type (
	// DataCell is a storage unit for text data.
	DataCell map[string]interface{}
	// FileCell is a storage unit for file data.
	// Stored path format: file/"Dir"/"RuleName"/"time"/"Name"
	FileCell map[string]interface{}
)

var (
	dataCellPool = &sync.Pool{
		New: func() interface{} {
			return DataCell{}
		},
	}
	fileCellPool = &sync.Pool{
		New: func() interface{} {
			return FileCell{}
		},
	}
)

// GetDataCell returns a DataCell from the pool with the given fields.
func GetDataCell(ruleName string, data map[string]interface{}, url string, parentUrl string, downloadTime string) DataCell {
	cell := dataCellPool.Get().(DataCell)
	cell["RuleName"] = ruleName
	cell["Data"] = data
	cell["Url"] = url
	cell["ParentUrl"] = parentUrl
	cell["DownloadTime"] = downloadTime
	return cell
}

// GetFileCell returns a FileCell from the pool with the given fields.
func GetFileCell(ruleName, name string, bytes []byte) FileCell {
	cell := fileCellPool.Get().(FileCell)
	cell["RuleName"] = ruleName
	cell["Name"] = name
	cell["Bytes"] = bytes
	return cell
}

// PutDataCell returns a DataCell to the pool.
func PutDataCell(cell DataCell) {
	cell["RuleName"] = nil
	cell["Data"] = nil
	cell["Url"] = nil
	cell["ParentUrl"] = nil
	cell["DownloadTime"] = nil
	dataCellPool.Put(cell)
}

// PutFileCell returns a FileCell to the pool.
func PutFileCell(cell FileCell) {
	cell["RuleName"] = nil
	cell["Name"] = nil
	cell["Bytes"] = nil
	fileCellPool.Put(cell)
}
