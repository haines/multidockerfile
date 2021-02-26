package command

type rootCommand struct {
	Version versionCommand `cmd:"" help:"Show the multidockerfile version information."`
}
