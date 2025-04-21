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
	targetLocalName := args[0]
	targetVersion := args[1]

	if verbose {
		fmt.Printf("Will update provider with local name '%s' to version '%s'\n", targetLocalName, targetVersion)
	}

	// TODO: Find "versions.tf" dynamicly
	err := updateProvider("versions.tf", targetLocalName, targetVersion)
	return err
}

func updateProvider(filename, localName, newVersion string) error {
	return patchFile(filename, func(hclFile *hclwrite.File) (*hclwrite.File, error) {

		bl, err := getBlockByTypeForWrite(hclFile.Body(), "terraform")
		if err != nil {
			return nil, err
		}

		bl, err = getBlockByTypeForWrite(hclFile.Body(), "required_providers")
		if err != nil {
			return nil, err
		}

		providerAttribute := bl.Body().GetAttribute(localName)
		if providerAttribute == nil {
			return nil, fmt.Errorf("Cannot find provider for localName '%s'", localName)
		}

		headlessBlock := providerAttribute.BuildTokens(nil).Bytes()
		fmt.Println(headlessBlock)

		bl.Body().SetAttributeValue("version", cty.StringVal(newVersion))
		return hclFile, nil
	})
}
