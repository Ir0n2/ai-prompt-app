package main

import (
        //"io"
        "fmt"
        "log"
        "net/http"
	"bufio"
        "strings"
        "os"
	//"os/exec"
        "github.com/blackestwhite/gopenai"
)

var global string

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
                input := r.FormValue("input")
                ip := r.RemoteAddr
		global = promptGPT(input)
		fmt.Println("Ip address = ", ip, "\n")
                //fmt.Println("Username = ", name, "\n")
                fmt.Println("Prompt = ", input, "\n")
		
		http.ServeFile(w, r, "Index.html")


		//http.ServeFile(w, r, "/home/zerocool/monkey/monkeyBuisness.html")
	} else {
                 http.ServeFile(w, r, "Index.html")		
	}

}

//checks if file exists
func doesFileExist(filename string) {

        _, err := os.Stat(filename)

        if err != nil {
                fmt.Println("log file doesnt exist")
                //return false
        } else {
                fmt.Println("fuck you log file exists")
                deleteFile("testlogfile")
		fmt.Println("left over log file deleted")
		//return true
        }


}

//takes input like scanln but with spaces
func getInput() string {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan() // use `for scanner.Scan()` to keep reading
    line := scanner.Text()

    return line
    //fmt.Println("captured:",line)
}

//One final function so everything works within one function call
//Now it returns the GPT output as a string. party party
func promptGPT(prompt string) string {
	
	doesFileExist("testlogfile")
	prompter(prompt)
	fuck := format()
	return fuck
}

func format() string {
	out := ""	
        num := printLines("testlogfile")
	for i := 0; num-2 >= i; i++ {
		if i == 0 {
			continue
		}			
		str := readList("testlogfile", i)
		//fmt.Println(str)
		//this if is to get around an error
		if str == "" {
			continue
		}
		//this if is also to get around errors
		if str == "} }]}" {
                        continue
                }

                //str is a line from testlogfile and test is str but split up at every curly bracket
                test := strings.Split(str, "{")

                //res2 removes the remaining bracket garbage.
                res2 := strings.TrimSuffix(test[3], "} }]}")

                //This will add whatever word is held by res2 to the out string, which is short for output. It's the chatgpt output.
		out += res2		
		//fmt.Println(out)
        }
	
	//delete the log file
	deleteFile("testlogfile")
	return out
}

func check(e error) {

        if e != nil {
                panic(e)
        }

}

//increments i through the list so i should = the amount of items in the list.
func printLines(path string) int {

        i := 0
        filePath := path
        readFile, err := os.Open(filePath)
        check(err)

        fileScanner := bufio.NewScanner(readFile)
        fileScanner.Split(bufio.ScanLines)
        var fileLines []string

        for fileScanner.Scan() {
                fileLines = append(fileLines, fileScanner.Text())
        }

        readFile.Close()

        var line string
        //this print statement is here because we need to do something with line or else it wont compile.
        fmt.Println(line)
        for _, line = range fileLines {

                i++

        }
        return i
}

//Reads list from file
func readList(path string, index int) string {

        filePath := path
        readFile, err := os.Open(filePath)

        check(err)

        fileScanner := bufio.NewScanner(readFile)
        fileScanner.Split(bufio.ScanLines)
        var fileLines []string

        for fileScanner.Scan() {
                fileLines = append(fileLines, fileScanner.Text())
        }

        readFile.Close()

        return fileLines[index]

}
//deletes file
func deleteFile(path string) {

        e := os.Remove(path)
        if e != nil {
                log.Fatal(e)
        }

}

func prompter(prompt string) {

	doesFileExist("testlogfile")

	key := "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"

        instance := gopenai.Setup(key)

        p := gopenai.ChatCompletionRequestBody{
                Model: "gpt-3.5-turbo",
                Messages: []gopenai.Message{
                        {Role: "user", Content: prompt},
                },
                Stream: true,
        }

        resultCh, err := instance.GenerateChatCompletion(p)
        if err != nil {
                log.Fatal(err)
        }

        f, err := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
        if err != nil {
            log.Fatalf("error opening file: %v", err)
        }
        defer f.Close()

        log.SetOutput(f)

        for chunk := range resultCh {

                log.Println(chunk)
        }
}

func result(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, global)

}

func main() {
	
	http.HandleFunc("/result", result)
	http.HandleFunc("/", login)
	http.ListenAndServe(":8080", nil)

	/*cmd := exec.Command("firefox", "localhost:8080")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}*/

}
