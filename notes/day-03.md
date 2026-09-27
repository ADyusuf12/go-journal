# Day 3 Notes

## Structs & Methods

Go structs group related fields. Methods attach behavior to structs via **receivers**.

### Basic Struct Example

```go
type Job struct {
    Title  string
    Status string
}
```

### Receivers: Value vs. Pointer

Memory and copying rules are strictly enforced by the compiler:

**Value Receiver** (`func (j Job) Summary() string`):
- Operates on a copy of the struct created on the stack frame
- Guarantees immutability; cannot alter the caller's original memory
- Use for read-only calculations or small value objects (e.g., Money)

**Pointer Receiver** (`func (j *Job) Publish() error`):
- Passes an 8-byte memory address
- Operates directly on the caller's instance to mutate state or avoid copying large structs
- **Rule of Thumb**: If one method on a struct requires a pointer receiver, use pointer receivers for all methods on that struct

### Automatic Pointer Dereferencing & Address-Taking

Go provides syntactic sugar for method calls so you don't manually write `(*j)` or `(&j)`:
```go
inv1.MarkAsPaidPointer() // Go automatically translates to (&inv1).MarkAsPaidPointer()
inv2.MarkAsPaidValue()   // Go automatically translates to (*inv2).MarkAsPaidValue()
```

### Visibility & Factory Constructor Pattern

- **Exported** (Capitalized SKU): Accessible outside the package
- **Unexported** (lowercase price): Package-private

Since Go lacks built-in class constructors, use the `New<StructName>` factory pattern to enforce domain invariants:
```go
func NewLineItem(sku string, quantity int) (*LineItem, error) {
    if quantity <= 0 {
        return nil, errors.New("quantity must be greater than zero")
    }
    return &LineItem{SKU: sku, Quantity: quantity}, nil
}
```

### Struct Embedding (Composition over Inheritance)

Go has no class inheritance (Invoice < ApplicationRecord). Instead, it embeds structs inside other structs:
```go
type AuditHeader struct {
    ID        int64
    CreatedAt time.Time
}

type Invoice struct {
    AuditHeader // Embedded struct (anonymous field)
    TotalCents  int64
}
```

**Field & Method Promotion**: Fields and methods of `AuditHeader` are accessible directly on `Invoice` (e.g., `inv.ID` or `inv.IsPersisted()`).

**Shadowing**: If `Invoice` defines its own `IsPersisted()` method, it overrides the embedded one.

### Core Mental Shifts from Rails

- **Receivers vs. Bang (!) Methods**: Pointer receivers are compiler-enforced memory semantics (passing reference addresses vs. copying values on the stack), not voluntary naming conventions like Ruby's bang methods.
- **Value Objects vs. Entities**: Value receivers act like immutable Rails Value Objects (Money). Pointer receivers act like stateful ActiveRecord models.
