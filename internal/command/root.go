package command

type rootCommand struct {
	Join    joinCommand    `cmd:"" help:"Join multiple Dockerfiles into a single multi-stage Dockerfile."`
	Version versionCommand `cmd:"" help:"Show the multidockerfile version information."`
}
