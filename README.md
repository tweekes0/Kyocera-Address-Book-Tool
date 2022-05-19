# Kyocera Address Book Tool

Kyocera Address Book tool or kyocera-ab-tool is an interractive CLI application
written in Go to facilitate the creation and maintainence of Kyocera scanner 
address books and their OneTouchKeys (scanner shortcuts). 

I wrote this tool to get a feel for Go, take advantage of Go's portability and 
the [Go Gopher](https://go.dev/blog/gopher) is also pretty cute. This tool creates
XML files to be imported into Kyocera scanners via the [Kyocera Net Viewer](https://www.kyoceradocumentsolutions.us/en/products/software/KYOCERANETVIEWER.html) tool.

## Usage

Download the release for your target OS and extract the folder to your desired location. 
Use **cmd.exe** and navigate to the extracted folder and run the executable. 

 ## Commands

    add_user        : add user to the current table. Fields must be separated by commas
    clear_table     : clears all users from the current table
    create_table    : creates new table and sets it to the current table
    delete_table    : deletes the specified table
    delete_user     : delete a single user from the current table
    exit            : exits the program
    export_table    : exports the current table to an xml file in the Address Books directory
    import_csv      : import users from csv file into current table
    list_tables     : list all tables
    show_users      : show all the users in the current table
    switch_table    : switch the current table
    update_user     : update user in the current table. Fields must be separated by commas


## Acknowledgements

This application uses these great libraries
- [go-sqlite3](https://github.com/mattn/go-sqlite3) -- Go SQLite3 driver
- [readline](https://github.com/chzyer/readline) -- Pure Go implementation of Readline for user input 
- [table](https://github.com/rodaine/table) -- Go CLI table generator
- [xgo](https://github.com/karalabe/xgo) -- CGO Cross Compiler

## Known Issues

- Kyocera Net Viewer seems to only be for Windows but the there are binaries for Linux
- The tool has libraries that use CGO so the application cannot be compiled for 
Windows using the Go build tool so a cross-compilation tool is necessary. 
- Running the application on Windows Powershell causes a weird graphical glitch 
but the application is fully functional.
