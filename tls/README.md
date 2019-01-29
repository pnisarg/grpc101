#### Generate CA key
` $openssl genrsa -out ca.key 2048 `

#### Generate self-signed certificate for CA
` $openssl req -x509 -new -sha256 -key ca.key -days 3650 -out ca.crt `

#### Generate key for server
`$openssl genrsa -out server.key 2048`

#### Generate CSR 
`$openssl req -new -key server.key -out server.csr`

#### Create config file needed to define the Subject Alternative Name (SAN) extension 
`$vi server.ext`
```
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = 127.0.0.1

```
#### Create signed certficate for server
`$openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 -extfile server.ext`


#### Generate self-signed certificate for server. 
We will use this to test our setup. Start server with self-signed certificate and make a request from client. You should see
`transport: authentication handshake failed: x509: certificate signed by unknown authority`
```
$openssl req -new -key server.key -out server1.csr
$openssl x509 -req -sha256 -in server1.csr -signkey server.key -out server1.crt -days 365
```