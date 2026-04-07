# bigfun
## A installer written in Go. This loader automates the download of this mods, configurations, and OneConfig settings directly from the GitHub repository to your `.minecraft` folder.

### How to Build

To generate the executable yourself, ensure you have [Go](https://go.dev/dl/) installed.
1. Open your terminal (CMD/PowerShell) in the project folder.
2. Initialize the module (if not already done).
3. Compile the Loader.
   ```bash
   go mod init loader
   go build -ldflags="-s -w" -o Loader.exe main.go
   ```
  
