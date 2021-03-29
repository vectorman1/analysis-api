rm -r certs/*.pem

mkdir certs

# 1. Generate API client's private key and certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout certs/ca-key.pem -out certs/ca-cert.pem -subj "/C=BG/ST=Sofia/L=Sofia/O=dystopia.systems/OU=N\/A/CN=*.dystopia.systems/emailAddress=master@dystopia.systems"

echo "CA's self-signed certificate"
openssl x509 -in certs/ca-cert.pem -noout -text

# 2. Generate GRPC server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout certs/server-key.pem -out certs/server-req.pem -subj "/C=BG/ST=Sofia/L=Sofia/O=dystopia.systems/OU=N\/A/CN=*.dystopia.systems/emailAddress=master@dystopia.systems"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in certs/server-req.pem -days 365 -CA certs/ca-cert.pem -CAkey certs/ca-key.pem -CAcreateserial -out certs/server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in certs/server-cert.pem -noout -text
