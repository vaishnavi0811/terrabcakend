package controllers

import (
	"context"
	"log"
	"os"
	"terracloud/app/functions"
	"terracloud/app/templates"

	"github.com/go-playground/validator"
	"github.com/hashicorp/go-tfe"
	"github.com/revel/revel"
)

type Convert struct {
	*revel.Controller
	mvm *templates.MVMVARS
	res *templates.Resources
	des *templates.Designs
}

func (c Convert) Collect() revel.Result {
	//c.Params.BindJSON(&c.res)
	designID := 48
	functions.Collect(designID)
	return nil
}
func (c Convert) AzureWindowsVM(workspaceName string, org string) revel.Result {
	userToken := c.Request.Header.Get("userToken")
	config := &tfe.Config{
		Token: userToken,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		return c.RenderText(err.Error())
	}
	ctx := context.Background()
	workspaceID, err := functions.GetWorkspaceID(ctx, client, workspaceName, org)
	if err != nil {
		return c.RenderText(err.Error())
	}
	path, err := os.Getwd()
	terraformfile := path + "\\" + workspaceID + "\\main.tf"

	c.Params.BindJSON(&c.mvm)
	v := validator.New()
	err = v.Struct(c.mvm)
	errlist := make(map[string][]string)
	var errors []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
			errors = append(errors, e.Field())
			errlist["errors"] = errors
		}
		return c.RenderText(err.Error())
	}
	err = functions.CreateAzureVM(c.mvm, terraformfile)
	if err != nil {
		return c.RenderText(err.Error())
	}
	gzipfile := functions.Gzip(terraformfile)
	log.Println(c.mvm)

	return c.RenderText(gzipfile)
}
