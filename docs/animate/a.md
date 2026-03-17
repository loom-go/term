---
weight: 2
title: A{}
---

```go {style=tokyonight-moon}
type A struct {
	Context  context.Context
	Duration time.Duration
	Tick     func(progress float64)
	Pacer    *Pacer
}
```
