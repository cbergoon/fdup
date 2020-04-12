<h1 align="center">fdup - file duplication utility</h1>
<p align="center">
<a href="https://travis-ci.org/cbergoon/fdup"><img src="https://travis-ci.org/cbergoon/fdup.svg?branch=master" alt="Build"></a>
<a href="https://goreportcard.com/report/github.com/cbergoon/fdup"><img src="https://goreportcard.com/badge/github.com/cbergoon/fdup?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/fdup"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

```fudp``` identifies duplicated files in a directory by hashing contents. Optionally, ```fdup``` can also identify directories with an exact duplicate within the 
specified scope. 

```fdup``` uses a hasing strategy to identify files (and/or directories) whose content is duplicated. Name and other metadata are not relevant to the results. 

#### Documentation 

See the docs [here](https://godoc.org/github.com/cbergoon/fdup).

#### Install
```
$ go get github.com/cbergoon/fdup
$ go install github.com/cbergoon/fdup
```

#### Example Usage
```sh
    $ fdup
```

#### License
This project is licensed under the MIT License.