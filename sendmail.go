package main

import (
	"github.com/go-ini/ini"
	"io/ioutil"
	"io"
	"strings"
	"os/exec"
	"os"
	"time"
	"log"
	"errors"
	"encoding/base64"
    "gopkg.in/gomail.v2"
)
type params struct {
	server string
	port int
	login string
	password string
	fromAddr string
	fromName string
	toAddr string
	toName string
	subject string
	body string
	directory string
	name string
	archive string
	archivator string
}

func getINI(filename string) (p params, err error) {
	cfg, err := ini.Load(filename)
    if err != nil {
        return
	}
	p.server = cfg.Section("smtp").Key("server").String()
	p.port, err = cfg.Section("smtp").Key("server").Int()
	if err != nil {
		p.port = 25
		err = nil
	}
	p.login = cfg.Section("smtp").Key("login").String()
	p.password = cfg.Section("smtp").Key("password").String()
	if !strings.HasPrefix(p.password, "?b") {

	}
	p.fromAddr = cfg.Section("mail").Key("from_addr").String()
	p.fromName = cfg.Section("mail").Key("from_name").String()
	p.toAddr = cfg.Section("mail").Key("to_addr").String()
	p.toName = cfg.Section("mail").Key("to_name").String()
	p.subject = cfg.Section("mail").Key("subject").String()
	p.body = cfg.Section("mail").Key("body").String()
	p.directory = cfg.Section("paths").Key("directory").String()
	p.name = cfg.Section("paths").Key("archive_name").String()
	p.archive = cfg.Section("paths").Key("archive_type").String()
	p.archivator = cfg.Section("paths").Key("7z").String()
	if strings.HasPrefix(p.password, "?b") {
		p.password, err = decodePassword(p.password)
	} else {
		cfg.Section("smtp").Key("password").SetValue(encodePassword(p.password))
		cfg.SaveTo(filename)
	}
	return
}
func encodePassword(s string) string {
	return "?b" + base64.StdEncoding.EncodeToString([]byte(s))
}

func decodePassword(s string) (string, error) {
	result, err := (base64.StdEncoding.DecodeString(s[2:]))
	return string(result), err
}

func getCurrentDate(separate string) string {
	if separate == "" || separate == "-" {
		return time.Now().Format("2006" + separate + "01" + separate + "02")
	}
	return time.Now().Format("02" + separate + "01" + separate + "2006")
}

func getCurrentTime(separate string) string {
	return time.Now().Format("15" + separate + "04")
}

func send(p params) (err error) {
	m := gomail.NewMessage()
    m.SetAddressHeader("From", p.fromAddr, p.fromName)
	m.SetAddressHeader("To", p.toAddr, p.toName)
    m.SetHeader("Subject", p.subject)
    m.SetBody("text/plain", p.body)
	if p.name != "" {
		m.Attach(p.name + "." + p.archive)
	}
	d := gomail.NewPlainDialer(p.server, p.port, p.login, p.password)
	err = d.DialAndSend(m)
	return
}
func replaceDateTime(s, dateSep, timeSep string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "<time>", getCurrentTime(timeSep)),
													"<date>", getCurrentDate(dateSep))
}
func createArchive(p params) (err error) {
	if p.archivator != "" && p.directory != "" && p.archive != "" {
		if _, err = os.Stat(p.directory); os.IsNotExist(err) {
			return
		}
		cmd := exec.Command(p.archivator, "a", "-t" + p.archive, "-mx7", p.name + "." + p.archive, p.directory)
		err = cmd.Run()
	} else {
		err = errors.New("The required fields are not filled in the ini-file")
	}
	return
}

func main() {
	l := log.New(os.Stdout,"",log.Ldate|log.Ltime|log.Lshortfile)
	if _, err := os.Stat("error.log"); os.IsNotExist(err) {
		file, err := os.Create("error.log")
		if err == nil {
			file.Close()
		}
	}
	file, err := os.OpenFile("error.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		defer file.Close()
		writer := io.Writer(file)
		l = log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	filenames, err := ioutil.ReadDir(".")
	if err != nil {
		l.Println(err)
		os.Exit(-1)
	}
	for _, filename := range filenames {
		if strings.HasSuffix(filename.Name(), ".ini") {
			p, err := getINI(filename.Name())
			if err != nil {
				l.Println(err)
				os.Exit(-1)
			}
			p.name = replaceDateTime(p.name, "-", "-")
			p.subject = replaceDateTime(p.subject, ".", ":")
			p.body = replaceDateTime(p.body, ".", ":")
			if err := createArchive(p); err != nil {
				l.Println(err)
				os.Exit(-1)
			}
			defer os.Remove(p.name + "." + p.archive)
			if err := send(p); err != nil {
				l.Println(err)
				os.Exit(-1)
			}
		}
	}
}