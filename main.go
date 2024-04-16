package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var hclOutput = hclwrite.NewEmptyFile()
var rootBody = hclOutput.Body()

type Config struct {
	inputPath      string
	outputPath     string
	merakiAPIKey   string
	organizationID string
}

func getConfig() *Config {
	inputPath := flag.String("infile", "", "Defines input ASA file")
	outputPath := flag.String("outfile", "", "Defines output Terraform file")
	merakiAPIKey := flag.String("api-key", "", "Defines API Key to be used for Meraki")
	organizationID := flag.String("org-id", "", "Defines Organizational ID used to create Infrastructure")
	flag.Parse()

	switch {
	case inputPath == nil:
		log.Fatal("Input path is empty. Please try again.")
	case outputPath == nil:
		log.Fatal("Ouput path is empty. Please try again.")
	case merakiAPIKey == nil:
		fmt.Println("Missing Meraki API Key. Checking OS env variable...")
		if os.Getenv("MERAKI_DASHBOARD_API_KEY") == "" {
			log.Fatal("Unable to find Meraki API key.")
		} else {
			*merakiAPIKey = os.Getenv("MERAKI_DASHBOARD_API_KEY")
		}
	}

	return &Config{
		inputPath:      *inputPath,
		outputPath:     *outputPath,
		merakiAPIKey:   *merakiAPIKey,
		organizationID: *organizationID,
	}
}

func main() {
	config := getConfig()

	/*terraformFile, err := os.Create(config.outputPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer terraformFile.Write(hclOutput.Bytes())*/

	terraformProvider := rootBody.AppendNewBlock("terraform", nil)
	terraformProviderBody := terraformProvider.Body()
	requiredProviders := terraformProviderBody.AppendNewBlock("required_providers", nil)
	requiredProvidersBody := requiredProviders.Body()
	requiredProvidersBody.SetAttributeValue("meraki", cty.ObjectVal(
		map[string]cty.Value{
			"source":  cty.StringVal("cisco-open/meraki"),
			"version": cty.StringVal("0.1.0-alpha"),
		}))
	rootBody.AppendNewline()

	merakiProvider := rootBody.AppendNewBlock("provider", []string{"meraki"})
	merakiProviderBody := merakiProvider.Body()
	merakiProviderBody.SetAttributeValue("meraki_dashboard_api_key", cty.StringVal(config.merakiAPIKey))
	rootBody.AppendNewline()

	GenerateNetworObjectskHCL(config.organizationID, config.inputPath)
	defer fmt.Printf("%s", hclOutput.Bytes())
}