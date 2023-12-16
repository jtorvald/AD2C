package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
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

type Device struct {
	Identifier  string
	Description string
}

func run() error {
	seenBefore = make(map[string]struct{}, 0)

	templateName := flag.String("t", "", "give a template name i.e. php, go, js")
	_ = flag.Bool("iphone", false, "pass -iphone if you want iPhone, iPad and iPod devices (default: false, but true if no other options given")
	_ = flag.Bool("watch", false, "pass -watch if you wan the Apple Watch models to be identified (default: false)")
	_ = flag.Bool("tv", false, "pass -tv (default: false")
	_ = flag.Bool("scan", false, "pass -scan to scan the Applications directory for traits databases")
	_ = flag.Bool("adamawolf", false, "pass -adamawolf to download the latest gist from https://gist.github.com/adamawolf/3048717")

	flag.Parse()

	iphone := isFlagPassed("iphone")
	watch := isFlagPassed("watch")
	tv := isFlagPassed("tv")
	scan := isFlagPassed("scan")
	adamawolf := isFlagPassed("adamawolf")

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

	xcodePaths, err := filepath.Glob("/Applications/Xcode*")

	if len(xcodePaths) == 0 && !adamawolf {
		scan = false
		adamawolf = true
	}

	if *templateName == "" {
		return errors.New("specify a template name. For example: -t=php")
	}

	devices := make([]Device, 0)

	if adamawolf {

		url := "https://gist.githubusercontent.com/adamawolf/3048717/raw/1ee7e1a93dff9416f6ff34dd36b0ffbad9b956e9/"
		client := http.Client{
			Timeout: 2 * time.Second,
		}
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		txtScanner := bufio.NewScanner(resp.Body)
		txtScanner.Split(bufio.ScanLines)

		for txtScanner.Scan() {
			line := txtScanner.Text()
			if strings.Index(line, ":") <= -1 {
				continue
			}

			foo := strings.Split(line, ":")
			identifier, description := strings.TrimSpace(foo[0]), strings.TrimSpace(foo[1])

			devices = appendToDevices(devices, Device{
				Identifier:  cleanupIdentifier(identifier),
				Description: description,
			})
		}
	} else {

		if !iphone && !watch && !tv {
			iphone = true
		}

		devices = appendToDevices(devices, Device{
			Identifier:  cleanupIdentifier("i386"),
			Description: "32-bit Simulator",
		}, Device{
			Identifier:  cleanupIdentifier("x86_64"),
			Description: "64-bit Simulator",
		}, Device{
			Identifier:  cleanupIdentifier("arm64"),
			Description: "64-bit Simulator",
		})

		if scan {

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
							devices = appendToDevices(devices, Device{
								Identifier:  cleanupIdentifier(k),
								Description: v,
							})
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
					devices = appendToDevices(devices, Device{
						Identifier:  cleanupIdentifier(k),
						Description: v,
					})
				}
			}
			if watch {
				devs, err := getDevices("/Applications/Xcode.app/Contents/Developer/Platforms/WatchOS.platform/usr/standalone/device_traits.db")
				if err != nil {
					return err
				}

				for k, v := range devs {
					devices = appendToDevices(devices, Device{
						Identifier:  cleanupIdentifier(k),
						Description: v,
					})
				}
			}

			if tv {
				devs, err := getDevices("/Applications/Xcode.app/Contents/Developer/Platforms/AppleTVOS.platform/usr/standalone/device_traits.db")
				if err != nil {
					return err
				}

				for k, v := range devs {
					devices = appendToDevices(devices, Device{
						Identifier:  cleanupIdentifier(k),
						Description: v,
					})
				}
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

var seenBefore map[string]struct{}

func appendToDevices(devices []Device, device ...Device) []Device {
	for _, d := range device {
		if _, ok := seenBefore[d.Identifier]; !ok {
			devices = append(devices, d)
			seenBefore[d.Identifier] = struct{}{}
		}
	}
	return devices
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

		devices[cleanupIdentifier(prodType)] = prodDesc
	}

	err = res.Err()
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func cleanupIdentifier(identifier string) string {
	if strings.HasSuffix(identifier, "-A") {
		return strings.TrimSuffix(identifier, "-A")
	}
	if strings.HasSuffix(identifier, "-B") {
		return strings.TrimSuffix(identifier, "-B")
	}
	return identifier
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
