package functions

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"terracloud/app/templates"

	"github.com/iancoleman/strcase"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Designs struct {
	ID            string
	Resourceid    string
	Dependson     pq.StringArray
	Configuration interface{}
}

const tagName = "required"
const (
	host     = "13.68.247.155"
	port     = 5432
	user     = "postgres"
	password = "October201!9"
	dbname   = "meg"
)

func Collect(designID int) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	//var designs []Designs
	var d Designs
	fmt.Println("Successfully connected!")

	//userSql := "SELECT resources FROM designs WHERE id = $1".
	getResource := "select json_array_elements(resource -> 'resources')->>'id' as id,json_array_elements(resource -> 'resources')->>'resourceid' as resourceid,json_array_elements(resource -> 'resources')->'dependson' as dependson from designs where id = $1"
	rows, err := db.Query(getResource, designID)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&d.ID, &d.Resourceid, pq.Array(&d.Dependson))
		if err != nil {
			log.Println("Failed to execute query: ", err)
		}
		inputs := reflect.ValueOf(d)
		typeofinput := inputs.Type()
		for i := 0; i < inputs.NumField(); i++ {
			moduleKey := strcase.ToSnake(typeofinput.Field(i).Name)
			moduleValue := inputs.Field(i).Interface()
			valueType := inputs.Field(i).Kind()
			log.Println(moduleKey, moduleValue, valueType)
		}
	}
}

func CreateAzureVM(mvmvars *templates.MVMVARS, terraformfile string) error {
	var subscription_id interface{}
	file, err := os.Create(terraformfile)
	if err != nil {
		return err
	}
	fmt.Fprintf(file, "variable \"client_id\" {\n\t type = string \n}")
	fmt.Fprintf(file, "\nvariable \"client_secret\" {\n\t type = string \n}")
	fmt.Fprintf(file, "\nvariable \"tenant_id\" {\n\t type = string \n}")
	fmt.Fprintf(file, "\nvariable \"subscription_id\" {\n\t type = string \n}\n")
	fmt.Fprintf(file, "\nmodule \"vm\" {\n")
	fmt.Fprintf(file, "\tversion =  \"1.0.4\"\n")
	fmt.Fprintf(file, "\tsource =  \"app.terraform.io/ClDevTeam/vm/azurerm\"\n")

	defer file.Close()
	inputs := reflect.ValueOf(*mvmvars)
	typeofinput := inputs.Type()
	log.Println("User inputs")
	for i := 0; i < inputs.NumField(); i++ {
		moduleKey := strcase.ToSnake(typeofinput.Field(i).Name)
		moduleValue := inputs.Field(i).Interface()
		log.Printf("%s = %v", moduleKey, moduleValue)
		if (inputs.Field(i).Kind() == reflect.Slice && inputs.Field(i).IsNil()) || (inputs.Field(i).Kind() == reflect.String && inputs.Field(i).IsZero()) || (inputs.Field(i).Kind() == reflect.Int && inputs.Field(i).IsZero()) {
			log.Printf("%s is empty or null", typeofinput.Field(i).Name)
			continue
		} else {
			if moduleKey == "os_data_disk_size_in_gb" {
				fmt.Fprintf(file, "\t%s = %v\n", moduleKey, moduleValue)
			} else if moduleKey == "data_disks" {
				len := inputs.Field(i).Len()
				value := moduleValue.([]int)
				for j := 0; j < len; j++ {
					if j == 0 {
						fmt.Fprintf(file, "\t%s = [ %v,", moduleKey, value[j])
					} else if j+1 == len {
						fmt.Fprintf(file, "%v ]\n", value[j])
					} else {
						fmt.Fprintf(file, "%v ,", value[j])
					}
				}
			} else if moduleKey == "tags" {
				fmt.Fprintf(file, "\t%s = { ", moduleKey)
				for key, element := range moduleValue.(map[string]string) {
					fmt.Fprintf(file, "\n\t\t%s = \"%s\"", key, element)
				}
				fmt.Fprintf(file, "\n\t")
			} else if moduleKey == "subscription_id" {
				subscription_id = moduleValue
				continue
			} else {
				fmt.Fprintf(file, "\t%s = \"%s\"\n", moduleKey, moduleValue)
			}
		}

	}
	fmt.Fprintf(file, "}\n")
	fmt.Fprintf(file, "\nprovider \"azurerm\" {\n")
	fmt.Fprintf(file, "\t version =  \"=2.4.0\"\n")
	fmt.Fprintf(file, "\t client_id = var.client_id\n")
	fmt.Fprintf(file, "\t client_secret = var.client_secret\n")
	fmt.Fprintf(file, "\t subscription_id = \"%v\"\n", subscription_id)
	fmt.Fprintf(file, "\t tenant_id = var.tenant_id\n")
	fmt.Fprintf(file, "\t features {}\n")
	fmt.Fprintf(file, "}\n")
	return nil
}
