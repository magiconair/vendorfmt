# Overview

vendorfmt rewrites a `vendor/vendor.json` file used by https://github.com/kardianos/govendor to a
more merge friendly format.

## Install

```
go get -u github.com/magiconair/vendorfmt/cmd/vendorfmt
```

## Usage

```shell
# format vendor/vendor.json
$ vendorfmt

# format other files
$ vendorfmt foo/bar/vendor.json
```

### Before

```json
{
	"comment": "",
	"ignore": "test",
	"package": [
		{
			"path": "appengine",
			"revision": ""
		},
		{
			"path": "appengine_internal",
			"revision": ""
		},
		{
			"path": "appengine_internal/base",
			"revision": ""
		}
	]
}
```

### After

```json
{
	"comment": "",
	"ignore": "test",
	"package": [
		{"path":"appengine","revision":""},
		{"path":"appengine_internal","revision":""},
		{"path":"appengine_internal/base","revision":""}
	]
}
```
