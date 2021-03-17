package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	t1 := time.Now()
	t2 := t1.UTC()
	t5 := time.Unix(t2.Unix(), 0).UTC()
	fmt.Println(t5)
}
