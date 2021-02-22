package model

import "auth-server/pkg/errors/error"

//OkResponce is model of OK responce
type OkResponce struct {
	Response map[string]interface{} `json:"response"`
}

//CreateOneOkResponce a constructor of "OK" responce
func CreateOneOkResponce(item interface{}) *OkResponce {
	return &OkResponce{
		Response: map[string]interface{}{
			"item": item,
		},
	}
}

//CreateOkResponce a constructor of "OK" responce
func CreateOkResponce(count int, items []interface{}) *OkResponce {
	return &OkResponce{Response: map[string]interface{}{
		"count": count,
		"items": items,
	}}
}

//BadRequest is model of BAD responce
type BadResponce struct {
	Error *error.HTTPError `json:"error"`
}

//CreateBadResponce a constructor of "BadResponce" struct
func CreateBadResponce(err *error.HTTPError) *BadResponce {
	return &BadResponce{err}
}
