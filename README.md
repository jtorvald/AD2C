# Apple Devices to Code (AD2C)
With this tool you can generate code to get the description for Apple device models.
It makes use of Go's text templates, so that you can make templates for your own use case.

This is an initial release to help myself out. 
If you like it, don't forget to star this repo. If you have feature requests in case you discovered a bug,
please file them under issues.

# Requirements #
The tool expects you to have Xcode on your system because it reads the devices directly from
the Xcode traits database on your system.

# How to run #
To write to a file:

`$ ad2bc -t=js > apple-devices.js`

To write to your clipboard:

`$ ad2bc -t=js | pbcopy`

Options are:
- `-iphone` generate code for iPhone and iPad/iPod models
- `-watch` generate code for Watch models
- `-tv` generate code for Apple TV models
- `-scan` scan your /Applications/Xcode* directories for databases

# Download releases #
The best is to [download the binary for your platform and run the executable](https://github.com/jtorvald/ad2c/releases).

# License #
MIT