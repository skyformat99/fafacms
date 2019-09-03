package model

import (
	"testing"
)

func TestGroup_GetById(t *testing.T) {
	Testxx()

	g := new(Group)
	g.Id = 2
	FafaRdb.Client.Id(g.Id).Get(g)
}
