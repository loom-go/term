module github.com/loom-go/loom/examples/term/console

go 1.25.6

require (
	github.com/loom-go/loom v0.0.0-20260309223821-57d50fb2517d
	github.com/loom-go/term v0.0.0-00010101000000-000000000000
)

require (
	github.com/AnatoleLucet/go-opentui v0.0.0-20260311124333-d904eb66f503 // indirect
	github.com/AnatoleLucet/sig v0.0.0-20260308162001-17251018b48a // indirect
	github.com/AnatoleLucet/tess v0.0.0-20260310111309-c8343f5a151d // indirect
	github.com/petermattis/goid v0.0.0-20251121121749-a11dd1a45f9a // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.38.0 // indirect
)

replace github.com/loom-go/term => ../..
