package framework

import (
	"errors"
	"reflect"
)

var (
	clearOnResetTransientObjs    = make([][]interface{}, 5)
	clearOnDeselectTransientObjs = make([][]interface{}, 5)
	k                            = -1
)

const (
	//ApiVersion = 513
	// PrivAccess thePrivAccess
	NotATransientObject = 0
	ClearOnReset        = 1
	ClearOnDeselect     = 2
)

func IsTransient(object []interface{}) int8 {
	if object != nil {
		if contains(clearOnResetTransientObjs, object) {
			return 1
		}
		if contains(clearOnDeselectTransientObjs, object) {
			return 2
		}
	}
	return 0
}
func contains(object [][]interface{}, item interface{}) bool {
	for i := 0; i < len(object); i++ {
		if cond := reflect.DeepEqual(object[i], item); cond == true {
			return true
		}
	}
	return false
}
func AddTransientArray(array []interface{}, event int8) error {
	if IsTransient(array) != 0 {
		//exception
		return errors.New("Error adding transient array")
	}
	switch event {
	case 1:
		k++
		clearOnResetTransientObjs[k] = array
	case 2:
		k++
		clearOnDeselectTransientObjs[k] = array
	default:
		//nothing
		return errors.New("Error adding transient array")
	}
	return nil
}
func MakeTransientArray(length int8, event int8, typ uint8) []interface{} {
	array := make([]interface{}, int(length))
	switch typ {
	case 0: //bool
		for i := 0; i < len(array); i++ {
			array[i] = bool(false)
		}
	case 1: //byte
		for i := 0; i < len(array); i++ {
			array[i] = byte(0)
		}
	case 2: //short
		for i := 0; i < len(array); i++ {
			array[i] = int16(0)
		}
	default: //object ie everything
		//nothing
	}
	AddTransientArray(array, event)
	return array
}
