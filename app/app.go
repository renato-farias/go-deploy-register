package main

import (
    "os"
    "fmt"
    "time"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "bytes"
    "reflect"
    "net/http"
    "github.com/gin-gonic/gin"
    "../app/modules"
)

// var config Config
var dbconn *sql.DB

// type Config struct {
//     Database struct {
//         Type string `json:"type",default:"mysql"`
//         Host string `json:"host",default:""`
//         Port int    `json:"port",default:""`
//         User string `json:"user",default:""`
//         Pass string `json:"pass",default:""`
//         Base string `json:"Base",default:""`
//     } `json:"database"`
// }



func DBConnect() {
    param := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", modules.CfgGetDatabaseUser(),
                         modules.CfgGetDatabasePass(),
                         modules.CfgGetDatabaseHost(),
                         modules.CfgGetDatabasePort(),
                         modules.CfgGetDatabaseBase())

    db, err := sql.Open(modules.CfgGetDatabaseDriver(), param)
    dbconn = db

    if err != nil {
        fmt.Println(err.Error())
    }

    err = db.Ping()

    if err != nil {
        fmt.Println(err.Error())
    }
}

func DBClose() {
    defer dbconn.Close()
}

func CreateTable(name string, ansi string) {

    stmt, err := dbconn.Prepare(ansi)

    if err != nil {
        fmt.Println(err.Error())
    }
    _, err = stmt.Exec()
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Printf("%s table successfully created.\n", name)
    }
}

func CheckTable(table string) {
    stmt, err := dbconn.Prepare(fmt.Sprintf("DESC %s;", table))

    if err != nil {
        fmt.Println(err.Error())
    }
    _, err = stmt.Exec()
    if err != nil {
        fmt.Printf("%s table doesn't exist.\n", table)
    } else {
        fmt.Printf("%s table is already exist. Nothing to do.\n", table)
    }
}

func DropTable(table string) {
    stmt, err := dbconn.Prepare(fmt.Sprintf("DROP TABLE %s;", table))

    if err != nil {
        fmt.Println(err.Error())
    }
    _, err = stmt.Exec()
    if err != nil {
        fmt.Printf("%s table could not be dropped.\n", table)
    } else {
        fmt.Printf("%s table was dropped.\n", table)
    }
}

func CreateApplications() {

    name := "applications"
    ansi := "CREATE TABLE IF NOT EXISTS `" + modules.CfgGetDatabaseBase() + "`.`" + name + "` (" +
            "  `id` INT(4) NOT NULL AUTO_INCREMENT," +
            "   `application_name` VARCHAR(30) NULL," +
            "   PRIMARY KEY (`id`))" +
            " ENGINE = InnoDB"
    CreateTable(name, ansi)
}

func CreateDeployments() {

    name := "deployments"
    ansi := "CREATE TABLE IF NOT EXISTS `" + modules.CfgGetDatabaseBase() + "`.`" + name + "` (" +
            " `id` INT(10) NOT NULL AUTO_INCREMENT," +
            " `start` TIMESTAMP NULL," +
            " `end` TIMESTAMP NULL," +
            " `version` VARCHAR(60) NULL," +
            " `status` ENUM('FAILED', 'OK', 'STARTED') NULL DEFAULT 'OK'," +
            " `application_id` INT(4) NOT NULL," +
            " PRIMARY KEY (`id`)," +
            " INDEX `fk_deployments_applications_idx` (`application_id` ASC)," +
            " CONSTRAINT `fk_deployments_applications`" +
            "   FOREIGN KEY (`application_id`)" +
            "   REFERENCES `" + modules.CfgGetDatabaseBase() + "`.`applications` (`id`)" +
            "   ON DELETE NO ACTION" +
            "   ON UPDATE NO ACTION)" +
            "ENGINE = InnoDB"
    CreateTable(name, ansi)
}

func CreateSchema() {
    DBConnect()
    CreateApplications()
    CreateDeployments()
    fmt.Println("")
}

func CheckSchema() {
    DBConnect()
    CheckTable("applications")
    CheckTable("deployments")
    fmt.Println("")
}

func DropSchema() {
    DBConnect()
    DropTable("deployments")
    DropTable("applications")
    fmt.Println("")
}

func in_array(val interface{}, array interface{}) bool {
    exists := false

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                exists = true
                return exists
            }
        }
    }
    return exists
}

func Api() {

    DBConnect()

    type Application struct {
        Id               int    `json:"id"`
        Application_Name string `json:"application_name"`
    }

    type Deployment_List struct {
        Id               int    `json:"id"`
        Version          string `json:"version"`
        Status           string `json:"status"`
        Application_Name string `json:"application_name"`
    }

    type Deployment struct {
        Id          int         `json:"id"`
        Start       time.Time   `json:"start"`
        End         time.Time   `json:"end"`
        Version     string      `json:"version"`
        Status      string      `json:"status"`
        Application Application `json:"application"`
    }

    router := gin.Default()

    authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
        "deploy": "secretpassword",
    }))

    // GET an application detail
    router.GET("/application/:id", func(c *gin.Context) {
        var (
            application Application
        )
        id := c.Param("id")
        row := dbconn.QueryRow(`SELECT id, application_name
                                FROM applications
                                WHERE id = ?;`, id)
        err := row.Scan(&application.Id, &application.Application_Name)
        if err != nil {
            c.IndentedJSON(404, gin.H{
                "error": nil,
            })
        } else {
            c.IndentedJSON(404, application)
        }

    })

    // GET all applications
    router.GET("/applications", func(c *gin.Context) {
        var (
            application  Application
            applications []Application
        )
        rows, err := dbconn.Query(`SELECT id, application_name
                                   FROM applications;`)
        if err != nil {
            fmt.Print(err.Error())
        }
        for rows.Next() {
            err = rows.Scan(&application.Id, &application.Application_Name)
            applications = append(applications, application)
            if err != nil {
                fmt.Print(err.Error())
            }
        }
        defer rows.Close()
        c.IndentedJSON(http.StatusOK, gin.H{
            "applications": applications,
            "count":        len(applications),
        })
    })

    // POST new application
    authorized.POST("/application", func(c *gin.Context) {
        var buffer bytes.Buffer
        application_name := c.PostForm("application_name")
        stmt, err := dbconn.Prepare(`INSERT into applications (application_name)
                                     VALUES(?);`)
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(application_name)

        if err != nil {
            fmt.Print(err.Error())
        }

        buffer.WriteString(application_name)
        defer stmt.Close()
        name := buffer.String()
        c.IndentedJSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf(" %s successfully created.", name),
        })
    })

    // GET all deployments
    router.GET("/deployments", func(c *gin.Context) {
        var (
            deployment  Deployment_List
            deployments []Deployment_List
        )
        rows, err := dbconn.Query(`SELECT d.id, d.version, d.status, a.application_name
                                   FROM deployments d, applications a
                                   WHERE a.id = d.application_id;`)
        if err != nil {
            fmt.Print(err.Error())
        }
        for rows.Next() {
            err = rows.Scan(&deployment.Id, &deployment.Version,
                            &deployment.Status, &deployment.Application_Name)
            deployments = append(deployments, deployment)
            if err != nil {
                fmt.Print(err.Error())
            }
        }
        defer rows.Close()
        c.IndentedJSON(http.StatusOK, gin.H{
            "deployments": deployments,
            "count":       len(deployments),
        })
    })

    // GET a deployment detail
    router.GET("/deployment/:id", func(c *gin.Context) {
        var (
            deployment  Deployment_List
            deployments []Deployment_List
        )
        rows, err := dbconn.Query(`SELECT d.id, d.version, d.status, a.application_name
                                   FROM deployments d, applications a
                                   WHERE a.id = d.application_id;`)
        if err != nil {
            fmt.Print(err.Error())
        }
        for rows.Next() {
            err = rows.Scan(&deployment.Id, &deployment.Version,
                            &deployment.Status, &deployment.Application_Name)
            deployments = append(deployments, deployment)
            if err != nil {
                fmt.Print(err.Error())
            }
        }
        defer rows.Close()
        c.IndentedJSON(http.StatusOK, gin.H{
            "deployments": deployments,
            "count":       len(deployments),
        })
    })

    // POST new deployment
    authorized.POST("/deployment", func(c *gin.Context) {
        start   := time.Now()
        status  := "STARTED"
        version := c.PostForm("version")
        application_id := c.PostForm("application_id")

        stmt, err := dbconn.Prepare(`INSERT into deployments
                                        (start, status, version, application_id)
                                     VALUES (?,?,?,?);`)
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(start, status, version, application_id)

        if err != nil {
            fmt.Print(err.Error())
        }

        defer stmt.Close()

        c.IndentedJSON(http.StatusOK, gin.H{
            "message": "Deployment successfully created.",
        })
    })

    // PUT status deployment
    authorized.PUT("/deployment/:id", func(c *gin.Context) {

        allowed_status := []string {"FAILED", "OK"}

        id := c.Param("id")
        end  := time.Now()
        status  := c.PostForm("status")
        application_id := c.PostForm("application_id")

        if in_array(status, allowed_status) == false {
            c.IndentedJSON(http.StatusBadRequest, gin.H{
                "message": "Wrong status. Please use FAILED or OK",
            })
            return
        }

        stmt, err := dbconn.Prepare(`UPDATE deployments
                                     SET status = ?,
                                         end = ?
                                     WHERE application_id = ?
                                       AND id = ?;`)
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(status, end, application_id, id)

        if err != nil {
            fmt.Print(err.Error())
        }

        defer stmt.Close()

        c.IndentedJSON(http.StatusOK, gin.H{
            "message": "Deployment successfully updated.",
        })
    })

    router.Run(":3000")
}

func main() {

    osargs := os.Args

    modules.LoadConfig()

    if len(osargs) > 1 {

        args := osargs[1:]

        switch arg1 := args[0]; arg1 {
            case "dropschema":
                DropSchema()
            case "checkschema":
                CheckSchema()
            case "createschema":
                CreateSchema()
            case "web":
                Api()
            default:
                modules.PrintHelp(osargs[0])
        }
    } else {
        modules.PrintHelp(osargs[0])
    }

}
