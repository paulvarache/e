// +build linux,darwin
package main

var SetVarPrefix string

func init() {
	SetVarPrefix = "export "
}
