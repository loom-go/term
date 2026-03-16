---
weight: 2
title: Attr{}
---

```go {style=tokyonight-moon}
type Attribute struct {
	Title       any // string | func() string
	Value       any // string | func() string
	Placeholder any // string | func() string
}

type Attr = Attribute
```
