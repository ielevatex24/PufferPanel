/*
 Copyright 2022 PufferPanel

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package steamcmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"pault.ag/go/debian/deb"
	"sync"
)

var client = &http.Client{}
var DownloadLink = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
var DebLink = "http://security.ubuntu.com/ubuntu/pool/main/g/glibc/libc6-i386_2.31-0ubuntu9.7_amd64.deb"
var PatchelfLink = "https://github.com/NixOS/patchelf/releases/download/0.15.0/patchelf-0.15.0-x86_64.tar.gz"

var downloader sync.Mutex

type Steamcmd struct {
}

func (op Steamcmd) Run(env pufferpanel.Environment) error {
	env.DisplayToConsole(true, "Downloading SteamCMD")
	return downloadSteamcmd()
}

func downloadSteamcmd() (err error) {
	downloader.Lock()
	defer downloader.Unlock()

	binaryFolder := filepath.Join(config.GetString("daemon.data.binaries"), "steamcmd")
	mainCommand := filepath.Join(config.GetString("daemon.data.binaries"), "steamcmd.sh")

	defer func(folder string) {
		if err != nil {
			_ = os.RemoveAll(folder)
		}
	}(binaryFolder)

	if _, err = os.Lstat(mainCommand); errors.Is(err, fs.ErrNotExist) {
		err = nil

		//we need to get the binaries and install them
		err = downloadTarGz(DownloadLink, binaryFolder)
		if err != nil {
			return err
		}

		//now... get deps we need to merge with this
		err = downloadTarGz(PatchelfLink, binaryFolder)
		if err != nil {
			return err
		}

		err = downloadDeb(DebLink, binaryFolder)
		if err != nil {
			return err
		}

		//clean up deb files that we really don't need....
		_ = os.RemoveAll(filepath.Join(binaryFolder, "etc"))
		_ = os.RemoveAll(filepath.Join(binaryFolder, "lib"))
		_ = os.RemoveAll(filepath.Join(binaryFolder, "share"))
		_ = os.RemoveAll(filepath.Join(binaryFolder, "usr"))

		lib32Path := filepath.Join(binaryFolder, "lib32")
		linux32Path := filepath.Join(binaryFolder, "linux32")
		_ = os.Rename(filepath.Join(lib32Path, "ld-2.31.so"), filepath.Join(linux32Path, "ld-linux-so.2"))
		_ = os.Rename(filepath.Join(lib32Path, "librt-2.31.so"), filepath.Join(linux32Path, "librt.so.1"))
		_ = os.Rename(filepath.Join(lib32Path, "libdl-2.31.so"), filepath.Join(linux32Path, "libdl.so.2"))
		_ = os.Rename(filepath.Join(lib32Path, "libpthread-2.31.so"), filepath.Join(linux32Path, "libpthread.so.0"))
		_ = os.Rename(filepath.Join(lib32Path, "libm-2.31.so"), filepath.Join(linux32Path, "libm.so.6"))
		_ = os.Rename(filepath.Join(lib32Path, "libc-2.31.so"), filepath.Join(linux32Path, "libc.so.6"))

		_ = os.RemoveAll(lib32Path)

		//bin/patchelf --set-interpreter linux32/ld-linux-so.2 linux32/steamcmd
		cmd := exec.Command("bin/patchelf", "--set-interpreter", "linux32/ld-linux-so.2", "linux32/steamcmd")
		cmd.Dir = binaryFolder
		err = cmd.Run()
		if err != nil && !errors.Is(err, &exec.ExitError{}) {
			return err
		}

		_ = os.RemoveAll(filepath.Join(binaryFolder, "bin"))
		err = os.Symlink(filepath.Join(binaryFolder, "steamcmd.sh"), mainCommand)
	}

	return err
}

func extractGz(gzipStream io.Reader, directory string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()
	return extractTar(uncompressedStream, directory)
}

func extractTar(stream io.Reader, directory string) error {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	var tarReader *tar.Reader
	if r, isGood := stream.(*tar.Reader); isGood {
		tarReader = r
	} else {
		tarReader = tar.NewReader(stream)
	}

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(directory, header.Name), 0755); err != nil {
				return err
			}
		case tar.TypeSymlink:
			continue
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Join(directory, filepath.Dir(header.Name)), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(filepath.Join(directory, header.Name))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				return err
			}
			_ = outFile.Close()
			err = os.Chmod(filepath.Join(directory, header.Name), header.FileInfo().Mode())
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("uknown type: %s in %s", string([]byte{header.Typeflag}), header.Name))
		}
	}
	return nil
}

func extractDeb(stream io.ReaderAt, directory string) error {
	file, err := deb.Load(stream, directory)
	if err != nil {
		return err
	}

	defer file.Close()

	return extractTar(file.Data, directory)
}

func downloadTarGz(link, folder string) error {
	response, err := client.Get(link)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	err = extractGz(response.Body, folder)
	return err
}

func downloadDeb(link, folder string) error {
	response, err := client.Get(link)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buff, response.Body)
	_ = response.Body.Close()

	if err != nil {
		return err
	}

	reader := bytes.NewReader(buff.Bytes())

	err = extractDeb(reader, folder)
	return err
}
