package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
        "os"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

// SPIFFE ID to allow requests from
var clientSpiffeID,socketPath string

func main() {

        clientSpiffeID = os.Getenv("CLIENT_SPIFFE_ID")
        socketPath = os.Getenv("SOCKET_PATH")

	if err := run(context.Background()); err != nil {
		log.Println(err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up a `/` resource handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received")
		_, _ = io.WriteString(w, "Top secret sensitive data from server workload!!!")
	})


	// Create a `workloadapi.X509Source`, it will connect to Workload API using provided socket.
	// If socket path is not defined using `workloadapi.SourceOption`, value from environment variable `SPIFFE_ENDPOINT_SOCKET` is used.
	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		log.Printf("unable to create X509Source: %w", err)
	}
	defer source.Close()

        //display the svid id and pem file of the source that we got from the workload api
        svid, err := source.GetX509SVID()

        if err != nil {
             log.Printf("Unable to retrieve X.509: %w", err)
        } else {
           pem, _, err := svid.Marshal()
            if err != nil {
              log.Printf("Unable to marshal X.509 SVID: %v", err)
            }
            log.Printf("Received X.509 SVID with ID %q: \n%s\n", svid.ID, string(pem))
        }

	// Allowed SPIFFE ID
	clientID := spiffeid.RequireFromString(clientSpiffeID)
        log.Println("Server listening for connections from", clientSpiffeID)

	// Create a `tls.Config` to allow mTLS connections, and verify that presented certificate has expected SPIFFE ID
	tlsConfig := tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeID(clientID))
	server := &http.Server{
		Addr:              ":8443",
		TLSConfig:         tlsConfig,
		ReadHeaderTimeout: time.Second * 10,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Printf("failed to serve: %w", err)
	}
	return nil
}
