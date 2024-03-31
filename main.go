/*
Author: PlaidSnowFrog
Date  : March 31 2024
*/


package main

import (
    "fmt"
    "time"
    "os"
    "io/ioutil"
    "path/filepath"
    "strconv"
    "os/exec"
    "runtime"
)

func main() {
    rawTime := time.Now()
    day := rawTime.Day()

    log, logErr := readLog()

    if logErr != nil {
        fmt.Println("Error:", logErr)
        return
    }

    if day == 1 || !log {
        urlErr := open("https://www.teacherspayteachers.com/items/download_items_stats")
        if urlErr != nil {
            fmt.Println("Error while opening url:", urlErr)
            return
        }
    }

    err := writeBat()
    if err != nil {
        fmt.Println("Error while adding to startup:", err)
    }
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}

func writeBat() error {
    ex, err := os.Executable()
    if err != nil {
        return err
    }

    startupPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
    batPath := filepath.Join(startupPath, "DownloadTpT.bat")
    batContent := fmt.Sprintf(`start "" "%s"`, ex)

    err = ioutil.WriteFile(batPath, []byte(batContent), 0644)
    return err
}

func writeLog(definer int) error {
    ex, err := os.Executable()
    if err != nil {
        return err
    }

    writePath := filepath.Join(filepath.Dir(ex), "tpt.log")

    err = ioutil.WriteFile(writePath, []byte(strconv.Itoa(definer)), 0644)
    return err
}

func readLog() (bool, error) {
    ex, err := os.Executable()
    if err != nil {
        return false, err
    }

    logPath := filepath.Join(filepath.Dir(ex), "tpt.log")

    data, err := os.ReadFile(logPath)
    if err != nil {
        return false, fmt.Errorf("error reading log file: %v", err)
    }

    strData := string(data)

    intData, err := strconv.Atoi(strData)
    if err != nil {
        return false, fmt.Errorf("error converting string to integer: %v", err)
    }

    switch intData {
    case 0:
        return false, nil
    case 1:
        return true, nil
    default:
        return false, fmt.Errorf("unexpected value in log file: %d", intData)
    }
}
