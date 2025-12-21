package protocol

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"

	"github.com/pkg/errors"

	"github.com/protodef-go/protodef-go/datatypes"
	"github.com/protodef-go/protodef-go/namespace"
	"github.com/tidwall/gjson"
)

func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := bytes.Fields(buf[:n])[1]
	id, err := strconv.ParseInt(string(idField), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

type Protocol struct {
	Types      []*datatypes.Type
	Namespaces map[string]*namespace.Namespace
}

func (p *Protocol) ReadJSON(d gjson.Result) error {
	types := d.Get("types")

	if !types.Exists() {
		return errors.New("protocol types is missing")
	}

	if !types.IsObject() {
		return errors.New("protocol type is not object")
	}

	fmt.Println("DEBUG [Protocol.ReadJSON - START]: p.Types slice ptr=", &p.Types, "cap=", cap(p.Types), "len=", len(p.Types))

	typeCount := 0
	for name, option := range types.Map() {
		t := datatypes.GetTypeFromJSON(name, option)
		if t == nil {
			// DEBUG
			if name == "ContainerID" || name == "optvarint" {
				fmt.Printf("DEBUG [Protocol.ReadJSON %d]: GetTypeFromJSON returned nil for %s\n", getGoroutineID(), name)
			}
			continue
		}
		// DEBUG
		if name == "ContainerID" || name == "optvarint" {
			fmt.Printf("DEBUG [Protocol.ReadJSON %d]: Adding type to p.Types: %s TypeName= %s current len=%d\n", getGoroutineID(), name, t.TypeName, len(p.Types))
		}
		p.Types = append(p.Types, t)
		typeCount++
		if name == "ContainerID" || name == "optvarint" {
			fmt.Printf("DEBUG [Protocol.ReadJSON %d]: After append, p.Types len=%d\n", getGoroutineID(), len(p.Types))
		}

		fmt.Printf("DEBUG [Protocol.ReadJSON %d]: t=%#v (%p)\n", getGoroutineID(), t, t)
		fmt.Printf("DEBUG [Protocol.ReadJSON %d]: After append => %#v\n", getGoroutineID(), p.Types)
	}
	println("DEBUG [Protocol.ReadJSON]: Added", typeCount, "types from top-level .types section. Total p.Types length=", len(p.Types))
	// DEBUG: Print all type names
	println("DEBUG: [Protocol.ReadJSON]All types in p.Types :")
	for i, t := range p.Types {
		fmt.Printf("\tDEBUG: [Protocol.ReadJSON %d] [%d] %#v (%p)\n", getGoroutineID(), i, t, t)
	}

	// DEBUG: Check if ContainerID exists BEFORE namespace processing
	hasContainerIDBefore := false
	for _, t := range p.Types {
		if t.Name == "ContainerID" {
			hasContainerIDBefore = true
			println("DEBUG [BEFORE namespace loop]: ContainerID exists at index")
			break
		}
	}
	if !hasContainerIDBefore {
		println("DEBUG [BEFORE namespace loop]: ContainerID DOES NOT exist")
	}

	p.Namespaces = make(map[string]*namespace.Namespace)
	for name, data := range d.Map() {
		if name == "types" {
			continue
		}

		namespace := &namespace.Namespace{Name: name}
		err := namespace.ReadJSON(data)
		if err != nil {
			return err
		}

		p.Namespaces[name] = namespace
	}

	println("DEBUG [Protocol.ReadJSON - END]: Final p.Types length=", len(p.Types))
	// DEBUG: Check if ContainerID is still there
	hasContainerID := false
	for _, t := range p.Types {
		if t.Name == "ContainerID" {
			hasContainerID = true
			break
		}
	}
	if hasContainerID {
		println("DEBUG [Protocol.ReadJSON - END]: ContainerID IS in p.Types")
	} else {
		println("DEBUG [Protocol.ReadJSON - END]: ContainerID IS NOT in p.Types")
	}
	return nil
}
