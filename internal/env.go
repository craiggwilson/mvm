package internal

func ExecuteEnv(cfg *EnvConfig) error {

	writef(cfg.Config, "set %s=%s", "MVM", cfg.MVMDirectory)

	return nil
}

type EnvConfig struct {
	*Config
}
