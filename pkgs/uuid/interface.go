package uuid

//go:generate mockgen -source interface.go -destination interface_mock.go -package uuid
type UUID interface {
	New() string
}
