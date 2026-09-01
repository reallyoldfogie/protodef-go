# protodef-go

A Go implementation of the [ProtoDef](https://github.com/ProtoDef-io/ProtoDef) protocol description system.

[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Overview

ProtoDef is a language-agnostic system for describing binary protocol structures. This Go implementation allows you to:
- Define protocol structures using JSON schemas
- Parse protocol definitions into type-safe Go structures
- Support for namespacing and type composition
- Compatible with the official ProtoDef specification

## Status

**Current Version:** 0.9.0-alpha  
**Spec Compliance:** ~85%

### Implementation Progress

#### ✅ Type System (100%)
All ProtoDef specification types are implemented and parseable:

**Numeric Types** (100%)
- ✅ Basic: `i8`, `u8`, `i16`, `u16`, `i32`, `u32`, `i64`, `u64`, `f32`, `f64`
- ✅ Little-endian: `li8`, `lu8`, `li16`, `lu16`, `li32`, `lu32`, `li64`, `lu64`, `lf32`, `lf64`
- ✅ Variable-length: `varint`
- ✅ Configurable size: `int`, `lint`

**Primitives** (100%)
- ✅ `bool`, `cstring`, `void`

**Structures** (100%)
- ✅ `array` - Lists with count prefix
- ✅ `container` - Named field structures (with anonymous field support)
- ✅ `count` - Count fields for arrays/buffers

**Conditional** (100%)
- ✅ `switch` - Conditional type selection
- ✅ `option` - Optional values

**Utils** (100%)
- ✅ `buffer` - Raw byte buffers
- ✅ `bitfield` - Bit-level fields
- ✅ `mapper` - Value-to-string mappings
- ✅ `pstring` - Length-prefixed strings

**Extensions** (Non-Standard)
- ✅ `varint64`, `varint128` - Extended variable-length integers
- ✅ `zigzag32`, `zigzag64` - ZigZag-encoded signed integers
- ✅ `bitflags` - Named boolean flags packed into an integer
- ✅ `registryEntryHolder`, `registryEntryHolderSet` - Registry-conditional values
- ✅ `entityMetadataLoop` - Loop that reads until a terminator value
- ✅ `topBitSetTerminatedArray` (alias `topbitsetalternative`) - Slot array terminated by a high-bit-set byte

These extensions are not part of the core ProtoDef spec; the last five are borrowed from Minecraft-protocol/`minecraft-data` definitions. See [Extensions](#extensions).

#### ⚠️ In Progress
- 🔨 Runtime encoding/decoding (planned)
- 🔨 Comprehensive test coverage (expanding)
- 🔨 Examples and tutorials (in progress)

## Installation

```bash
go get github.com/protodef-go/protodef-go
```

## Quick Start

```go
import (
    "github.com/protodef-go/protodef-go/protodef"
)

// Parse a protocol definition file
proto, err := protodef.ReadProtocolFile("protocol.json")
if err != nil {
    panic(err)
}

// proto.Types holds every top-level type definition (*datatypes.Type)
for _, t := range proto.Types {
    fmt.Println(t.Name, t.TypeName)
}

// proto.Namespaces holds any nested namespace sections, keyed by name
handshake := proto.Namespaces["handshake"]
```

## Documentation

### Core Documentation
- [Implementation Fixes Plan](docs/implementation-fixes.md) - Detailed implementation roadmap
- [Phase Summaries](docs/) - Phase 1-4 implementation summaries
- [Extensions Guide](docs/extensions.md) - Non-standard numeric type extensions
- [TypeExtras Implementation](docs/TYPEEXTRAS_IMPLEMENTATION.md) - How per-type argument parsing (`TypeExtras`) works
- [RawDefinition Feature](docs/rawdefinition.md) - Debugging with original type definitions
- [Implementation Verification](PROTODEF_IMPLEMENTATION_VERIFICATION.md) - Coverage check against the ProtoDef spec
- [TODO](TODO.md) - Outstanding work tracked against the upstream spec

### ProtoDef Specification
- [Official Spec](ProtoDef/README.md)
- [Datatypes Reference](ProtoDef/doc/datatypes.md)
- [Protocol Format](ProtoDef/doc/protocol.md)

## Features

### ✅ Completed
- Full type system implementation
- JSON protocol definition parsing
- Namespace support
- Type composition and references
- Little-endian type support
- Anonymous container fields
- Configurable integer sizes (`int`, `lint`)
- Value mapping (`mapper` type)
- Bit-level flags (`bitfield`, `bitflags`)
- Registry-conditional types (`registryEntryHolder`, `registryEntryHolderSet`)
- Terminator-based loops (`entityMetadataLoop`, `topBitSetTerminatedArray`)
- JSON Schema validation of protocol definitions (`protocol.ValidateJSONSchema`, with optional `ajv-cli` support for advanced regex)
- RawDefinition preservation for debugging

### 🔜 Planned
- Runtime serialization/deserialization
- Code generation from schemas
- Performance optimizations
- More comprehensive examples

## Project Structure

```
protodef-go/
├── datatypes/               # Type definitions and parsing
│   ├── primitives.go        # Basic types (bool, void, cstring)
│   ├── numbers.go           # Numeric types (all variants)
│   ├── intExtras.go         # Configurable-size int/lint arguments
│   ├── container.go         # Container structures (incl. anonymous fields)
│   ├── array.go             # Array type
│   ├── count.go             # Count type
│   ├── switch.go            # Switch conditional
│   ├── option.go            # Option type
│   ├── conditional.go       # Conditional-type shared support
│   ├── buffer.go            # Buffer type
│   ├── bitfield.go          # Bitfield type
│   ├── bitflags.go          # Bitflags type
│   ├── mapper.go            # Mapper type
│   ├── pstring.go           # Prefixed string type
│   ├── registryEntryHolder.go     # Registry-conditional single value
│   ├── registryEntryHolderSet.go  # Registry-conditional value set
│   ├── entitymetadataloop.go      # Terminator-based read loop
│   ├── topbitsetterminatedarray.go # High-bit-terminated slot array
│   ├── type.go              # Core type system / JSON dispatch
│   ├── typeExtras.go        # Per-type argument-parsing interface
│   └── types.go             # Native type registry
├── protocol/                 # Protocol parsing and JSON Schema validation
├── namespace/                # Namespace management
├── protodef/                 # Top-level file-loading API (ReadProtocolFile)
├── examples/                 # Example programs (e.g. rawdefinition_example.go)
├── docs/                     # Documentation
│   ├── implementation-fixes.md
│   ├── extensions.md
│   ├── rawdefinition.md
│   ├── TYPEEXTRAS_IMPLEMENTATION.md
│   ├── phase1-summary.md
│   ├── phase2-summary.md
│   ├── phase3-summary.md
│   └── phase4-summary.md
└── ProtoDef/                 # Official spec (git submodule)
```

## Examples

### Defining a Protocol

```json
{
  "types": {
    "packet": [
      "container",
      [
        {"name": "id", "type": "varint"},
        {"name": "data", "type": "buffer", "countType": "varint"}
      ]
    ],
    "position": [
      "container",
      [
        {"name": "x", "type": "i32"},
        {"name": "y", "type": "i32"},
        {"name": "z", "type": "i32"}
      ]
    ]
  }
}
```

### Using Little-Endian Types

```json
{
  "types": {
    "file_header": [
      "container",
      [
        {"name": "magic", "type": "lu32"},
        {"name": "version", "type": "lu16"},
        {"name": "size", "type": "lu64"}
      ]
    ]
  }
}
```

### Using Mapper for Enums

```json
{
  "types": {
    "packet_type": [
      "mapper",
      {
        "type": "u8",
        "mappings": {
          "0": "handshake",
          "1": "data",
          "2": "disconnect"
        }
      }
    ]
  }
}
```

## Compatibility

### ProtoDef Implementations

This implementation aims for compatibility with:
- [node-protodef](https://github.com/ProtoDef-io/node-protodef) - Reference implementation
- [elixir-protodef](https://github.com/ProtoDef-io/elixir-protodef)
- [protodefc](https://github.com/ProtoDef-io/protodefc) - Rust compiler

### Extensions

Some types are non-standard extensions to the core ProtoDef spec:
- Numeric extensions `varint64`, `varint128`, `zigzag32`, `zigzag64` — documented in [docs/extensions.md](docs/extensions.md)
- Minecraft-protocol-derived types `bitflags`, `registryEntryHolder`, `registryEntryHolderSet`, `entityMetadataLoop`, `topBitSetTerminatedArray` (alias `topbitsetalternative`) — implemented for compatibility with `minecraft-data`/`node-minecraft-protocol` protocol definitions, not part of the official ProtoDef spec

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./datatypes
```

## Contributing

Contributions are welcome! Areas that need work:

1. **Runtime Encoding/Decoding** - Implement serialization for all types
2. **Test Coverage** - Add comprehensive tests for all types
3. **Examples** - Real-world protocol implementations
4. **Performance** - Optimization and benchmarking
5. **Documentation** - Tutorials and guides

See [docs/implementation-fixes.md](docs/implementation-fixes.md) for the detailed roadmap.

## Compliance Matrix

| Category | Types | Status |
|----------|-------|--------|
| Numeric (incl. `int`/`lint`) | 23 types | ✅ 100% |
| Primitives | 3 types | ✅ 100% |
| Structures | 3 types | ✅ 100% |
| Conditional | 2 types | ✅ 100% |
| Utils | 4 types | ✅ 100% |
| **ProtoDef spec subtotal** | **35 types** | **✅ 100%** |
| Non-standard extensions | 9 types | ✅ 100% |
| **Total** | **44 types** | **✅ 100%** |

## Projects Using ProtoDef

ProtoDef is used by various projects for protocol definitions:
- [minecraft-protocol](https://github.com/PrismarineJS/node-minecraft-protocol) - Minecraft protocol
- [prismarine-nbt](https://github.com/PrismarineJS/prismarine-nbt) - NBT format
- [node-raknet](https://github.com/mhsjlw/node-raknet) - RakNet protocol

## License

MIT - see [LICENSE](LICENSE)

## Acknowledgments

- [ProtoDef](https://github.com/ProtoDef-io/ProtoDef) - Original specification
- [PrismarineJS](https://github.com/PrismarineJS) - Protocol definitions and inspiration

## Links

- [ProtoDef Specification](https://github.com/ProtoDef-io/ProtoDef)
- [ProtoDef Documentation](https://github.com/ProtoDef-io/ProtoDef/tree/master/doc)
- [Issue Tracker](https://github.com/protodef-go/protodef-go/issues)
