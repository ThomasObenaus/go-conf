package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func isFieldExported(typeOfField reflect.StructField) bool {
	return typeOfField.PkgPath == ""
}

// castToPrimitive supports casting of primitive types (such as int, string,...) to the given target type.
func castToPrimitive(rawValue interface{}, targetType reflect.Type) (interface{}, error) {
	typeOfValue := reflect.TypeOf(rawValue)

	if targetType == reflect.TypeOf(time.Second) {
		dur, err := cast.ToDurationE(rawValue)
		if err != nil {
			return nil, errors.Wrapf(err, "Casting '%v' to duration", rawValue)
		}
		rawValue = dur
		typeOfValue = reflect.TypeOf(rawValue)
	}

	if !typeOfValue.ConvertibleTo(targetType) {
		return nil, fmt.Errorf("Can't convert %v to %v", typeOfValue, targetType)
	}
	return reflect.ValueOf(rawValue).Convert(targetType).Interface(), nil
}

// castToStruct supports casting of structs (also nested) to the given target type.
func castToStruct(rawValue interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Can't cast to struct since the target type is not a struct. Instead it is %v", targetType)
	}

	parsedDefinitionCastedToMap, ok := rawValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to cast %v (type=%T) to %v. Type must be map[string]interface{}", rawValue, rawValue, targetType)
	}
	castedToTargetType, err := createAndFillStruct(targetType, parsedDefinitionCastedToMap)
	if err != nil {
		return nil, errors.Wrap(err, "Handling default value for element in a slice of structs")
	}
	return castedToTargetType.Interface(), nil
}

// castToSlice supports casting of slices (of primitives and structs) to the given target type.
func castToSlice(rawValue interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Can't cast to slice since the target type is not a slice. Instead it is %v", targetType)
	}

	typeOfRawValue := reflect.TypeOf(rawValue)
	if typeOfRawValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Can't cast to slice since the given raw value is no slice. Instead it is %v", typeOfRawValue)
	}

	sliceValue := reflect.ValueOf(rawValue)

	// obtain the type of the slices elements
	elementType := targetType.Elem()
	sliceInTargetType := reflect.MakeSlice(targetType, 0, sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		rawDefaultValueElement := sliceValue.Index(i).Interface()
		switch castedRawElement := rawDefaultValueElement.(type) {
		case map[string]interface{}:
			// handles structs
			castedToTargetType, err := createAndFillStruct(elementType, castedRawElement)
			if err != nil {
				return nil, errors.Wrap(err, "Handling default value for element in a slice of structs")
			}
			sliceInTargetType = reflect.Append(sliceInTargetType, castedToTargetType)
		default:
			// handles primitive elements (int, string, ...)
			castedToTargetType, err := castToPrimitive(rawDefaultValueElement, elementType)
			if err != nil {
				return nil, errors.Wrap(err, "Casting default value to element type")
			}
			sliceInTargetType = reflect.Append(sliceInTargetType, reflect.ValueOf(castedToTargetType))
		}

	}
	return sliceInTargetType.Interface(), nil
}

// castToTargetType casts the given raw value to the given target type.
func castToTargetType(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.Struct:
		return castToStruct(rawUntypedValue, targetType)
	case reflect.Slice:
		return castToSlice(rawUntypedValue, targetType)
	default:
		return castToPrimitive(rawUntypedValue, targetType)
	}
}

func getConfigTagDefinition(fieldDeclaration reflect.StructField) (string, bool) {
	return fieldDeclaration.Tag.Lookup("cfg")
}

// createAndFillStruct creates a struct based on the given type and fills its fields based on the given data.
// For being able to fill the struct the given datas keys have to match the config tags that are defined on the target type.
//
// e.g. for type
//
//	type my struct {
//		Field1 string `cfg:"{'name':'field_1'}"`
//	}
//
// the data map should contain an entry with name 'field_1'
// 	data := map[string]interface{}{"field_1":"a value"}
func createAndFillStruct(targetTypeOfStruct reflect.Type, data map[string]interface{}) (reflect.Value, error) {
	if targetTypeOfStruct.Kind() != reflect.Struct {
		return reflect.Zero(targetTypeOfStruct), fmt.Errorf("The target type must be a struct")
	}

	newStruct := reflect.New(targetTypeOfStruct)
	newStructValue := newStruct.Elem()

	for i := 0; i < targetTypeOfStruct.NumField(); i++ {
		fieldDeclaration := targetTypeOfStruct.Field(i)
		fieldValue := newStructValue.FieldByName(fieldDeclaration.Name)
		fieldType := fieldDeclaration.Type
		configTag, hasConfig := getConfigTagDefinition(fieldDeclaration)
		if !hasConfig {
			continue
		}

		entry, err := parseConfigTagDefinition(configTag, fieldType, "")
		if err != nil {
			return reflect.Zero(targetTypeOfStruct), errors.Wrapf(err, "Parsing configTag '%s'", configTag)
		}
		val, ok := data[entry.Name]
		if !ok {
			if entry.IsRequired() {
				return reflect.Zero(targetTypeOfStruct), fmt.Errorf("Missing value for required field (struct-field='%s',expected-key='%s')", fieldDeclaration.Name, entry.Name)
			}

			// take the default value
			val = entry.Def
		}

		// cast the parsed default value to the target type
		castedToTargetTypeIf, err := castToTargetType(val, fieldType)
		if err != nil {
			return reflect.Zero(targetTypeOfStruct), errors.Wrapf(err, "Casting to target type")
		}
		castedToTargetType := reflect.ValueOf(castedToTargetTypeIf)

		// ensure that the casted value can be set
		if !isFieldExported(fieldDeclaration) {
			return reflect.Zero(targetTypeOfStruct), fmt.Errorf("Can't set value for unexported field (struct-field='%s',key='%s').", fieldDeclaration.Name, entry.Name)
		}
		if !fieldValue.CanSet() {
			return reflect.Zero(targetTypeOfStruct), fmt.Errorf("Can't set value for field (struct-field='%s',key='%s').", fieldDeclaration.Name, entry.Name)
		}
		fieldValue.Set(castedToTargetType)
	}

	return newStructValue, nil
}

func fullFieldName(nameOfParent string, fieldName string) string {
	if len(nameOfParent) == 0 {
		return fieldName
	}
	return fmt.Sprintf("%s.%s", nameOfParent, fieldName)
}

func parseStringContainingSliceOfX(sliceString string, targetSliceType reflect.Type) ([]interface{}, error) {
	sliceString = strings.ReplaceAll(sliceString, "'", `"`)

	slice := reflect.MakeSlice(targetSliceType, 0, 0).Interface()
	err := json.Unmarshal([]byte(sliceString), &slice)
	if err != nil {
		return nil, errors.Wrapf(err, "Parsing string that contains a slice of %v", targetSliceType)
	}

	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return nil, fmt.Errorf("Unmarshalling did not produce the expected slice type instead it produced '%T'", slice)
	}

	return slice.([]interface{}), nil
}

// handleViperWorkarounds viper does not handle all types correctly. e.g. a slice of structs or booleans is not supported and just returned as
// a jsonstring. handleViperWorkarounds casts those jsonstrings into the correct golang types.
func handleViperWorkarounds(val interface{}, targetType reflect.Type, hasMappingfunc bool) (interface{}, error) {
	if val == nil {
		return val, nil
	}

	if hasMappingfunc {
		return val, nil
	}

	typeOfValue := reflect.TypeOf(val)

	// immediately return / do nothing in case we have no string
	if typeOfValue.Kind() != reflect.String {
		return val, nil
	}

	// immediately return / do nothing in case the target type is no slice
	if targetType.Kind() != reflect.Slice {
		return val, nil
	}

	valueAsString, err := cast.ToStringE(val)
	if err != nil {
		return nil, errors.Wrapf(err, "Casting %v (type=%T) to string", val, val)
	}
	return parseStringContainingSliceOfX(valueAsString, targetType)
}

// handleYamlElementListInput ensures that the value is coming in as an yaml element list are in the right format because usually they are given as []interface{}.
// Internally they are converted to a json string.
func handleYamlElementListInput(val interface{}, targetType reflect.Type) (interface{}, error) {
	if val == nil {
		return val, nil
	}
	typeOfValue := reflect.TypeOf(val)

	// do nothing if the target type is no slice
	if targetType.Kind() != reflect.Slice {
		return val, nil
	}
	// do nothing if the given type is no slice
	if typeOfValue.Kind() != reflect.Slice {
		return val, nil
	}

	targetTypeElementKind := targetType.Elem().Kind()
	valueTypeElementKind := typeOfValue.Elem().Kind()
	// do nothing if the element types match
	if targetTypeElementKind == valueTypeElementKind {
		return val, nil
	}

	return cfgValueToStructuredString(val)
}

// cfgValueToStructuredString handles structured config values that are either valid json or yaml element lists.
// The function returns a json string that represents the input data.
func cfgValueToStructuredString(cfgValue interface{}) (string, error) {
	cfgValueAsStr := ""

	// We have to handle two cases.
	// 1. The config parameter was provided as one string (e.g. via env-parameter or via command-line).
	// 2. The config parameter was provided via yaml as a list of objects (will be type []interface{})
	switch typedIface := cfgValue.(type) {
	case string:
		// This handles the case that config parameter was provided as one string.
		// example [{'key1':'value1','key2':'value2'}]
		cfgValueAsStr = cfgValue.(string)
	case []interface{}:
		// This handles the case that config parameter was provided via yaml as a list of objects.
		return yamlElementListToJsonString(typedIface)
	default:
		return "", fmt.Errorf("Unable to handle parameter of type %T", cfgValue)
	}
	return cfgValueAsStr, nil
}

// yamlElementListToJsonString converts contents of a yaml list into a json string
func yamlElementListToJsonString(yamlData []interface{}) (string, error) {
	// This handles the case that config parameter was provided via yaml as a list of objects.
	// example :
	// my-config:
	// - key1: "value1"
	//   key2: 111
	// - key1: "value2"
	//   key2: 222
	pairList := make([]string, 0)

	for _, element := range yamlData {
		pair := ""
		switch castedElement := element.(type) {
		case map[string]string:
			internalPair, err := json.Marshal(castedElement)
			if err != nil {
				return "", errors.Wrap(err, "Handling map[string]string")
			}
			pair = string(internalPair)
		case map[string]interface{}, map[interface{}]interface{}:
			cElement, err := cast.ToStringMapE(element)
			if err != nil {
				return "", errors.Wrap(err, "Casting map[string]interface{}")
			}
			internalPair, err := json.Marshal(cElement)
			if err != nil {
				return "", errors.Wrap(err, "Handling map[string]interface{}")
			}
			pair = string(internalPair)
		default:
			return "", fmt.Errorf("Unable to handle element %#v of type %T", element, element)
		}

		pairList = append(pairList, pair)
	}
	// concatenate all elements to [{'key1':'value1','key2':'value2'}]
	return fmt.Sprintf(`[%s]`, strings.Join(pairList, ",")), nil
}
