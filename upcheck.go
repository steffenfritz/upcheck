package main

import (
    "bufio"
    "log"
    "net"
    _ "net/smtp"
    "os"
    "time"
    "crypto/tls"
    "gopkg.in/gomail.v2"
)

// globals
var smtpUser string = ""
var smtpPass string = ""
var smtpHost string = ""
var smtpPort int = 25

// main function
func main(){

 secs := 3
 timeOut := time.Duration(secs) * time.Second

 hosts := readSourceFile("hostfile.txt")

 for _, host := range hosts {
    checkUp(host, timeOut)
 }
}


// readSourceFile with hosts and ports. Format per line HOST:PORT 
func readSourceFile(hostFile string) []string{

    f, err := os.Open(hostFile)
    if err != nil {
        log.Fatal("Could not read hosts file. Qutting.")
    }   

    defer f.Close()

    var hosts []string
    scanner := bufio.NewScanner(f) 
    for scanner.Scan() {
        hosts = append(hosts, scanner.Text())
    }
    return hosts
}

// check if a connection can be established 
func checkUp(host string, timeOut time.Duration) {
    _, err := net.DialTimeout("tcp", host, timeOut)


    if err != nil {
        println("CRITICAL. Service " + host + " not running.")
        sendMail(host)
    }else {
        println("OK. Service is running.")
    }
}

// sendMail if connection fails
func sendMail(host string) {
    m := gomail.NewMessage()
    m.SetHeader("From", "")
    m.SetHeader("To", "")
    m.SetHeader("Subject", "* * * ALARM * * *")
    m.SetBody("text/html", "CRITICAL! " + host + "ist nicht verfügbar!")


    d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
    d.TLSConfig = &tls.Config{InsecureSkipVerify: true}


    if err := d.DialAndSend(m); err != nil {
        log.Fatal(err)
    }
}
