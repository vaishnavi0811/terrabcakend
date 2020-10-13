package controllers

import (
	"context"
	"log"
	"os"
	"terracloud/app/functions"
	"terracloud/app/templates"

	"github.com/hashicorp/go-tfe"
	"github.com/revel/revel"
)

type Deployment struct {
	*revel.Controller
	apply *templates.ApplyPlan
	token string
}

func (c Deployment) ConfigAndPlan(workspaceID string) revel.Result {
	path, err := os.Getwd()
	filepath := path + "\\" + workspaceID + "\\"
	//filepath := "/home/ubuntu/terraform/apicall"

	token := os.Getenv("userToken")
	////userToken := c.Request.Header.Get("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return c.RenderText(err.Error())
	}
	ctx := context.Background()
	autoQueueRuns := true
	speculative := false
	configVersionOptions := tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: &autoQueueRuns,
		Speculative:   &speculative,
	}
	configVersionID, err := functions.ConfigAndPlan(ctx, client, &configVersionOptions, workspaceID, filepath)
	if err != nil {
		return c.RenderText(err.Error())
	}
	planID, runID := functions.GetPlanID(ctx, client, workspaceID, configVersionID)
	runDetails := make(map[string]string)
	runDetails["runID"] = runID
	runDetails["planID"] = planID
	runDetails["configVersionID"] = configVersionID
	return c.RenderJSON(runDetails)
}
func (c Deployment) PrintPlan(planID string) revel.Result {
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	message := functions.PrintPlan(ctx, client, planID)
	return c.RenderText(message)
}
func (c Deployment) ApplyPlan(runID string) revel.Result {
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return c.RenderText(err.Error())
	}
	ctx := context.Background()
	c.Params.BindJSON(&c.apply)
	applyComment := c.apply.ApplyMessage
	err = functions.Apply(ctx, client, &applyComment, runID)
	if err != nil {
		log.Fatal(err)
		return c.RenderText(err.Error())
	}
	return c.RenderText("Plan Applied")
}
func (c Deployment) CreateVariable(org string, workspaceName string) revel.Result {
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")

	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	var variableOptions *tfe.VariableCreateOptions
	c.Params.BindJSON(&variableOptions)
	workspaceID, err := functions.GetWorkspaceID(ctx, client, workspaceName, org)
	variable, err := functions.CreateVariables(ctx, client, workspaceID, variableOptions)
	if err != nil {
		c.Response.Status = 409
		return c.RenderText(err.Error())
	}
	c.Response.Status = 201
	return c.RenderJSON(variable)
}
func (c Deployment) CreateWorkspace(org string, workspaceName string) revel.Result {
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	var workspaceCreateOptions *tfe.WorkspaceCreateOptions
	c.Params.BindJSON(&workspaceCreateOptions)
	workspace, err := functions.CreateWorkspace(ctx, client, org, workspaceCreateOptions)
	if err != nil {
		c.Response.Status = 409
		return c.RenderText(err.Error())
	}
	c.Response.Status = 201
	return c.RenderJSON(workspace)
}
func (c Deployment) GetWorkspace(org string, workspaceName string) revel.Result {
	//var secureParams *functions.SecureParams
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	workspaceID, err := functions.GetWorkspaceID(ctx, client, workspaceName, org)
	if err != nil {
		c.Response.Status = 404
		return c.RenderText(err.Error())
	}
	return c.RenderText(workspaceID)
}
func (c Deployment) GetRuns(org string, workspaceName string) revel.Result {
	//var secureParams *functions.SecureParams
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	workspaceID, err := functions.GetWorkspaceID(ctx, client, workspaceName, org)
	runs, err := functions.GetRuns(ctx, client, workspaceID)
	if err != nil {
		c.Response.Status = 404
		return c.RenderText(err.Error())
	}
	return c.RenderJSON(runs)
}
func (c Deployment) GetVars(org string, workspaceName string) revel.Result {
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	workspaceID, err := functions.GetWorkspaceID(ctx, client, workspaceName, org)
	vars, err := functions.GetVariables(ctx, client, workspaceID)
	if err != nil {
		c.Response.Status = 404
		return c.RenderText(err.Error())
	}
	c.Response.Status = 200
	return c.RenderJSON(vars)
}
func (c Deployment) GetRun(runID string) revel.Result {
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err.Error())
		return c.RenderText(err.Error())
	}
	ctx := context.Background()
	run, err := functions.GetRun(ctx, client, runID)
	if err != nil {
		log.Fatal(err.Error())
		return c.RenderText(err.Error())
	}
	return c.RenderJSON(run)
}
func (c Deployment) PrintApplyLog(runID string) revel.Result {
	//userToken := c.Request.Header.Get("userToken")

	token := os.Getenv("userToken")
	config := &tfe.Config{
		Token: token,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err.Error())
		return c.RenderText(err.Error())
	}
	ctx := context.Background()
	run, err := functions.GetRun(ctx, client, runID)
	applyID := run.Apply.ID
	applyLog, err := functions.GetApplyLog(ctx, client, applyID)
	if err != nil {
		log.Fatal(err.Error())
		return c.RenderText(err.Error())
	}
	return c.RenderText(applyLog)
}
