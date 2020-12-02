//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

var (
	app *portapps.App
)

const (
	vmOptionsFile = "idea.vmoptions"
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("intellij-idea-community-portable", "IntelliJ IDEA Community"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	ideaExe := "idea.exe"
	ideaVmOptionsFile := "idea.exe.vmoptions"
	if runtime.GOARCH == "amd64" {
		ideaExe = "idea64.exe"
		ideaVmOptionsFile = "idea64.exe.vmoptions"
	}

	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "bin", ideaExe)
	app.WorkingDir = utl.PathJoin(app.AppPath, "bin")

	// override idea.properties
	ideaPropContent := strings.Replace(`# DO NOT EDIT! AUTOMATICALLY GENERATED BY PORTAPPS.
idea.config.path={{ DATA_PATH }}/config
idea.system.path={{ DATA_PATH }}/system
idea.plugins.path={{ DATA_PATH }}/plugins
idea.log.path={{ DATA_PATH }}/log`, "{{ DATA_PATH }}", utl.FormatUnixPath(app.DataPath), -1)

	ideaPropPath := utl.PathJoin(app.DataPath, "idea.properties")
	if err := utl.CreateFile(ideaPropPath, ideaPropContent); err != nil {
		log.Fatal().Err(err).Msg("Cannot write idea.properties")
	}

	// https://www.jetbrains.com/help/idea/tuning-intellij-idea.html#configure-platform-properties
	os.Setenv("IDEA_PROPERTIES", ideaPropPath)

	// https://www.jetbrains.com/help/idea/tuning-the-ide.html#configure-jvm-options
	os.Setenv("IDEA_VM_OPTIONS", utl.PathJoin(app.DataPath, vmOptionsFile))
	if !utl.Exists(utl.PathJoin(app.DataPath, vmOptionsFile)) {
		utl.CopyFile(utl.PathJoin(app.AppPath, "bin", ideaVmOptionsFile), utl.PathJoin(app.DataPath, vmOptionsFile))
	} else {
		utl.CopyFile(utl.PathJoin(app.DataPath, vmOptionsFile), utl.PathJoin(app.AppPath, "bin", ideaVmOptionsFile))
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
