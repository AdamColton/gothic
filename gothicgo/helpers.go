package gothicgo

type errStr string

func (err errStr) Error() string {
	return string(err)
}
