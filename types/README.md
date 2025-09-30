# Type System and Code Generation

Automatic TypeScript definition generation for cross-language type safety.

## Features

- **Automatic Generation**: Parse Python type hints, Rust structs, Go interfaces
- **TypeScript AST**: Uses ts-morph for TypeScript AST manipulation
- **Cross-language Mapping**: Python mypy → TypeScript, Rust serde → TypeScript
- **Protocol Buffers**: Define cross-language interfaces
- **IntelliSense**: Full IDE support with cross-language autocomplete

## Supported Type Mappings

### Python → TypeScript
- `int` → `number`
- `str` → `string`
- `List[T]` → `T[]`
- `Dict[K, V]` → `Record<K, V>`
- `Optional[T]` → `T | null`

### Rust → TypeScript
- `i32` → `number`
- `String` → `string`
- `Vec<T>` → `T[]`
- `HashMap<K, V>` → `Record<K, V>`
- `Option<T>` → `T | null`

### Go → TypeScript
- `int` → `number`
- `string` → `string`
- `[]T` → `T[]`
- `map[K]V` → `Record<K, V>`

## Implementation Status

🚧 **Planning Phase** - Type system architecture design in progress.
