/**
*
* TODO
* ----
* . fix bug with reading input in goroutine (it's empty -- use channels?)
* . show current transfers
* . communicate p2p without needing to know port...possible?
*
* file transfers in a separate thread
* ---
* . as currently implemented, incoming transfer will interrupt
*   whatever's going on at the time
**/

package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "strconv"
)

var config struct {
    recvPort int
}

func main() {
    // NOTE: This whole port detection thing is bad.  As currently structured,
    // 1 port needed for each file transfer.  Think there is a way to share a
    // port?
    config.recvPort = getPortForTransfers()
    if config.recvPort == -1 {
        fmt.Println("Couldn't find open port for transfers. Exiting")
        return
    } else {
        fmt.Println("Listening for transfers on port " +
            strconv.Itoa(config.recvPort))
    }

    waitForTransfers()

    var input int
    for {
        printInstructions()
        fmt.Scanf("%d", &input)
        switch input {
            case 1: sendFile()
            case 2: showCurrentTransfers()
        }
    }
}

func getPortForTransfers() int {
    for checkPort := 28321; checkPort < 28332; checkPort++ {
        if isPortOpen(checkPort) {
            return checkPort
        }
    }

    return -1
}

func waitForTransfers() {
    go receiveFile();
}

func printInstructions() {
    fmt.Println("GoP2P")
    fmt.Println("-----")
    fmt.Println("[1] Send file")
    fmt.Println("[2] Show current transfers")
}

func sendFile() {
    fmt.Println("Enter an IP address (incl. port):")
    var ipAddress string
    fmt.Scanf("%s", &ipAddress)

    fmt.Println("Enter the path to the file:")
    var filePath string
    fmt.Scanf("%s", &filePath)

    conn, err := net.Dial("tcp", ipAddress)

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
    listener, err := net.Listen("tcp", ":" + strconv.Itoa(config.recvPort))
    if err != nil {
        log.Print(err)
        return
    }

    defer listener.Close()

    conn, _ := listener.Accept()

    var accept string
    fmt.Println("Incoming file transfer from " + conn.RemoteAddr().String() +
        ".  Accept? [yn]")
    fmt.Scanf("%s", &accept)
    fmt.Println("---------------")
    fmt.Print(accept)
    fmt.Println("---------------")

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

func isPortOpen(port int) bool {
    listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))

    if err != nil {
        return false
    } else {
        listener.Close()
        return true
    }
}

func showCurrentTransfers() {

}
