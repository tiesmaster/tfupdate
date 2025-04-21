package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

var updateModuleCmd = &cobra.Command{
	Use:   "module MODULE_ID VERSION",
	Short: "Update a Terraform module to a new version",
	RunE:  runUpdateModuleCommand,
	Args:  cobra.ExactArgs(2),
}

var attributeName string

func init() {
	updateCmd.AddCommand(updateModuleCmd)
	updateModuleCmd.Flags().StringVarP(&attributeName, "attribute-name", "n", "version", "attribute name to use for the version")
}

func runUpdateModuleCommand(cmd *cobra.Command, args []string) error {
	targetModule := args[0]
	targetVersion := args[1]

	if verbose {
		fmt.Printf("Will update module '%s' to version '%s'\n", args[0], args[1])
	}

	err := updateModuleForAllFiles(targetModule, targetVersion)
	return err
}

func updateModuleForAllFiles(moduleId, newVersion string) error {
	err := ensureTargetDir()
	if err != nil {
		return err
	}

	tfFiles, err := getTerraformFiles()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("Discovered TF files: ")
		for _, m := range tfFiles {
			fmt.Printf("\t%v\n", m)
		}
	}

	for _, f := range tfFiles {
		err = updateModule(f, moduleId, newVersion)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateModule(filename, moduleId, newVersion string) error {
	return patchFile(filename, func(hclFile *hclwrite.File) (*hclwrite.File, error) {
		mods := getModuleBlocksBySourceForWrite(hclFile.Body(), moduleId)
		for _, m := range mods {
			m.Body().SetAttributeValue(attributeName, cty.StringVal(newVersion))
		}
		return hclFile, nil
	})
}

func isSource(bl *hclwrite.Block, moduleId string) bool {
	srcAttr := bl.Body().GetAttribute("source")
	expr := srcAttr.Expr()
	tokens := expr.BuildTokens(nil)

	moduleSource := string(tokens.Bytes())
	s := strings.TrimSpace(moduleSource)
	s = strings.Trim(s, `"`)
	return s == moduleId
}
