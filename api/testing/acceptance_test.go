package testing

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

//go:embed testdata/http/**/*.json
var embedHttpTestdata embed.FS

type testRunner struct {
	serverAddress string
	resp          *http.Response
	reqBody       io.Reader
}

var opts = &godog.Options{
	Concurrency: 6,
	Randomize:   -1,
	Format:      "pretty",
	Paths:       []string{"features"},
	Output:      colors.Colored(os.Stdout),
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, opts)
}

func TestFeatures(t *testing.T) {
	if opts.Tags == "" {
		t.Skip()
		return
	}
	opts.TestingT = t
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func (t *testRunner) serverAddressIs(address string) error {
	t.serverAddress = address
	return nil
}

func (t *testRunner) requestBodyIs(filename string) error {
	filename = filepath.Join("testdata", "http", "input", filename+".json")
	fileBytes, err := embedHttpTestdata.ReadFile(filename)
	if err != nil {
		return err
	}
	t.reqBody = bytes.NewReader(fileBytes)
	return nil
}

func (t *testRunner) iSendrequestTo(method, endpoint string) (err error) {
	url := fmt.Sprintf("%s%s", t.serverAddress, endpoint)
	req, err := http.NewRequest(method, url, t.reqBody)
	if err != nil {
		return err
	}

	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 1 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 2 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Second * 3,
		Transport: netTransport,
	}

	t.resp, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (t *testRunner) theResponseCodeShouldBe(code int) error {
	if code != t.resp.StatusCode {
		return fmt.Errorf("http status code mismatch, want: %d, got: %d", code, t.resp.StatusCode)
	}
	return nil
}

type StringAnyMap map[string]any

// wtf
type gotJSON struct {
	PublicID string `json:"publicId,omitempty"`
	StringAnyMap
}

func (t *testRunner) theResponseShouldMatchJSON(filename string) error {
	var want, got gotJSON
	filename = filepath.Join("testdata", "http", "output", filename+".json")
	fileBytes, err := embedHttpTestdata.ReadFile(filename)
	if err != nil {
		return err
	}

	// re-encode expected response
	if err = json.Unmarshal(fileBytes, &want); err != nil {
		return err
	}
	defer t.resp.Body.Close()
	bodyBytes, err := io.ReadAll(t.resp.Body)
	if err != nil {
		return err
	}

	// re-encode actual response too
	if err = json.Unmarshal(bodyBytes, &got); err != nil {
		return err
	}
	opts := cmpopts.IgnoreFields(gotJSON{}, "PublicID")

	// the matching may be adapted per different requirements.
	if diff := cmp.Diff(got, want, opts); diff != "" {
		return fmt.Errorf("\nresponseBody mismatch got(-) want(+)\n%s", diff)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	t := &testRunner{}
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		t = &testRunner{}
		return ctx, nil
	})
	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, t.iSendrequestTo)
	ctx.Step(`^the server address is "([^"]*)"$`, t.serverAddressIs)
	ctx.Step(`^the request body is "([^"]*)"$`, t.requestBodyIs)
	ctx.Step(`^the response code should be (\d+)$`, t.theResponseCodeShouldBe)
	ctx.Step(`^the response body should match "([^"]*)"$`, t.theResponseShouldMatchJSON)
}
