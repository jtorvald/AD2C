package apple

/*
Automatically generated by AD2C (Apple Devices to Code)
https://github.com/jtorvald/AD2C
*/


// GetReadableDeviceModel returns the description for a given device model
func GetReadableDeviceModel(model string) string {
    switch model { {{range $identifier, $description := . }}
        case "{{$identifier}}": return "{{$description}}"{{end}}
    }
    return "Unknown"
}