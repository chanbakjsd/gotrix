# gotrix

[![Examples](https://img.shields.io/badge/Example-__examples%2F-blueviolet?style=flat-square)][example-link]
[![Report Card](https://goreportcard.com/badge/github.com/chanbakjsd/gotrix?style=flat-square)][report-link]
[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue?style=flat-square)][doc-link]

[example-link]: https://github.com/chanbakjsd/gotrix/tree/master/_examples
[report-link]: https://goreportcard.com/report/github.com/chanbakjsd/gotrix
[doc-link]: https://pkg.go.dev/github.com/chanbakjsd/gotrix

Gotrix is a work-in-progress implementation of the client portion of [Matrix's client-server API][spec-link].
It is still actively being developed and the API may change especially when they are recently introduced.

It currently implements all of the parts mandated by specification but does not implement all modules available.
A list of available modules that are done can be found at [TODO.md][todo-link].

If you require the use of features that have not been implemented yet,
[gomatrix][gomatrix-link] and [mautrix-go][mautrixgo-link] are alternative clients of Matrix in Go.

[spec-link]: https://spec.matrix.org/v1.1/client-server-api/
[todo-link]: https://github.com/chanbakjsd/gotrix/blob/master/TODO.md
[gomatrix-link]: https://github.com/matrix-org/gomatrix
[mautrixgo-link]: https://github.com/tulir/mautrix-go
