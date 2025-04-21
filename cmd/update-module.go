package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

var updateModuleCmd = &cobra.Command{
	Use:   "update-module MODULE_ID VERSION",
	Short: "Update module",
	RunE:  runUpdateModuleCommand,
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(updateModuleCmd)
}

func runUpdateModuleCommand(cmd *cobra.Command, args []string) error {
	targetModule := args[0]
	targetVersion := args[1]

	fmt.Printf("Will update module '%s' to version '%s'\n", args[0], args[1])

	err := updateModuleForAllFiles(targetModule, targetVersion)
	return err
}

func updateModuleForAllFiles(moduleId, newVersion string) error {
	tfFiles, err := getTerraformFiles()
	if err != nil {
		return err
	}

	for _, f := range tfFiles {
		updateModule(f, moduleId, newVersion)
	}

	return nil
}

func getTerraformFiles() ([]string, error) {
	dir := os.DirFS(".")
	matches, err := fs.Glob(dir, "*.tf")
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, errors.New("no TF files detected")
	}

	return matches, nil
}

func updateModule(filename, moduleId, newVersion string) error {
	return patchFile(filename, func(hclFile *hclwrite.File) (*hclwrite.File, error) {
		mods := getModuleBlocksBySourceForWrite(hclFile, moduleId)
		for _, m := range mods {
			m.Body().SetAttributeValue("version", cty.StringVal(newVersion))
		}
		return hclFile, nil
	})
}

func patchFile(filename string, patch func(hclFile *hclwrite.File) (*hclwrite.File, error)) error {
	input, _ := os.ReadFile(filename)
	hclFile, diags := hclwrite.ParseConfig(input, filename, hcl.Pos{Line: 1, Column: 1})

	if diags.HasErrors() {
		return errors.New("failed to parse TF file: " + diags.Error())
	}

	newHclFile, err := patch(hclFile)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filename, newHclFile.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}

func getModuleBlocksBySourceForWrite(hclFile *hclwrite.File, moduleId string) []*hclwrite.Block {
	var modBlocks []*hclwrite.Block
	for _, bl := range hclFile.Body().Blocks() {
		if bl.Type() == "module" && isSource(bl, moduleId) {
			modBlocks = append(modBlocks, bl)
		}
	}
	return modBlocks
}

func isSource(bl *hclwrite.Block, moduleId string) bool {
	srcAttr := bl.Body().GetAttribute("source")
	expr := srcAttr.Expr()
	tokens := expr.BuildTokens(nil)
	for _, t := range tokens {
		fmt.Println(t)
	}

	moduleSource := getStringValue(expr)
	moduleSource = string(tokens.Bytes())
	s := strings.TrimSpace(moduleSource)
	s = strings.Trim(s, `"`)
	return s == moduleId
}

func getStringValue(expr *hclwrite.Expression) string {
	for _, t := range expr.BuildTokens(nil) {
		if t.Type == hclsyntax.TokenQuotedLit {
			return string(t.Bytes)
		}
	}

	panic("cannot reach")
	// or just return error
}

// func getStringValue2(expr *hclwrite.Expression) string {
// 	trav, diags := hcl.AbsTraversalForExpr(nil)
// 	for _, tr := range trav {
// 		tr.

// 	}
// }
