package server

import (
  "fmt"
  "io/ioutil"
  "net/http"
  _ "net/http/pprof" // http server profiling
  "os"
  "strings"

  "github.com/joshuaprunier/trite/common"
)

// RunServer receives a port number and a directory path for create definitions output by trite in dump mode and another directory path with an xtrabackup processed with the --export flag
func RunServer(tablePath string, backupPath string, port string) {
  // Make sure directory passed in has trailing slash
  if strings.HasSuffix(backupPath, "/") == false {
    backupPath = backupPath + "/"
  }

  // Ensure the backup has been prepared for transporting with --export
  check := dirWalk(backupPath, false)
  if check == false {
    fmt.Println()
    fmt.Println()
    fmt.Println("It appears that --export has not be run on your backups!")
    fmt.Println()
    fmt.Println()
    os.Exit(1)
  }

  // Start HTTP server listener
  fmt.Println()
  fmt.Println("Starting server listening on port", port)
  http.Handle("/tables/", http.StripPrefix("/tables/", http.FileServer(http.Dir(tablePath))))
  http.Handle("/backups/", http.StripPrefix("/backups/", http.FileServer(http.Dir(backupPath))))
  err := http.ListenAndServe(":"+port, nil)

  // Check if port is already in use
  if err != nil {
    if err.Error() == "listen tcp :"+port+": bind: address already in use" {
      fmt.Println()
      fmt.Println()
      fmt.Println("ERROR: Port", port, "is already in use!")
      fmt.Println()
      fmt.Println()
      os.Exit(1)
    } else {
      common.CheckErr(err)
    }
  }
}

// dirWalk the backup directory and confirm there are .exp files which is proof --export was run
func dirWalk(dir string, flag bool) bool {
  files, ferr := ioutil.ReadDir(dir)
  common.CheckErr(ferr)
  for _, file := range files {
    // Check if file is a .exp, that means --export has been performed on the backup
    _, ext := common.ParseFileName(file.Name())

    // Handle sub dirs recursive function
    if file.IsDir() {
      flag := dirWalk(dir+file.Name()+"/", flag)
      if flag == true {
        return flag
      }
    } else {
      if ext == "exp" {
        flag = true
        break
      }
    }
  }
  return flag
}
