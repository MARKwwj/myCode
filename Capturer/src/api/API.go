package api

import (
	"flag"
	"fmt"
	"strings"
)

// IModule IModule
type IModule interface {
	Name() string
	Run() error
}

var (
	mDebug       = flag.Bool("debug", false, "Use debug mode")
	mTest        = flag.Bool("test", false, "Use test mode")
	mModules     []IModule
	mStartModule = flag.String("start", "", "Enable the specified module")
)

func init() {
	flag.Parse()
	*mStartModule = strings.ToLower(*mStartModule)
	StartService()
}

// RegisterCapturer RegisterCapturer
func RegisterCapturer(module IModule) {
	if IsDebug() {
		Info("RegisterCapturer -> %s", module.Name())
	}
	mModules = append(mModules, module)
}

// IsDebug IsDebug
func IsDebug() bool {
	return *mDebug
}

// IsTest IsTest
func IsTest() bool {
	return *mTest
}

// StartAllCapturer StartAllCapturer
func StartAllCapturer(defaultModule, currentVersion string) {
	defer CloseService()
	// if err := initDownloader(); err != nil {
	// 	Error("Init downloader error: %s", err.Error())
	// 	return
	// }
	if len(*mStartModule) <= 0 {
		*mStartModule = defaultModule
	}
	SetConsoleTitle(fmt.Sprintf("Capturer - %s", currentVersion))
	*mStartModule = strings.ToLower(*mStartModule)
	mStartModuleList := strings.Split(*mStartModule, ",")

	for _, module := range mModules {
		if module != nil {
			startModule := false
			moduleName := strings.ToLower(module.Name())
			for _, startName := range mStartModuleList {
				if moduleName == startName {
					startModule = true
					break
				}
			}
			if startModule {
				Info("Run capturer -> %s", module.Name())
				go func(module IModule) {
					if err := module.Run(); err != nil {
						Error("[%s]%s", module.Name(), err.Error())
					}
				}(module)
			} else if IsDebug() {
				Debug("Skip capturer -> %s", module.Name())
			}
		}
	}
	select {}
}
