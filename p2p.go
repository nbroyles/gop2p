/**
*
* TODO
* ----
* . specify output file
* . do file transfers in separate thread
* . show current transfers
*
*
* file transfers in a separate thread
* ---
* . open listener in thread and wait for transfers
* . how to interrupt current display with request to accept transfer?
**/

package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

func main() {
    var input int
    for {
//        waitForTransfers()
        printInstructions()
        fmt.Scanf("%d", &input)
        switch input {
            case 1: sendFile()
            case 2: receiveFile()
        }
    }
}

/*func waitForTransfers() {
    go func() {
    }()
}*/

func printInstructions() {
    fmt.Println("GoP2P")
    fmt.Println("-----")
    fmt.Println("[1] Send file")
    fmt.Println("[2] Wait for transfer")
}

func sendFile() {
    fmt.Println("Enter an IP address:")
    var ipAddress string
    fmt.Scanf("%s", &ipAddress)

    fmt.Println("Enter the path to the file:")
    var filePath string
    fmt.Scanf("%s", &filePath)

    // TODO: this port could be in use, then what?
    conn, err := net.Dial("tcp", ipAddress + ":28321")

    // TODO: When reading and writing in chunks, can spit out progress
    // TODO: Also considering firing up send in a separate thread/go routine
    fileHandle, err := os.Open(filePath)
    if err != nil {
        log.Print("File doesn't exist: " + filePath)
    }

    defer fileHandle.Close()

    file := bufio.NewReader(fileHandle)

    buffer := make([]byte, 1024)
    for { 
        bytesRead, _ := file.Read(buffer) 
        if bytesRead == 0 { break }

        conn.Write(buffer)
    }
    conn.Close()
}

func receiveFile() {
    listener, err := net.Listen("tcp", ":28321")
    if err != nil {
        log.Print(err)
        return
    }

    defer listener.Close()

    fmt.Println("Waiting for connection...")
    conn, _ := listener.Accept()

    var accept string
    fmt.Println("Incoming file transfer from " + conn.RemoteAddr().String() + 
        ".  Accept? [yn]")
    fmt.Scanf("%s", &accept)

    if accept != "y" {
        return  
    }

    var fileName string
    fmt.Println("Enter a file name:")
    fmt.Scanf("%s", &fileName)


    // have a connection -- open another thread to wait for transfers
    //waitForTransfers()

    fileHandle, _ := os.Create(fileName)
    if err != nil {
        log.Print("Couldn't create file: " + fileName)
    }

    defer fileHandle.Close()

    file := bufio.NewWriter(fileHandle)

    buffer := make([]byte, 1024)
    for {
        bytesRead, _ := conn.Read(buffer)
        if bytesRead == 0 { break }

        file.Write(buffer)
    }

    file.Flush()
}
