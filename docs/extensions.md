# ProtoDef-Go Extensions

**Status:** Documented  
**Version:** 1.0

## Overview

This document describes the non-standard type extensions implemented in protodef-go that are not part of the official ProtoDef specification. These extensions provide additional functionality for specific use cases but may not be compatible with other ProtoDef implementations.

## Purpose

Extensions are provided to support:
- Advanced variable-length integer encoding schemes
- Optimized integer representations for specific protocols
- Protocol-specific requirements not covered by the base spec

## ⚠️ Compatibility Warning

**Important:** These extensions are **not part of the official ProtoDef specification** and will not be recognized by other ProtoDef implementations (Node.js, Elixir, Rust, etc.).

If you need cross-platform protocol compatibility, use only the standard ProtoDef types documented in the [official specification](../ProtoDef/doc/datatypes.md).

---

## Extension Types

### 1. `varint64` - 64-bit Variable-Length Integer

**Category:** Numeric (Non-Standard)  
**Size:** Variable (1-10 bytes)

#### Description
Extended varint encoding that supports full 64-bit integer range. Based on Protocol Buffers varint encoding but extended to handle 64-bit values.

#### Encoding
- Uses continuation bit scheme (MSB of each byte)
- Little-endian byte order
- Can represent values from 0 to 2^64-1
- More compact than fixed 8-byte representation for small values

#### Use Cases
- Protocols requiring large integers with space optimization
- Database record IDs, timestamps with high precision
- File offsets in large files
- Network protocols with 64-bit counters

#### Example
```json
{
  "recordId": "varint64"
}
```

**Value Examples:**
- `127` → 1 byte
- `16383` → 2 bytes
- `2^63` → 10 bytes

#### Trade-offs
- **Pros:** Space-efficient for small values, supports full 64-bit range
- **Cons:** Variable size complicates random access, non-standard encoding

---

### 2. `varint128` - 128-bit Variable-Length Integer

**Category:** Numeric (Non-Standard)  
**Size:** Variable (1-19 bytes)

#### Description
Extremely large variable-length integer encoding supporting up to 128-bit values. Useful for cryptographic identifiers, UUIDs stored as integers, or scientific computing.

#### Encoding
- Extended varint encoding with continuation bits
- Supports values from 0 to 2^128-1
- Little-endian byte order
- Up to 19 bytes for maximum value

#### Use Cases
- IPv6 addresses as integers (128-bit)
- Cryptographic key identifiers
- UUID representation
- Large-scale distributed system IDs
- Scientific computing with very large numbers

#### Example
```json
{
  "sessionId": "varint128"
}
```

**Practical Values:**
- Small values: 1-2 bytes
- UUID-sized values (~128 bits): ~18-19 bytes
- Maximum theoretical: 19 bytes

#### Trade-offs
- **Pros:** Supports extremely large values, still compact for small values
- **Cons:** Maximum size larger than fixed 16-byte representation, complex encoding

---

### 3. `zigzag32` - ZigZag-Encoded 32-bit Integer

**Category:** Numeric (Non-Standard)  
**Size:** Variable (1-5 bytes)

#### Description
ZigZag encoding combined with varint for efficient representation of signed integers, especially those that frequently oscillate around zero.

#### Encoding
- Maps signed integers to unsigned: `(n << 1) ^ (n >> 31)`
- Then applies varint encoding
- Efficiently encodes both positive and negative values near zero
- Based on Protocol Buffers sint32

#### Mapping Examples
```
 0 → 0
-1 → 1
 1 → 2
-2 → 3
 2 → 4
```

#### Use Cases
- Signed integers that vary around zero (deltas, differences)
- Coordinate systems (relative positions)
- Financial data (profits/losses)
- Temperature readings (positive/negative)

#### Example
```json
{
  "deltaX": "zigzag32",
  "deltaY": "zigzag32",
  "profit": "zigzag32"
}
```

**Encoding Efficiency:**
- Values -64 to 63: 1 byte
- Values -8192 to 8191: 2 bytes
- Much more efficient than storing negative numbers in standard varint

#### Trade-offs
- **Pros:** Very efficient for small signed integers, handles negatives well
- **Cons:** Requires encoding/decoding step, non-standard format

---

### 4. `zigzag64` - ZigZag-Encoded 64-bit Integer

**Category:** Numeric (Non-Standard)  
**Size:** Variable (1-10 bytes)

#### Description
64-bit version of ZigZag encoding. Efficiently represents signed 64-bit integers, particularly those close to zero.

#### Encoding
- Maps signed integers to unsigned: `(n << 1) ^ (n >> 63)`
- Then applies varint encoding
- Supports full 64-bit signed range: -2^63 to 2^63-1
- Based on Protocol Buffers sint64

#### Use Cases
- Large signed deltas (timestamps, offsets)
- 64-bit coordinate systems
- Financial calculations requiring high precision
- Scientific data with positive/negative large values

#### Example
```json
{
  "timeDelta": "zigzag64",
  "accountBalance": "zigzag64",
  "altitude": "zigzag64"
}
```

**Encoding Efficiency:**
- Small values (-64 to 63): 1 byte
- Medium values: 2-5 bytes
- Large values near max: up to 10 bytes

#### Trade-offs
- **Pros:** Optimal for 64-bit signed integers, efficient near zero
- **Cons:** Encoding overhead, non-standard, larger max size than zigzag32

---

## Comparison Matrix

| Type | Signed | Size Range | Best Use Case | Standard |
|------|--------|------------|---------------|----------|
| `varint` | No | 1-5 bytes | Unsigned 32-bit | ✅ Yes |
| `varint64` | No | 1-10 bytes | Unsigned 64-bit | ❌ Extension |
| `varint128` | No | 1-19 bytes | 128-bit values | ❌ Extension |
| `zigzag32` | Yes | 1-5 bytes | Signed 32-bit near zero | ❌ Extension |
| `zigzag64` | Yes | 1-10 bytes | Signed 64-bit near zero | ❌ Extension |

---

## Encoding Details

### Varint Encoding (All Types)
```
Value: 300 (0x012C)
Binary: 10101100 00000010
Bytes: [0xAC, 0x02]
       └─┬──┘  └─┬──┘
    continuation  final
         bit       byte
```

### ZigZag Encoding
```
Original: -1
ZigZag: ((-1) << 1) ^ ((-1) >> 31) = 1
Varint: 0x01 (1 byte)

Original: -2
ZigZag: ((-2) << 1) ^ ((-2) >> 31) = 3
Varint: 0x03 (1 byte)
```

---

## Usage Recommendations

### When to Use Extensions

✅ **Use extensions when:**
- Building Go-only protocols
- Performance optimization is critical
- You need features not in standard ProtoDef
- Cross-platform compatibility is not required

❌ **Avoid extensions when:**
- Building multi-language systems
- Protocol must be ProtoDef-standard compliant
- Interoperability with other ProtoDef implementations needed
- Protocol documentation should be universal

### Migration Path

If you need to migrate from extensions to standard types:

| Extension | Standard Alternative | Trade-off |
|-----------|---------------------|-----------|
| `varint64` | `i64` / `u64` | Fixed 8 bytes, no space optimization |
| `varint128` | Two `i64` fields | More complex, but standard |
| `zigzag32` | `i32` | Fixed 4 bytes, no space optimization |
| `zigzag64` | `i64` | Fixed 8 bytes, no space optimization |

---

## Implementation Status

### Current Implementation

| Type | Parsing | Type Definition | Encoding | Decoding |
|------|---------|----------------|----------|----------|
| `varint64` | ✅ | ✅ | ⚠️ Pending | ⚠️ Pending |
| `varint128` | ✅ | ✅ | ⚠️ Pending | ⚠️ Pending |
| `zigzag32` | ✅ | ✅ | ⚠️ Pending | ⚠️ Pending |
| `zigzag64` | ✅ | ✅ | ⚠️ Pending | ⚠️ Pending |

**Note:** Type definitions and parsing are complete. Runtime encoding/decoding will be implemented in future releases.

---

## Code Examples

### Using Extensions in Protocol Definitions

```json
{
  "types": {
    "player_position": [
      "container",
      [
        {"name": "entityId", "type": "varint64"},
        {"name": "x", "type": "zigzag32"},
        {"name": "y", "type": "zigzag32"},
        {"name": "z", "type": "zigzag32"}
      ]
    ],
    "session": [
      "container",
      [
        {"name": "sessionId", "type": "varint128"},
        {"name": "timestamp", "type": "varint64"},
        {"name": "balance", "type": "zigzag64"}
      ]
    ]
  }
}
```

### Go Code Usage

```go
// Type definitions (available now)
var (
    VarInt64  = &Type{Name: "varint64", TypeName: "varint64"}
    VarInt128 = &Type{Name: "varint128", TypeName: "varint128"}
    ZigZag32  = &Type{Name: "zigzag32", TypeName: "zigzag32"}
    ZigZag64  = &Type{Name: "zigzag64", TypeName: "zigzag64"}
)

// Encoding/decoding (future implementation)
// encoder.WriteVarint64(value)
// value := decoder.ReadZigZag32()
```

---

## References

### Related Standards
- **Protocol Buffers:** [Varint Encoding](https://developers.google.com/protocol-buffers/docs/encoding#varints)
- **Protocol Buffers:** [Signed Integers (ZigZag)](https://developers.google.com/protocol-buffers/docs/encoding#signed-ints)

### ProtoDef Specification
- **Official Spec:** [ProtoDef Repository](https://github.com/ProtoDef-io/ProtoDef)
- **Standard Types:** [datatypes.md](../ProtoDef/doc/datatypes.md)

---

## Contributing

If you have ideas for additional extensions or improvements to existing ones:

1. Ensure the extension doesn't duplicate standard ProtoDef functionality
2. Document the use case clearly
3. Provide encoding/decoding specifications
4. Consider cross-platform implications
5. Submit as a pull request with tests

---

## Version History

### Version 1.0 (2025-10-31)
- Initial documentation of varint64, varint128, zigzag32, zigzag64
- Type definitions and parsing implemented
- Encoding/decoding marked as future work

---

## License

This extension documentation is part of the protodef-go project and follows the same license as the main project.
