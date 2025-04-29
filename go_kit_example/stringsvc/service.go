package main

import (
	"fmt"
	"strings"
)

var _ StringService = (*stringService)(nil)

type StringService interface {
	UpperCase(string) (string, error)
	Count(string) int
}

type stringService struct{}

func (stringService) UpperCase(s string) (string, error) {
	if s == "" {
		return "", EmptyErr
	}

	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

var EmptyErr = fmt.Errorf("empty string")
