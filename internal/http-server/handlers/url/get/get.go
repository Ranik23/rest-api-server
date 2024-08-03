package get


import (
	//"fmt"
	//"net/http"
)


type URLGetter interface {
	GetURL(alias string) (string, error)
}

