# RawDefinition Feature

## Overview

The `RawDefinition` field in the `Type` struct preserves the original protodef JSON snippet that was used to create the type. This is useful for debugging to compare the original definition with the parsed structure.

## Purpose

When working with complex protocol definitions, it can be helpful to:
- See the exact original JSON that defined a type
- Compare the original definition with the parsed Type structure
- Debug parsing issues by tracing back to the source definition
- Generate human-readable debug output that shows both representations

## Usage

Every `Type` instance now includes a `RawDefinition` string field that contains the original JSON snippet.

### Example 1: Simple Type Alias

```go
jsonStr := `"varint"`
result := gjson.Parse(jsonStr)
typ := datatypes.GetTypeFromJSON("PlayerCount", result)

fmt.Println(typ.RawDefinition) // Output: "varint"
```

### Example 2: Container Type

```go
jsonStr := `["container", [
    {"name": "x", "type": "i32"},
    {"name": "y", "type": "i32"},
    {"name": "z", "type": "i32"}
]]`
result := gjson.Parse(jsonStr)
typ := datatypes.GetTypeFromJSON("Position", result)

// The top-level type has the full container definition
fmt.Println(typ.RawDefinition) 
// Output: ["container", [{"name": "x", "type": "i32"}, ...]]

// Individual fields also preserve their raw definitions
container := typ.Extras.(*datatypes.Container)
for _, field := range container.Fields {
    fmt.Printf("%s: %s\n", field.Name, field.Type.RawDefinition)
}
// Output:
// x: "i32"
// y: "i32"
// z: "i32"
```

### Example 3: Nested Containers

For nested types (like containers with anonymous fields), each level preserves its own raw definition:

```go
jsonStr := `["container", [
    {"name": "header", "type": "u32"},
    {"anon": true, "type": ["container", [
        {"name": "bitFlags", "type": "u8"},
        {"name": "metadata", "type": "u16"}
    ]]},
    {"name": "footer", "type": "u32"}
]]`
result := gjson.Parse(jsonStr)
typ := datatypes.GetTypeFromJSON("Packet", result)

// Top-level container
fmt.Println(typ.RawDefinition) // Full packet definition

// Anonymous nested container
container := typ.Extras.(*datatypes.Container)
anonField := container.Fields[1]
fmt.Println(anonField.Type.RawDefinition) 
// Output: ["container", [{"name": "bitFlags", "type": "u8"}, ...]]
```

## Implementation Details

### Where RawDefinition is Captured

The `RawDefinition` field is populated in these locations:

1. **`GetTypeFromJSON()`** - When creating a type from JSON with a name
2. **`GetType()`** - When parsing type definitions directly
3. **Nested types** - Container fields, array elements, switch cases, etc. all preserve their raw definitions

### What's Stored

The `RawDefinition` field stores the exact JSON string from the `gjson.Result.Raw` field, which includes:
- Complete JSON structure (arrays, objects, strings)
- Whitespace as it appears in the original
- Nested structures for complex types

### Type Cloning

When types are cloned (e.g., in `Container.Clone()`), the `RawDefinition` is automatically preserved through struct copying.

## Use Cases

### 1. Debugging Type Parsing

When a type doesn't parse as expected, you can inspect the raw definition to see what the parser received:

```go
if typ.TypeName != expectedType {
    fmt.Printf("Expected %s but got %s from definition: %s\n", 
        expectedType, typ.TypeName, typ.RawDefinition)
}
```

### 2. Error Messages

Include the original definition in error messages for better context:

```go
if err := validateType(typ); err != nil {
    return fmt.Errorf("invalid type definition %s: %w", typ.RawDefinition, err)
}
```

### 3. Documentation Generation

Generate documentation showing both the original protodef and the parsed structure:

```go
fmt.Printf("Type: %s\n", typ.Name)
fmt.Printf("Definition: %s\n", typ.RawDefinition)
fmt.Printf("Parsed as: %s\n", typ.TypeName)
```

### 4. Testing and Validation

Verify that types round-trip correctly by comparing the parsed structure back to the original:

```go
func TestTypeRoundtrip(t *testing.T) {
    original := `["container", [{"name": "x", "type": "i32"}]]`
    typ := parseType(original)
    
    // Verify the definition was preserved
    assert.Equal(t, original, typ.RawDefinition)
}
```

## Performance Considerations

- The `RawDefinition` field stores the complete JSON string, which may increase memory usage for large protocol definitions
- For most use cases, the memory overhead is negligible compared to the debugging benefits
- If memory is a critical concern, you could choose to omit `RawDefinition` in production builds using build tags

## Example Program

See `examples/rawdefinition_example.go` for a complete working example demonstrating various use cases.

## See Also

- [Type System Documentation](../ProtoDef/doc/datatypes.md)
- [Container Type](../datatypes/container.go)
- [Testing Examples](../datatypes/rawdefinition_test.go)
