# Internet Relay Chat Secure
Minimalist and Simplified Internet Relay Chat (IRC) via TLS (RFC 7194)

```
   +-----------------------+       +---------------------+
   |   Certificado de      |       |      Servidor       |
   |   Autoridade (CA)     |       |                     |
   |                       |       |                     |
   |       Assina          |       |       Gera          |
   |       CSR             |       |       CSR           |
   |                       |       |                     |
   V                       V       V                     V
 +------+                 +-------+                 +--------+
 | AKID |                 | AKID  |                 | AKID   |
 +------+                 +-------+                 +--------+
   |                        |                        |
   |                        |                        |
   V                        |                        V
 Chat Acessível             |                Acesso Negado
 a partir do AKID           |                Sem AKID válido
                            |
                            V
                      +------------+
                      |    Chat    |
                      +------------+
```
