package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
	"github.com/xuri/excelize/v2"
)

type Wallet struct {
	Address    string
	Mnemonic   string
	PrivateKey string
}

const ProgramName = "hdwalletCli"

var (
	filename string
	number   int
)

func Cmd(programName string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   programName,
		Short: fmt.Sprintf("generate %s program", programName),
		Long:  fmt.Sprintf("generate %s program with address & privateKey & Mnemonic", programName),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Program %s params with filename=[%s], number=[%d]\n", programName, filename, number)
			if len(args) > 0 {
				return fmt.Errorf("unexpected args are met: %v", args)
			}
			isExist := exists(filename + ".xlsx")
			if isExist {
				return errors.New(filename + ".xlsx is Exist")
			}
			wallets := []*Wallet{}
			for i := 0; i < number; i++ {
				walt, err := newWallet()
				if err != nil {
					log.Panicln(err)
				}
				wallets = append(wallets, walt)

			}
			if len(wallets) == 0 {
				log.Panicln("wallet is empty")
			}
			err := saveExcel(wallets, filename)
			if err != nil {
				return err
			}
			fmt.Println(filename + ".xlsx wallet file save succ")
			return nil
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "mywallet", "Wallet Excel Filename")
	cmd.Flags().IntVarP(&number, "number", "n", 10, "Wallet Address Number")

	return cmd
}

func newWallet() (*Wallet, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, "")
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	address := account.Address.Hex()
	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return nil, err
	}

	walt := &Wallet{
		Address:    address,
		Mnemonic:   mnemonic,
		PrivateKey: privateKey,
	}
	return walt, nil
}

func saveExcel(wallets []*Wallet, filename string) error {
	f := excelize.NewFile()
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#fc5531"},
	})
	if err != nil {
		return err
	}

	index := f.NewSheet("Sheet1")
	// header
	f.SetCellValue("Sheet1", "A1", "Address")
	f.SetCellValue("Sheet1", "B1", "PrivateKey")
	f.SetCellValue("Sheet1", "C1", "Mnemonic")
	f.SetCellValue("Sheet1", "D1", "Remark")
	f.SetCellValue("Sheet1", "E1", "【审阅】-【保护工作表】")
	f.SetCellStyle("Sheet1", "A1", "E1", style)

	for k, v := range wallets {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", k+2), v.Address)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", k+2), v.PrivateKey)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", k+2), v.Mnemonic)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", k+2), "")
	}

	f.SetActiveSheet(index)
	err = f.SaveAs(filename + ".xlsx")
	if err != nil {
		return err
	}
	return err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func main() {
	cmd := Cmd(ProgramName)
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}
