---
weight: 3
title: Pacer{}
---

```go {style=tokyonight-moon}
type Pacer struct {
	rate     time.Duration
}

func NewPacer(rate time.Duration) *Pacer
```
