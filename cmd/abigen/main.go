















package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

var (
	
	gitCommit = ""
	gitDate   = ""

	app *cli.App

	
	abiFlag = cli.StringFlag{
		Name:  "abi",
		Usage: "Path to the Ethereum contract ABI json to bind, - for STDIN",
	}
	binFlag = cli.StringFlag{
		Name:  "bin",
		Usage: "Path to the Ethereum contract bytecode (generate deploy method)",
	}
	typeFlag = cli.StringFlag{
		Name:  "type",
		Usage: "Struct name for the binding (default = package name)",
	}
	jsonFlag = cli.StringFlag{
		Name:  "combined-json",
		Usage: "Path to the combined-json file generated by compiler",
	}
	solFlag = cli.StringFlag{
		Name:  "sol",
		Usage: "Path to the Ethereum contract Solidity source to build and bind",
	}
	solcFlag = cli.StringFlag{
		Name:  "solc",
		Usage: "Solidity compiler to use if source builds are requested",
		Value: "solc",
	}
	vyFlag = cli.StringFlag{
		Name:  "vy",
		Usage: "Path to the Ethereum contract Vyper source to build and bind",
	}
	vyperFlag = cli.StringFlag{
		Name:  "vyper",
		Usage: "Vyper compiler to use if source builds are requested",
		Value: "vyper",
	}
	excFlag = cli.StringFlag{
		Name:  "exc",
		Usage: "Comma separated types to exclude from binding",
	}
	pkgFlag = cli.StringFlag{
		Name:  "pkg",
		Usage: "Package name to generate the binding into",
	}
	outFlag = cli.StringFlag{
		Name:  "out",
		Usage: "Output file for the generated binding (default = stdout)",
	}
	langFlag = cli.StringFlag{
		Name:  "lang",
		Usage: "Destination language for the bindings (go, java, objc)",
		Value: "go",
	}
	aliasFlag = cli.StringFlag{
		Name:  "alias",
		Usage: "Comma separated aliases for function and event renaming, e.g. foo=bar",
	}
)

func init() {
	app = flags.NewApp(gitCommit, gitDate, "ethereum checkpoint helper tool")
	app.Flags = []cli.Flag{
		abiFlag,
		binFlag,
		typeFlag,
		jsonFlag,
		solFlag,
		solcFlag,
		vyFlag,
		vyperFlag,
		excFlag,
		pkgFlag,
		outFlag,
		langFlag,
		aliasFlag,
	}
	app.Action = utils.MigrateFlags(abigen)
	cli.CommandHelpTemplate = flags.OriginCommandHelpTemplate
}

func abigen(c *cli.Context) error {
	utils.CheckExclusive(c, abiFlag, jsonFlag, solFlag, vyFlag) 
	if c.GlobalString(pkgFlag.Name) == "" {
		utils.Fatalf("No destination package specified (--pkg)")
	}
	var lang bind.Lang
	switch c.GlobalString(langFlag.Name) {
	case "go":
		lang = bind.LangGo
	case "java":
		lang = bind.LangJava
	case "objc":
		lang = bind.LangObjC
		utils.Fatalf("Objc binding generation is uncompleted")
	default:
		utils.Fatalf("Unsupported destination language \"%s\" (--lang)", c.GlobalString(langFlag.Name))
	}
	
	var (
		abis    []string
		bins    []string
		types   []string
		sigs    []map[string]string
		libs    = make(map[string]string)
		aliases = make(map[string]string)
	)
	if c.GlobalString(abiFlag.Name) != "" {
		
		var (
			abi []byte
			err error
		)
		input := c.GlobalString(abiFlag.Name)
		if input == "-" {
			abi, err = ioutil.ReadAll(os.Stdin)
		} else {
			abi, err = ioutil.ReadFile(input)
		}
		if err != nil {
			utils.Fatalf("Failed to read input ABI: %v", err)
		}
		abis = append(abis, string(abi))

		var bin []byte
		if binFile := c.GlobalString(binFlag.Name); binFile != "" {
			if bin, err = ioutil.ReadFile(binFile); err != nil {
				utils.Fatalf("Failed to read input bytecode: %v", err)
			}
			if strings.Contains(string(bin), "
				utils.Fatalf("Contract has additional library references, please use other mode(e.g. --combined-json) to catch library infos")
			}
		}
		bins = append(bins, string(bin))

		kind := c.GlobalString(typeFlag.Name)
		if kind == "" {
			kind = c.GlobalString(pkgFlag.Name)
		}
		types = append(types, kind)
	} else {
		
		exclude := make(map[string]bool)
		for _, kind := range strings.Split(c.GlobalString(excFlag.Name), ",") {
			exclude[strings.ToLower(kind)] = true
		}
		var err error
		var contracts map[string]*compiler.Contract

		switch {
		case c.GlobalIsSet(solFlag.Name):
			contracts, err = compiler.CompileSolidity(c.GlobalString(solcFlag.Name), c.GlobalString(solFlag.Name))
			if err != nil {
				utils.Fatalf("Failed to build Solidity contract: %v", err)
			}
		case c.GlobalIsSet(vyFlag.Name):
			output, err := compiler.CompileVyper(c.GlobalString(vyperFlag.Name), c.GlobalString(vyFlag.Name))
			if err != nil {
				utils.Fatalf("Failed to build Vyper contract: %v", err)
			}
			contracts = make(map[string]*compiler.Contract)
			for n, contract := range output {
				name := n
				
				
				if !strings.Contains(n, ":") {
					
					name = abi.ToCamelCase(strings.TrimSuffix(filepath.Base(name), ".vy"))
				}
				contracts[name] = contract
			}

		case c.GlobalIsSet(jsonFlag.Name):
			jsonOutput, err := ioutil.ReadFile(c.GlobalString(jsonFlag.Name))
			if err != nil {
				utils.Fatalf("Failed to read combined-json from compiler: %v", err)
			}
			contracts, err = compiler.ParseCombinedJSON(jsonOutput, "", "", "", "")
			if err != nil {
				utils.Fatalf("Failed to read contract information from json output: %v", err)
			}
		}
		
		for name, contract := range contracts {
			if exclude[strings.ToLower(name)] {
				continue
			}
			abi, err := json.Marshal(contract.Info.AbiDefinition) 
			if err != nil {
				utils.Fatalf("Failed to parse ABIs from compiler output: %v", err)
			}
			abis = append(abis, string(abi))
			bins = append(bins, contract.Code)
			sigs = append(sigs, contract.Hashes)
			nameParts := strings.Split(name, ":")
			types = append(types, nameParts[len(nameParts)-1])

			libPattern := crypto.Keccak256Hash([]byte(name)).String()[2:36]
			libs[libPattern] = nameParts[len(nameParts)-1]
		}
	}
	
	if c.GlobalIsSet(aliasFlag.Name) {
		
		
		
		
		re := regexp.MustCompile(`(?:(\w+)[:=](\w+))`)
		submatches := re.FindAllStringSubmatch(c.GlobalString(aliasFlag.Name), -1)
		for _, match := range submatches {
			aliases[match[1]] = match[2]
		}
	}
	
	code, err := bind.Bind(types, abis, bins, sigs, c.GlobalString(pkgFlag.Name), lang, libs, aliases)
	if err != nil {
		utils.Fatalf("Failed to generate ABI binding: %v", err)
	}
	
	if !c.GlobalIsSet(outFlag.Name) {
		fmt.Printf("%s\n", code)
		return nil
	}
	if err := ioutil.WriteFile(c.GlobalString(outFlag.Name), []byte(code), 0600); err != nil {
		utils.Fatalf("Failed to write ABI binding: %v", err)
	}
	return nil
}

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}