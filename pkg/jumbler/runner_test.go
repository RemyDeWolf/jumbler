package jumbler

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"
	"github.com/remydewolf/jumbler/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestGetMode(t *testing.T) {
	mode, err := GetMode("encrypt")
	require.NoError(t, err)
	require.Equal(t, mode, ModeEncrypt)

	mode, err = GetMode("decrypt")
	require.NoError(t, err)
	require.Equal(t, mode, ModeDecrypt)

	_, err = GetMode("unknown")
	require.ErrorContains(t, err, "invalid mode, expect encrypt or decrypt")
}

func TestRun(t *testing.T) {
	// copy all the test data to a temp directory
	tmpDir := t.TempDir()
	require.NoError(t, cp.Copy("./test-data", tmpDir))

	cfg := config.Config{
		Path:        tmpDir,
		Password:    "correct-pwd",
		AutoApprove: true,
	}

	isEncrypted := func(t *testing.T, name string) {
		ext := filepath.Ext(name)
		require.Equal(t, JumblerExt, ext)
	}
	isClear := func(t *testing.T, name string) {
		ext := filepath.Ext(name)
		require.NotEqual(t, JumblerExt, ext)
	}

	//run and check filenames have been encrypted
	require.NoError(t, Run(context.Background(), cfg, ModeEncrypt))
	checkFiles(t, tmpDir, isEncrypted)

	//try to decrypt with the wrong password
	require.ErrorContains(t, Run(
		context.Background(),
		config.Config{
			Path:        tmpDir,
			Password:    "wrong-pwd",
			AutoApprove: true,
		},
		ModeDecrypt,
	), "cipher: message authentication failed")

	//run and check filenames have been decrypted
	require.NoError(t, Run(context.Background(), cfg, ModeDecrypt))
	checkFiles(t, tmpDir, isClear)

	//run in dry mode check nothing is changed
	require.NoError(t, Run(
		context.Background(),
		config.Config{
			Path:        tmpDir,
			Password:    "correct-pwd",
			AutoApprove: true,
			DryRun:      true,
		}, ModeEncrypt))
	checkFiles(t, tmpDir, isClear)

	//only encrypt for a given extension
	require.NoError(t, Run(
		context.Background(),
		config.Config{
			Path:        tmpDir,
			Password:    "correct-pwd",
			AutoApprove: true,
			Ext:         ".png",
		},
		ModeEncrypt,
	))
	isPngEncrypted := func(t *testing.T, name string) {
		ext := filepath.Ext(name)
		if ext == ".png" {
			t.Errorf("%v shoud have been encrypted", name)
		}
		//allow other extensions
		require.Contains(t, []string{".pdf", JumblerExt}, ext)
	}
	checkFiles(t, tmpDir, isPngEncrypted)

	//only decrypt for a given extension
	require.NoError(t, Run(
		context.Background(),
		config.Config{
			Path:        tmpDir,
			Password:    "correct-pwd",
			AutoApprove: true,
			Ext:         ".png",
		},
		ModeDecrypt,
	))
	isPngDecrypted := func(t *testing.T, name string) {
		ext := filepath.Ext(name)
		//allow only decrypted files
		require.Contains(t, []string{".pdf", ".png"}, ext)
	}
	checkFiles(t, tmpDir, isPngDecrypted)

	// test that the "jumbler" file is never encrypted
	tmpDir = t.TempDir()
	jmblrFile, err := os.Create(path.Join(tmpDir, "jumbler"))
	require.NoError(t, err)
	require.ErrorContains(t, Run(
		context.Background(),
		config.Config{
			Path:        tmpDir,
			Password:    "correct-pwd",
			AutoApprove: true,
		},
		ModeEncrypt,
	), "no file to consider")
	require.FileExists(t, jmblrFile.Name())

}

func checkFiles(t *testing.T, dir string, fn func(*testing.T, string)) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			fn(t, f.Name())
		}
	}
}
