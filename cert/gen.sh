rm *.pem

openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=AR/ST=Corrientes/L=Corrientes/O=DEV/OU= TUTORIAL/CN=*.apisgo.com.ar/emailAddress=soporte.divcor@gmail.com" 

openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=AR/ST=Corrientes/L=Corrientes/O=DEV/OU=BLOG/CN=*.apisgo.com.ar/emailAddress=soporte.divcor@gmail.com"

openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.conf

openssl x509 -in server-cert.pem -noout -text