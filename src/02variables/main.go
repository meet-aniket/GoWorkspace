package main

import "fmt"

// constant public variable
const LoginToken string = "random-public-token" // Public Token

func main() {
	var username string = "aniket"
	fmt.Println(username)
	fmt.Printf("Variable is of type: %T \n", username)

	var isLoggedIn bool = true
	fmt.Println(isLoggedIn)
	fmt.Printf("Variable is of type: %T \n", isLoggedIn)

	var smallValue uint8 = 255
	fmt.Println(smallValue)
	fmt.Printf("Variable is of type: %T \n", smallValue)

	var smallFloat float32 = 255.53795473975
	fmt.Println(smallFloat)
	fmt.Printf("Variable is of type: %T \n", smallFloat)

	// default values and some aliases
	var anotherVariable int
	var anotherVariableString string
	fmt.Println(anotherVariable)
	fmt.Println(anotherVariableString)
	fmt.Printf("Variable is of type: %T \n", anotherVariable)

	// implicit type
	var website = "test-website.com"
	fmt.Println("Website is: ", website)

	// no var style
	numberOfUsers := 100000
	fmt.Println("NumberOfUsers are: ", numberOfUsers)
}
