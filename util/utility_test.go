package util

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type UtilSuite struct{}

var _ = Suite(&UtilSuite{})

func (s *UtilSuite) TestRandom(c *C) {
	r1 := String(64)
	c.Check(r1, HasLen, 64)

	r2 := String(64)
	c.Check(r1, Not(Equals), r2)
}
