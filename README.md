# soda
Soda is a tiny CLI application that helps establishing a one-to-one encrypted and authenticated channel for private communication, driven by NaCl's public-key cryptosystem [Box](https://nacl.cr.yp.to/box.html) (Curve25519 + XSalsa20 + Poly1305).

## Install
Recommended way ([dep](https://github.com/golang/dep) is required):
```bash
$ git clone https://github.com/Equim-chan/soda.git $GOPATH/src/ekyu.moe/soda
$ dep ensure
$ go install
$ $GOPATH/bin/soda
```

Traditional way:
```bash
$ go get ekyu.moe/soda
$ $GOPATH/bin/soda
```

## License
[Apache-2.0](https://github.com/Equim-chan/soda/blob/master/LICENSE)
