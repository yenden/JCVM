package api

import (
	"log"
	"reflect"
)

var (
	/*ClearOnResetTransientObjs represent objects
	* that need to be cleared on reset
	 */
	ClearOnResetTransientObjs = make([]interface{}, 30)

	/*ClearOnDeselectTransientObjs represent objects
	* that need to be cleared when the applet is deselect
	 */
	ClearOnDeselectTransientObjs = make([]interface{}, 30)

	kr = -1 //indices clearonReset elements
	kd = -1 //indices clearonDeselect elements
)

const (
	//To identify the type of transient objects
	NotATransientObject = 0
	ClearOnReset        = 1
	ClearOnDeselect     = 2
)

/*IsTransient checks if an object is transient*/
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

/*AddTransientArray adds an array as transient*/
func AddTransientArray(array interface{}, event int8) {
	if IsTransient(array) != 0 {
		//exception
		log.Fatal("Error adding transient array")
	}
	switch event {
	case 1:
		kr++
		ClearOnResetTransientObjs[kr] = array
	case 2:
		kd++
		ClearOnDeselectTransientObjs[kd] = array
	default:
		//nothing
		log.Fatal("Error adding transient array")
	}
}

/*MakeTransientBooleanArray creates a transient boolean array */
func MakeTransientBooleanArray(length int16, event int8) []bool {
	array := make([]bool, int(length))
	AddTransientArray(array, event)
	return array
}

/*MakeTransientByteArray creates a transient byte array */
func MakeTransientByteArray(length int16, event int8) []byte {
	array := make([]byte, int(length))
	AddTransientArray(array, event)
	return array
}

/*MakeTransientShortArray creates a transient short array */
func MakeTransientShortArray(length int16, event int8) []int16 {
	array := make([]int16, int(length))
	AddTransientArray(array, event)
	return array
}

/*MakeTransientObjectArray creates a transient array of  objects*/
func MakeTransientObjectArray(length int16, event int8) interface{} {
	array := make([]interface{}, int(length))
	AddTransientArray(array, event)
	return array
}
