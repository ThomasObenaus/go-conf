package config

import (
	"reflect"
	"testing"

	"github.com/ThomasObenaus/go-conf/interfaces"
	mock_provider "github.com/ThomasObenaus/go-conf/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ignoreAllCallsToLogger(mockedProvider *mock_provider.MockProvider) {
	// just please the calls to the logger
	mockedProvider.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockedProvider.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockedProvider.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

func Test_getTargetTypeAndValue(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	m := my{
		Field1: "field1",
	}

	// WHEN
	t1, v1, err1 := getTargetTypeAndValue(&m)

	// THEN
	require.NoError(t, err1)
	assert.Equal(t, reflect.TypeOf(m), t1)
	assert.NotEqual(t, reflect.Ptr, t1.Kind())
	assert.True(t, v1.IsValid())
}

func Test_getTargetTypeAndValue_Fail(t *testing.T) {
	// GIVEN
	type my struct {
		Field1 string `cfg:"{'name':'field_1'}"`
	}
	m := my{
		Field1: "field1",
	}

	// WHEN
	t1, v1, err1 := getTargetTypeAndValue(m)

	// THEN
	require.Error(t, err1)
	assert.Nil(t, t1)
	assert.True(t, v1.IsZero())

	// WHEN
	t2, v2, err2 := getTargetTypeAndValue(nil)

	// THEN
	require.Error(t, err2)
	assert.Nil(t, t2)
	assert.True(t, v2.IsZero())

	// GIVEN
	var nilVal *my

	// WHEN
	t3, v3, err3 := getTargetTypeAndValue(nilVal)

	// THEN
	require.Error(t, err3)
	assert.Nil(t, t3)
	assert.True(t, v3.IsZero())
}

func Test_applyConfig_Empty(t *testing.T) {
	// GIVEN
	type empty struct {
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockedProvider := mock_provider.NewMockProvider(mockCtrl)
	myTestCfg := empty{}
	ignoreAllCallsToLogger(mockedProvider)

	// WHEN
	err := applyConfig(mockedProvider, &myTestCfg, "", configTag{}, map[string]interfaces.MappingFunc{})

	// THEN
	assert.NoError(t, err)
}

func Test_applyConfig(t *testing.T) {
	// GIVEN
	type myNestedConfig struct {
		FieldA int `cfg:"{'name':'field-a'}"`
		FieldB int `cfg:"{'name':'field-b'}"`
	}

	type myTestConfig struct {
		Field1 string           `cfg:"{'name':'field-1'}"`
		Field2 myNestedConfig   `cfg:"{'name':'field-2'}"`
		Field3 []int            `cfg:"{'name':'field-3'}"`
		Field4 []myNestedConfig `cfg:"{'name':'field-4'}"`
		Field5 []myNestedConfig `cfg:"{'name':'field-5'}"`
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockedProvider := mock_provider.NewMockProvider(mockCtrl)
	myTestCfg := myTestConfig{}
	ignoreAllCallsToLogger(mockedProvider)

	mockedProvider.EXPECT().IsSet("field-1").Return(true)
	mockedProvider.EXPECT().Get("field-1").Return("value field-1")
	mockedProvider.EXPECT().IsSet("field-2.field-a").Return(true)
	mockedProvider.EXPECT().Get("field-2.field-a").Return(11)
	mockedProvider.EXPECT().IsSet("field-2.field-b").Return(true)
	mockedProvider.EXPECT().Get("field-2.field-b").Return(22)
	mockedProvider.EXPECT().IsSet("field-3").Return(true)
	mockedProvider.EXPECT().Get("field-3").Return([]int{1, 2, 3})
	mockedProvider.EXPECT().IsSet("field-4").Return(true)
	field4Value := []myNestedConfig{
		{FieldA: 11, FieldB: 12},
		{FieldA: 21, FieldB: 22},
	}
	mockedProvider.EXPECT().Get("field-4").Return(field4Value)
	mockedProvider.EXPECT().IsSet("field-5").Return(true)
	mockedProvider.EXPECT().Get("field-5").Return("[{'field-a':33,'field-b':44}]")

	// WHEN
	err := applyConfig(mockedProvider, &myTestCfg, "", configTag{}, map[string]interfaces.MappingFunc{})

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "value field-1", myTestCfg.Field1)
	assert.Equal(t, 11, myTestCfg.Field2.FieldA)
	assert.Equal(t, 22, myTestCfg.Field2.FieldB)
	assert.Len(t, myTestCfg.Field3, 3)
	assert.Equal(t, 1, myTestCfg.Field3[0])
	assert.Equal(t, 2, myTestCfg.Field3[1])
	assert.Equal(t, 3, myTestCfg.Field3[2])
	assert.Len(t, myTestCfg.Field4, 2)
	assert.Equal(t, 11, myTestCfg.Field4[0].FieldA)
	assert.Equal(t, 12, myTestCfg.Field4[0].FieldB)
	assert.Equal(t, 21, myTestCfg.Field4[1].FieldA)
	assert.Equal(t, 22, myTestCfg.Field4[1].FieldB)
	assert.Len(t, myTestCfg.Field5, 1)
	assert.Equal(t, 33, myTestCfg.Field5[0].FieldA)
	assert.Equal(t, 44, myTestCfg.Field5[0].FieldB)
}

func Test_applyConfig_Fail(t *testing.T) {
	// GIVEN - Wrong type returned
	type myTestConfigWrongType struct {
		Field1 int `cfg:"{'name':'field-1'}"`
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockedProviderWrongType := mock_provider.NewMockProvider(mockCtrl)
	myTestCfgWrongType := myTestConfigWrongType{}
	ignoreAllCallsToLogger(mockedProviderWrongType)

	mockedProviderWrongType.EXPECT().IsSet("field-1").Return(true)
	mockedProviderWrongType.EXPECT().Get("field-1").Return("that is not an int")

	// WHEN
	errWrongType := applyConfig(mockedProviderWrongType, &myTestCfgWrongType, "", configTag{}, map[string]interfaces.MappingFunc{})

	// THEN
	assert.Error(t, errWrongType)

	// GIVEN - Wrong type returned for slice of structs
	type myNestedConfig struct {
		FieldA int `cfg:"{'name':'field-a'}"`
	}
	type myTestConfigWrongTypeSliceOfStructs struct {
		Field1 []myNestedConfig `cfg:"{'name':'field-1'}"`
	}

	mockedProviderWrongTypeSliceOfStructs := mock_provider.NewMockProvider(mockCtrl)
	myTestCfgWrongTypeSliceOfStructs := myTestConfigWrongTypeSliceOfStructs{}
	ignoreAllCallsToLogger(mockedProviderWrongTypeSliceOfStructs)

	mockedProviderWrongTypeSliceOfStructs.EXPECT().IsSet("field-1").Return(true)
	mockedProviderWrongTypeSliceOfStructs.EXPECT().Get("field-1").Return("that is not an int")

	// WHEN
	errWrongTypeSliceOfStructs := applyConfig(mockedProviderWrongTypeSliceOfStructs, &myTestCfgWrongTypeSliceOfStructs, "", configTag{}, map[string]interfaces.MappingFunc{})

	// THEN
	assert.Error(t, errWrongTypeSliceOfStructs)

	// GIVEN - target is no pointer to a struct
	mockedProvider := mock_provider.NewMockProvider(mockCtrl)
	myTestCfg := myTestConfigWrongTypeSliceOfStructs{}
	ignoreAllCallsToLogger(mockedProvider)

	// WHEN
	err := applyConfig(mockedProvider, myTestCfg, "", configTag{}, map[string]interfaces.MappingFunc{})

	// THEN
	assert.Error(t, err)
}

func Test_applyConfig_slices(t *testing.T) {
	// GIVEN
	type myTestConfig struct {
		Field1 []int    `cfg:"{'name':'field-1','default':[1,2,3]}"`
		Field2 []bool   `cfg:"{'name':'field-2','default':[true,false]}"`
		Field3 []string `cfg:"{'name':'field-3','default':['a','b','c']}"`
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockedProvider := mock_provider.NewMockProvider(mockCtrl)
	myTestCfg := myTestConfig{}
	ignoreAllCallsToLogger(mockedProvider)

	mockedProvider.EXPECT().IsSet("field-1").Return(true)
	mockedProvider.EXPECT().Get("field-1").Return([]int{3, 2, 1})
	mockedProvider.EXPECT().IsSet("field-2").Return(true)
	mockedProvider.EXPECT().Get("field-2").Return([]bool{false, true})
	mockedProvider.EXPECT().IsSet("field-3").Return(true)
	mockedProvider.EXPECT().Get("field-3").Return([]string{"c", "b", "a"})

	// WHEN
	err := applyConfig(mockedProvider, &myTestCfg, "", configTag{}, map[string]interfaces.MappingFunc{})

	// THEN
	assert.NoError(t, err)
	assert.Len(t, myTestCfg.Field1, 3)
	assert.Len(t, myTestCfg.Field2, 2)
	assert.Len(t, myTestCfg.Field3, 3)
}
