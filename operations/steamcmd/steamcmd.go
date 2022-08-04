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
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var DownloadLink = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
var DebLink = "http://security.ubuntu.com/ubuntu/pool/main/g/glibc/libc6-i386_2.31-0ubuntu9.7_amd64.deb"
var PatchelfLink = "https://github.com/NixOS/patchelf/releases/download/0.15.0/patchelf-0.15.0-x86_64.tar.gz"

var downloader sync.Mutex

type Steamcmd struct {
}

func (Steamcmd) Run(env pufferpanel.Environment) (err error) {
	env.DisplayToConsole(true, "Downloading SteamCMD")

	binaryFolder := filepath.Join(config.GetString("daemon.data.binaries"), "steamcmd")
	mainCommand := filepath.Join(config.GetString("daemon.data.binaries"), "steamcmd.sh")

	downloader.Lock()
	defer downloader.Unlock()

	if _, err = os.Lstat(mainCommand); errors.Is(err, fs.ErrNotExist) {
		//we need to download, but if we can't, we should purge what we downloaded
		defer func(folder string) {
			if err != nil {
				_ = os.RemoveAll(folder)
			}
		}(binaryFolder)

		//we need to get the binaries and install them
		err = pufferpanel.HttpGetTarGz(DownloadLink, binaryFolder)
		if err != nil {
			return
		}

		//now... get deps we need to merge with this
		err = pufferpanel.HttpGetTarGz(PatchelfLink, binaryFolder)
		if err != nil {
			return
		}

		err = pufferpanel.HttpDownloadDeb(DebLink, binaryFolder)
		if err != nil {
			return
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
			return
		}

		_ = os.RemoveAll(filepath.Join(binaryFolder, "bin"))
		err = os.Symlink(filepath.Join(binaryFolder, "steamcmd.sh"), mainCommand)
	}

	return
}
