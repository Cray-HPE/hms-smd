/*
 * MIT License
 *
 * (C) Copyright [2018-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	hmshttp "stash.us.cray.com/HMS/hms-go-http-lib"
	"syscall"
	"time"
)

var ctx context.Context
var baseRequest hmshttp.HTTPRequest

func readFile(file string) (jsonBytes []byte) {
	// Open and parse the file.
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	jsonBytes, _ = ioutil.ReadAll(jsonFile)
	jsonString := string(jsonBytes)

	_ = jsonFile.Close()

	if jsonString == "" {
		panic("JSON file empty...are you sure you created the ConfigMap with the file " + file + "?\n")
	}

	fmt.Printf("%s file contents:\n%s\n", file, string(jsonString))

	return
}

func nodeNids(hsmURL string) {
	nodeNidMapFile, nodeNidMapOk := os.LookupEnv("NODE_NID_MAP_FILE")
	if !nodeNidMapOk {
		panic("Value not set for NODE_NID_MAP_FILE")
	}

	fmt.Println("Deleting existing NodeMaps collection...")
	deleteRequest := baseRequest
	deleteRequest.FullURL = hsmURL + "/hsm/v1/Defaults/NodeMaps"
	deleteRequest.Method = "DELETE"

	// Don't bother checking whether the delete worked.
	_, _, _ = deleteRequest.DoHTTPAction()

	// Now do the upload
	fmt.Printf("Uploading NodeMaps file (%s) to HSM (%s)...\n", nodeNidMapFile, hsmURL)

	jsonBytes := readFile(nodeNidMapFile)

	postRequest := baseRequest
	postRequest.FullURL = hsmURL + "/hsm/v1/Defaults/NodeMaps"
	postRequest.Method = "POST"
	postRequest.Payload = jsonBytes

	postResponsePayload, postStatusCode, postErr := postRequest.DoHTTPAction()
	if postErr != nil {
		errorString := fmt.Sprintf("\n\nLoader FAILED (%d), err: %s\n", postStatusCode, postErr)
		if postResponsePayload != nil {
			errorString += fmt.Sprintf("Response body:\n%s", postResponsePayload)
		}
		panic(errorString)
	}

	fmt.Println("\n\nSMD successfully loaded with above map.")
}

func main() {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c

		// Cancel the context to cancel any in progress HTTP requests.
		cancel()
	}()

	hsmURL, hsmURLOk := os.LookupEnv("HSM_URL")
	if !hsmURLOk {
		panic("Value not set for HSM_URL")
	}

	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 100

	baseRequest = hmshttp.HTTPRequest{
		Client:              httpClient,
		Context:             ctx,
		Timeout:             10 * time.Second,
		ExpectedStatusCodes: []int{http.StatusOK},
		ContentType:         "application/json",
		CustomHeaders:       make(map[string]string),
	}

	// Set a custom header on the base request so we can later identify all connections coming from this loader.
	baseRequest.CustomHeaders["HMS-Service"] = "smd-loader"

	// Upload NodeMaps to HSM to set default NIDs to be used for discovery.
	nodeNids(hsmURL)
}
