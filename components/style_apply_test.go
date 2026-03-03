package components

// style priorization:
// - on first render, each style property (width, color, etc) is applied top-down. last in the tree wins.
// - after that, properties are prioritised winthin each layer.
//
// "blue" wins:
// Box(
// 	Apply(Style{Color: "red"}),
// 	Apply(Style{Color: "blue"}),
// )
//
// "blue" wins:
// Box(
// 	Apply(
// 		Style{Color: "red"},
// 		Style{Color: "blue"},
//  ),
// )
//
// "blue" wins (even when signal changes):
// Box(
// 	Apply(
// 		Style{Color: getColor},
// 		Style{Color: "blue"},
//  ),
// )
//
// "blue" wins, then getColor signal wins on change:
// Box(
// 	Apply(Style{Color: getColor}),
// 	Apply(Style{Color: "blue"}),
// )
//
// "blue" wins, then "red" wins on hover:
// Box(
// 	ApplyOn("hover", Style{Color: "red"}),
// 	Apply(Style{Color: "blue"}),
// )
