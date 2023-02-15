package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
        "os"
        "time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)
//expecting URL in the form https://localhost:8443/ from SERVER_URL env
var serverURL,serverSpiffeID,socketPath string

func main() {

        serverURL = os.Getenv("SERVER_URL")
        serverSpiffeID = os.Getenv("SERVER_SPIFFE_ID")
        socketPath = os.Getenv("SOCKET_PATH")

        for {
	  if err := run(context.Background()); err != nil {
	  	log.Print(err)
	  }
          time.Sleep(10 * time.Second)
        }
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()


        log.Print("Requesting my SVID..")
        //Note that this step is not a prerequisite for using the workloadapi.NewX509Source function
        //We are only doing it for the purpose of getting and displaying the SVID contents

        svid, err := workloadapi.FetchX509SVID(ctx,workloadapi.WithAddr(socketPath))
	if err != nil {
		log.Printf("unable to retrieve svid: %w", err)
	} else {
            pem, _, err := svid.Marshal()
            if err != nil {
              log.Printf("Unable to marshal X.509 SVID: %v", err)
            }
            log.Printf("Received SVID with ID %q: \n%s\n", svid.ID, string(pem))
        }

	// Create a `workloadapi.X509Source`, it will connect to Workload API using provided socket path
	// If socket path is not defined using `workloadapi.SourceOption`, value from environment variable `SPIFFE_ENDPOINT_SOCKET` is used.
	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		log.Printf("unable to create X509Source: %w", err)
	}
	defer source.Close()

	// Allowed SPIFFE ID
	serverID := spiffeid.RequireFromString(serverSpiffeID)

        log.Print("Requesting web page from " , serverURL)

	// Create a `tls.Config` to allow mTLS connections, and verify that presented certificate has SPIFFE ID `spiffe://example.org/server`
	tlsConfig := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeID(serverID))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	r, err := client.Get(serverURL)
	if err != nil {
		return fmt.Errorf("error connecting to %q: %w", serverURL, err)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read body: %w", err)
	}

	log.Printf("%s", body)
	return nil
}
