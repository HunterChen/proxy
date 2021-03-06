package bak

// import (
// 	"bytes"
// 	// "fmt"
// 	// "strings"
// 	"reflect"
// 	"testing"
// 	// "time"
// )

// func SerializeTest(t *testing.T, doc interface{}, expectedBuffer []byte) {
// 	bson := NewBSON()
// 	// Validate the size of the bson array
// 	size, _ := bson.CalculateObjectSize(reflect.ValueOf(doc))
// 	if size != len(expectedBuffer) {
// 		t.Errorf("size comparison failed %v != %v", size, len(expectedBuffer))
// 	}

// 	// Serialize the document allowing self allocation of buffer
// 	b, err := bson.Serialize(doc, nil, 0)
// 	// Ensure the buffers match
// 	if err != nil || len(b) != len(expectedBuffer) || bytes.Compare(b, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned \nexp: %v:%v\ngot: %v:%v", expectedBuffer, len(expectedBuffer), b, len(b))
// 	}

// 	// Serialize into pre-allocated buffer
// 	b = make([]byte, len(expectedBuffer))
// 	// Serialize the document
// 	b, err = bson.Serialize(doc, b, 0)
// 	// Ensure the buffers match
// 	if err != nil || len(b) != len(expectedBuffer) || bytes.Compare(b, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned \nexp: %v:%v\ngot: %v:%v", expectedBuffer, len(expectedBuffer), bson, len(b))
// 	}
// }

// func DeserializeTest(t *testing.T, b []byte, empty interface{}, expected interface{}) {
// 	bson := NewBSON()
// 	// Deserialize the data into the type
// 	err := bson.Deserialize(b, empty)
// 	if err != nil {
// 		t.Errorf("[%v] Failed to deserialize %v into type %v", err, b, expected)
// 	}

// 	// Check if this is a document
// 	switch doc := empty.(type) {
// 	case *Document:
// 		switch doc1 := expected.(type) {
// 		case Document:
// 			if doc1.Equal(doc) == false {
// 				t.Errorf("failed to deserialize document correctly 3")
// 			}
// 		case *Document:
// 			if doc1.Equal(doc) == false {
// 				t.Errorf("failed to deserialize document correctly 4")
// 			}
// 		}
// 	default:
// 		if reflect.DeepEqual(empty, expected) == false {
// 			t.Errorf("failed to deserialize struct correctly 5")
// 		}
// 	}
// }

func TestEmptyDocumentSerialization(t *testing.T) {
	type T struct{}
	// Expected buffer from serialization
	var expectedBuffer = []byte{5, 0, 0, 0, 0}
	// Serialize tests
	SerializeTest(t, NewDocument(), expectedBuffer)
	SerializeTest(t, &T{}, expectedBuffer)
}

func TestDocumentWithInt32Serialization(t *testing.T) {
	type T struct {
		Int int32 `bson:"int,omitempty"`
	}
	// Expected buffer from serialization
	var expectedBuffer = []byte{14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	// Create document
	document := NewDocument()
	document.Add("int", int32(10))

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)
	SerializeTest(t, &T{10}, expectedBuffer)
}

func TestSimpleStringSerialization(t *testing.T) {
	type T struct {
		String string `bson:"string,omitempty"`
	}
	// Expected buffer from serialization
	var expectedBuffer = []byte{29, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 0}
	document := NewDocument()
	document.Add("string", "hello world")

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)
	SerializeTest(t, &T{"hello world"}, expectedBuffer)
}

func TestSimpleStringAndIntSerialization(t *testing.T) {
	type T struct {
		String string `bson:"string,omitempty"`
		Int    int32  `bson:"int,omitempty"`
	}
	// Expected buffer from serialization
	var expectedBuffer = []byte{38, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0}
	document := NewDocument()
	document.Add("string", "hello world")
	document.Add("int", int32(10))

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)
	SerializeTest(t, &T{"hello world", 10}, expectedBuffer)
}

func TestSimpleNestedDocumentSerialization(t *testing.T) {
	// Expected buffer from serialization
	var expectedBuffer = []byte{48, 0, 0, 0, 2, 115, 116, 114, 105, 110, 103, 0, 12, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0, 3, 100, 111, 99, 0, 14, 0, 0, 0, 16, 105, 110, 116, 0, 10, 0, 0, 0, 0, 0}
	document := NewDocument()
	subdocument := NewDocument()
	subdocument.Add("int", int32(10))
	document.Add("string", "hello world")
	document.Add("doc", subdocument)

	type T1 struct {
		Int int32 `bson:"int,omitempty"`
	}

	type T2 struct {
		String string `bson:"string,omitempty"`
		Doc    *T1    `bson:"doc,omitempty"`
	}

	// Serialize tests
	SerializeTest(t, document, expectedBuffer)
	SerializeTest(t, &T2{"hello world", &T1{10}}, expectedBuffer)

	// De serializing tests
	DeserializeTest(t, expectedBuffer, NewDocument(), document)
	DeserializeTest(t, expectedBuffer, &T2{}, &T2{"hello world", &T1{10}})
}

func TestSimpleArraySerialization(t *testing.T) {
	var expectedBuffer = []byte{35, 0, 0, 0, 4, 97, 114, 114, 97, 121, 0, 23, 0, 0, 0, 2, 48, 0, 2, 0, 0, 0, 97, 0, 2, 49, 0, 2, 0, 0, 0, 98, 0, 0, 0}
	document := NewDocument()
	array := make([]interface{}, 0)
	array = append(array, "a")
	array = append(array, "b")
	document.Add("array", array)
	bson, err := Serialize(document, nil, 0)

	t.Logf("[%v]", len(bson))
	t.Logf("[%v]", bson)
	t.Logf("[%v]", expectedBuffer)

	if err != nil {
		t.Fatalf("Failed to create bson document %v", err)
	}

	if len(bson) != len(expectedBuffer) {
		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
	}

	if bytes.Compare(bson, expectedBuffer) != 0 {
		t.Fatalf("Illegal BSON returned")
	}

	// Deserialize the object
	obj, err := Deserialize(expectedBuffer)
	if err != nil {
		t.Fatalf("Failed to deserialize the bson array")
	}

	validateObjectSize(t, obj, 1)
	a, _ := obj.FieldAsArray("array")
	validateString(t, a[0], "a")
	validateString(t, a[1], "b")
}

// func TestSimpleBinarySerialization(t *testing.T) {
// 	var expectedBuffer = []byte{26, 0, 0, 0, 5, 98, 105, 110, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0}
// 	document := NewDocument()
// 	bin := &Binary{0, []byte("hello world")}
// 	document.Add("bin", bin)
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateBinaryField(t, obj, "bin", bin)
// }

// func TestMixedDocumentSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{51, 0, 0, 0, 4, 97, 114, 114, 97, 121, 0, 39, 0, 0, 0, 5, 48, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 3, 49, 0, 12, 0, 0, 0, 16, 97, 0, 1, 0, 0, 0, 0, 0, 0}
// 	subdocument := NewDocument()
// 	subdocument.Add("a", int32(1))
// 	document := NewDocument()
// 	array := make([]interface{}, 0)
// 	bin := &Binary{0, []byte("hello world")}
// 	array = append(array, bin)
// 	array = append(array, subdocument)
// 	document.Add("array", array)
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	a, _ := obj.FieldAsArray("array")
// 	validateBinary(t, a[0], bin)
// 	validateIntField(t, toDocument(t, a[1]), "a", 1)
// }

// func TestObjectIdSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{21, 0, 0, 0, 7, 105, 100, 0, 49, 50, 51, 52, 53, 54, 55, 56, 49, 50, 51, 52, 0}
// 	document := NewDocument()
// 	objectid := &ObjectId{[]byte("123456781234")}
// 	document.Add("id", objectid)
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateObjectIdField(t, obj, "id", objectid)
// }

// func TestJavascriptNoScopeSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{34, 0, 0, 0, 13, 106, 115, 0, 21, 0, 0, 0, 118, 97, 114, 32, 97, 32, 61, 32, 102, 117, 110, 99, 116, 105, 111, 110, 40, 41, 123, 125, 0, 0}
// 	document := NewDocument()
// 	js := &Javascript{"var a = function(){}"}
// 	document.Add("js", js)

// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateJavascriptField(t, obj, "js", js)
// }

// func TestJavascriptWithScopeSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{50, 0, 0, 0, 15, 106, 115, 0, 41, 0, 0, 0, 21, 0, 0, 0, 118, 97, 114, 32, 97, 32, 61, 32, 102, 117, 110, 99, 116, 105, 111, 110, 40, 41, 123, 125, 0, 12, 0, 0, 0, 16, 97, 0, 1, 0, 0, 0, 0, 0}
// 	scope := NewDocument()
// 	scope.Add("a", int32(1))
// 	document := NewDocument()
// 	js := &JavascriptWScope{"var a = function(){}", scope}
// 	document.Add("js", js)

// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	j := validateJavascriptFieldWScope(t, obj, "js", js)
// 	validateIntField(t, j, "a", int32(1))
// }

// func TestMinMaxSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{15, 0, 0, 0, 255, 109, 105, 110, 0, 127, 109, 97, 120, 0, 0}

// 	// serializeAndPrint('min and max', {min: new MinKey(), max: new MaxKey()});
// 	document := NewDocument()
// 	document.Add("min", &Min{})
// 	document.Add("max", &Max{})

// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 2)
// 	validateMaxField(t, obj, "max")
// 	validateMinField(t, obj, "min")
// }

// func TestDateAndTimeSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{31, 0, 0, 0, 9, 111, 110, 101, 0, 160, 134, 1, 0, 0, 0, 0, 0, 9, 116, 119, 111, 0, 160, 134, 1, 0, 0, 0, 0, 0, 0}

// 	// Actual document
// 	document := NewDocument()
// 	document.Add("one", &Date{100000})
// 	document.Add("two", time.Unix(100000, 0))
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 2)
// 	validateTimeField(t, obj, "one", time.Unix(100000, 0))
// 	validateTimeField(t, obj, "two", time.Unix(100000, 0))
// }

// func TestBufferSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{24, 0, 0, 0, 5, 98, 0, 11, 0, 0, 0, 0, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 0}

// 	// Actual document
// 	document := NewDocument()
// 	buffer := []byte("hello world")
// 	document.Add("b", buffer)
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateBufferField(t, obj, "b", buffer)
// }

// func TestTimestampSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{16, 0, 0, 0, 17, 116, 0, 160, 134, 1, 0, 0, 0, 0, 0, 0}

// 	// Actual document
// 	document := NewDocument()
// 	document.Add("t", &Timestamp{100000})
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateTimestampField(t, obj, "t", &Timestamp{100000})
// }

// func TestInt64AndUInt64Serialization(t *testing.T) {
// 	var expectedBuffer = []byte{27, 0, 0, 0, 18, 111, 0, 255, 255, 255, 255, 255, 255, 255, 255, 18, 116, 0, 160, 134, 1, 0, 0, 0, 0, 0, 0}

// 	// Actual document
// 	document := NewDocument()
// 	document.Add("o", int64(-1))
// 	document.Add("t", uint64(100000))
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 2)
// 	validateInt64Field(t, obj, "o", int64(-1))
// 	validateUInt64Field(t, obj, "t", uint64(100000))
// }

// func TestFloat64Serialization(t *testing.T) {
// 	var expectedBuffer = []byte{16, 0, 0, 0, 1, 111, 0, 31, 133, 235, 81, 184, 30, 9, 64, 0}

// 	// Actual document
// 	document := NewDocument()
// 	document.Add("o", float64(3.14))
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateFloat64Field(t, obj, "o", float64(3.14))
// }

// func TestFloat32Serialization(t *testing.T) {
// 	var expectedBuffer = []byte{16, 0, 0, 0, 1, 111, 0, 102, 102, 102, 102, 102, 102, 246, 191, 0}

// 	// Actual document
// 	document := NewDocument()
// 	document.Add("o", float32(-1.4))
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateFloat32Field(t, obj, "o", float32(-1.4))
// }

// func TestBooleanSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{13, 0, 0, 0, 8, 111, 0, 1, 8, 116, 0, 0, 0}

// 	// Actual document
// 	document := NewDocument()
// 	document.Add("o", true)
// 	document.Add("t", false)
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 2)
// 	validateBooleanField(t, obj, "o", true)
// 	validateBooleanField(t, obj, "t", false)
// }

// func TestNilSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{8, 0, 0, 0, 10, 111, 0, 0}

// 	// Actual document
// 	document := NewDocument()
// 	document.Add("o", nil)
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateNilField(t, obj, "o")
// }

// func TestRegExpSerialization(t *testing.T) {
// 	var expectedBuffer = []byte{17, 0, 0, 0, 11, 111, 0, 91, 116, 101, 115, 116, 93, 0, 105, 0, 0}

// 	// Actual document
// 	document := NewDocument()
// 	regexp := &RegExp{"[test]", "i"}
// 	document.Add("o", regexp)
// 	bson, err := Serialize(document, nil, 0)

// 	t.Logf("[%v]", len(bson))
// 	t.Logf("[%v]", bson)
// 	t.Logf("[%v]", expectedBuffer)

// 	if err != nil {
// 		t.Fatalf("Failed to create bson document %v", err)
// 	}

// 	if len(bson) != len(expectedBuffer) {
// 		t.Fatalf("Illegal BSON length returned %v = %v", len(bson), len(expectedBuffer))
// 	}

// 	if bytes.Compare(bson, expectedBuffer) != 0 {
// 		t.Fatalf("Illegal BSON returned")
// 	}

// 	// Deserialize the object
// 	obj, err := Deserialize(expectedBuffer)
// 	if err != nil {
// 		t.Fatalf("Failed to deserialize the bson array")
// 	}

// 	validateObjectSize(t, obj, 1)
// 	validateRegExpField(t, obj, "o", regexp)
// }
