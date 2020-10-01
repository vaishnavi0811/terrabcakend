package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

/* func (c App) GetWorkspace() revel.Result {

	config := &tfe.Config{
		Token: "JAq94WOOoR9s3g.atlasv1.jj40mNZ3GFoCGH8FVNrD7xvmQHnHaQmUZwtMKY7UVsXv0eo6f6nCC13BD3WtAyK1VBY",
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	//variables
	workspaceName := "space1"
	orgName := "ClDevTeam"
	ctx := context.Background()
	//autoqueueruns := true
	listOptions := tfe.ListOptions{
		PageSize: 20,
	}
	workspaceOptions := tfe.WorkspaceListOptions{
		ListOptions: listOptions,
		Search:      &workspaceName,
	}
	workspaces, err := client.Workspaces.List(ctx, orgName, workspaceOptions)
	if err != nil {
		log.Fatal(err)
	}
	workspaceID := workspaces.Items[0].ID
	return c.RenderText(workspaceID)
	//fmt.Println("workspace ID: ", workspaceID)
} */
