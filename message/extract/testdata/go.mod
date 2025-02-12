module github.com/issue9/localeutil/testdata

require github.com/issue9/localeutil v1.0.0

// NOTE: 需要保证与根目录中 go.mod 的 text 具有相同的版本，否则测试会失败！
require golang.org/x/text v0.22.0 // indirect

replace github.com/issue9/localeutil => ../../..

go 1.23.0
