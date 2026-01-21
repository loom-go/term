module github.com/AnatoleLucet/loom-term

go 1.25.6

require (
	github.com/AnatoleLucet/tess v0.0.0-20260203151116-72dc74b8dcb6
	golang.org/x/term v0.38.0
)

require (
	github.com/AnatoleLucet/go-opentui v0.0.0-20260209134342-7a88a7c05506 // indirect
	github.com/clipperhouse/uax29/v2 v2.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/AnatoleLucet/go-opentui => ../../../go-opentui

replace github.com/AnatoleLucet/tess => ../../../tess
