//+build mage

package main

import (
	"fmt"
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/mg"
)

// Hello prints a message (shows that you can define custom Mage targets).
func Hello() {
	fmt.Println("hello plugin developer!")
}

func BuildAllExceptDarwinARM64() { //revive:disable-line
	b := build.Build{}
	mg.Deps(b.Linux, b.Windows, b.Darwin, b.LinuxARM64, b.LinuxARM)
}

// Default configures the default target.
var Default = BuildAllExceptDarwinARM64
