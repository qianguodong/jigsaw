package main

import (
	"github.com/guodongq/jigsaw/pkg/operator"
	"github.com/guodongq/jigsaw/pkg/util/profile"
)

func main() {
	defer profile.Profile().Stop()

	operator.Execute()
}
