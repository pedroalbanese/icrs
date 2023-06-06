# Internet Relay Chat Secure
Minimalist and Simplified Internet Relay Chat (IRC) via TLS (RFC 7194)

Internet Relay Chat (IRC) is a text-based chat protocol that enables real-time exchange of text messages between users connected to an IRC server. It functions as a network of virtual chat rooms where users can chat in public or private channels, create and manage their own channels, and interact with other users in real time. The server acts as a router, ensuring that all messages are delivered to the correct recipients.

```
   +-----------------------+       +---------------------+
   |   Certificado de      |       |      Servidor       |
   |   Autoridade (CA)     |       |                     |
   |                       |       |                     |
   |       Assina          |       |       Gera          |
   |       CSR             |       |       CSR           |
   |                       |       |                     |
   V                       V       V                     V
 +------+                 +-------+                  +--------+
 | AKID |                 | AKID  |                  | AKID   |
 +------+                 +-------+                  +--------+
   |                        |                            |
   |                        |                            |
   V                        |                            V
 Chat Acessível             |                    Acesso Negado
 a partir do AKID           |                    Sem AKID válido
                            |
                            V
                      +------------+
                      |    Chat    |
                      +------------+
```


