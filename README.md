# IRCS Daemon (KDS)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/pedroalbanese/ircs/blob/master/LICENSE.md) 
[![GoDoc](https://godoc.org/github.com/pedroalbanese/ircs?status.png)](http://godoc.org/github.com/pedroalbanese/ircs)
[![GitHub downloads](https://img.shields.io/github/downloads/pedroalbanese/ircs/total.svg?logo=github&logoColor=white)](https://github.com/pedroalbanese/ircs/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/pedroalbanese/ircs)](https://goreportcard.com/report/github.com/pedroalbanese/ircs)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pedroalbanese/ircs)](https://golang.org)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/pedroalbanese/ircs)](https://github.com/pedroalbanese/ircs/releases)
### Internet Relay Chat Secure (Key Distribution System)
Minimalist Internet Relay Chat (IRC) via Transport Layer Security (RFC 7194).

Internet Relay Chat (IRC) is a text-based chat protocol that enables real-time exchange of text messages between users connected to an IRC server. It functions as a network of virtual chat rooms where users can interact with other users in real time. The server acts as a router, ensuring that all messages are delivered to the correct recipients.

This chat application provides a high level of security through robust encryption mechanisms. All communication within the chat rooms is encrypted using industry-standard protocols. Additionally, this chat system does not support asynchronous operations, ensuring that all interactions occur in real time. It does not provide any advantages over other chat applications, except for the ability to enforce usage policies based on digital certification.

```
   +-----------------------+     +----------------------+
   |   Certificate         |     |        Server        |
   |   Authority (CA)      |     |                      |
   |                       |     |                      |
   |   Sign                |     |    Generate          |
   |   CSR                 |     |    CSR               |
   |                       |     |                      |
   V                       V     V                      V
+-------+                 +-------+                  +--------+
| AKID  |                 | AKID  |                  | AKID   |
+-------+                 +-------+                  +--------+
   |                        |                           |
   |                        |                           |
   V                        |                           V
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

## Usage
```
Usage of ircs:
  -cert string
        Certificate file path.
  -crl string
        Certificate revocation list.
  -ipport string
        Server address. (default "localhost:8000")
  -key string
        Private key file path.
  -mode string
        Mode: <server|client> (default "client")
  -pwd string
        Password. (for Private key PEM decryption)
  -strict
        Restrict users.
```

## Examples
This program requires either the [EDGE Toolkit](https://github.com/pedroalbanese/edgetk) or OpenSSL to generate keys and certificates.

#### Asymmetric RSA keypair generation:
```sh
./edgetk -pkey keygen -bits 4096 [-priv private.pem] [-pwd "pass"]
```
#### Generate Self Signed Certificate:
```sh
./edgetk -pkey certgen -key private.pem [-pwd "pass"] [-cert "cacert.pem"]
```
#### Generate Certificate Signing Request:
```sh
./edgetk -pkey req -key private.pem [-pwd "pass"] [-cert certificate.csr]
```
#### Sign CSR with CA Certificate:
```sh
./edgetk -pkey x509 -key private.pem -root cacert.pem -cert certificate.csr > signedcert.crt
```
#### Generate Certificate Revocation List:
```sh
./edgetk -pkey crl -cert cacert.pem -key private.pem -crl old.crl serials.txt > NewCRL.crl
```
## Daemon
### Server
```sh
./ircs -mode server -key private.pem -cert cacert.pem [-strict]
```
### Client
```sh
./ircs -key clientpriv.pem -cert signedcert.crt [-ipport localhost:8000]
```

## Client Commands
There are only four commands for the client to interact with the server:
```
 1. JOIN <room_name>:
        Description: This command allows the user to enter a specific chat room.
        Example: JOIN Chat_Room

 2. LEAVE:
        Description: This command allows the user to exit the current chat room.
        Example: LEAVE

 3. LIST:
        Description: When executed inside a chat room, this command lists the 
        participants currently present in the room. When executed outside of a 
        room, it lists the available chat rooms for the user to choose from.
        Example: LIST

 4. QUIT:
        Description: This command allows the user to exit the chat system entirely.
        Example: QUIT
```

(TODO)
