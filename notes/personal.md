# When you use ActiveRecord, you're operating at a high level of abstraction:

## Ruby: Beautiful, expressive, but hides everything

```ruby
Product.transaction do
  product = Product.lock.find_by!(sku: "PRD-101")
  product.decrement!(:stock, 5)
end
```

Underneath that simple Ruby block, Rails is doing a dizzying amount of heavy lifting: grabbing a connection from connection pool thread-locals, issuing a BEGIN, building the SQL string, managing prepared statement caches, type-casting the PG adapter responses into Ruby objects, tracking dirty attributes, issuing the UPDATE, and firing callbacks—all while hiding the memory and I/O mechanics from you.

When you transition to Go, the curtain gets pulled completely back.

## What Changes in the "Programming Mindset"

### You See the I/O Cost

You realize that every database call isn't just a method on an object—it's a socket write over TCP, a stream of binary bytes coming back down the wire, and an explicit set of memory allocations to parse those bytes.

### Explicit Memory Control

Instead of an ORM instantiating heavy Ruby objects with internal state tracking, you hand Go a memory address pointer (`&product.Stock`), and the standard library streams raw bytes directly into that memory slot.

### Deterministic Failure Modes

In Rails, an unhandled exception inside a transaction block rolls back automatically via control-flow magic. In Go, `defer tx.Rollback()` makes resource cleanup explicit, deterministic, and visible right at the call site.

You aren't just writing backend applications anymore—you're managing system resources, network streams, and memory layouts. That mechanical clarity makes you a drastically better engineer in any language, including Ruby!

---

## Rails: The Opinionated Mansion vs. Go: Unfiltered Freedom

In Rails, you are working inside a massive, highly opinionated mansion ("Convention over Configuration"). When you want to do things the way Rails intends, it is lightning-fast. But the moment you want to customize the plumbing—like tweaking how middleware passes request context or intercepting status codes—you find yourself fighting heavy abstractions, monkey-patching, or digging through ActionDispatch source code.

In Go, because everything boils down to simple interfaces and functions (`http.Handler`), you own every single wire. You aren't fighting a framework; you're just writing code. You have total architectural freedom to build exactly the pipeline you need with zero hidden magic.
