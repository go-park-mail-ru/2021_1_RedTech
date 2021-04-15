package jsonerrors

import "fmt"

func JSONMessage(m string) string {
	return fmt.Sprintf(`{"message":"%s"}`, m)
}

var (
	JSONEncode = JSONMessage("json encode")
	JSONDecode = JSONMessage("json decode")
	URLParams  = JSONMessage("url params")
	Session    = JSONMessage("session")
	CSRF       = JSONMessage("no csrf-token")
)
