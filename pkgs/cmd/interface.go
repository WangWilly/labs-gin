package cmd

//go:generate mockgen -source=interface.go -destination=cmd_mock.go -package=cmd
type Cmd interface {
	Run() error
}
