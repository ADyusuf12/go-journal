# Day 4: Interfaces, Type Assertions & Type Switches

## Key Concepts Covered

1. **Implicit Interface Satisfaction**:
   - Interfaces define contracts (method signatures only).
   - Go has no `implements` keyword. A concrete struct satisfies an interface automatically if it implements all required methods with matching parameter and return signatures.

2. **Interfaces Belong to the Consumer**:
   - Go convention: "Accept interfaces, return structs".
   - Interfaces are defined by the consumer/caller package to keep abstractions decoupled and lightweight.

3. **Method Sets & Pointer Receivers**:
   - A value instance `T` only possesses value-receiver methods `(t T)`.
   - A pointer instance `*T` possesses both value-receiver `(t T)` and pointer-receiver `(t *T)` methods.
   - If an interface method uses a pointer receiver `(t *T)`, you must pass a pointer address (`&Struct{}`) to satisfy the interface.

4. **Type Assertions (`val, ok := i.(ConcreteType)`)**:
   - Allows runtime extraction or inspection of the underlying concrete value held by an interface variable.
   - Always use the two-value "comma, ok" idiom (`if !ok`) to avoid runtime panics on type mismatches.

5. **Type Switches (`switch v := i.(type)`)**:
   - A control structure to branch logic based on the concrete type inside an interface variable.
   - Inside each `case` block, `v` is automatically scoped and typed as that specific concrete type.

6. **The Empty Interface (`any` / `interface{}`)**:
   - Satisfied by every type in Go.
   - Used for dynamic data structures (e.g., JSON parsing, generic logging), but should be avoided when static typing can be preserved.

---

## Code Example Summary (`interfaces/main.go`)
