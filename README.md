# cli-filesharer
CLI tool in GO to share files through LAN

### Usage: 
To build the app:
`go build -o sendfiles .` 

To run it: 
`sendfiles file1.txt file2.txt file7.rar dir/file4.zip \FullPath\To\File.jpeg`


It opens a web-server which serves all the requested files.
To access a list of all the served files open `localhost/` on the server, or the given IP by the app to the user. 

Installers may be created shortly, but really it's just 2 commands. 

I did this in a day because I had no USB drive to share some 4GB files from a laptop to a PC. So suggestions are welcome. :) 
