# soda
Soda is a tiny CLI application that helps establishing a one-to-one encrypted and authenticated channel for private communication, driven by NaCl's public-key cryptosystem [Box](https://nacl.cr.yp.to/box.html) (Curve25519 + XSalsa20 + Poly1305).

## Install
[dep](https://github.com/golang/dep) and [goversioninfo](https://github.com/josephspurrier/goversioninfo) are required.
```bash
$ git clone https://github.com/Equim-chan/soda.git $GOPATH/src/ekyu.moe/soda
$ cd $GOPATH/src/ekyu.moe/soda
$ make
$ make install
$ $GOPATH/bin/soda
```

## License
[Apache-2.0](https://github.com/Equim-chan/soda/blob/master/LICENSE)
