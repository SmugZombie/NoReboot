package main

// NoReboot
// A simple script designed to disable automatic reboots by windows updates. 
// Smugzombie - github.com/smugzombie
// Version 1.2
// We're not saying that you should skip installing updates, as they're important to keep your device secure and up to date. 
// However, there are scenarios where you make want to take full control and decide exactly when to restart your computer 
// to apply new updates, and this is when knowing how to stop automatic reboots comes in handy.

import "fmt"
import "os"
import "io"
import "bufio"
import "strings"

// Global Variables
var rootpath = "C:\\WINDOWS\\System32\\Tasks\\Microsoft\\Windows\\UpdateOrchestrator\\"
var rebootpath = rootpath + "Reboot"
var rebootpathbak = rebootpath + ".bak"
var aboutpath = rootpath + ".NoReboot"
var undo = ""
var version = "1.2"

func main() {
	// Show the splash screen
	splash()

	checkIfAdmin()

	// Check to see if NoReboot is already enabled
	if(checkForNoReboot()){
		// If so, tell user and ask if they wish to remove it
		fmt.Println("NoReboot is already installed \n")
		if(readUserInput("Would you like to uninstall it now? (y|yes): ")){
			// If yes, remove it
			removeNoReboot()
		}else{
			// Otherwise, exit gracefully
			fmt.Println("Ok, Maybe next time. Goodbye!")
		}
	}else{
		// If not, ask the user if they would like to install it.
		fmt.Println("NoReboot is not installed \n")
		if(readUserInput("Would you like to install it now? (y|yes): ")){
			// If yes, install
			installNoReboot()
		}else{
			// Else, exit gracefully
			fmt.Println("Ok, Maybe next time. Goodbye!")
		}
	}
	userWait()
}

// A simple function to ask a user a question, and only really care if they answer positively, otherwise return false
func readUserInput(message string) bool{
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text2, _ := reader.ReadString('\n')
	// Move to lowercase so contains matching is easy
	text := strings.ToLower(text2) 
	// Check for a y, even "no i dont want to because its really foolish" should match, as it has a y
	if(strings.Contains(text, "y")){
		return true
	}else{
		return false
	}
}

// Simply create a user prompt to ensure the console window stays open long enough to read
func userWait(){
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n" + "Press Enter to Continue")
    text, _ := reader.ReadString('\n')
    _ = text
}

// Display the splash for the app
func splash(){
	fmt.Println(" _   _      ______     _                 _   ")
	fmt.Println("| \\ | |     | ___ \\   | |               | |  ")
	fmt.Println("|  \\| | ___ | |_/ /___| |__   ___   ___ | |_ ")
	fmt.Println("| . ` |/ _ \\|    // _ \\ '_ \\ / _ \\ / _ \\| __|")
	fmt.Println("| |\\  | (_) | |\\ \\  __/ |_) | (_) | (_) | |_ ")
	fmt.Println("\\_| \\_/\\___/\\_| \\_\\___|_.__/ \\___/ \\___/ \\__| v"+version)
	fmt.Println("                                             ")
	fmt.Println("Brought to you by: Smugzombie (github.com/smugzombie)")
	fmt.Println("")
}

// Check to see if NoReboot is already enabled
func checkForNoReboot() bool{
	if(fileExists(rebootpath)){
		if(isDirorFile(rebootpath) == "dir"){
			return true
		}
	}
	return false
}

// Install NoReboot
func installNoReboot(){
	// Check if current reboot file exists
	fmt.Println("Looking for: " + rebootpath)
	if(fileExists(rebootpath)){
		// Double check if NoReboot isn't actually installed..
		if(isDirorFile(rebootpath) == "file"){
			// If not, lets see if the backup file exists
			fmt.Println("Looking for: " + rebootpathbak)
			if(fileExists(rebootpathbak)){
				// If the file exists, consider it backed up and move on
				fmt.Println("Already Backed Up!")
			}else{
				// If the file does not exist.. create it.. then move on
				fmt.Println("Backing up: " + rebootpath	+ " to " + rebootpathbak)
				moveFile(rebootpath, rebootpathbak)
			}
				// Delete Original File
				fmt.Println("Deleting Original Reboot File")
				deleteFile(rebootpath)
				// Create Directory in Place
				fmt.Println("Creating Reboot Directory")
				CreateDirIfNotExist(rebootpath)
				// Create about file
				createAbout(aboutpath)

				fmt.Println("\nNoReboot Successfully Applied!")
		}else{
			// If the directory does exist, NoReboot is already applied
			fmt.Println("\nNoReboot already applied")
		}
	}else{
		// Uhoh, we don't know where the reboot file is, but shouldn't proceed
		fmt.Println("File does not exist. Cannot Proceed")
	}
}

// Uninstall NoReboot
func removeNoReboot(){
	// Delete our directory
	fmt.Println("Removing NoReboot Directory")
	deleteFile(rebootpath)
	// Restore the backup file to original
	fmt.Println("Restoring: " + rebootpathbak	+ " to " + rebootpath)
	moveFile(rebootpathbak, rebootpath)
	// Delete the backup file
	fmt.Println("\nNoReboot has restored your reboot configuration.")
	deleteFile(rebootpathbak)
	// Delete the about file
	deleteFile(aboutpath)
}

// Check if file exists
func fileExists(filepath string) bool{
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}else{
		return true
	}
}

// Copys a file from one place to another
func moveFile(filepath string, newfilepath string) bool{
	Copy(filepath, newfilepath)
	return true
}

// Deletes a file
func deleteFile(filepath string) bool{
	var err = os.Remove(filepath)
	if err != nil { return false } // If something went wrong, simply return false
	return true
}

// Create a directory, if not already existing
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
				panic(err)
		}
	}
}

// Copies a file from src to dst location
func Copy(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil { return }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil { return }
    defer func() {
        cerr := out.Close()
        if err == nil { err = cerr }
    }()
    if _, err = io.Copy(out, in); err != nil { return }
    err = out.Sync()
    return
}

// Checks to see if the filepath in question is actually a file or directory
func isDirorFile(filepath string) string{
	fi, err := os.Stat(filepath)
    if err != nil {
        fmt.Println(err)
        return ""
    }
    switch mode := fi.Mode(); {
    case mode.IsDir():
        return "dir"
    case mode.IsRegular():
        return "file"
    }
    return ""
}

// Creates the about file that explains how to disable
func createAbout(filepath string){
	var _, err = os.Stat(filepath)

	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		if err != nil {
        	fmt.Println(err)
    	}
		defer file.Close()
	}

	var file, err2 = os.OpenFile(filepath, os.O_RDWR, 0644)
	if err2 != nil {
        fmt.Println(err)
    }
	defer file.Close()

	_, err = file.WriteString("This file exists to show that NoReboot has been applied to this machine. To remove NoReboot simply run the script again and say yes to uninstalling or to do so manually delete the 'Reboot' directory and move 'Reboot.bak' back to 'Reboot'")
	if err != nil { fmt.Println(err) }

	err = file.Sync()
	if err != nil { fmt.Println(err) }
}

func checkIfAdmin(){
    Block{
        Try: func() {
            fo, err := os.Create("c:\\test.txt")
            if err != nil { 
                Throw("Oops, NoReboot Needs to be run with Administrative Priviledges. Run-As Administrator to continue.") 
            }
            defer fo.Close()
        },
        Catch: func(e Exception) {
            fmt.Printf("%v\n", e)
            userWait()
            os.Exit(1)
        },
        Finally: func() {
        },
    }.Do()
}

/// Try Block ///
type Block struct {
    Try     func()
    Catch   func(Exception)
    Finally func()
}

type Exception interface{}

func Throw(up Exception) {
    panic(up)
}
func (tcf Block) Do() {
    if tcf.Finally != nil {
            defer tcf.Finally()
    }
    if tcf.Catch != nil {
        defer func() {
            if r := recover(); r != nil {
                    tcf.Catch(r)
            }
        }()
    }
    tcf.Try()
}
/// End Try Block ///
