package pdcl

type LocalStorageConfig struct {
	Directory string `envconfig:"PDCL_LOCAL_DIRECTORY" required:"true"`
}
