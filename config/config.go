package config

var strs = make(map[string]map[string]string) //[env][key] => val

func Environments(environments ...string) {
	for _, env := range environments {
		strs[env] = make(map[string]string)
	}
	if activeEnv == "" && len(environments) > 0 {
		activeEnv = environments[0]
	}
}

var activeEnv string

func SetEnvironment(env string) {
	activeEnv = env
}

type StrSetter interface {
	As(string, ...string) StrSetter
}

type strSetter string

func (s strSetter) As(val string, envs ...string) StrSetter {
	key := string(s)
	if len(envs) == 0 {
		for _, env := range strs {
			env[key] = val
		}
		return s
	}

	for _, envStr := range envs {
		if env, ok := strs[envStr]; ok {
			env[key] = val
		}
	}
	return s
}

func SetString(key string) StrSetter { return strSetter(key) }
func GetString(key string) string {
	env, ok := strs[activeEnv]
	if !ok {
		return ""
	}
	return env[key]
}
