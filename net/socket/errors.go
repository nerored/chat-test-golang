/*
	网络错误区分：用于分辨是否需要终止服务
*/
package socket

import (
	"net"
)

type errortype int

const (
	error_type_noerror errortype = iota
	error_type_temporary
	error_type_othererro
)

func errcheck(err error) errortype {
	if err == nil {
		return error_type_noerror
	}

	switch realErr := err.(type) {
	case net.Error:
		if realErr.Temporary() {
			return error_type_temporary
		}

		return error_type_othererro
	default:
		return error_type_othererro
	}
}
