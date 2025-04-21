package cmd

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

var updateProviderCmd = &cobra.Command{
	Use:   "provider SOURCE_ADDRESS VERSION",
	Short: "Update a Terraform provider to a new version",
	RunE:  runUpdateProviderCommand,
	Args:  cobra.ExactArgs(2),
}

func init() {
	updateCmd.AddCommand(updateProviderCmd)
}

func runUpdateProviderCommand(cmd *cobra.Command, args []string) error {
	targetProviderAddress := args[0]
	targetVersion := args[1]
	localName, err := parseSourceAddressToLocalName(targetProviderAddress)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Will update provider '%s' ('%s') to version '%s'\n", targetProviderAddress, localName, targetVersion)
	}

	// TODO: Find "versions.tf" dynamicly
	err = updateProvider("versions.tf", localName, targetProviderAddress, targetVersion)
	return err
}

func updateProvider(filename, localName, providerAddress, newVersion string) error {
	return patchFile(filename, func(hclFile *hclwrite.File) (*hclwrite.File, error) {

		bl, _ := getBlockByTypeForWrite(hclFile.Body(), "terraform")

		bl.Body().SetAttributeValue("version", cty.StringVal(newVersion))
		return hclFile, nil
	})
}
