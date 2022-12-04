namespace "Apple"

/*
Automatically generated by AD2C (Apple Devices to Code)
https://github.com/jtorvald/AD2C
*/

static class AppleDevices
{
    static function GetReadableDeviceModel($model): string
    {
        switch ($model) { {{range $identifier, $description := . }}
            case "{{$identifier}}": return "{{$description}}";{{end}}
        }
        return "Unknown";
    }

}

