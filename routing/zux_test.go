package routing

import (
	"testing"
)

func Test_Zuxer(t *testing.T) {
	go func() { runMuxer() }()

	select {}
}
