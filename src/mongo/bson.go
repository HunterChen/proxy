package mongo

import (
	// "fmt"
	"reflect"
	"strings"
)

type bsonType byte

const (
	bsonDouble           bsonType = 0x01
	bsonString           bsonType = 0x02
	bsonDocument         bsonType = 0x03
	bsonArray            bsonType = 0x04
	bsonBinary           bsonType = 0x05
	bsonUndefined        bsonType = 0x06
	bsonObjectId         bsonType = 0x07
	bsonBoolean          bsonType = 0x08
	bsonDateTime         bsonType = 0x09
	bsonNull             bsonType = 0x0a
	bsonRegExp           bsonType = 0x0b
	bsonDBPointer        bsonType = 0x0c
	bsonJavaScript       bsonType = 0x0d
	bsonJavaScriptWScope bsonType = 0x0f
	bsonInt32            bsonType = 0x10
	bsonTimestamp        bsonType = 0x11
	bsonInt64            bsonType = 0x12
	bsonMinKey           bsonType = 0xff
	bsonMaxKey           bsonType = 0x7f
)

type ObjectId struct {
	Id []byte
}

type Binary struct {
	SubType byte
	Data    []byte
}

type Javascript struct {
	Code string
}

type JavascriptWScope struct {
	Code  string
	Scope *Document
}

type Date struct {
	Value int64
}

type RegExp struct {
	Pattern string
	Options string
}

type Timestamp struct {
	Value int64
}

type Min struct {
}

type Max struct {
}

type DBPointer struct {
}

type TypeInfos struct {
	Types map[string]*TypeInfo
}

type TypeInfo struct {
	Fields        map[string]*FieldInfo
	FieldsByIndex []*FieldInfo
	NumberOfField int
	HasGetBSON    bool
}

type FieldInfo struct {
	Name         string
	MetaDataName string
}

func writeU32(buffer []byte, index int, value uint32) {
	buffer[index+3] = byte((value >> 24) & 0xff)
	buffer[index+2] = byte((value >> 16) & 0xff)
	buffer[index+1] = byte((value >> 8) & 0xff)
	buffer[index] = byte(value & 0xff)
}

func writeU64(buffer []byte, index int, value uint64) {
	buffer[index+7] = byte((value >> 56) & 0xff)
	buffer[index+6] = byte((value >> 48) & 0xff)
	buffer[index+5] = byte((value >> 40) & 0xff)
	buffer[index+4] = byte((value >> 32) & 0xff)
	buffer[index+3] = byte((value >> 24) & 0xff)
	buffer[index+2] = byte((value >> 16) & 0xff)
	buffer[index+1] = byte((value >> 8) & 0xff)
	buffer[index] = byte(value & 0xff)
}

func readUInt64(buffer []byte, index int) uint64 {
	return (uint64(buffer[index]) << 0) |
		(uint64(buffer[index+1]) << 8) |
		(uint64(buffer[index+2]) << 16) |
		(uint64(buffer[index+3]) << 24) |
		(uint64(buffer[index+4]) << 32) |
		(uint64(buffer[index+5]) << 40) |
		(uint64(buffer[index+6]) << 48) |
		(uint64(buffer[index+7]) << 56)
}

func readUInt32(buffer []byte, index int) uint32 {
	return (uint32(buffer[index]) << 0) |
		(uint32(buffer[index+1]) << 8) |
		(uint32(buffer[index+2]) << 16) |
		(uint32(buffer[index+3]) << 24)
}

type Getter interface {
	GetBSON() (interface{}, error)
}

type BSON struct {
	typeInfos *TypeInfos
}

func NewBSON() *BSON {
	return &BSON{&TypeInfos{make(map[string]*TypeInfo)}}
}

func parseTypeInformation(typeInfos *TypeInfos, originalValue reflect.Value, value reflect.Value) *TypeInfo {
	// We have a pointer get the underlying value
	if value.Type().Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Check if we have a cached type
	cachedType := typeInfos.Types[value.Type().Name()]
	if cachedType != nil {
		return cachedType
	}

	// Get the number of fields
	numberOfFields := value.NumField()

	// Create typeInfo box
	typeInfo := TypeInfo{}
	// Pre-allocate a map with the entries we need
	typeInfo.Fields = make(map[string]*FieldInfo, numberOfFields*2)
	typeInfo.FieldsByIndex = make([]*FieldInfo, numberOfFields)
	typeInfo.NumberOfField = numberOfFields

	// Iterate over all the fields and collect the metadata
	for index := 0; index < numberOfFields; index++ {
		// Get the field information
		fieldType := value.Type().Field(index)
		// Get the field name
		key := fieldType.Name
		// Get the tag for the field
		tag := fieldType.Tag.Get("bson")

		// Split the tag into parts
		parts := strings.Split(tag, ",")

		// Override the key if the metadata has one
		if len(parts) > 0 && parts[0] != "" {
			key = parts[0]
		}

		// Create a new fieldInfo instance
		fieldInfo := FieldInfo{fieldType.Name, key}
		// Add to the map
		typeInfo.Fields[fieldType.Name] = &fieldInfo
		typeInfo.Fields[key] = &fieldInfo
		typeInfo.FieldsByIndex[index] = &fieldInfo
	}

	if originalValue.Type().Kind() == reflect.Ptr {
		// Iterate over all the
		numberOfMethods := originalValue.NumMethod()

		// Iterate over all the fields and collect the metadata
		for index := 0; index < numberOfMethods; index++ {
			// Method type
			methodType := originalValue.Type().Method(index)
			if methodType.Name == "GetBSON" {
				typeInfo.HasGetBSON = true
				break
			}
		}
	}

	// We need to save the type information of the GetBSON method aswell
	if typeInfo.HasGetBSON {
		if vi, ok := originalValue.Interface().(Getter); ok {
			getv, err := vi.GetBSON()
			if err != nil {
				panic(err)
			}

			// Add the type information to our cache
			parseTypeInformation(typeInfos, reflect.ValueOf(getv), reflect.ValueOf(getv))
		}
	}

	// Save type
	typeInfos.Types[value.Type().Name()] = &typeInfo
	// Return the type information
	return &typeInfo
}
