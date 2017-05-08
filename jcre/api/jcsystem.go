package api

import (
	"log"
	"reflect"
)

var (
	ClearOnResetTransientObjs    = make([]interface{}, 30)
	ClearOnDeselectTransientObjs = make([]interface{}, 30)
	Kr                           = -1
	Kd                           = -1
)

const (
	//ApiVersion = 513
	// PrivAccess thePrivAccess
	NotATransientObject = 0
	ClearOnReset        = 1
	ClearOnDeselect     = 2
)

func IsTransient(object interface{}) int8 {
	if object != nil {
		if contains(ClearOnResetTransientObjs, object) {
			return 1
		}
		if contains(ClearOnDeselectTransientObjs, object) {
			return 2
		}
	}
	return 0
}
func contains(object []interface{}, item interface{}) bool {
	for i := 0; i < len(object); i++ {
		if cond := reflect.DeepEqual(object[i], item); cond == true {
			return true
		}
	}
	return false
}
func AddTransientArray(array interface{}, event int8) {
	if IsTransient(array) != 0 {
		//exception
		log.Fatal("Error adding transient array")
	}
	switch event {
	case 1:
		Kr++
		ClearOnResetTransientObjs[Kr] = array
	case 2:
		Kd++
		ClearOnDeselectTransientObjs[Kd] = array
	default:
		//nothing
		log.Fatal("Error adding transient array")
	}
}
func MakeTransientBooleanArray(length int16, event int8) []bool {
	array := make([]bool, int(length))
	AddTransientArray(array, event)
	return array
}
func MakeTransientByteArray(length int16, event int8) []byte {
	array := make([]byte, int(length))
	AddTransientArray(array, event)
	return array
}
func MakeTransientShortArray(length int16, event int8) []int16 {
	array := make([]int16, int(length))
	AddTransientArray(array, event)
	return array
}
func MakeTransientObjectArray(length int16, event int8) interface{} {
	array := make([]interface{}, int(length))
	AddTransientArray(array, event)
	return array
}

/*
func MakeTransientArray(array interface{}, length int16, event int8) {
	switch array.(type) {
	case *[]bool:
		array = make([]bool, int(length))
	case *[]byte:
		array = make([]byte, int(length))
	case *[]int16: //short
		array = make([]int16, int(length))
	default: //any other object
		array = make([]interface{}, int(length))
	}
	AddTransientArray(array, event)
	length++
}*/
