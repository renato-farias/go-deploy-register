package modules

import (
    "fmt"
    "github.com/jinzhu/configor"
)

var config = struct {
    Database struct {
        Driver string `default:"mysql"`
        Host   string `default:"127.0.0.1"`
        Port   int    `default:"3306"`
        User   string `default:"root"`
        Pass   string `default:"password"`
        Base   string `default:"base"`
    }
}{}

func LoadConfiguration(file string) {
    configor.Load(&config, file)
}

func CfgGetDatabaseDriver() string {
    return config.Database.Driver
}

func CfgGetDatabaseHost() string {
    return config.Database.Host
}

func CfgGetDatabasePort() int {
    return config.Database.Port
}

func CfgGetDatabaseUser() string {
    return config.Database.User
}

func CfgGetDatabasePass() string {
    return config.Database.Pass
}

func CfgGetDatabaseBase() string {
    return config.Database.Base
}

func LoadConfig() {
    LoadConfiguration("config/config.yml")

    fmt.Println("")
    fmt.Println("=======================")
    fmt.Println("Database Configuration:")
    fmt.Println("=======================")
    fmt.Printf("Type: %s\n", config.Database.Driver)
    fmt.Printf("Host: %s\n", config.Database.Host)
    fmt.Printf("Port: %d\n", config.Database.Port)
    fmt.Printf("User: %s\n", config.Database.User)
    fmt.Println("Pass: ********")
    fmt.Printf("Base: %s\n", config.Database.Base)
    fmt.Println("")
}

func PrintHelp(prog string) {
    fmt.Printf("Use: %s [option]:\n", prog)
    fmt.Println("")
    fmt.Println("    options list")
    fmt.Println("    help           - To print this message.")
    fmt.Println("    dropschema     - To print this message.")
    fmt.Println("    checkschema    - To print this message.")
    fmt.Println("    createschema   - To print this message.")
}
