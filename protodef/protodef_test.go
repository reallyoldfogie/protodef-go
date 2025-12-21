package protodef

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadProtocolFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_protocol.json")

	jsonContent := `{
		"types": {
			"packet_id": "varint",
			"string": ["pstring", {"countType": "varint"}]
		},
		"handshake": {
			"types": {
				"handshake_packet": ["container", [
					{"name": "protocol_version", "type": "varint"},
					{"name": "server_address", "type": "string"}
				]]
			}
		}
	}`

	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)
	if err != nil {
		t.Fatalf("ReadProtocolFile() error = %v", err)
	}

	if protocol == nil {
		t.Fatal("protocol is nil")
	}

	if protocol.Types == nil {
		t.Fatal("protocol.Types is nil")
	}

	if len(protocol.Types) != 2 {
		t.Errorf("len(protocol.Types) = %d, want 2", len(protocol.Types))
	}

	if len(protocol.Namespaces) != 1 {
		t.Errorf("len(protocol.Namespaces) = %d, want 1", len(protocol.Namespaces))
	}

	if _, ok := protocol.Namespaces["handshake"]; !ok {
		t.Error("handshake namespace not found")
	}
}

func TestReadProtocolFileMinimal(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "minimal.json")

	jsonContent := `{
		"types": {
			"u8": "native"
		}
	}`

	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)
	if err != nil {
		t.Fatalf("ReadProtocolFile() error = %v", err)
	}

	if protocol == nil {
		t.Fatal("protocol is nil")
	}

	if len(protocol.Types) != 1 {
		t.Errorf("len(protocol.Types) = %d, want 1", len(protocol.Types))
	}
}

func TestReadProtocolFileNotFound(t *testing.T) {
	protocol, err := ReadProtocolFile("/nonexistent/file.json")

	if err == nil {
		t.Fatal("ReadProtocolFile() should return error for nonexistent file")
	}

	if protocol != nil {
		t.Error("protocol should be nil on error")
	}
}

func TestReadProtocolFileInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid.json")

	// Write invalid JSON (not an object)
	invalidJSON := `["this", "is", "an", "array"]`

	err := os.WriteFile(testFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)

	if err == nil {
		t.Fatal("ReadProtocolFile() should return error for invalid JSON structure")
	}

	expectedError := "protocol file is obviously wrong"
	if err.Error() != expectedError {
		t.Errorf("error = %q, want %q", err.Error(), expectedError)
	}

	if protocol != nil {
		t.Error("protocol should be nil on error")
	}
}

func TestReadProtocolFileMissingTypes(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "no_types.json")

	jsonContent := `{
		"handshake": {
			"types": {
				"test": "u8"
			}
		}
	}`

	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)

	if err == nil {
		t.Fatal("ReadProtocolFile() should return error when types field is missing")
	}

	if protocol != nil {
		t.Error("protocol should be nil on error")
	}
}

func TestReadProtocolFileEmptyTypes(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty_types.json")

	jsonContent := `{
		"types": {}
	}`

	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)
	if err != nil {
		t.Fatalf("ReadProtocolFile() error = %v", err)
	}

	if protocol == nil {
		t.Fatal("protocol is nil")
	}

	if len(protocol.Types) != 0 {
		t.Errorf("len(protocol.Types) = %d, want 0", len(protocol.Types))
	}
}

func TestReadProtocolFileComplexStructure(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "complex.json")

	jsonContent := `{
		"types": {
			"varint": "native",
			"string": ["pstring", {"countType": "varint"}],
			"position": ["container", [
				{"name": "x", "type": "i32"},
				{"name": "y", "type": "i32"},
				{"name": "z", "type": "i32"}
			]],
			"entity_metadata": ["array", {
				"countType": "varint",
				"type": ["switch", {
					"compareTo": "type",
					"fields": {
						"0": "byte",
						"1": "varint",
						"2": "f32"
					},
					"default": "void"
				}]
			}]
		},
		"handshaking": {
			"types": {
				"handshake": ["container", [{"name": "version", "type": "varint"}]]
			}
		},
		"play": {
			"types": {
				"spawn_entity": ["container", [{"name": "entity_id", "type": "varint"}]]
			},
			"client": {
				"types": {
					"player_position": ["container", [{"name": "pos", "type": "position"}]]
				}
			}
		}
	}`

	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)
	if err != nil {
		t.Fatalf("ReadProtocolFile() error = %v", err)
	}

	if protocol == nil {
		t.Fatal("protocol is nil")
	}

	if len(protocol.Types) != 4 {
		t.Errorf("len(protocol.Types) = %d, want 4", len(protocol.Types))
	}

	if len(protocol.Namespaces) != 2 {
		t.Errorf("len(protocol.Namespaces) = %d, want 2", len(protocol.Namespaces))
	}

	play, ok := protocol.Namespaces["play"]
	if !ok {
		t.Fatal("play namespace not found")
	}

	client, ok := play.Namespaces["client"]
	if !ok {
		t.Fatal("client namespace not found in play")
	}

	if len(client.Types) != 1 {
		t.Errorf("len(client.Types) = %d, want 1", len(client.Types))
	}
}

func TestReadProtocolFileMalformedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "malformed.json")

	// Write malformed JSON that gjson will parse but won't have types field
	malformedJSON := `{"notypes": "value"}`

	err := os.WriteFile(testFile, []byte(malformedJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)

	// Should error because types field is missing
	if err == nil {
		t.Fatal("ReadProtocolFile() should return error for JSON without types")
	}

	if protocol != nil {
		t.Error("protocol should be nil on error")
	}
}

func TestReadProtocolFileEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.json")

	err := os.WriteFile(testFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	protocol, err := ReadProtocolFile(testFile)

	if err == nil {
		t.Fatal("ReadProtocolFile() should return error for empty file")
	}

	if protocol != nil {
		t.Error("protocol should be nil on error")
	}
}

func TestReadProtocolFilePermissionDenied(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "noperm.json")

	jsonContent := `{"types": {"test": "u8"}}`
	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Remove read permissions
	err = os.Chmod(testFile, 0000)
	if err != nil {
		t.Fatalf("Failed to change permissions: %v", err)
	}

	// Restore permissions after test
	defer os.Chmod(testFile, 0644)

	protocol, err := ReadProtocolFile(testFile)

	if err == nil {
		t.Fatal("ReadProtocolFile() should return error for file without read permissions")
	}

	if protocol != nil {
		t.Error("protocol should be nil on error")
	}
}
