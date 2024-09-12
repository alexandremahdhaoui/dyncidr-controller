package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/alexandremahdhaoui/dyncidr-controller/pkg/apis/v1alpha1"
	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/yaml"
)

// ------------------------------------ CONFIG ------------------------------------ //

const (
	configEnvKey = "DYNCIDR_CONTROLLER_CONFIG_PATH"
)

var defaultConfigPath = []string{
	"./dynccidr-controller.yaml",
	"/etc/dynccidr-controller.yaml",
}

// ------------------------------------ MAIN ------------------------------------ //

// What is the code suppose to do?
// - 1. Get a public ipv6 of one of the node/router.
// - 2. Parse /64 subnet.
// - 3. Compute set of CIDR ranges in that subnet for LB IPAM. (Based on config)
// - 4. Update CRDs with new CIDR ranges. (Based on config)
// - OPTIONAL: Delete & Recreate service type loadbalancers associated with those CIDR ranges.

// TODO: move the config to a CRD.

func main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ipNet, err := GetV6Subnet(config.IPV6Lookup.URL, config.IPV6Lookup.JSONPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(ipNet.String())
}

func AllocateCIDR() {}
func UpdateCRDs()   {}

// ------------------------------------ GET IPV6 ------------------------------------ //

var errGettingIPV6 = errors.New("getting IPV6")

func GetV6Subnet(url string, jp *jsonpath.JSONPath) (*net.IPNet, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Join(err, errGettingIPV6)
	}

	if code := resp.StatusCode; code != 200 {
		return nil, errors.Join(fmt.Errorf("received status code: %d", code), errGettingIPV6)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(err, errGettingIPV6)
	}

	if err = resp.Body.Close(); err != nil {
		return nil, errors.Join(err, errGettingIPV6)
	}

	var v any
	if err = json.Unmarshal(b, &v); err != nil {
		return nil, errors.Join(err, errGettingIPV6)
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	if err = jp.Execute(buf, v); err != nil {
		return nil, errors.Join(err, errGettingIPV6)
	}

	if ip := net.ParseIP(buf.String()); ip == nil {
		return nil, errors.Join(fmt.Errorf("cannot parse ip: %s", buf.String()), errGettingIPV6)
	}

	_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/64", buf.String()))
	if err != nil {
		return nil, errors.Join(err, errGettingIPV6)
	}

	return ipNet, nil
}

// ------------------------------------ READ CONFIG ------------------------------------ //

var errReadingConfig = errors.New("reading config")

func readConfig() (v1alpha1.DynamicNetwork, error) {
	out := v1alpha1.DynamicNetwork{}

	fp := os.Getenv(configEnvKey)
	if fp == "" {
		for _, s := range defaultConfigPath {
			if _, err := os.Stat(s); err != nil {
				continue
			}

			fp = s
			break
		}
	}

	b, err := os.ReadFile(fp)
	if err != nil {
		return v1alpha1.DynamicNetwork{}, errors.Join(err, errReadingConfig)
	}

	if err := yaml.Unmarshal(b, &out); err != nil {
		return v1alpha1.DynamicNetwork{}, errors.Join(err, errReadingConfig)
	}

	return out, nil
}
