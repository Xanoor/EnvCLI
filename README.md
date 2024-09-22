
# 🗒️EnvCLI - Enveloppe Command Line Interface

Project created to learn more about Golang and assist with .env file creation.
It may contain some issues; if you find any, please contact me.

Use -help, -man, or -h to access all commands and learn how to use them !
## Language
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
## Run Locally

Clone the project

```bash
  git clone https://github.com/Xanoor/EnvCLI
```

Go to the project directory

```bash
  cd EnvCLI
```

Create the executable:

```bash
  go mod init EnvCLI
  go build
```

You can now run your executable !
## Commands example
Here are some examples of commands for the "test.env" file (use -help for all commands):

Create test.env:
```bash
-create test
```

Add variable(s):
```bash
-add test -var VAR1 VAR2
```

Update variable(s):
```bash
-update test -var VAR1 VAR2
```

Remove variable:
```bash
-remove test -var VAR1 VAR2
```

Delete file:
```bash
-delete test
```
## Support

For support, discord -> xanoor1

