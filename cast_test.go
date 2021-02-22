package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createAndFillStruct_Fail(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}

	// WHEN
	_, errMissingRequired := createAndFillStruct(reflect.TypeOf(my{}), map[string]interface{}{})

	// THEN
	assert.Error(t, errMissingRequired)

	// GIVEN
	type myInvalidConfig struct {
		Field1 string `cfg:"{}"`
	}

	// WHEN
	_, errInvalidConfig := createAndFillStruct(reflect.TypeOf(myInvalidConfig{}), map[string]interface{}{})

	// THEN
	assert.Error(t, errInvalidConfig)

	// GIVEN
	type myUnexportedField struct {
		// nolint
		field1 string `cfg:"{'name':'field_1','default':'value'}"`
	}

	// WHEN
	_, errUnexportedField := createAndFillStruct(reflect.TypeOf(myUnexportedField{}), map[string]interface{}{})

	// THEN
	assert.Error(t, errUnexportedField)

	// WHEN
	_, errNoStruct := createAndFillStruct(reflect.TypeOf(int(0)), map[string]interface{}{})

	// THEN
	assert.Error(t, errNoStruct)
}

func Test_createAndFillStruct(t *testing.T) {
	// GIVEN
	type nested struct {
		FieldA          float64 `cfg:"{'name':'field_a'}"`
		ShouldBeIgnored bool
	}

	type my struct {
		Field1          string `cfg:"{'name':'field_1'}"`
		Field2          nested `cfg:"{'name':'field_2'}"`
		ShouldBeIgnored bool
	}

	expected := my{Field1: "field-1", Field2: nested{FieldA: 22.22}}

	data := map[string]interface{}{
		"field_1": expected.Field1,
		"field_2": map[string]interface{}{
			"field_a": expected.Field2.FieldA,
		},
	}

	// WHEN
	structVal, err := createAndFillStruct(reflect.TypeOf(my{}), data)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, structVal.Interface())
}

func Test_castToSlice_Fail(t *testing.T) {
	// GIVEN
	vIntSlice := []int{11, 22}
	vIntSliceIf := []interface{}{vIntSlice[0], vIntSlice[1]}

	// WHEN
	casted1, err1 := castToSlice(vIntSliceIf, reflect.TypeOf(int(0)))

	// THEN
	assert.Error(t, err1)
	assert.Nil(t, casted1)

	// GIVEN
	vInt := int(10)

	// WHEN
	casted2, err2 := castToSlice(vInt, reflect.TypeOf([]int{}))

	// THEN
	assert.Error(t, err2)
	assert.Nil(t, casted2)
}

func Test_castToSlice(t *testing.T) {
	// GIVEN
	vIntSlice := []int{11, 22}
	vIntSliceIf := []interface{}{vIntSlice[0], vIntSlice[1]}

	// WHEN
	casted1, err1 := castToSlice(vIntSliceIf, reflect.TypeOf([]int{}))

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, vIntSlice, casted1)

	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	expected := []my{
		{Field1: "field-1"},
		{Field1: "field-2"},
	}

	vStructIf := []interface{}{
		map[string]interface{}{"field_1": expected[0].Field1},
		map[string]interface{}{"field_1": expected[1].Field1},
	}

	// WHEN
	casted2, err2 := castToSlice(vStructIf, reflect.TypeOf([]my{}))

	// THEN
	assert.NoError(t, err2)
	assert.Equal(t, expected, casted2)
}

func Test_castToStruct_Fail(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	vStructMissingRequired := map[string]interface{}{}
	vNoStruct := int(11)

	// WHEN
	casted1, err1 := castToStruct(vStructMissingRequired, reflect.TypeOf(my{}))
	casted2, err2 := castToStruct(vStructMissingRequired, reflect.TypeOf(int(0)))
	casted3, err3 := castToStruct(vNoStruct, reflect.TypeOf(my{}))

	// THEN
	assert.Error(t, err1)
	assert.Nil(t, casted1)
	assert.Error(t, err2)
	assert.Nil(t, casted2)
	assert.Error(t, err3)
	assert.Nil(t, casted3)
}

func Test_castToStruct(t *testing.T) {
	// GIVEN
	type nested struct {
		FieldA float64 `cfg:"{'name':'field_a'}"`
	}
	type my struct {
		Field1 string   `cfg:"{'name':'field_1'}"`
		Field2 int      `cfg:"{'name':'field_2'}"`
		Field3 nested   `cfg:"{'name':'field_3'}"`
		Field4 []int    `cfg:"{'name':'field_4'}"`
		Field5 []nested `cfg:"{'name':'field_5'}"`
	}

	expected := my{
		Field1: "a field",
		Field2: 11,
		Field3: nested{
			FieldA: 22.22,
		},
		Field4: []int{11, 22},
		Field5: []nested{{FieldA: 22.22}},
	}

	vStruct := map[string]interface{}{
		"field_1": expected.Field1,
		"field_2": expected.Field2,
		"field_3": map[string]interface{}{
			"field_a": expected.Field3.FieldA,
		},
		"field_4": []interface{}{
			expected.Field4[0],
			expected.Field4[1],
		},
		"field_5": []interface{}{
			map[string]interface{}{"field_a": expected.Field5[0].FieldA},
		},
	}

	// WHEN
	casted1, err1 := castToStruct(vStruct, reflect.TypeOf(my{}))

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, reflect.TypeOf(my{}), reflect.TypeOf(casted1))
	assert.Equal(t, expected, casted1)
}

func Test_castToPrimitive(t *testing.T) {
	// GIVEN
	vInt := int(11)
	vIntSlice := []int{11, 22}
	vString := "something"

	// WHEN
	castedInt, errInt := castToPrimitive(vInt, reflect.TypeOf(int(0)))
	castedIntSlice, errIntSlice := castToPrimitive(vIntSlice, reflect.TypeOf([]int{}))
	castedWrongType, errWrongType := castToPrimitive(vString, reflect.TypeOf(int(0)))

	// THEN
	assert.NoError(t, errInt)
	assert.Equal(t, 11, castedInt)
	assert.NoError(t, errIntSlice)
	assert.Equal(t, []int{11, 22}, castedIntSlice)
	assert.Error(t, errWrongType)
	assert.Nil(t, castedWrongType)
}

func Test_isFieldExported(t *testing.T) {
	// GIVEN
	type my struct {
		ExportedField string
		// nolint
		unExportedField string
	}
	reflectedType := reflect.TypeOf(my{})
	exportedField, _ := reflectedType.FieldByName("ExportedField")
	unExportedField, _ := reflectedType.FieldByName("unExportedField")

	// WHEN + THEN
	assert.True(t, isFieldExported(exportedField))
	assert.False(t, isFieldExported(unExportedField))
}

func Test_fullFieldName(t *testing.T) {
	assert.Equal(t, "root", fullFieldName("", "root"))
	assert.Equal(t, "root.child", fullFieldName("root", "child"))
	assert.Equal(t, "root.children.child", fullFieldName("root.children", "child"))
}

func Test_parseStringContainingSliceOfMaps(t *testing.T) {
	v1 := `[{"name":"name1","key":"key1","count":1},{"name":"name2","key":"key2","count":2}]`
	v2 := `[]`
	v3 := `invalid`
	v4 := `{}`
	v5 := `[1 2 3 4]`
	mapType := reflect.TypeOf([]map[string]interface{}{})

	// WHEN
	r1, err1 := parseStringContainingSliceOfX(v1, mapType)
	r2, err2 := parseStringContainingSliceOfX(v2, mapType)
	r3, err3 := parseStringContainingSliceOfX(v3, mapType)
	r4, err4 := parseStringContainingSliceOfX(v4, mapType)
	r5, err5 := parseStringContainingSliceOfX(v5, mapType)

	// THEN
	require.NoError(t, err1)
	require.Len(t, r1, 2)

	assert.Equal(t, "name1", cast.ToStringMap(r1[0])["name"])
	assert.Equal(t, "key1", cast.ToStringMap(r1[0])["key"])
	assert.Equal(t, float64(1), cast.ToStringMap(r1[0])["count"])
	assert.Equal(t, "name2", cast.ToStringMap(r1[1])["name"])
	assert.Equal(t, "key2", cast.ToStringMap(r1[1])["key"])
	assert.Equal(t, float64(2), cast.ToStringMap(r1[1])["count"])
	require.NoError(t, err2)
	assert.Empty(t, r2)
	assert.Error(t, err3)
	assert.Empty(t, r3)
	assert.Error(t, err4)
	assert.Empty(t, r4)
	assert.Error(t, err5)
	assert.Empty(t, r5)
}

func Test_handleViperWorkarounds(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string
		Field2 int
	}

	// WHEN
	valNil, errNil := handleViperWorkarounds(nil, reflect.TypeOf(0), false)
	valNoString, errNoString := handleViperWorkarounds(1, reflect.TypeOf(0), false)
	valNoSlice, errNoSlice := handleViperWorkarounds("1", reflect.TypeOf("0"), false)
	valBoolSlice, errBoolSlice := handleViperWorkarounds("[true,false,true]", reflect.TypeOf([]bool{}), false)
	valMapSlice, errMapSlice := handleViperWorkarounds(`[{"field1":"hello 1","field2":11},{"field1":"hello 2","field2":22}]`, reflect.TypeOf([]my{}), false)
	valDurationSlice, errDurationSlice := handleViperWorkarounds("", reflect.TypeOf([]time.Duration{}), false)
	valHasMapfunc, errHasMapfunc := handleViperWorkarounds("1", reflect.TypeOf(0), true)

	// THEN
	assert.NoError(t, errNil)
	assert.Nil(t, valNil)
	assert.NoError(t, errNoString)
	assert.Equal(t, 1, valNoString)
	assert.NoError(t, errNoSlice)
	assert.Equal(t, "1", valNoSlice)
	assert.NoError(t, errBoolSlice)
	bSlice := cast.ToBoolSlice(valBoolSlice)
	assert.Equal(t, []bool{true, false, true}, bSlice)
	assert.NoError(t, errMapSlice)

	mapSlice := []map[string]interface{}{}
	for _, m := range valMapSlice.([]interface{}) {
		mapSlice = append(mapSlice, cast.ToStringMap(m))
	}

	assert.Equal(t, []map[string]interface{}{
		{"field1": "hello 1", "field2": float64(11)},
		{"field1": "hello 2", "field2": float64(22)},
	}, mapSlice)
	assert.Error(t, errDurationSlice)
	assert.Nil(t, valDurationSlice)
	assert.NoError(t, errHasMapfunc)
	assert.Equal(t, "1", valHasMapfunc)
}

func Test_yamlElementListToJsonString(t *testing.T) {
	// GIVEN
	element1_1 := map[string]interface{}{"name": "hans", "age": 12}
	element1_2 := map[string]interface{}{"name": "benno", "age": 22}
	elements1 := []interface{}{element1_1, element1_2}

	element2_1 := map[string]string{"firstname": "hans", "lastname": "wurst"}
	element2_2 := map[string]string{"firstname": "benno", "lastname": "benni"}
	elements2 := []interface{}{element2_1, element2_2}

	elements3 := []interface{}{1, 2, 3}

	// WHEN
	str1, err1 := yamlElementListToJsonString(elements1)
	str2, err2 := yamlElementListToJsonString(elements2)
	_, err3 := yamlElementListToJsonString(elements3)
	str4, err4 := yamlElementListToJsonString(nil)

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, `[{"age":12,"name":"hans"},{"age":22,"name":"benno"}]`, str1)
	assert.NoError(t, err2)
	assert.Equal(t, `[{"firstname":"hans","lastname":"wurst"},{"firstname":"benno","lastname":"benni"}]`, str2)
	assert.Error(t, err3)
	assert.NoError(t, err4)
	assert.Equal(t, `[]`, str4)
}

func Test_cfgValueToStructuredString(t *testing.T) {
	// GIVEN
	element1_1 := map[string]interface{}{"name": "hans", "age": 12}
	element1_2 := map[string]interface{}{"name": "benno", "age": 22}
	elements1 := []interface{}{element1_1, element1_2}

	// WHEN
	str1, err1 := cfgValueToStructuredString(elements1)
	_, err2 := cfgValueToStructuredString(162)
	str3, err3 := cfgValueToStructuredString("hello world")

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, `[{"age":12,"name":"hans"},{"age":22,"name":"benno"}]`, str1)
	assert.Error(t, err2)
	assert.NoError(t, err3)
	assert.Equal(t, "hello world", str3)
}

func Test_handleYamlElementListInput(t *testing.T) {
	// GIVEN
	element1_1 := map[string]interface{}{"name": "hans", "age": 12}
	element1_2 := map[string]interface{}{"name": "benno", "age": 22}
	elements1 := []interface{}{element1_1, element1_2}

	type myStruct struct {
		Name string
		Age  int
	}

	elements2 := []myStruct{{}, {}}

	// WHEN
	str1, err1 := handleYamlElementListInput(elements1, reflect.TypeOf([]myStruct{}))
	str2, err2 := handleYamlElementListInput(element1_1, reflect.TypeOf(myStruct{}))
	str3, err3 := handleYamlElementListInput(element1_1, reflect.TypeOf([]myStruct{}))
	_, err4 := handleYamlElementListInput(nil, reflect.TypeOf(myStruct{}))
	str5, err5 := handleYamlElementListInput(elements2, reflect.TypeOf([]myStruct{}))

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, `[{"age":12,"name":"hans"},{"age":22,"name":"benno"}]`, str1)
	assert.NoError(t, err2)
	assert.Equal(t, element1_1, str2)
	assert.NoError(t, err3)
	assert.Equal(t, element1_1, str3)
	assert.NoError(t, err4)
	assert.NoError(t, err5)
	assert.Equal(t, elements2, str5)
}

func Test_castSimple(t *testing.T) {
	// GIVEN

	// WHEN
	v1, err1 := castSimple("a", reflect.TypeOf(""))
	v2, err2 := castSimple("true", reflect.TypeOf(bool(false)))
	v3, err3 := castSimple("12.34", reflect.TypeOf(float32(0)))
	v4, err4 := castSimple("12.34", reflect.TypeOf(float64(0)))
	v5, err5 := castSimple("1234", reflect.TypeOf(int(0)))
	v6, err6 := castSimple("123", reflect.TypeOf(int8(0)))
	v7, err7 := castSimple("1234", reflect.TypeOf(int16(0)))
	v8, err8 := castSimple("1234", reflect.TypeOf(int32(0)))
	v9, err9 := castSimple("1234", reflect.TypeOf(int64(0)))
	v10, err10 := castSimple("1234", reflect.TypeOf(uint(0)))
	v11, err11 := castSimple("123", reflect.TypeOf(uint8(0)))
	v12, err12 := castSimple("1234", reflect.TypeOf(uint16(0)))
	v13, err13 := castSimple("1234", reflect.TypeOf(uint32(0)))
	v14, err14 := castSimple("1234", reflect.TypeOf(uint64(0)))
	v15, err15 := castSimple("not possible", reflect.TypeOf([]int{}))

	// THEN
	assert.NoError(t, err1)
	assert.Equal(t, "a", v1)
	assert.NoError(t, err2)
	assert.Equal(t, true, v2)
	assert.NoError(t, err3)
	assert.Equal(t, float32(12.34), v3)
	assert.NoError(t, err4)
	assert.Equal(t, float64(12.34), v4)
	assert.NoError(t, err5)
	assert.Equal(t, int(1234), v5)
	assert.NoError(t, err6)
	assert.Equal(t, int8(123), v6)
	assert.NoError(t, err7)
	assert.Equal(t, int16(1234), v7)
	assert.NoError(t, err8)
	assert.Equal(t, int32(1234), v8)
	assert.NoError(t, err9)
	assert.Equal(t, int64(1234), v9)
	assert.NoError(t, err10)
	assert.Equal(t, uint(1234), v10)
	assert.NoError(t, err11)
	assert.Equal(t, uint8(123), v11)
	assert.NoError(t, err12)
	assert.Equal(t, uint16(1234), v12)
	assert.NoError(t, err13)
	assert.Equal(t, uint32(1234), v13)
	assert.NoError(t, err14)
	assert.Equal(t, uint64(1234), v14)
	assert.Error(t, err15)
	assert.Nil(t, v15)
}
