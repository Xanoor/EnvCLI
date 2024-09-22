// Env - Command Line Interface by Xanoor
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Gray = "\033[37m"

// List of commands.
var cmds = []string{"-help", "-get", "-create", "-rename", "-delete", "-read", "-add", "-remove", "-update", "-man", "-h", "-quit", "-q"}

func isCommand(array []string, cmdName string) (int, int, bool) {
	var index int
	var maxIndex int
	var hasBeenFound = false

	// Iterate through the array to find the command
	for i, v := range array {
		if v == cmdName {
			index = i + 1 // Set the starting index after the command
			hasBeenFound = true
		}

		// If the command has been found, check for the end of arguments
		if hasBeenFound && ((strings.HasPrefix(v, "-") && v != cmdName) || i == len(array)-1) {
			if strings.HasPrefix(v, "-") {
				maxIndex = i - 1 // Set max index to the last argument before the next command
			} else {
				maxIndex = i // If it's the last item, set max index to current index
			}
			return index, maxIndex, true // Return the indices and found status
		}
	}

	return 0, 0, false // Return false if the command was not found
}

func main() {
	fmt.Println(`
  ______             _____ _      _____ 
 |  ____|           / ____| |    |_   _|
 | |__   _ ____   _| |    | |      | |  
 |  __| | '_ \ \ / / |    | |      | |  
 | |____| | | \ V /| |____| |____ _| |_ 
 |______|_| |_|\_/  \_____|______|_____|		  
 ENVelope Command Line Interface By Xanoor
	`)

	for {
		fmt.Print(Red + "[EnvCLI] [" + time.Now().Format("15:04:05") + "] : " + Reset)
		reader := bufio.NewReader(os.Stdin)
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred!")
			return
		}

		cmd = strings.TrimSpace(cmd)     // Remove any leading/trailing whitespace
		input := strings.Split(cmd, " ") // Split input into command and arguments

		// Check if the command is valid
		if _, _, found := isCommand(cmds, input[0]); found {
			if len(input) > 1 {
				// Handle commands with additional arguments
				switch input[0] {
				case "-create":
					fmt.Println(create(input))
				case "-get":
					fmt.Println(get(input))
				case "-rename":
					fmt.Println(rename(input))
				case "-delete":
					fmt.Println(delete(input))
				case "-help", "-man", "-h":
					fmt.Println(help(input))
				case "-add":
					fmt.Println(add(input))
				case "-read":
					fmt.Println(read(input))
				case "-remove":
					fmt.Println(remove(input))
				case "-update":
					fmt.Println(update(input))
				default:
					fmt.Println(Red + "Unknown command." + Reset)
				}
			} else {
				// If only the command is provided
				if input[0] == "-help" || input[0] == "-man" || input[0] == "-h" {
					fmt.Println(help(input))
				} else if input[0] == "-quit" || input[0] == "-q" {
					break
				} else {
					fmt.Println(Red + "Expected one argument!" + Reset)
				}
			}
		} else {
			// If the command is invalid
			fmt.Println(Red + "Invalid command!" + Yellow + "\nList of commands: -help\nDon't forget to add \"-\" in front of the command name!\n[EXAMPLE]: -help, -quit, -create, -delete" + Reset)
		}
	}
}

// Adds the .env extension to the file name if it's not already present.
func addExtension(fileName string) string {
	if !strings.HasSuffix(fileName, ".env") {
		return fileName + ".env"
	}
	return fileName
}

func verify(sentence string) bool {
	fmt.Println(Yellow + sentence + Reset)
	var res string
	fmt.Scanln(&res)

	// Check the user's response
	if res == "y" || res == "Y" {
		return true // User confirmed with "yes"
	} else if res == "n" || res == "N" {
		return false // User declined with "no"
	} else {
		fmt.Println(Yellow + "Unknown response! Response may be \"y\" or \"n\", not \"" + res + "\"" + Reset)
		return false
	}
}

func writeData(fileName, variable, data string) bool {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer f.Close() // Ensure the file is closed at the end of the function

	_, err = fmt.Fprintln(f, variable+"="+data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return false
	}

	return true
}

func overwriteFile(filePath string, content []string) bool {
	// Create or open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close() // Ensure the file is closed after use

	// Write the new content to the file
	_, err = file.WriteString(strings.Join(content, "\n"))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return false
	}
	return true
}

func getVariable(variable string, content string) string {
	result := ""
	state := false
	// Iterate through each line in the content
	for _, line := range strings.Split(content, "\n") {
		// Check if the line contains the specified variable/value
		if strings.Contains(strings.ToUpper(line), variable) {
			result += "-" + line + "\n" // Append the line to the result
			state = true
		}
	}

	if state {
		return Green + "Variable(s)/Value found!\n" + result + Reset
	} else {
		return Red + variable + " variable/value doesn't exist!" + Reset
	}
}

func getFileData(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "Error", err
	}
	return string(data), nil
}

// Verify if the file exists / is valid
func isFileValid(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func updateVar(index int, maxIndex int, fileName string, command []string, prompt bool) string {
	// Iterate over the specified variable indices to update their values
	for i := index; i <= maxIndex; i++ {
		var newData []string
		// Retrieve the current data from the file
		data, err := getFileData(fileName)
		if err != nil {
			return Red + "Error occurred during operation!" + Reset
		}

		// Check if the variable exists in the file
		if !strings.Contains(data, command[i]+"=") {
			return Red + command[i] + " wasn't found!" + Reset
		}

		// Prompt the user for a new value for the variable
		fmt.Println(Yellow + "Insert new value for variable " + command[i] + ":" + Reset)
		reader := bufio.NewReader(os.Stdin)
		value, err := reader.ReadString('\n')
		if err != nil {
			return Red + "An error occurred!" + Reset
		}

		value = strings.TrimSpace(value) // Remove any extra whitespace

		var varToUpdate []string
		// If prompting is enabled, confirm the value change with the user
		if prompt {
			if state := verify(Yellow + "Are you sure you want to change the value of variable " + command[i] + " to " + value + "? (y/n)" + Reset); !state {
				fmt.Println(Red + "Variable not updated!" + Reset)
				continue // Skip to the next variable if the user declines
			}
		}
		varToUpdate = append(varToUpdate, command[i])

		// Split the current data into lines to find the variable to update
		splittedData := strings.Split(data, "\n")
		for l := 0; l < len(splittedData); l++ {
			splittedVar := strings.Split(splittedData[l], "=")
			// Check if the current line corresponds to the variable being updated
			if splittedVar[0] == command[i] && slices.Contains(varToUpdate, command[i]) {
				newData = append(newData, splittedVar[0]+"="+value) // Update the value
			} else {
				newData = append(newData, splittedData[l]) // Keep the existing line
			}
		}

		// Attempt to overwrite the file with the updated data
		if result := overwriteFile(fileName, newData); result {
			fmt.Println(Green + command[i] + " successfully updated!" + Reset)
		} else {
			fmt.Println(Red + "Error when updating variable " + command[i] + "!" + Reset)
		}
	}
	return ""
}

func update(command []string) string {
	command[1] = addExtension(command[1])

	if isFileValid(command[1]) {
		// Check for the presence of the -var option to update variables
		index, maxIndex, found := isCommand(command, "-var")
		if found {
			// Check for the presence of the -p option (prompt for confirmation)
			_, _, prompt := isCommand(command, "-p")
			// Call the updateVar function to update the specified variables
			return updateVar(index, maxIndex, command[1], command, prompt)
		} else {
			return Red + "Incorrect use of command -update!\n-help -update for more info!" + Reset
		}
	} else {
		return Red + command[1] + " not found!" + Reset
	}
}

func create(command []string) string {
	command[1] = addExtension(command[1])

	if isFileValid(command[1]) {
		// Ask the user if they want to overwrite the existing file
		if !verify(Yellow + "File already exists! Do you want to overwrite it? (y/n)" + Reset) {
			return Red + "Action cancelled!" + Reset
		}
	}

	// Create the new file
	file, err := os.Create(command[1])
	if err != nil {
		return Red + "Error when opening the file." + Reset
	}
	defer file.Close() // Ensure the file is closed after we're done

	// Check for the presence of the -var option to add variables immediately
	index, maxIndex, found := isCommand(command, "-var")
	if found {
		fmt.Println(addVariable(index, maxIndex, command)) // Add variables if -var is present
		return Green + "File and variable(s) created successfully!" + Reset
	} else {
		// If the -s option is not found, ask if the user wants to add variables
		if _, _, found := isCommand(command, "-s"); !found {
			if response := verify(Green + "File created, do you want to add variable(s)? (y/n)"); response {
				fmt.Println("Write the first variable name: (stop with \"q\" or \"quit\")")
				for {
					var varName string
					fmt.Scanln(&varName)
					if varName == "q" || varName == "quit" {
						break // Exit the loop if the user wants to stop
					}
					fmt.Println(Yellow + "Value of variable " + varName + ":" + Reset)
					var varValue string
					fmt.Scanln(&varValue)
					if varValue == "q" || varValue == "quit" {
						break // Exit the loop if the user wants to stop
					}
					// Attempt to write the variable to the file
					if hasWroteData := writeData(command[1], varName, varValue); hasWroteData {
						fmt.Println(Green + "Variable " + varName + " added to env.\n" + Reset)
					} else {
						break // Exit the loop on failure to write data
					}
					fmt.Println("Write the next variable name: (stop with \"q\" or \"quit\")")
				}
				return Green + "File and variable(s) created successfully!" + Reset
			}
		}
		return Green + "File created successfully!" + Reset
	}
}

func addVariable(index int, maxIndex int, command []string) string {
	// Loop through each variable from the specified index to the maximum index
	for i := index; i <= maxIndex; i++ {
		fmt.Println(Yellow + "Insert value for variable " + command[i] + ":" + Reset)
		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString('\n') // Read user input until newline
		if err != nil {
			return Red + "An error occurred while reading the input!" + Reset
		}
		data = strings.TrimSpace(data) // Remove any leading/trailing whitespace

		// Attempt to write the variable data to the file
		if hasWroteData := writeData(command[1], command[i], data); !hasWroteData {
			return Red + "An error occurred when adding " + command[i] + " with value: " + data + Reset
		}
	}
	return Green + "\nVariable(s) added!" + Reset
}

func add(command []string) string {
	command[1] = addExtension(command[1])

	if isFileValid(command[1]) {
		// Check for the presence of the -var option and get the indices
		index, maxIndex, found := isCommand(command, "-var")
		if found {
			// Call the function to add variables to the file
			return addVariable(index, maxIndex, command)
		} else {
			return Red + "Incorrect use of the -add command!\n-help -add for more info!" + Reset
		}
	} else {
		return Red + command[1] + " not found!" + Reset
	}
}

func remove(command []string) string {
	command[1] = addExtension(command[1])

	if isFileValid(command[1]) {
		index, _, found := isCommand(command, "-var") // Check for the presence of the -var option
		if found {
			content, err := getFileData(command[1])
			if err != nil {
				return Red + "Error when reading file!" + Reset
			}

			contentArray := strings.Split(content, "\n") // Split the file content into an array of lines
			varToRemove := command[index:]               // Variables to remove start from the index of -var

			// Iterate over the content array from the end to avoid index issues while removing
			for i := len(contentArray) - 1; i >= 0; i-- {
				v := contentArray[i]
				if len(v) > 2 {
					for _, z := range varToRemove {
						if strings.Contains(v, z) { // Check if the line contains the variable to remove
							// Remove the variable from the content array
							if i == len(contentArray) {
								contentArray = contentArray[:i]
							} else {
								contentArray = append(contentArray[:i], contentArray[i+1:]...)
							}
						}
					}
				}
			}

			if response := verify(Yellow + "Are you sure you want to remove these variable(s)? (y/n)" + Reset); response {
				// Overwrite the file with the updated content
				if hasOverwrote := overwriteFile(command[1], contentArray); hasOverwrote {
					return Green + "\nVariable(s) removed!" + Reset
				}
			} else {
				return Green + "Variable(s) not removed!" + Reset
			}

		} else {
			return Red + "Incorrect use of the -remove command!\n-help -remove for more info!" + Reset
		}
	} else {
		return Yellow + command[1] + " not found!" + Reset
	}
	return Red + "A problem occurred!" + Reset
}

func get(command []string) string {
	command[1] = addExtension(command[1])

	fileContent, err := getFileData(command[1])
	if err != nil {
		return Red + "Error: File " + command[1] + " doesn't exist or cannot be read!" + Reset
	}

	// Check if a variable name is provided as a command argument
	if len(command) >= 3 {
		command[2] = strings.ToUpper(command[2])    // Convert variable name to uppercase
		return getVariable(command[2], fileContent) // Retrieve the variable value
	} else {
		fmt.Println(Yellow + "Enter variable name:" + Reset)
		var res string
		fmt.Scanln(&res)                     // Read the variable name from user input
		res = strings.ToUpper(res)           // Convert to uppercase
		return getVariable(res, fileContent) // Retrieve the variable value
	}
}

func read(command []string) string {
	command[1] = addExtension(command[1])

	// Attempt to get the content of the specified file
	fileContent, err := getFileData(command[1])
	if err != nil {
		return Red + "Error: File " + command[1] + " doesn't exist or cannot be read!" + Reset
	}

	// Return the content of the file
	return Green + "Here is the content of " + command[1] + ":\n" + Reset + fileContent + "\n"
}

func renameFile(oldName string, newName string) string {
	newName = addExtension(newName)

	if isFileValid(newName) {
		return Red + "A file with the name \"" + newName + "\" already exists!" + Reset
	}

	// Attempt to rename the old file to the new name
	err := os.Rename(oldName, newName)
	if err != nil {
		return Red + "Error: Unable to rename " + oldName + " to " + newName + "!" + Reset
	}

	return Green + oldName + " has been renamed to " + newName + Reset
}

func rename(command []string) string {
	command[1] = addExtension(command[1])

	if isFileValid(command[1]) {
		if len(command) >= 3 { // Check if a new name is provided
			return renameFile(command[1], command[2]) // Rename the file with the provided name
		} else {
			// Prompt the user to enter a new file name
			fmt.Println(Yellow + "Enter a new file name." + Reset)
			var res string
			fmt.Scanln(&res) // Read the new name from user input
			if len(res) > 0 {
				return renameFile(command[1], res) // Rename the file if a valid name is given
			} else {
				return Red + "No name has been given!" + Reset
			}
		}
	} else {
		return Yellow + command[1] + " not found!" + Reset
	}
}

func deleteFile(fileName string) string {
	err := os.Remove(fileName)
	if err != nil {
		return Red + "Error: " + fileName + " has not been deleted!" + Reset
	}
	return Green + fileName + " has been successfully deleted!" + Reset
}

func delete(command []string) string {
	command[1] = addExtension(command[1])

	if isFileValid(command[1]) {
		if _, _, found := isCommand(command, "-v"); found { // Check if the -v option is present
			return deleteFile(command[1]) // Delete the file without confirmation
		} else {
			// Prompt the user for confirmation before deletion
			if response := verify(Yellow + "Are you sure you want to delete " + command[1] + " (y/n)? " + Reset); response {
				return deleteFile(command[1]) // Delete the file if confirmed
			} else {
				return Green + "File not deleted!" + Reset
			}
		}
	} else {
		return Red + command[1] + " not found!" + Reset
	}
}

// HELP COMMAND
func help(command []string) string {
	if len(command) > 1 {
		switch command[1] {
		case "create", "-create":
			fmt.Println(Gray + `
[HELP - CREATE COMMAND - EnvCLI]
Usage: -create [FILE NAME] [OPTIONS]
Options:
	-var []  | List of default variables to add.
	-s       | Skip variable(s) prompt.
	
-> Create a .env file.
			` + Reset)
		case "rename", "-rename":
			fmt.Println(Gray + `
[HELP - RENAME COMMAND - EnvCLI]
Usage: -rename [FILE NAME] [NEW NAME]

-> Rename a file.
			` + Reset)
		case "delete", "-delete":
			fmt.Println(Gray + `
[HELP - DELETE COMMAND - EnvCLI]
Usage: -delete [FILE NAME] [OPTIONS]
Options: 
		-v   | Skip validation.

-> Delete a file.
			` + Reset)
		case "read", "-read":
			fmt.Println(Gray + `
[HELP - READ COMMAND - EnvCLI]
Usage: -read [FILE NAME]

-> Return the content of the .env file.
			` + Reset)
		case "get", "-get":
			fmt.Println(Gray + `
[HELP - GET COMMAND - EnvCLI]
Usage: -get [FILE NAME] [VARIABLE(S)]

-> Return a list of occurrences of the given variable(s).
			` + Reset)
		case "remove", "-remove":
			fmt.Println(Gray + `
[HELP - REMOVE COMMAND - EnvCLI]
Usage: -remove [FILE NAME] -var [VARIABLE(S)]

-> Remove variable(s) from the .env file.
			` + Reset)
		case "add", "-add":
			fmt.Println(Gray + `
[HELP - ADD COMMAND - EnvCLI]
Usage: -add [FILE NAME] -var [VAR(S)]

-> Add variable(s) to a .env file.
			` + Reset)
		case "update", "-update":
			fmt.Println(Gray + `
[HELP - UPDATE COMMAND - EnvCLI]
Usage: -update [FILE NAME] -var [VARIABLE(S)] [OPTIONS]
Options:
	-p       | Create confirmation message for every variable.
	
-> Update .env file variable values.
			` + Reset)
		default:
			fmt.Println(Red + command[1] + " is an unknown command!")
			fmt.Println(Yellow + "List of commands: -help, -create, -read, -get, -delete, -remove" + Reset)
		}
	} else {
		return Gray + `
[HELP - CREATE COMMAND - EnvCLI]
Usage: -create [FILE NAME] [OPTIONS]
Options:
	-var []  | List of default variables to add.
	-s       | Skip variable(s) prompt.
	
-> Create a .env file.
-------------------------------------
[HELP - UPDATE COMMAND - EnvCLI]
Usage: -update [FILE NAME] -var [VARIABLE(S)] [OPTIONS]
Options:
	-p       | Create confirmation message for every variable.
	
-> Update .env file variable values.
-------------------------------------
[HELP - DELETE COMMAND - EnvCLI]
Usage: -delete [FILE NAME] [OPTIONS]
Options: 
		-v   | Skip validation.

-> Delete a file.
-------------------------------------
[HELP - RENAME COMMAND - EnvCLI]
Usage: -rename [FILE NAME] [NEW NAME]

-> Rename a file.
-------------------------------------
[HELP - REMOVE COMMAND - EnvCLI]
Usage: -remove [FILE NAME] -var [VARIABLE(S)]

-> Remove variable(s) from the .env file.
-------------------------------------
[HELP - GET COMMAND - EnvCLI]
Usage: -get [FILE NAME] [VARIABLE(S)]

-> Return a list of occurrences of the given variable(s).
-------------------------------------
[HELP - READ COMMAND - EnvCLI]
Usage: -read [FILE NAME]

-> Return the content of the .env file.
-------------------------------------
[HELP - ADD COMMAND - EnvCLI]
Usage: -add [FILE NAME] -var [VAR(S)]

-> Add variable(s) to a .env file.
					` + Reset
	}
	return ""
}
