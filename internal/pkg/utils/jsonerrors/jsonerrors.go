package jsonerrors

import "fmt"

func JSONMessage(m string) string {
	return fmt.Sprintf(`{"message":"%s"}`, m)
}
