package converter

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestToLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		input    interface{}
		checkFn  func(*lua.LState, lua.LValue) bool
		expected bool
	}{
		{
			name:  "nil",
			input: nil,
			checkFn: func(L *lua.LState, v lua.LValue) bool {
				return v == lua.LNil
			},
			expected: true,
		},
		{
			name:  "string",
			input: "hello",
			checkFn: func(L *lua.LState, v lua.LValue) bool {
				return v.String() == "hello"
			},
			expected: true,
		},
		{
			name:  "int",
			input: 42,
			checkFn: func(L *lua.LState, v lua.LValue) bool {
				return v.(lua.LNumber) == 42
			},
			expected: true,
		},
		{
			name:  "float64",
			input: 3.14,
			checkFn: func(L *lua.LState, v lua.LValue) bool {
				return v.(lua.LNumber) == 3.14
			},
			expected: true,
		},
		{
			name:  "bool true",
			input: true,
			checkFn: func(L *lua.LState, v lua.LValue) bool {
				return v.(lua.LBool) == lua.LTrue
			},
			expected: true,
		},
		{
			name:  "bool false",
			input: false,
			checkFn: func(L *lua.LState, v lua.LValue) bool {
				return v.(lua.LBool) == lua.LFalse
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLua(L, tt.input)
			if !tt.checkFn(L, result) {
				t.Errorf("ToLua() = %v, expected different value", result)
			}
		})
	}
}

func TestMapToTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	m := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"active": true,
	}

	table := MapToTable(L, m)

	// Check values
	if table.RawGetString("name").String() != "John" {
		t.Errorf("Expected name='John', got %v", table.RawGetString("name"))
	}
	if table.RawGetString("age").(lua.LNumber) != 30 {
		t.Errorf("Expected age=30, got %v", table.RawGetString("age"))
	}
	if table.RawGetString("active").(lua.LBool) != lua.LTrue {
		t.Errorf("Expected active=true, got %v", table.RawGetString("active"))
	}
}

func TestSliceToTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	arr := []interface{}{"a", "b", "c"}

	table := SliceToTable(L, arr)

	// Lua arrays are 1-indexed
	if table.RawGetInt(1).String() != "a" {
		t.Errorf("Expected [1]='a', got %v", table.RawGetInt(1))
	}
	if table.RawGetInt(2).String() != "b" {
		t.Errorf("Expected [2]='b', got %v", table.RawGetInt(2))
	}
	if table.RawGetInt(3).String() != "c" {
		t.Errorf("Expected [3]='c', got %v", table.RawGetInt(3))
	}
}

func TestFromLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		setup    func(*lua.LState) lua.LValue
		checkFn  func(interface{}) bool
		expected bool
	}{
		{
			name: "nil",
			setup: func(L *lua.LState) lua.LValue {
				return lua.LNil
			},
			checkFn: func(v interface{}) bool {
				return v == nil
			},
			expected: true,
		},
		{
			name: "string",
			setup: func(L *lua.LState) lua.LValue {
				return lua.LString("hello")
			},
			checkFn: func(v interface{}) bool {
				return v.(string) == "hello"
			},
			expected: true,
		},
		{
			name: "number",
			setup: func(L *lua.LState) lua.LValue {
				return lua.LNumber(42)
			},
			checkFn: func(v interface{}) bool {
				return v.(float64) == 42.0
			},
			expected: true,
		},
		{
			name: "bool",
			setup: func(L *lua.LState) lua.LValue {
				return lua.LBool(true)
			},
			checkFn: func(v interface{}) bool {
				return v.(bool) == true
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv := tt.setup(L)
			result := FromLua(L, lv)
			if !tt.checkFn(result) {
				t.Errorf("FromLua() = %v, expected different value", result)
			}
		})
	}
}

func TestTableToGo_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a Lua array table
	table := L.NewTable()
	table.RawSetInt(1, lua.LString("a"))
	table.RawSetInt(2, lua.LString("b"))
	table.RawSetInt(3, lua.LString("c"))

	result := TableToGo(L, table)

	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected array, got %T", result)
	}

	if len(arr) != 3 {
		t.Errorf("Expected length 3, got %d", len(arr))
	}

	if arr[0].(string) != "a" || arr[1].(string) != "b" || arr[2].(string) != "c" {
		t.Errorf("Expected [a, b, c], got %v", arr)
	}
}

func TestTableToGo_Map(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a Lua map table
	table := L.NewTable()
	table.RawSetString("name", lua.LString("John"))
	table.RawSetString("age", lua.LNumber(30))

	result := TableToGo(L, table)

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", result)
	}

	if m["name"].(string) != "John" {
		t.Errorf("Expected name='John', got %v", m["name"])
	}

	if m["age"].(float64) != 30.0 {
		t.Errorf("Expected age=30, got %v", m["age"])
	}
}

func TestRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test round trip: Go -> Lua -> Go
	original := map[string]interface{}{
		"name": "John",
		"age":  30,
		"tags": []interface{}{"developer", "golang"},
	}

	luaTable := ToLua(L, original)
	goValue := FromLua(L, luaTable)

	result, ok := goValue.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", goValue)
	}

	if result["name"].(string) != "John" {
		t.Errorf("Expected name='John', got %v", result["name"])
	}

	if result["age"].(float64) != 30.0 {
		t.Errorf("Expected age=30, got %v", result["age"])
	}

	tags, ok := result["tags"].([]interface{})
	if !ok {
		t.Fatalf("Expected tags array, got %T", result["tags"])
	}

	if len(tags) != 2 || tags[0].(string) != "developer" || tags[1].(string) != "golang" {
		t.Errorf("Expected tags=[developer, golang], got %v", tags)
	}
}
