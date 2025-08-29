package hint

// Envs returns the environments for command line autocompletion.
func Envs() []string {
	return []string{"dev", "stg", "prod"}
}
