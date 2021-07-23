package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ThomasObenaus/go-conf/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testParseConfigTag struct {
	configTagStr string
	typeOfEntry  reflect.Type
	nameOfParent string
}

func Test_parseConfigTagDefinition_Fail(t *testing.T) {
	// GIVEN
	invalidType1 := "{'name':'string-slice','default':['default1','default2']}"
	invalidType2 := "{'name':'string-slice','default':['default1','default2']}"
	invalid1 := "{}"
	invalid2 := "just invalid [][}{"

	// WHEN
	_, errInvalidType1 := parseConfigTagDefinition(invalidType1, reflect.TypeOf(int(0)), "")
	_, errInvalidType2 := parseConfigTagDefinition(invalidType2, reflect.TypeOf([]int{}), "")
	_, errInvalid1 := parseConfigTagDefinition(invalid1, reflect.TypeOf([]int{}), "")
	_, errInvalid2 := parseConfigTagDefinition(invalid2, reflect.TypeOf([]int{}), "")

	// THEN
	assert.Error(t, errInvalidType1)
	assert.Error(t, errInvalidType2)
	assert.Error(t, errInvalid1)
	assert.Error(t, errInvalid2)
}

func Test_parseConfigTagDefinition_Struct(t *testing.T) {
	type mystruct struct {
		Field1 string `cfg:"{'name':'f1','default':'default'}"`
		Field2 int    `cfg:"{'name':'f2','default':111}"`
	}

	// GIVEN
	simpleStruct := testParseConfigTag{
		configTagStr: "{'name':'string-slice','default':{'f1':'value1'}}",
		typeOfEntry:  reflect.TypeOf(mystruct{}),
	}

	// WHEN + THEN
	tag, err := parseConfigTagDefinition(simpleStruct.configTagStr, simpleStruct.typeOfEntry, simpleStruct.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, mystruct{Field1: "value1", Field2: 111}, tag.Def)
	assert.Equal(t, simpleStruct.typeOfEntry, reflect.TypeOf(tag.Def))
	assert.Equal(t, simpleStruct.typeOfEntry, tag.desiredType)
	assert.False(t, tag.IsRequired())
}

func Test_parseConfigTagDefinition_Slices(t *testing.T) {
	// GIVEN
	stringSlice := testParseConfigTag{
		configTagStr: "{'name':'string-slice','default':['default1','default2']}",
		typeOfEntry:  reflect.TypeOf([]string{}),
	}

	// WHEN + THEN
	tag, err := parseConfigTagDefinition(stringSlice.configTagStr, stringSlice.typeOfEntry, stringSlice.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, []string{"default1", "default2"}, tag.Def)
	assert.Equal(t, stringSlice.typeOfEntry, reflect.TypeOf(tag.Def))
	assert.Equal(t, stringSlice.typeOfEntry, tag.desiredType)
	assert.False(t, tag.IsRequired())

	type mystruct struct {
		Field1 string `cfg:"{'name':'f1','default':'default'}"`
	}

	// GIVEN
	structSlice := testParseConfigTag{
		configTagStr: "{'name':'struct-slice','default':[{'f1':'value1'},{}]}",
		typeOfEntry:  reflect.TypeOf([]mystruct{}),
	}

	// WHEN + THEN
	tag, err = parseConfigTagDefinition(structSlice.configTagStr, structSlice.typeOfEntry, structSlice.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, []mystruct{{Field1: "value1"}, {Field1: "default"}}, tag.Def)
	assert.Equal(t, structSlice.typeOfEntry, reflect.TypeOf(tag.Def))
	assert.Equal(t, structSlice.typeOfEntry, tag.desiredType)
	assert.False(t, tag.IsRequired())
}

func Test_parseConfigTagDefinition_Simple(t *testing.T) {
	// GIVEN
	simpleString := testParseConfigTag{
		configTagStr: "{'name':'field-string','desc':'string field','default':'default'}",
		typeOfEntry:  reflect.TypeOf(""),
	}

	// WHEN + THEN
	tagStr, err := parseConfigTagDefinition(simpleString.configTagStr, simpleString.typeOfEntry, simpleString.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-string", tagStr.Name)
	assert.Equal(t, "string field", tagStr.Description)
	assert.Equal(t, "default", tagStr.Def)
	assert.Equal(t, simpleString.typeOfEntry, reflect.TypeOf(tagStr.Def))
	assert.Equal(t, simpleString.typeOfEntry, tagStr.desiredType)
	assert.False(t, tagStr.IsRequired())

	// GIVEN
	simpleInt := testParseConfigTag{
		configTagStr: "{'name':'field-int','desc':'int field','default':1111}",
		typeOfEntry:  reflect.TypeOf(int(0)),
	}

	// WHEN + THEN
	tagInt, err := parseConfigTagDefinition(simpleInt.configTagStr, simpleInt.typeOfEntry, simpleInt.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-int", tagInt.Name)
	assert.Equal(t, "int field", tagInt.Description)
	assert.Equal(t, 1111, tagInt.Def)
	assert.Equal(t, simpleInt.typeOfEntry, reflect.TypeOf(tagInt.Def))
	assert.Equal(t, simpleInt.typeOfEntry, tagInt.desiredType)
	assert.False(t, tagInt.IsRequired())

	// GIVEN
	simpleFloat := testParseConfigTag{
		configTagStr: "{'name':'field-float','desc':'float field','default':22.22}",
		typeOfEntry:  reflect.TypeOf(float64(0)),
	}

	// WHEN + THEN
	tagFloat, err := parseConfigTagDefinition(simpleFloat.configTagStr, simpleFloat.typeOfEntry, simpleFloat.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-float", tagFloat.Name)
	assert.Equal(t, "float field", tagFloat.Description)
	assert.Equal(t, 22.22, tagFloat.Def)
	assert.Equal(t, simpleFloat.typeOfEntry, reflect.TypeOf(tagFloat.Def))
	assert.Equal(t, simpleFloat.typeOfEntry, tagFloat.desiredType)
	assert.False(t, tagFloat.IsRequired())

	// GIVEN
	simpleBool := testParseConfigTag{
		configTagStr: "{'name':'field-bool','desc':'bool field','default':true}",
		typeOfEntry:  reflect.TypeOf(bool(true)),
	}

	// WHEN + THEN
	tagBool, err := parseConfigTagDefinition(simpleBool.configTagStr, simpleBool.typeOfEntry, simpleBool.nameOfParent)
	assert.NoError(t, err)
	assert.Equal(t, "field-bool", tagBool.Name)
	assert.Equal(t, "bool field", tagBool.Description)
	assert.Equal(t, true, tagBool.Def)
	assert.Equal(t, simpleBool.typeOfEntry, reflect.TypeOf(tagBool.Def))
	assert.Equal(t, simpleBool.typeOfEntry, tagBool.desiredType)
	assert.False(t, tagBool.IsRequired())
}

func Test_parseConfigTagDefinition_Required(t *testing.T) {
	// GIVEN
	simpleOptional := testParseConfigTag{
		configTagStr: "{'name':'field-string','desc':'string field','default':'default'}",
		typeOfEntry:  reflect.TypeOf(""),
	}

	// WHEN + THEN
	tag, err := parseConfigTagDefinition(simpleOptional.configTagStr, simpleOptional.typeOfEntry, simpleOptional.nameOfParent)
	assert.NoError(t, err)
	assert.False(t, tag.IsRequired())

	// GIVEN
	simpleRequired := testParseConfigTag{
		configTagStr: "{'name':'field-string','desc':'string field'}",
		typeOfEntry:  reflect.TypeOf(""),
	}

	// WHEN + THEN
	tag, err = parseConfigTagDefinition(simpleRequired.configTagStr, simpleRequired.typeOfEntry, simpleRequired.nameOfParent)
	assert.NoError(t, err)
	assert.True(t, tag.IsRequired())
}

func Test_extractConfigTagsOfStruct_Primitives(t *testing.T) {

	// GIVEN
	type primitives struct {
		ShouldBeSkipped string
		SomeFieldStr    string  `cfg:"{'name':'field-str','desc':'a string field','default':'default value'}"`
		SomeFieldInt    int     `cfg:"{'name':'field-int','desc':'a int field','default':11}"`
		SomeFieldFloat  float64 `cfg:"{'name':'field-float','desc':'a float field','default':22.22}"`
		SomeFieldBool   bool    `cfg:"{'name':'field-bool','desc':'a bool field','default':true}"`
	}
	prims := primitives{}

	// WHEN
	entries, err := extractConfigTagsOfStruct(&prims, interfaces.NoLogging, "", configTag{})

	// THEN
	assert.NoError(t, err)
	require.Len(t, entries, 4)
	assert.Equal(t, "field-str", entries[0].Name)
	assert.Equal(t, "a string field", entries[0].Description)
	assert.Equal(t, "default value", entries[0].Def)
	assert.Equal(t, reflect.TypeOf(""), reflect.TypeOf(entries[0].Def))
	assert.False(t, entries[0].IsRequired())

	assert.Equal(t, "field-int", entries[1].Name)
	assert.Equal(t, "a int field", entries[1].Description)
	assert.Equal(t, reflect.TypeOf(int(0)), reflect.TypeOf(entries[1].Def))
	assert.Equal(t, 11, entries[1].Def)
	assert.False(t, entries[1].IsRequired())

	assert.Equal(t, "field-float", entries[2].Name)
	assert.Equal(t, "a float field", entries[2].Description)
	assert.Equal(t, reflect.TypeOf(float64(0)), reflect.TypeOf(entries[2].Def))
	assert.Equal(t, 22.22, entries[2].Def)
	assert.False(t, entries[2].IsRequired())

	assert.Equal(t, "field-bool", entries[3].Name)
	assert.Equal(t, "a bool field", entries[3].Description)
	assert.Equal(t, reflect.TypeOf(bool(true)), reflect.TypeOf(entries[3].Def))
	assert.Equal(t, true, entries[3].Def)
	assert.False(t, entries[3].IsRequired())
}

func Test_extractConfigTagsOfStruct_Required(t *testing.T) {

	// GIVEN
	type primitives struct {
		SomeFielOptional  string `cfg:"{'name':'field-str','desc':'a string field','default':'default value'}"`
		SomeFieldRequired string `cfg:"{'name':'field-str','desc':'a string field'}"`
	}
	prims := primitives{}

	// WHEN
	entries, err := extractConfigTagsOfStruct(&prims, interfaces.NoLogging, "", configTag{})

	// THEN
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.False(t, entries[0].IsRequired())
	assert.True(t, entries[1].IsRequired())
}

func Test_isOfPrimitiveType(t *testing.T) {
	type my struct {
	}

	is1, err1 := isOfPrimitiveType(reflect.TypeOf(int(0)))
	is2, err2 := isOfPrimitiveType(reflect.TypeOf(my{}))
	is3, err3 := isOfPrimitiveType(reflect.TypeOf([]int{}))
	is4, err4 := isOfPrimitiveType(reflect.TypeOf([]my{}))
	i := 2
	is5, err5 := isOfPrimitiveType(reflect.TypeOf(&i))

	assert.NoError(t, err1)
	assert.True(t, is1)
	assert.NoError(t, err2)
	assert.False(t, is2)
	assert.NoError(t, err3)
	assert.True(t, is3)
	assert.NoError(t, err4)
	assert.True(t, is4)
	assert.NoError(t, err5)
	assert.True(t, is5)
}

func Test_processAllConfigTagsOfStruct(t *testing.T) {
	// GIVEN
	type primitives struct {
		NoConfigTag       string
		SomeFielOptional  string `cfg:"{'name':'field-1','desc':'a string field','default':'default value'}"`
		SomeFieldRequired string `cfg:"{'name':'field-2','desc':'a string field'}"`
	}
	prims := primitives{}

	// WHEN
	obtainedConfigTags := make([]configTag, 0)
	err := processAllConfigTagsOfStruct(&prims, interfaces.NoLogging, "", configTag{}, func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error {
		obtainedConfigTags = append(obtainedConfigTags, cfgTag)
		return nil
	})

	// THEN
	require.NoError(t, err)
	require.Len(t, obtainedConfigTags, 2)
	assert.Equal(t, "field-1", obtainedConfigTags[0].Name)
	assert.Equal(t, "field-2", obtainedConfigTags[1].Name)
}

func Test_processAllConfigTagsOfStruct_Fail(t *testing.T) {
	// GIVEN
	type primitives struct {
		SomeFieldRequired string `cfg:"{'name':'field-2','desc':'a string field'}"`
	}
	prims := primitives{}

	// WHEN
	errNoPointer := processAllConfigTagsOfStruct(prims, interfaces.NoLogging, "", configTag{}, func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error {
		return nil
	})
	errNil := processAllConfigTagsOfStruct(nil, interfaces.NoLogging, "", configTag{}, func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error {
		return nil
	})
	errFailHandler := processAllConfigTagsOfStruct(&prims, interfaces.NoLogging, "", configTag{}, func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error {
		return fmt.Errorf("FAILURE")
	})

	// THEN
	require.Error(t, errNoPointer)
	require.Error(t, errNil)
	require.Error(t, errFailHandler)
}

func Test_extractConfigTagsOfStruct(t *testing.T) {
	// GIVEN
	type nestedAgain struct {
		FieldI  string `cfg:"{'name':'field-i','desc':'a string field','default':'field-i default'}"`
		FieldII string `cfg:"{'name':'field-ii','desc':'a nested field','default':'field-ii default'}"`
	}

	type nested struct {
		FieldA string      `cfg:"{'name':'field-a','desc':'a string field','default':'field-a default'}"`
		FieldB nestedAgain `cfg:"{'name':'field-b','desc':'a nested field','default':{'field-i':'field-i value','field-ii':'field-ii value'}}"`
	}

	type my struct {
		NoConfigTag string
		Field1      string   `cfg:"{'name':'field-1','desc':'a string field','default':'field-1 default'}"`
		Field3      nested   `cfg:"{'name':'field-3','desc':'a nested field'}"`
		Field4      []nested `cfg:"{'name':'field-4','desc':'a list of nested field','default':[{'field-a':'field-a value','field-b':{}}]}"`
	}
	strct := my{}

	// WHEN
	cfgTags, errNoPointer := extractConfigTagsOfStruct(&strct, interfaces.NoLogging, "", configTag{})

	// THEN
	require.NoError(t, errNoPointer)
	require.Len(t, cfgTags, 5)
}

func Test_extractConfigTagsOfStruct_NoFieldAnnotation(t *testing.T) {

	// GIVEN
	type fieldsNotAnnotated struct {
		Field1 string
	}
	type config struct {
		SomeFieldStr fieldsNotAnnotated `cfg:"{'name':'field-1','desc':'A field of a complex type whose fields are NOT annotated'}"`
	}

	cfg := config{}

	// WHEN
	entries, err := extractConfigTagsOfStruct(&cfg, interfaces.NoLogging, "", configTag{})

	// THEN
	assert.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "field-1", entries[0].Name)
	assert.True(t, entries[0].IsRequired())
}

func Test_hasAnnotatedFields(t *testing.T) {

	// GIVEN
	type fieldsNotAnnotated struct {
		Field1 string
		Field2 string
		Field3 string
	}

	type fieldsAnnotated struct {
		Field1 string
		Field2 string `cfg:"{'name':'field-2'}"`
		Field3 string
	}

	// WHEN
	has1 := hasAnnotatedFields(reflect.TypeOf(fieldsNotAnnotated{}))
	has2 := hasAnnotatedFields(reflect.TypeOf(fieldsAnnotated{}))
	has3 := hasAnnotatedFields(reflect.TypeOf(nil))
	has4 := hasAnnotatedFields(reflect.TypeOf(""))
	has5 := hasAnnotatedFields(reflect.TypeOf(1))
	i := 6
	has6 := hasAnnotatedFields(reflect.TypeOf(i))
	n := fieldsNotAnnotated{}
	has7 := hasAnnotatedFields(reflect.TypeOf(&n))
	m := fieldsAnnotated{}
	has8 := hasAnnotatedFields(reflect.TypeOf(&m))

	// THEN
	assert.False(t, has1)
	assert.True(t, has2)
	assert.False(t, has3)
	assert.False(t, has4)
	assert.False(t, has5)
	assert.False(t, has6)
	assert.False(t, has7)
	assert.True(t, has8)
}

func Test_extractConfigTagsOfStruct_NonPrimitiveWithoutAnnotation(t *testing.T) {
	// GIVEN
	type nested struct {
		FieldA string
	}

	type my struct {
		Field1 string `cfg:"{'name':'field-1','desc':'a string field','default':'field-1 default'}"`
		Field3 nested
	}
	strct := my{}

	// WHEN
	cfgTags, errNoPointer := extractConfigTagsOfStruct(&strct, interfaces.NoLogging, "", configTag{})

	// THEN
	require.NoError(t, errNoPointer)
	require.Len(t, cfgTags, 5)
}
