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

Unlike other IRCs, this application does not generate logs and does not store messages on the hard drive. This means that the sender of the messages can be absolutely certain that only those present in the room will read the message, and it will no longer be accessible thereafter.

```
   +-----------------------+     +----------------------+
   |   Certificate         |     |        Server        |
   |   Authority (CA)      |     |                      |
   |                       |     |                      |
   |   Sign                |     |    Generate          |
   |   CSR                 |     |    CSR               |
   |                       |     |                      |
   V                       V     V                      V
+-------+                 +-------+                 +--------+
| AKID  |                 | AKID  |                 | AKID   |
+-------+                 +-------+                 +--------+
   |                        |                           |
   |                        |                           |
   V                        |                           V
 Client                     |                     Access Denied
 with Certificate           |                     No valid AKID
   |                        |
   |                        V
   |                  +------------+
   +----------------> |  CRL       |
                      |  Check     |
                      +------------+
                            |
                            |     +---------------------+
                            |     |    Revocation List  |
                            +---->|                     |
                                  |    Not after XXXX   |
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

## Q Code

The Q Code is a three-letter abbreviation system used in radio communications to transmit messages more efficiently and concisely. It was widely used by amateur radio operators as well as professional and military radio operators.

The Q Code system was originally developed for use in radio telegraph communications, where message transmission can be slow and costly. It was created to standardize and simplify communications by allowing common information to be transmitted quickly using short combinations of letters.

Each Q code starts with the letter "Q" followed by two additional letters, forming a three-letter combination. Each code has a specific meaning associated with it. For example, QTH is used to inquire about location, QSL is used to confirm receipt of a message, and QRM is used to indicate the presence of radio interference.

Although the Q Code was widely used in the past, its popularity has waned over time due to the advancement of technology and the development of more sophisticated communication systems. However, it is still recognized and understood by many radio operators, especially in emergency situations or when communicating with more experienced operators.

### Codes
```
QRA: What is the name of your station?
    Description: It's used to ask the other station about the name of their station.
    Example: QRA

QRB: What is the distance to my station?
    Description: It's used to ask the other station about the distance to your station.
    Example: QRB

QRG: What is my exact frequency?
    Description: It's used to ask the other station about the exact frequency of your transmission.
    Example: QRG

QRH: Does my frequency vary?
    Description: It's used to ask the other station if your frequency is varying.
    Example: QRH

QRI: What is the tone of my transmission?
    Description: It's used to ask the other station about the tone of your transmission.
    Example: QRI

QRJ: Can you copy me well?
    Description: It's used to ask the other station if they can copy your signals well.
    Example: QRJ

QRK: What is the readability of my signals?
    Description: It's used to ask the other station about the readability of your signals.
    Example: QRK

QRL: Are you busy?
    Description: It's used to ask the other station if they are busy.
    Example: QRL

QRM: Are you being interfered with?
    Description: It's used to ask the other station if they are experiencing interference.
    Example: QRM

QRN: Are you troubled by static noise?
    Description: It's used to ask the other station if they are experiencing static noise.
    Example: QRN

QRO: Shall I increase power?
    Description: It's used to ask the other station if you should increase your transmission power.
    Example: QRO

QRP: Shall I decrease power?
    Description: It's used to ask the other station if you should decrease your transmission power.
    Example: QRP

QRQ: Shall I send faster?
    Description: It's used to ask the other station if you should increase the sending speed.
    Example: QRQ

QRS: Shall I send slower?
    Description: It's used to ask the other station if you should decrease the sending speed.
    Example: QRS

QRT: Shall I stop transmitting?
    Description: It's used to ask the other station if you should stop transmitting.
    Example: QRT

QRU: Do you have anything for me?
    Description: It's used to ask the other station if they have any messages or information for you.
    Example: QRU

QRZ: Who is calling me?
    Description: It's used to ask the other station who is trying to establish communication.
    Example: QRZ

QSA: What is the strength of my signals?
    Description: It's used to ask the other station about the strength of your signals.
    Example: QSA

QSB: Are my signals fading?
    Description: It's used to ask the other station if your signals are fading.
    Example: QSB

QSD: Is my keying defective?
    Description: It's used to ask the other station if your keying is defective.
    Example: QSD

QSG: Shall I send __ messages at a time?
    Description: It's used to ask the other station if you should send a specific number of messages at a time.
    Example: QSG <number of messages>

QSK: Can you hear me between your signals?
    Description: It's used to ask the other station if they can hear you between their signals.
    Example: QSK

QSL: Can you acknowledge receipt?
    Description: It's used to ask the other station if they can acknowledge receipt of your message.
    Example: QSL

QSM: Shall I repeat the last message?
    Description: It's used to ask the other station if you should repeat the last message.
    Example: QSM

QSN: Did you hear me?
    Description: It's used to ask the other station if they heard your transmission.
    Example: QSN

QSO: Can you communicate with __ direct?
    Description: It's used to ask the other station if they can communicate with a specific station directly.
    Example: QSO <station name>

QSP: Will you relay a message to __?
    Description: It's used to ask the other station if they will relay a message to a specific station.
    Example: QSP <station name>?

QSR: Do you want me to repeat my call?
    Description: It's used to ask the other station if they want you to repeat your call.
    Example: QSR

QSS: Shall I send slower?
    Description: It's used to ask the other station if you should send at a slower speed.
    Example: QSS

QSU: Shall I send more slowly?
    Description: It's used to ask the other station if you should send at an even slower speed.
    Example: QSU

QSV: Shall I send __ long dashes?
    Description: It's used to ask the other station if you should send a specific number of long dashes.
    Example: QSV <number of long dashes>

QSX: Will you listen on __ frequency?
    Description: It's used to ask the other station if they will listen on a specific frequency.
    Example: QSX <frequency>

QSY: Shall I change frequency?
    Description: It's used to ask the other station if you should change frequency.
    Example: QSY

QSZ: Shall I send each word or group more than once?
    Description: It's used to ask the other station if you should send each word or group more than once.
    Example: QSZ

QTC: How many messages have you to send?
    Description: It's used to ask the other station how many messages they have to send.
    Example: QTC

QTH: What is your location?
    Description: It's used to ask the other station about their location.
    Example: QTH

QTR: What is the time?
    Description: It's used to ask the other station about the current time.
    Example: QTR
```
(TODO)
