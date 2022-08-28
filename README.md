[![Go Report Card](https://goreportcard.com/badge/github.com/mskrha/imap-checker)](https://goreportcard.com/report/github.com/mskrha/imap-checker)

## imap-checker

### Description
Simple tool to check IMAP server and print number of unread and all messages for specified accounts.

### Important notice
It is required to checkout commit 4132e15 of library github.com/emersion/go-sasl, because after this commit the XOAUTH2 support was removed.

### Build
```shell
git clone https://github.com/mskrha/imap-checker.git
cd imap-checker
make
```

### Configuration
```shell
mkdir ~/.imap-checker
cp /etc/imap-checker/config.json ~/.imap-checker/
```

### Usage
```shell
$ imap-checker
Account 1: 0/0, Account 2: 7/12, Account 3: 2/5
```
