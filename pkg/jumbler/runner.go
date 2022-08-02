package jumbler

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/remydewolf/jumbler/pkg/config"
)

type Mode string

const (
	ModeEncrypt Mode = "encrypt"
	ModeDecrypt Mode = "decrypt"
	JumblerExt       = ".jmb"
)

func GetMode(s string) (Mode, error) {
	mode := Mode(s)
	if mode == ModeEncrypt || mode == ModeDecrypt {
		return mode, nil
	}
	return mode, fmt.Errorf("invalid mode, expect %v or %v", ModeEncrypt, ModeDecrypt)
}

func Run(ctx context.Context, cfg config.Config, mode Mode) error {

	if cfg.DryRun {
		fmt.Println("Running with dry run, no change will be made")
	}

	var files []string
	err := filepath.Walk(cfg.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				//ignore directory, only rename files
				return nil
			}
			//never encrypt these files
			if filepath.Base(path) == "jumbler" {
				return nil
			}
			if mode == ModeEncrypt && (cfg.Ext == "" || cfg.Ext == filepath.Ext(path)) {
				files = append(files, path)
			}
			if mode == ModeDecrypt && filepath.Ext(path) == JumblerExt {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no file to consider")
	}
	fmt.Printf("Found %v file(s):\n", len(files))
	for _, f := range files {
		if !cfg.Quiet && mode == ModeEncrypt {
			fmt.Println(f)
		}
	}
	if !cfg.DryRun && !cfg.AutoApprove {
		fmt.Printf("Press 'Enter' to confirm and %v these file names ... (%v total)", string(mode), len(files))
		if _, err = bufio.NewReader(os.Stdin).ReadBytes('\n'); err != nil {
			return err
		}
	}

	cryptoFuncs := map[Mode]func(string, string) (string, error){
		ModeEncrypt: encryptFilename,
		ModeDecrypt: decryptFilename,
	}

	//proceed with encrypt or decript the filename for each file
	start := time.Now()
	for _, file := range files {
		newFile, err := cryptoFuncs[mode](file, cfg.Password)
		if err != nil {
			return err
		}
		if !cfg.Quiet {
			fmt.Printf("[%v] %v to %v\n", mode, file, newFile)
		}
		if !cfg.DryRun {
			if err := os.Rename(file, newFile); err != nil {
				if strings.Contains(err.Error(), "file name too long") {
					//ignore this file
					continue
				}
				return err
			}
		}
	}
	duration := time.Since(start)
	if cfg.DryRun {
		fmt.Printf("Would have updated %v files\n", len(files))
	} else {
		fmt.Printf("%ved %v files in %v\n", mode, len(files), duration)
	}
	return nil
}

//hash your passphrase using a hashing algorithm that produces 32 byte hashes.
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
