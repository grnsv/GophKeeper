package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
)

func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("error_message")
	if code == http.StatusInternalServerError {
		log.Println(err)
		e.StrEscape("Internal Server Error")
	} else {
		e.StrEscape(err.Error())
	}
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
