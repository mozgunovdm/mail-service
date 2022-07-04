package response

import "net/http"

type Response struct{}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
