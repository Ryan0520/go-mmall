package backend

import (
	"github.com/Ryan0520/go-mmall/models"
	"testing"
)

func TestProduct_Count(t *testing.T) {
	models.Setup()
	p := &Product{}
	count, err := p.Count()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
	}
}
