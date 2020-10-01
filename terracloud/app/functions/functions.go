package functions

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

type SecureParams struct {
	UserToken string
}

func GetWorkspaceID(ctx context.Context, client *tfe.Client, workspaceName string, orgName string) (string, error) {
	listOptions := tfe.ListOptions{
		PageSize: 20,
	}
	workspaceOptions := tfe.WorkspaceListOptions{
		ListOptions: listOptions,
		Search:      &workspaceName,
	}
	workspaces, err := client.Workspaces.List(ctx, orgName, workspaceOptions)
	workspaceID := workspaces.Items[0].ID
	return workspaceID, err
	//fmt.Println("workspace ID: ", workspaceID)
}
func GetPlanID(ctx context.Context, client *tfe.Client, workspaceID string, configVersionID string) (string, string) {
	var planID, runID string
	//fmt.Println("enter get run id")
	listOptions := tfe.ListOptions{
		PageSize: 20,
	}
	runlistOptions := tfe.RunListOptions{
		ListOptions: listOptions,
	}

	for {
		getRuns, err := client.Runs.List(ctx, workspaceID, runlistOptions)
		requiredRuns := getRuns.Items
		//fmt.Println(requiredRuns)
		if err != nil {
			return "Error in listing plans", err.Error()
		}
		for _, run := range requiredRuns {
			if run.ConfigurationVersion.ID == configVersionID && run.Status == "planned" {
				runID = run.ID
				planID = run.Plan.ID
				//runStatus = "planned"
				break
			} else if run.ConfigurationVersion.ID == configVersionID && run.Status == "errored" {
				runID = run.ID
				planID = run.Plan.ID
				//runStatus = "errored"
				break
			} else if run.ConfigurationVersion.ID == configVersionID && run.Status == "discarded" {
				runID = run.ID
				planID = run.Plan.ID
				//runStatus = "discarded"
				break
			}
		}
		if planID != "" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return planID, runID
}
func GetRun(ctx context.Context, client *tfe.Client, runID string) (*tfe.Run, error) {
	fmt.Println("we are in get run")
	runStruct, err := client.Runs.Read(ctx, runID)
	return runStruct, err
}
func GetRuns(ctx context.Context, client *tfe.Client, workspaceID string) (*tfe.RunList, error) {
	listOptions := tfe.ListOptions{
		PageSize: 20,
	}
	runlistOptions := tfe.RunListOptions{
		ListOptions: listOptions,
	}
	getRuns, err := client.Runs.List(ctx, workspaceID, runlistOptions)

	return getRuns, err
}
func GetApplyLog(ctx context.Context, client *tfe.Client, applyID string) (string, error) {
	getApplyLogs, err := client.Applies.Logs(ctx, applyID)
	applyLogByte, err := ioutil.ReadAll(getApplyLogs)
	applyLog := string(applyLogByte)
	return applyLog, err
}
func PrintPlan(ctx context.Context, client *tfe.Client, planID string) string {
	getPlanLog, err := client.Plans.Logs(ctx, planID)
	planLogByte, err := ioutil.ReadAll(getPlanLog)
	if err != nil {
		return err.Error()
	}
	planLog := string(planLogByte)
	return planLog
}
func Apply(ctx context.Context, client *tfe.Client, applyComment *string, runID string) error {
	fmt.Println("enter apply block")
	runApplyOptions := tfe.RunApplyOptions{
		Comment: applyComment,
	}
	err := client.Runs.Apply(ctx, runID, runApplyOptions)
	if err != nil {
		return err
	} else {
		return nil
	}
}
func createConfigVersion(ctx context.Context, client *tfe.Client, autoqueueruns *bool, workspaceID string, uploadPATH string) string {
	configurationVersionOptions := tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: autoqueueruns,
	}
	newConfigurationVersion, err := client.ConfigurationVersions.Create(ctx, workspaceID, configurationVersionOptions)
	if err != nil {
		log.Fatal(err)
	}
	UploadURL := newConfigurationVersion.UploadURL
	configVersionID := newConfigurationVersion.ID
	err = client.ConfigurationVersions.Upload(ctx, UploadURL, uploadPATH)
	if err != nil {
		fmt.Println("the upload block returns error")
		log.Fatal(err)
	}
	//fmt.Println("your run is created")
	return configVersionID
}
func CreateVariables(ctx context.Context, client *tfe.Client, workspaceID string, variableCreateOptions *tfe.VariableCreateOptions) (*tfe.Variable, error) {
	//create terraform variable
	variable, err := client.Variables.Create(ctx, workspaceID, *variableCreateOptions)
	return variable, err
}
func CreateWorkspace(ctx context.Context, client *tfe.Client, org string, workspaceCreateOptions *tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	workspace, err := client.Workspaces.Create(ctx, org, *workspaceCreateOptions)
	return workspace, err
}
func ConfigAndPlan(ctx context.Context, client *tfe.Client, configVersionOptions *tfe.ConfigurationVersionCreateOptions, workspaceID string, uploadPATH string) (string, error) {
	/* configurationVersionOptions := tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: autoqueueruns,
	} */
	/* var filename string
	if err != nil {
		//fmt.Println("the upload block returns error")
		log.Fatal(err)
	}

	for _, file := range files {
		filename = uploadPATH + "/" + file.Name()
	} */
	newConfigurationVersion, err := client.ConfigurationVersions.Create(ctx, workspaceID, *configVersionOptions)
	if err != nil {
		return "", err
	}
	UploadURL := newConfigurationVersion.UploadURL
	configVersionID := newConfigurationVersion.ID
	err = client.ConfigurationVersions.Upload(ctx, UploadURL, uploadPATH)
	if err != nil {
		//fmt.Println("the upload block returns error")
		return "", err
	}
	//fmt.Println("your run is created")

	return configVersionID, err
}
func WriteFileToDisk(file string, filecontent map[string]interface{}, filepath string) error {
	err := os.MkdirAll(filepath, 766)
	if err != nil {
		return err
	} else {
		terraformfile := filepath + file
		terraformjson, err := json.MarshalIndent(filecontent, "", "")
		if err != nil {
			return err
		} else {
			ioutil.WriteFile(terraformfile, terraformjson, os.ModePerm)
			return err
		}
	}
	/* if err != nil {
		return err
	} else {
		localfile, err := os.Create(path.Join(filepath, filename))
		if err != nil {
			return err
		}
		_, err = localfile.Write(file)
		if err != nil {
			return err
		}
		return nil
	} */
}
func Gzip(terraformfile string) string {

	f, _ := os.Open(terraformfile)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	terraformfile = strings.Replace(terraformfile, ".tf", ".tar.gz", -1)
	f, _ = os.Create(terraformfile)
	w := gzip.NewWriter(f)
	w.Write(content)
	w.Close()
	return terraformfile
}
