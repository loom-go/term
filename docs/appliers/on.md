---
weight: 3
title: On{}
---

```go {style=tokyonight-moon}
type On struct {
	Hover        func(*core.EventMouse)
	Click        func(*core.EventMouse)
	MouseEnter   func(*core.EventMouse)
	MouseLeave   func(*core.EventMouse)
	MousePress   func(*core.EventMouse)
	MouseRelease func(*core.EventMouse)
	MouseMove    func(*core.EventMouse)
	MouseScroll  func(*core.EventMouse)
	MouseDrag    func(*core.EventMouse)

	KeyPress   func(*core.EventKey)
	KeyRelease func(*core.EventKey)
	Paste      func(*core.EventPaste)

	Focus func(*core.EventFocus)
	Blur  func(*core.EventBlur)

	Input  func(*core.EventInput)
	Submit func(*core.EventSubmit)
}
```
