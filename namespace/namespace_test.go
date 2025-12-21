package namespace

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestNamespaceReadJSON(t *testing.T) {
	jsonStr := `{
		"types": {
			"type1": "i32",
			"type2": "string",
			"type3": "varint"
		}
	}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if ns.Types == nil {
		t.Fatal("ns.Types is nil")
	}

	if len(ns.Types) != 3 {
		t.Fatalf("len(ns.Types) = %d, want 3", len(ns.Types))
	}
}

func TestNamespaceReadJSONWithTypesArray(t *testing.T) {
	jsonStr := `{
		"types": ["type1", "type2", "type3"]
	}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if ns.Types == nil {
		t.Fatal("ns.Types is nil")
	}

	if len(ns.Types) != 3 {
		t.Fatalf("len(ns.Types) = %d, want 3", len(ns.Types))
	}
}

func TestNamespaceReadJSONWithNestedNamespace(t *testing.T) {
	jsonStr := `{
		"types": {
			"base_type": "i32"
		},
		"sub_namespace": {
			"types": {
				"nested_type": "string"
			}
		}
	}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if ns.Types == nil {
		t.Fatal("ns.Types is nil")
	}

	if len(ns.Types) != 1 {
		t.Errorf("len(ns.Types) = %d, want 1", len(ns.Types))
	}

	if ns.Namespaces == nil {
		t.Fatal("ns.Namespaces is nil")
	}

	if len(ns.Namespaces) != 1 {
		t.Errorf("len(ns.Namespaces) = %d, want 1", len(ns.Namespaces))
	}

	subNs, ok := ns.Namespaces["sub_namespace"]
	if !ok {
		t.Fatal("sub_namespace not found")
	}

	if subNs.Name != "sub_namespace" {
		t.Errorf("subNs.Name = %q, want %q", subNs.Name, "sub_namespace")
	}

	if len(subNs.Types) != 1 {
		t.Errorf("len(subNs.Types) = %d, want 1", len(subNs.Types))
	}
}

func TestNamespaceReadJSONMultipleNamespaces(t *testing.T) {
	jsonStr := `{
		"types": {
			"main_type": "u8"
		},
		"namespace1": {
			"types": ["type_a", "type_b"]
		},
		"namespace2": {
			"types": {
				"type_c": "i64"
			}
		}
	}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(ns.Namespaces) != 2 {
		t.Fatalf("len(ns.Namespaces) = %d, want 2", len(ns.Namespaces))
	}

	if _, ok := ns.Namespaces["namespace1"]; !ok {
		t.Error("namespace1 not found")
	}

	if _, ok := ns.Namespaces["namespace2"]; !ok {
		t.Error("namespace2 not found")
	}
}

func TestNamespaceReadJSONEmpty(t *testing.T) {
	jsonStr := `{}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if ns.Namespaces == nil {
		t.Error("ns.Namespaces should be initialized")
	}
}

func TestNamespaceReadJSONOnlyTypes(t *testing.T) {
	jsonStr := `{
		"types": {
			"packet_id": "varint",
			"player_name": "string",
			"position": ["container", [{"name": "x", "type": "f64"}]]
		}
	}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(ns.Types) != 3 {
		t.Errorf("len(ns.Types) = %d, want 3", len(ns.Types))
	}

	if len(ns.Namespaces) != 0 {
		t.Errorf("len(ns.Namespaces) = %d, want 0", len(ns.Namespaces))
	}
}

func TestNamespaceReadJSONDeepNesting(t *testing.T) {
	jsonStr := `{
		"types": {"root_type": "i8"},
		"level1": {
			"types": {"level1_type": "i16"},
			"level2": {
				"types": {"level2_type": "i32"}
			}
		}
	}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	level1, ok := ns.Namespaces["level1"]
	if !ok {
		t.Fatal("level1 namespace not found")
	}

	if len(level1.Types) != 1 {
		t.Errorf("len(level1.Types) = %d, want 1", len(level1.Types))
	}

	level2, ok := level1.Namespaces["level2"]
	if !ok {
		t.Fatal("level2 namespace not found")
	}

	if len(level2.Types) != 1 {
		t.Errorf("len(level2.Types) = %d, want 1", len(level2.Types))
	}
}

func TestNamespaceReadJSONNoTypes(t *testing.T) {
	jsonStr := `{
		"namespace1": {
			"types": {"type1": "u8"}
		}
	}`

	ns := &Namespace{}
	result := gjson.Parse(jsonStr)
	err := ns.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(ns.Types) != 0 {
		t.Errorf("len(ns.Types) = %d, want 0", len(ns.Types))
	}

	if len(ns.Namespaces) != 1 {
		t.Errorf("len(ns.Namespaces) = %d, want 1", len(ns.Namespaces))
	}
}
