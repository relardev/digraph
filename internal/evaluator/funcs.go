package evaluator

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/maja42/goval"
)

func Prepare(expression string) string {
	return strings.ReplaceAll(expression, "url.Parse(", "urlParseQQQ(")
}

var Funcs = map[string]goval.ExpressionFunction{
	"urlParseQQQ": func(args ...interface{}) (interface{}, error) {
		u, err := url.Parse(args[0].(string))
		if err != nil {
			return nil, err
		}
		bytes, err := json.Marshal(u)
		if err != nil {
			return nil, err
		}
		var final map[string]interface{}
		err = json.Unmarshal(bytes, &final)
		if err != nil {
			return nil, err
		}
		return final, nil
	},
	"len": func(args ...interface{}) (interface{}, error) {
		arg := args[0]
		switch arg := arg.(type) {
		case []interface{}:
			return len(arg), nil
		case map[string]interface{}:
			return len(arg), nil
		case []string:
			return len(arg), nil
		case string:
			return len(arg), nil
		case []int:
			return len(arg), nil
		default:
			return nil, fmt.Errorf("unknown type in len %T", arg)
		}
	},
}
