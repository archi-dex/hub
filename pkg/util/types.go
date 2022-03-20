package util

type Map interface {
	Get(string) string
	Has(string) bool
}
