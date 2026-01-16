package converter

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// ToLua converts a Go value to a Lua value.
// Supports: string, numbers (int, int32, int64, float32, float64), bool,
// map[string]interface{}, []interface{}, and nil.
// Unknown types are converted to string representation.
func ToLua(L *lua.LState, v interface{}) lua.LValue {
	if v == nil {
		return lua.LNil
	}

	switch val := v.(type) {
	case string:
		return lua.LString(val)
	case float64:
		return lua.LNumber(val)
	case float32:
		return lua.LNumber(val)
	case int:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case int32:
		return lua.LNumber(val)
	case bool:
		return lua.LBool(val)
	case map[string]interface{}:
		return MapToTable(L, val)
	case []interface{}:
		return SliceToTable(L, val)
	default:
		// For unknown types, convert to string
		return lua.LString(fmt.Sprintf("%v", val))
	}
}

// MapToTable converts a Go map[string]interface{} to a Lua table.
func MapToTable(L *lua.LState, m map[string]interface{}) *lua.LTable {
	table := L.NewTable()

	for k, v := range m {
		table.RawSetString(k, ToLua(L, v))
	}

	return table
}

// SliceToTable converts a Go []interface{} to a Lua table.
// Lua arrays are 1-indexed, so the slice is converted accordingly.
func SliceToTable(L *lua.LState, arr []interface{}) *lua.LTable {
	table := L.NewTable()

	for i, item := range arr {
		table.RawSetInt(i+1, ToLua(L, item)) // Lua arrays are 1-indexed
	}

	return table
}

// FromLua converts a Lua value to a Go value.
// Supports: nil, bool, string, number, table (array or map).
func FromLua(L *lua.LState, lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		return TableToGo(L, v)
	default:
		return v.String()
	}
}

// TableToGo converts a Lua table to a Go value (map or slice).
// It detects if the table is an array (consecutive integer keys starting from 1)
// or a map (mixed or string keys).
func TableToGo(L *lua.LState, tbl *lua.LTable) interface{} {
	// Check if it's an array (consecutive integer keys starting from 1)
	arr := make([]interface{}, 0)
	isArray := true
	maxKey := 0

	tbl.ForEach(func(key, val lua.LValue) {
		if num, ok := key.(lua.LNumber); ok {
			keyNum := int(num)
			if keyNum > 0 {
				if keyNum > maxKey {
					maxKey = keyNum
				}
				// Check if keys are consecutive
				if keyNum <= len(arr)+1 {
					for len(arr) < keyNum {
						arr = append(arr, nil)
					}
					arr[keyNum-1] = FromLua(L, val)
				} else {
					isArray = false
				}
			} else {
				isArray = false
			}
		} else {
			isArray = false
		}
	})

	// If it's a proper array (consecutive keys 1..n), return as array
	if isArray && maxKey == len(arr) {
		return arr
	}

	// Otherwise, treat as map
	result := make(map[string]interface{})
	tbl.ForEach(func(key, val lua.LValue) {
		keyStr := key.String()
		result[keyStr] = FromLua(L, val)
	})
	return result
}
