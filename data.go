package main

type Row interface {
	Value(key string) string
	Values() []string
}
