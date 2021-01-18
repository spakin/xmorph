xmorph
======

[![Go Reference](https://pkg.go.dev/badge/github.com/spakin/xmorph.svg)](https://pkg.go.dev/github.com/spakin/xmorph)
[![Go Report Card](https://goreportcard.com/badge/github.com/spakin/xmorph)](https://goreportcard.com/report/github.com/spakin/xmorph)

Description
-----------

`xmorph` is a package for the [Go programming language](https://golang.org/) that provides functions for warping and morping images.  In this context, *warping* means distorting a single image, and *morphing* means gradually distorting and dissolving one image into another.  In both cases, the caller provides two meshes: a source mesh and a destination mesh.  The functions provided by the package distort images such that points that originally lie above the source mesh are made to lie above the corresponding points on the destination mesh.

The `xmorph` package provides the following features:

* Warping and morphing work on any image type that implements [`image.Image`](https://golang.org/pkg/image/#Image).

* A mesh can be created as empty (all-zero coordinates), with coordinates regularly spaced over a specified area, or converted from a given 2-D slice of coordinates.

* Coordinates can be described with either an [`image.Point`](https://golang.org/pkg/image/#Point) (integer-valued) or an analogous `xmorph.Point` (floating-point-valued).

* Meshes can be read from and written to files in the same format used by `morph`, `xmorph`, and `gtkmorph`, facilitating interoperability.

The package itself is primarily a Go interface to the venerable [`libmorph` library](http://xmorph.sourceforge.net/).  `libmorph` provides the foundation for the `morph` command-line program and the `xmorph` and `gtkmorph` graphical user interfaces.

Installation
------------

You will first need to install the [`libmorph`](http://xmorph.sourceforge.net/) library and header files.  If your operating system does not provide these via its usual software-installation mechanism, the `libmorph` source code is available from SourceForge under https://sourceforge.net/projects/xmorph/.  (Download the latest `xmorph-*.tar.gz` file.)  The Go package expects to find the header files in an `xmorph` subdirectory in your C preprocessor's include path and to be able to link to `libmorph` by passing `-lmorph` to the linker.

The `xmorph` package has opted into the [Go module system](https://golang.org/ref/mod).  Hence, once `libmorph` is installed properly, `xmorph` should install automatically when your Go code does an
```Go
import "github.com/spakin/xmorph"
```

`xmorph` can also be installed from the command line by running
```bash
go get -u github.com/spakin/xmorph
```
or by manually downloading the code from GitHub and building and installing it.

Documentation
-------------

Descriptions and examples of the `xmorph` API can be found online in the [pkg.go.dev `xmorph` documentation](https://pkg.go.dev/github.com/spakin/xmorph).

Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott+xmorph@pakin.org*
