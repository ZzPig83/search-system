package main

import (
	"search-system/service"
	"testing"
)

func TestSearch(t *testing.T) {
	query := "城市刑警普通案件"

	service.Search(query)
}
