# IRCS Daemon (KDS)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/pedroalbanese/ircs/blob/master/LICENSE.md) 
[![GoDoc](https://godoc.org/github.com/pedroalbanese/ircs?status.png)](http://godoc.org/github.com/pedroalbanese/ircs)
[![Go Report Card](https://goreportcard.com/badge/github.com/pedroalbanese/ircs)](https://goreportcard.com/report/github.com/pedroalbanese/ircs)
### Internet Relay Chat Secure (Key Distribution System)
Minimalist Internet Relay Chat (IRC) via Transport Layer Security (RFC 7194).

Internet Relay Chat (IRC) is a text-based chat protocol that enables real-time exchange of text messages between users connected to an IRC server. It functions as a network of virtual chat rooms where users can interact with other users in real time. The server acts as a router, ensuring that all messages are delivered to the correct recipients.

```
   +-----------------------+      +----------------------+
   |   Certificate         |      |        Server        |
   |   Authority (CA)      |      |                      |
   |                       |      |                      |
   |   Sign                |      |    Generate          |
   |   CSR                 |      |    CSR               |
   |                       |      |                      |
   V                       V      V                      V
+-------+                 +-------+                  +--------+
| AKID  |                 | AKID  |                  | AKID   |
+-------+                 +-------+                  +--------+
   |                        |                            |
   |                        |                            |
   V                        |                            V
 Client                     |                     Access Denied
 with Certificate           |                     No valid AKID
                            |
                            V
                      +------------+
                      |  CRL       |
                      |  Check     |
                      +------------+
                            |
                            |     +---------------------+
                            |     |    Revocation List  |
                            +---->|                     |
                                  |    Not after XXX    |
                                  +---------------------+
                                  |
                                  V
                            +------------+
                            |    Chat    |
                            +------------+
```

### Documentation
```
GOST R 34.10-2012 public key signature function (RFC 7091)
GOST R 34.11-2012 Streebog hash function (RFC 6986)
GOST R 34.12-2015 128-bit block cipher Kuznechik (RFC 7801)
GOST R 50.1.114-2016 GOST R 34.10-2012 and GOST R 34.11-2012 
RFC 5280: Internet X.509 PKI Certificate Revocation List (CRL)
RFC 7194: Internet Relay Chat (IRC) via TLS
RFC 7539: ChaCha20-Poly1305 AEAD Stream cipher
RFC 8032: Ed25519 Signature a.k.a. EdDSA (Daniel J. Bernstein)
RFC 8446: Transport Layer Security (TLS) Protocol Version 1.3
RFC 9058: MGM AEAD mode for 64 and 128 bit ciphers (E. Griboedova)
```

(TODO)
