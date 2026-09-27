# Day 2

## Functions

Functions can accept typed parameters:

```go
func greet(name string)
```

## Return Values

Single return value:

```go
func greet() string
```

Multiple return values:

```go
func getUser() (string, int)
```

## Error Handling

Common pattern:

```go
result, err := someFunction()
```

## Nil

No error:

```go
return result, nil
```

## Error Creation

```go
errors.New("message")
```

## Common Go Pattern

```go
result, err := someFunction()

if err != nil {
	return err
}
```

## Biggest Insight

Go prefers explicit error handling instead of relying primarily on exceptions.
``