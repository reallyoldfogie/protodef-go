package main

import (
	"fmt"

	"github.com/reallyoldfogie/protodef-go/datatypes"
	"github.com/tidwall/gjson"
)

func main() {
	// Example 1: Simple type alias
	fmt.Println("=== Example 1: Simple Type Alias ===")
	simpleJSON := `"varint"`
	result := gjson.Parse(simpleJSON)
	typ := datatypes.GetTypeFromJSON("PlayerCount", result)
	
	fmt.Printf("Type Name: %s\n", typ.Name)
	fmt.Printf("Type TypeName: %s\n", typ.TypeName)
	fmt.Printf("Raw Definition: %s\n\n", typ.RawDefinition)

	// Example 2: Container with fields
	fmt.Println("=== Example 2: Container Type ===")
	containerJSON := `["container", [
		{"name": "x", "type": "i32"},
		{"name": "y", "type": "i32"},
		{"name": "z", "type": "i32"}
	]]`
	result = gjson.Parse(containerJSON)
	typ = datatypes.GetTypeFromJSON("Position", result)
	
	fmt.Printf("Type Name: %s\n", typ.Name)
	fmt.Printf("Type TypeName: %s\n", typ.TypeName)
	fmt.Printf("Raw Definition: %s\n", typ.RawDefinition)
	
	container, ok := typ.Extras.(*datatypes.Container)
	if ok {
		fmt.Println("\nFields:")
		for _, field := range container.Fields {
			fmt.Printf("  - %s: type=%s, raw=%s\n", 
				field.Name, 
				field.Type.TypeName,
				field.Type.RawDefinition)
		}
	}
	fmt.Println()

	// Example 3: Nested container with anonymous field
	fmt.Println("=== Example 3: Nested Container (with anon field) ===")
	nestedJSON := `["container", [
		{"name": "packetId", "type": "varint"},
		{"anon": true, "type": ["container", [
			{"name": "bitFlags", "type": "u8"},
			{"name": "metadata", "type": "u16"}
		]]},
		{"name": "checksum", "type": "u32"}
	]]`
	result = gjson.Parse(nestedJSON)
	typ = datatypes.GetTypeFromJSON("Packet", result)
	
	fmt.Printf("Type Name: %s\n", typ.Name)
	fmt.Printf("Type TypeName: %s\n", typ.TypeName)
	fmt.Printf("Raw Definition:\n%s\n\n", typ.RawDefinition)
	
	container, ok = typ.Extras.(*datatypes.Container)
	if ok {
		fmt.Println("Fields:")
		for i, field := range container.Fields {
			if field.Anon {
				fmt.Printf("  [%d] <anonymous field>\n", i)
				fmt.Printf("      Type: %s\n", field.Type.TypeName)
				fmt.Printf("      Raw Definition: %s\n", field.Type.RawDefinition)
				
				// Show nested container fields if available
				if nestedContainer, ok := field.Type.Extras.(*datatypes.Container); ok {
					fmt.Println("      Nested fields:")
					for _, nestedField := range nestedContainer.Fields {
						fmt.Printf("        - %s: %s\n", nestedField.Name, nestedField.Type.TypeName)
					}
				}
			} else {
				fmt.Printf("  [%d] %s: %s\n", i, field.Name, field.Type.TypeName)
				fmt.Printf("      Raw: %s\n", field.Type.RawDefinition)
			}
		}
	}
	fmt.Println()

	// Example 4: Array type
	fmt.Println("=== Example 4: Array Type ===")
	arrayJSON := `["array", {
		"countType": "varint",
		"type": "u64"
	}]`
	result = gjson.Parse(arrayJSON)
	typ = datatypes.GetTypeFromJSON("EntityList", result)
	
	fmt.Printf("Type Name: %s\n", typ.Name)
	fmt.Printf("Type TypeName: %s\n", typ.TypeName)
	fmt.Printf("Raw Definition: %s\n\n", typ.RawDefinition)

	fmt.Println("=== Use Case: Debugging ===")
	fmt.Println("When debugging protocol parsing issues, you can now:")
	fmt.Println("1. Compare the original protodef JSON with the parsed Type structure")
	fmt.Println("2. See exactly how nested types were defined")
	fmt.Println("3. Trace back from a Type instance to its original definition")
	fmt.Println("4. Generate human-readable debug output showing both representations")
}
