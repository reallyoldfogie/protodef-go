package protocol

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestProtocolReadJSON(t *testing.T) {
	jsonStr := `{
		"types": {
			"packet": "i32",
			"string": ["pstring", {"countType": "varint"}],
			"varint": "native"
		}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if protocol.Types == nil {
		t.Fatal("protocol.Types is nil")
	}

	if len(protocol.Types) != 3 {
		t.Fatalf("len(protocol.Types) = %d, want 3", len(protocol.Types))
	}
}

func TestProtocolReadJSONWithNamespace(t *testing.T) {
	jsonStr := `{
		"types": {
			"base_type": "u8"
		},
		"handshake": {
			"types": {
				"server_address": "string",
				"server_port": "u16"
			}
		}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(protocol.Types) != 1 {
		t.Errorf("len(protocol.Types) = %d, want 1", len(protocol.Types))
	}

	if protocol.Namespaces == nil {
		t.Fatal("protocol.Namespaces is nil")
	}

	if len(protocol.Namespaces) != 1 {
		t.Errorf("len(protocol.Namespaces) = %d, want 1", len(protocol.Namespaces))
	}

	handshake, ok := protocol.Namespaces["handshake"]
	if !ok {
		t.Fatal("handshake namespace not found")
	}

	if handshake.Name != "handshake" {
		t.Errorf("handshake.Name = %q, want %q", handshake.Name, "handshake")
	}

	if len(handshake.Types) != 2 {
		t.Errorf("len(handshake.Types) = %d, want 2", len(handshake.Types))
	}
}

func TestProtocolReadJSONMultipleNamespaces(t *testing.T) {
	jsonStr := `{
		"types": {
			"packet_id": "varint"
		},
		"handshaking": {
			"types": {
				"handshake": ["container", [{"name": "protocol_version", "type": "varint"}]]
			}
		},
		"status": {
			"types": {
				"status_request": "void"
			}
		},
		"login": {
			"types": {
				"login_start": ["container", [{"name": "name", "type": "string"}]]
			}
		}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(protocol.Namespaces) != 3 {
		t.Fatalf("len(protocol.Namespaces) = %d, want 3", len(protocol.Namespaces))
	}

	expectedNamespaces := []string{"handshaking", "status", "login"}
	for _, name := range expectedNamespaces {
		if _, ok := protocol.Namespaces[name]; !ok {
			t.Errorf("namespace %q not found", name)
		}
	}
}

func TestProtocolReadJSONMissingTypes(t *testing.T) {
	jsonStr := `{
		"handshake": {
			"types": {
				"test": "i32"
			}
		}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err == nil {
		t.Fatal("ReadJSON() should return error when types is missing")
	}

	expectedError := "protocol types is missing"
	if err.Error() != expectedError {
		t.Errorf("error = %q, want %q", err.Error(), expectedError)
	}
}

func TestProtocolReadJSONTypesNotObject(t *testing.T) {
	jsonStr := `{
		"types": ["not", "an", "object"]
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err == nil {
		t.Fatal("ReadJSON() should return error when types is not an object")
	}

	expectedError := "protocol type is not object"
	if err.Error() != expectedError {
		t.Errorf("error = %q, want %q", err.Error(), expectedError)
	}
}

func TestProtocolReadJSONTypesEmpty(t *testing.T) {
	jsonStr := `{
		"types": {}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(protocol.Types) != 0 {
		t.Errorf("len(protocol.Types) = %d, want 0", len(protocol.Types))
	}

	if protocol.Namespaces == nil {
		t.Error("protocol.Namespaces should be initialized")
	}
}

func TestProtocolReadJSONComplexTypes(t *testing.T) {
	jsonStr := `{
		"types": {
			"entity_metadata": ["array", {
				"countType": "varint",
				"type": ["container", [
					{"name": "key", "type": "u8"},
					{"name": "value", "type": ["switch", {
						"compareTo": "key",
						"fields": {
							"0": "byte",
							"1": "varint"
						}
					}]}
				]]
			}]
		}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(protocol.Types) != 1 {
		t.Fatalf("len(protocol.Types) = %d, want 1", len(protocol.Types))
	}

	if protocol.Types[0].Name != "entity_metadata" {
		t.Errorf("protocol.Types[0].Name = %q, want %q", protocol.Types[0].Name, "entity_metadata")
	}
}

func TestProtocolReadJSONOnlyTypes(t *testing.T) {
	jsonStr := `{
		"types": {
			"u8": "native",
			"i32": "native",
			"string": ["pstring", {"countType": "varint"}]
		}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(protocol.Types) != 3 {
		t.Errorf("len(protocol.Types) = %d, want 3", len(protocol.Types))
	}

	if len(protocol.Namespaces) != 0 {
		t.Errorf("len(protocol.Namespaces) = %d, want 0", len(protocol.Namespaces))
	}
}

func TestProtocolReadJSONNestedNamespaces(t *testing.T) {
	jsonStr := `{
		"types": {
			"root": "u8"
		},
		"game": {
			"types": {
				"game_type": "u8"
			},
			"player": {
				"types": {
					"player_info": "string"
				}
			}
		}
	}`

	protocol := &Protocol{}
	result := gjson.Parse(jsonStr)
	err := protocol.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	game, ok := protocol.Namespaces["game"]
	if !ok {
		t.Fatal("game namespace not found")
	}

	player, ok := game.Namespaces["player"]
	if !ok {
		t.Fatal("player namespace not found in game")
	}

	if len(player.Types) != 1 {
		t.Errorf("len(player.Types) = %d, want 1", len(player.Types))
	}
}
