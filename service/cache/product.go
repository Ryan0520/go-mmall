package cache

import (
	"github.com/Ryan0520/go-mmall/pkg/e"
	"strconv"
	"strings"
)

type Product struct {
	ID int

	FilterOffSale bool
	Keyword       string
	OrderBy       string
	PageNum       int
	PageSize      int
}

func (p *Product) GetProductKey() string {
	if p.FilterOffSale {
		return e.CacheProduct + "_" + strconv.Itoa(p.ID) + "_" + "on_sale"
	}
	return e.CacheProduct + "_" + strconv.Itoa(p.ID)
}

func (p *Product) GetProductsKey() string {
	keys := []string{
		e.CacheProduct,
		"LIST",
	}

	if p.ID > 0 {
		keys = append(keys, strconv.Itoa(p.ID))
	}
	if p.FilterOffSale {
		keys = append(keys, "on_sale")
	}
	if p.Keyword != "" {
		keys = append(keys, p.Keyword)
	}
	if p.OrderBy != "" {
		keys = append(keys, p.OrderBy)
	}
	if p.PageNum > 0 {
		keys = append(keys, strconv.Itoa(p.PageNum))
	}
	if p.PageSize > 0 {
		keys = append(keys, strconv.Itoa(p.PageSize))
	}

	return strings.Join(keys, "_")
}
