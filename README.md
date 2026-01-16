# go-lua-converter

A GopherLua utility package for converting between Go values and Lua values.

## Installation

```bash
go get github.com/kerimovok/go-lua-converter
```

## Usage

### Go to Lua Conversion

```go
import (
    lua "github.com/yuin/gopher-lua"
    converter "github.com/kerimovok/go-lua-converter"
)

L := lua.NewState()
defer L.Close()

// Convert Go value to Lua
goValue := map[string]interface{}{
    "name": "John",
    "age": 30,
    "tags": []interface{}{"developer", "golang"},
}
luaTable := converter.ToLua(L, goValue)

// Or use specific converters
luaMap := converter.MapToTable(L, map[string]interface{}{"key": "value"})
luaArray := converter.SliceToTable(L, []interface{}{1, 2, 3})
```

### Lua to Go Conversion

```go
// Get a Lua value from the stack or table
luaValue := L.Get(-1) // or from table

// Convert to Go value
goValue := converter.FromLua(L, luaValue)

// Convert Lua table to Go map or slice
luaTable := L.GetGlobal("mytable").(*lua.LTable)
goValue := converter.TableToGo(L, luaTable)
```

## Functions

### `ToLua(L *lua.LState, v interface{}) lua.LValue`

Converts a Go value to a Lua value.

- **Supported types:** `string`, `int`, `int32`, `int64`, `float32`, `float64`, `bool`, `map[string]interface{}`, `[]interface{}`, `nil`
- **Unknown types:** Converted to string representation

### `MapToTable(L *lua.LState, m map[string]interface{}) *lua.LTable`

Converts a Go map to a Lua table.

### `SliceToTable(L *lua.LState, arr []interface{}) *lua.LTable`

Converts a Go slice to a Lua table (1-indexed).

### `FromLua(L *lua.LState, lv lua.LValue) interface{}`

Converts a Lua value to a Go value.

- **Supported types:** `nil`, `bool`, `string`, `number`, `table`
- **Unknown types:** Converted to string

### `TableToGo(L *lua.LState, tbl *lua.LTable) interface{}`

Converts a Lua table to a Go value (map or slice).

- **Array detection:** If table has consecutive integer keys starting from 1, returns `[]interface{}`
- **Otherwise:** Returns `map[string]interface{}`

## Notes

- Lua arrays are 1-indexed, Go slices are 0-indexed - conversion handles this automatically
- Tables with mixed keys (both integer and string) are treated as maps
- Numbers are converted to `float64` when converting from Lua to Go
