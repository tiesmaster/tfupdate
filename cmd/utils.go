package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ensureTargetDir() error {
	if targetDir != "" {
		err := os.Chdir(targetDir)
		if err != nil {
			return err
		}
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

func getModuleBlocksBySourceForWrite(hclBody *hclwrite.Body, moduleId string) []*hclwrite.Block {
	var modBlocks []*hclwrite.Block
	for _, bl := range hclBody.Blocks() {
		if bl.Type() == "module" && isSource(bl, moduleId) {
			modBlocks = append(modBlocks, bl)
		}
	}
	return modBlocks
}

func getBlockByTypeForWrite(hclBody *hclwrite.Body, typeId string) (*hclwrite.Block, error) {
	for _, bl := range hclBody.Blocks() {
		if bl.Type() == typeId {
			return bl, nil
		}
	}

	return nil, errors.New("Cannot find block with type " + typeId)
}

func parseSourceAddressToLocalName(sourceAddress string) (string, error) {
	s := strings.Split(sourceAddress, "/")
	if len(s) != 2 {
		return "", errors.New("Unable to parse source address " + sourceAddress + "as valid Terraform provider")
	}

	return s[1], nil
}
