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
	On(string, string) StrSetter
	All(string) StrSetter
}

type strSetter string

func (s strSetter) On(env, val string) StrSetter {
	key := string(s)
	strs[env][key] = val
	return s
}

func (s strSetter) All(val string) StrSetter {
	key := string(s)
	for _, env := range strs {
		env[key] = val
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
