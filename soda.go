// Soda is a tiny CLI application that helps establishing a one-to-one
// encrypted and authenticated channel for private communication, driven
// by NaCl's public-key cryptosystem Box (Curve25519 + XSalsa20 + Poly1305).
package main // import "ekyu.moe/soda"

//go:generate goversioninfo -icon=icon.ico

var (
	Version   = "(dev)"
	GitHash   = "unknwon-dirty"
	BuildDate = "(unknown)"
)
