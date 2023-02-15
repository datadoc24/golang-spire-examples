# golang-spire-examples
Example Spire worloads to accompany article in Admin Magazine

Based on the mTLS examples at https://github.com/spiffe/go-spiffe/blob/main/v2/examples/spiffe-tls/README.md. The client program sends a request to the server every few seconds and logs the result, so you can easily see the impact of the changes you are making to your SPIRE infrastructure. Both programs display their SVID bundle information, which you can decode with openssl x509 -text -in <pem contents> to see the SVID's details.
