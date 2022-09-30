package albatross

import "regexp"

func verifyUrl(url string) (bool, error) {
	regex := `^(https|http|ws|wss):\/\/`
	return regexp.Match(regex, []byte(url))
}

func addOptionalParam[T any, D any](params []interface{}, optional []T, defaultValue D) []interface{} {
	if len(optional) > 0 {
		return append(params, optional[0])
	}
	return append(params, defaultValue)
}
