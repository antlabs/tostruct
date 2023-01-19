// antlabs, guonaihong 2022
// apache 2.0
package fromhead

import (
	"bytes"
	"net/http"

	"github.com/antlabs/tostruct/internal/map2struct"
)

type FromHead struct {
	buf        bytes.Buffer
	structName string
	h          http.Header
}

func New(structName string) *FromHead {
	return &FromHead{structName: structName, h: make(http.Header)}
}

func (h *FromHead) AppendHeader(hk, hv []byte) *FromHead {
	h.h.Add(string(hk), string(hv))
	return h
}

func (h *FromHead) Gen() (string, error) {
	return map2struct.MapGenStruct(h.h, h.structName, "header")
}
