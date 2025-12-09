# Common Errors and Lessons Learned

## Go Testing Issues

### Empty Loop Bug in Tests
**Problem:** Test was passing when it should fail because the loop never executed.

**Example:**
```go
// BAD: This test passes even when cleanInput() returns empty slice
for i := range actual {
    // This never runs if actual is empty
    if actual[i] != expected[i] {
        t.Errorf("...")
    }
}

// GOOD: Check lengths first
if len(actual) != len(expected) {
    t.Errorf("Expected length of %d, got %d", len(expected), len(actual))
    continue
}
for i := range actual {
    // Now we know both slices have same length
    if actual[i] != expected[i] {
        t.Errorf("...")
    }
}
```

**Lesson:** Always validate assumptions in tests - check lengths, nil values, etc. before iterating.

## Package Declaration Issues

### Package Name Mismatch
**Problem:** `found packages main (main.go) and abs (repl_test.go)` 

**Solution:** Test files must declare the same package as the code they're testing.

```go
// main.go
package main

// repl_test.go - must match
package main  // not package abs
```

## Go Idioms

### Comma OK Pattern for Map Lookups
**Problem:** Need to distinguish between a key that exists with zero value vs a key that doesn't exist.

**Example from repl.go:27:**
```go
// GOOD: Use comma ok idiom
command, exists := getCommands()[commandInput]
if exists {
    err := command.callback()
    // handle command
} else {
    fmt.Println("Unknown command")
}

// BAD: Direct access can't distinguish missing vs zero value
command := getCommands()[commandInput]
if command != nil { // This won't work for structs
    // ...
}
```

**Lesson:** Always use `value, ok := map[key]` when you need to check if a key exists in a map.