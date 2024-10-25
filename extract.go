package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ThomasObenaus/go-conf/interfaces"
	"github.com/pkg/errors"
)

// configTag represents the definition for a config read from the type tag.
// A config tag on a type is expected to be defined as:
//
//	`cfg:"{'name':'<name of the config>','desc':'<description>','default':<default value>,'short':'<shorthand name for the flag>','mapfun':<name of the mapping function>}"`
type configTag struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"desc,omitempty"`
	Def         interface{} `json:"default,omitempty"`
	MapFunName  string      `json:"mapfun,omitempty"`
	ShortName   string      `json:"short,omitempty"`
	desiredType reflect.Type

	// isComplexTypeWithoutAnnotatedFields is true if this configTag represents a field whose type is a non primitive one and whose fields are not annotated with the cfg tag.
	isComplexTypeWithoutAnnotatedFields bool
}

func (e configTag) String() string {
	return fmt.Sprintf(`name:"%s",desc:"%s",default:%v (%T),required=%t,isComplexTypeWithoutAnnotatedFields=%t`, e.Name, e.Description, e.Def, e.Def, e.IsRequired(), e.isComplexTypeWithoutAnnotatedFields)
}

func (e configTag) IsRequired() bool {
	return e.Def == nil
}

func (e configTag) HasMapfunc() bool {
	return len(e.MapFunName) > 0
}

// parseConfigTagDefinition parses a definition like
//
//	`cfg:"{'name':'<name of the config>','desc':'<description>','default':<default value>,'mapfun':<name of the mapping function>}"`
//
// to a configTag
func parseConfigTagDefinition(configTagStr string, typeOfEntry reflect.Type, nameOfParent string) (configTag, error) {
	configTagStr = strings.TrimSpace(configTagStr)
	// replace all single quotes by double quotes to get a valid json
	configTagStr = strings.ReplaceAll(configTagStr, "'", `"`)

	// parse the config tag
	parsedDefinition := configTag{}
	if err := json.Unmarshal([]byte(configTagStr), &parsedDefinition); err != nil {
		return configTag{}, errors.Wrapf(err, "Parsing configTag from '%s'", configTagStr)
	}

	if len(parsedDefinition.Name) == 0 {
		return configTag{}, fmt.Errorf("Missing required config tag field 'name' on '%s'", configTagStr)
	}

	result := configTag{
		// update name to reflect the hierarchy
		Name:        fullFieldName(nameOfParent, parsedDefinition.Name),
		Description: parsedDefinition.Description,
		desiredType: typeOfEntry,
		MapFunName:  parsedDefinition.MapFunName,
		ShortName:   parsedDefinition.ShortName,
	}

	// only in case a default value is given
	if parsedDefinition.Def != nil {

		// This handles the case where the type of a field defined in the config annotation does not match
		// the type of the field that is annotated.
		// Example:
		// type cfg struct {
		// 	F1 zerolog.Level `cfg:"{'name':'logl','default':'info'}"`
		// }
		// Here F1 is of type zerolog.Level (int8) and the defined type in the annotation is string (based on the default value)
		//
		// In order to support this situation we take the type of the config annotation to cast the default value.
		if reflect.TypeOf(parsedDefinition.Def) == reflect.TypeOf("") {
			typeOfEntry = reflect.TypeOf("")
		}

		castedValue, err := castToTargetType(parsedDefinition.Def, typeOfEntry)
		if err != nil {
			return configTag{}, errors.Wrap(err, "Casting parsed default value to target type")
		}
		result.Def = castedValue
	}
	return result, nil
}

// extractConfigTagFromStructField extracts the configTag from the given StructField.
// Beside the extracted configTag a bool value indicating if the given type is a primitive type is returned.
func extractConfigTagFromStructField(field reflect.StructField, parent configTag) (isPrimitive bool, tag *configTag, err error) {
	fType := field.Type

	// find out if we have a primitive type
	isPrimitive, err = isOfPrimitiveType(fType)
	if err != nil {
		return false, nil, errors.Wrapf(err, "Checking for primitive type failed for field '%v'", field)
	}

	configTagDefinition, hasCfgTag := getConfigTagDefinition(field)
	if !hasCfgTag {
		return isPrimitive, nil, nil
	}

	cfgTag, err := parseConfigTagDefinition(configTagDefinition, fType, parent.Name)
	if err != nil {
		return isPrimitive, nil, errors.Wrapf(err, "Parsing the config definition ('%s') failed for field '%v'", configTagDefinition, field)
	}

	return isPrimitive, &cfgTag, nil
}

// CreateEntriesFromStruct creates Entries based on the annotations provided at the given target struct.
//
// Only fields with annotations of the form
//
//	`cfg:"{'name':<name>,'desc':<description>,'default':<default value>}"`
//
// will be regarded.
//
// For example for the struct below
//
//	type Cfg struct {
//		Name string `cfg:"{'name':'name','desc':'the name of the config','default':'the name'}"`
//	}
//
// A config entry
//
//	e := NewEntry("name","the name of the config",Default("the name"))
//
// will be created.
func CreateEntriesFromStruct(target interface{}, logger interfaces.LoggerFunc) ([]Entry, error) {

	entries := make([]Entry, 0)

	configTags, err := extractConfigTagsOfStruct(target, logger, "", configTag{})
	if err != nil {
		return nil, err
	}
	for _, configTag := range configTags {
		desiredType := configTag.desiredType
		// Use string as desired type in case a mapping function is defined.
		if configTag.HasMapfunc() {
			desiredType = reflect.TypeOf("")
		}

		// create and append the new config entry
		entry := NewEntry(configTag.Name, configTag.Description, Default(configTag.Def), DesiredType(desiredType), ShortName(configTag.ShortName))
		entries = append(entries, entry)
		logger(interfaces.LogLevelInfo, "Added new config new entry=%v\n", entry)
	}

	return entries, nil
}

// extractConfigTagsOfStruct extracts recursively all configTags from the given struct.
// Fields of the target struct that are not annotated with a configTag are ignored.
//
// target - the target that should be processed (has to be a pointer to a struct)
// nameOfParentField - the name of the targets parent field. This is needed since this function runs recursively through the given target struct.
// parent - the configTag of the targets parent field. This is needed since this function runs recursively through the given target struct.
func extractConfigTagsOfStruct(target interface{}, logger interfaces.LoggerFunc, nameOfParentField string, parent configTag) ([]configTag, error) {

	entries := make([]configTag, 0)

	targetType := reflect.TypeOf(target)

	logger(interfaces.LogLevelDebug, "[Extract-(%s)] structure-type=%v definition=%v\n", nameOfParentField, targetType, parent)

	err := processAllConfigTagsOfStruct(target, logger, nameOfParentField, parent, func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error {
		logPrefix := fmt.Sprintf("[Extract-(%s)]", fieldName)

		if !isPrimitive && !cfgTag.isComplexTypeWithoutAnnotatedFields {
			fieldValueIf := fieldValue.Addr().Interface()
			subEntries, err := extractConfigTagsOfStruct(fieldValueIf, logger, fieldName, cfgTag)
			if err != nil {
				return errors.Wrap(err, "Extracting subentries")
			}
			entries = append(entries, subEntries...)

			logger(interfaces.LogLevelDebug, "%s added %d configTags.\n", logPrefix, len(entries))
			return nil
		}

		// This is a non primitive type which is annotated with a cfg struct tag.
		// But none of this types fields is annotated.
		// This is usually the case for external types, which are types that can't be annotated.
		// To allow the usage of external types, the entry itself is added.
		if !isPrimitive && cfgTag.isComplexTypeWithoutAnnotatedFields {
			entries = append(entries, cfgTag)
			logger(interfaces.LogLevelInfo, "%s complex type has no annotated fields. Hence a config entry for this field is added.\n", logPrefix)
			return nil
		}

		entries = append(entries, cfgTag)
		logger(interfaces.LogLevelDebug, "%s added configTag entry=%v.\n", logPrefix, cfgTag)

		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "Extracting config tags from type %v", targetType)
	}
	return entries, nil
}

// handleConfigTagFunc function type for handling an extracted configTag of a given field
type handleConfigTagFunc func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error

// processAllConfigTagsOfStruct finds the configTag on each field of the given struct. Each of this configTags will handled by the given handleConfigTagFunc.
// Fields of the target struct that are not annotated with a configTag are ignored (handleConfigTagFunc won't be called).
//
// target - the target that should be processed (has to be a pointer to a struct)
// nameOfParentField - the name of the targets parent field. This is needed since this function runs recursively through the given target struct.
// parent - the configTag of the targets parent field. This is needed since this function runs recursively through the given target struct.
// handleConfigTagFun - a function that should be used to handle each of the targets struct fields.
func processAllConfigTagsOfStruct(target interface{}, logger interfaces.LoggerFunc, nameOfParentField string, parent configTag, handleConfigTagFun handleConfigTagFunc) error {
	if target == nil {
		return fmt.Errorf("The target must not be nil")
	}

	targetType, targetValue, err := getTargetTypeAndValue(target)
	if err != nil {
		return errors.Wrapf(err, "Obtaining target type and -value for target='%v',nameOfParentField='%s',parent='%s'", target, nameOfParentField, parent)
	}

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := targetValue.Field(i)
		fType := field.Type

		fieldName := fullFieldName(nameOfParentField, field.Name)
		logPrefix := fmt.Sprintf("[Process-(%s)]", fieldName)
		logger(interfaces.LogLevelDebug, "%s field-type=%s\n", logPrefix, fType)

		isPrimitive, cfgTag, err := extractConfigTagFromStructField(field, parent)
		if err != nil {
			return errors.Wrap(err, "Extracting config tag")
		}

		// skip the field in case there is no config tag
		if cfgTag == nil {
			logger(interfaces.LogLevelInfo, "%s no tag found entry will be skipped.\n", logPrefix)
			continue
		}

		// This is a non primitive type whose fields are not annotated
		if !isPrimitive && !hasAnnotatedFields(fType) {
			cfgTag.isComplexTypeWithoutAnnotatedFields = true
		}

		logger(interfaces.LogLevelDebug, "%s parsed config entry=%v. Is primitive=%t.\n", logPrefix, cfgTag, isPrimitive)

		err = handleConfigTagFun(fieldName, isPrimitive, fType, fieldValue, *cfgTag)
		if err != nil {
			return errors.Wrapf(err, "Handling configTag %s for field '%s'", *cfgTag, fieldName)
		}
	}
	return nil
}

// hasAnnotatedFields returns true if the given complex type (struct) has at least one field with a cfg tag annotation.
func hasAnnotatedFields(t reflect.Type) bool {
	if t == nil {
		return false
	}

	result := false

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		_, hasCfgTag := getConfigTagDefinition(field)
		if hasCfgTag {
			return true
		}
	}
	return result
}

// isOfPrimitiveType returns true if the given type is a primitive one (can be easily casted).
// This is also the case for slices.
func isOfPrimitiveType(fieldType reflect.Type) (bool, error) {
	kind := fieldType.Kind()
	switch kind {
	case reflect.Struct:
		return false, nil
	case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return true, nil
	case reflect.Ptr:
		elementType := fieldType.Elem()
		return isOfPrimitiveType(elementType)
	case reflect.Slice:
		return true, nil
	default:
		return false, fmt.Errorf("Kind '%s' with type '%s' is not supported", kind, fieldType)
	}
}
