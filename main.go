package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {

	err := run()
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, err.Error())
		if err != nil {
			panic(err)
		}
	}
}

func run() error {

	templateName := flag.String("t", "", "give a template name i.e. php, go, js")
	_ = flag.Bool("iphone", false, "pass -iphone if you want iPhone, iPad and iPod devices (default: false, but true if no other options given")
	_ = flag.Bool("watch", false, "pass -watch if you wan the Apple Watch models to be identified (default: false)")
	_ = flag.Bool("tv", false, "pass -tv (default: false")
	_ = flag.Bool("scan", false, "pass -scan to scan the Applications directory for traits databases")

	flag.Parse()

	iphone := isFlagPassed("iphone")
	watch := isFlagPassed("watch")
	tv := isFlagPassed("tv")
	scan := isFlagPassed("scan")

	filename, err := os.Executable()
	if err != nil {
		return err
	}

	applicationDirectory := filepath.Dir(filename)

	if _, err = os.Stat(filepath.Join(applicationDirectory, "templates")); errors.Is(err, fs.ErrNotExist) {
		applicationDirectory, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	if *templateName == "" {
		return errors.New("specify a template name. For example: -t=php")
	}

	if !iphone && !watch && !tv {
		iphone = true
	}

	devices := make(map[string]string, 0)
	devices["i386"] = "32-bit Simulator"
	devices["x86_64"] = "64-bit Simulator"
	devices["arm64"] = "64-bit Simulator"
	
	if scan {
		xcodePaths, err := filepath.Glob("/Applications/Xcode*")
		if err != nil {
			return err
		}
		for _, path := range xcodePaths {

			err = filepath.WalkDir(filepath.Join(path, "/Contents/Developer/Platforms/"), func(path string, d fs.DirEntry, err error) error {
				if strings.Index(path, ".platform") > 0 &&
					strings.HasSuffix(path, "traits.db") {

					if strings.Index(path, "AppleTVOS.platform") > 0 && !tv {
						return nil
					}
					if strings.Index(path, "WatchOS.platform") > 0 && !watch {
						return nil
					}
					if strings.Index(path, "iPhoneOS.platform") > 0 && !iphone {
						return nil
					}

					devs, err := getDevices(path)
					if err != nil {
						return err
					}

					for k, v := range devs {
						devices[k] = v
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

	} else {

		if iphone {
			devs, err := getDevices("/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/usr/standalone/device_traits.db")
			if err != nil {
				return err
			}

			for k, v := range devs {
				devices[k] = v
			}
		}
		if watch {
			devs, err := getDevices("/Applications/Xcode.app/Contents/Developer/Platforms/WatchOS.platform/usr/standalone/device_traits.db")
			if err != nil {
				return err
			}

			for k, v := range devs {
				devices[k] = v
			}
		}

		if tv {
			devs, err := getDevices("/Applications/Xcode.app/Contents/Developer/Platforms/AppleTVOS.platform/usr/standalone/device_traits.db")
			if err != nil {
				return err
			}

			for k, v := range devs {
				devices[k] = v
			}
		}
	}

	if len(devices) == 0 {
		return errors.New("no devices found")
	}

	tpl, err := template.ParseFiles(applicationDirectory + "/templates/" + *templateName + ".tpl")
	if err != nil {
		return err
	}
	err = tpl.Execute(os.Stdout, devices)
	if err != nil {
		return err
	}

	return nil
}

func getDevices(databaseFile string) (map[string]string, error) {

	db, err := sql.Open("sqlite", databaseFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res, err := db.Query(`SELECT ProductType, ProductDescription FROM Devices`)
	if err != nil {
		return nil, err
	}

	devices := make(map[string]string, 0)
	var prodType, prodDesc string
	for res.Next() {
		err = res.Scan(&prodType, &prodDesc)
		if err != nil {
			return nil, err
		}

		devices[prodType] = prodDesc
	}

	err = res.Err()
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func isFlagPassed(name string) (found bool) {
	flag.Visit(func(f *flag.Flag) {
		if strings.EqualFold(f.Name, name) {
			found = true
			if f.Value != nil {
				if strings.ToLower(f.Value.String()) == "false" {
					found = false
				}
			}
		}
	})
	return found
}
