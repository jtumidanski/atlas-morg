package handlers

import (
	"atlas-morg/attributes"
	"context"
	"net/http"
)

func (w *World) MiddlewareValidateMonster(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		cs := &attributes.MonsterInputDataContainer{}
		err := attributes.FromJSON(cs, r.Body)
		if err != nil {
			w.l.Println("[ERROR] deserializing monster input", err)
			rw.WriteHeader(http.StatusBadRequest)
			attributes.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyWorld{}, *cs)
		r = r.WithContext(ctx)

		f(rw, r)
	})
}
