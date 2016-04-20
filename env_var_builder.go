package piper

import "strings"

type EnvVarBuilder struct{}

func (b EnvVarBuilder) Build(environment []string, params map[string]string) []DockerEnv {
	env := make(map[string]string)
	for _, variable := range environment {
		parts := strings.Split(variable, "=")
		env[parts[0]] = parts[1]
	}

	var envVars []DockerEnv
	for key, value := range params {
		if env[key] != "" {
			value = env[key]
		}
		envVars = append(envVars, DockerEnv{
			Key:   key,
			Value: value,
		})
	}

	return envVars
}
