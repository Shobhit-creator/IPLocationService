package utils

import (
	"net"
	"net/http"
)

func GetIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header (used by proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header (sometimes used instead of X-Forwarded-For)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // return as-is if parsing fails
	}

	return ip
}
