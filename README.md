# Apple Devices to Code (AD2C)
With this tool you can generate code to get the description for Apple device models.
It makes use of Go's text templates, so that you can make templates for your own use case.

This is an initial release to help myself out. 
If you like it, don't forget to star this repo. If you have feature requests in case you discovered a bug,
please file them under issues.

# Requirements #
If you don't have Xcode installed, the tool will automatically use the gist from @adamawolf 
[gist.github.com/adamawolf](https://gist.github.com/adamawolf/3048717)

When you have Xcode installed it will use the database from Xcode to generate the code.

# How to run #
To write to a file:

`$ ad2bc -t=js -iphone -tv -watch > apple-devices.js`

To write to your clipboard:

`$ ad2bc -t=js -iphone -tv -watch | pbcopy`

Force the use of the gist:

`$ ad2bc -t=js -iphone -tv -watch -adamwolf | pbcopy`

Options are:
- `-iphone` generate code for iPhone and iPad/iPod models
- `-watch` generate code for Watch models
- `-tv` generate code for Apple TV models
- `-scan` scan your /Applications/Xcode* directories for databases
- `-adamawolf` when you don't have Xcode installed, use the gist from [gist.github.com/adamawolf](https://gist.github.com/adamawolf/3048717)

# Download releases #
The best is to [download the binary for your platform and run the executable](https://github.com/jtorvald/ad2c/releases).

# License #
MIT