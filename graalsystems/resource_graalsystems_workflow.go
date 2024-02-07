package graalsystems

import (
	"context"
	"encoding/json"
	"fmt"
	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"slices"
)

// resourceGraalSystemsWorkflow defines the schema for the workflow resource
func resourceGraalSystemsWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsWorkflowCreate,
		ReadContext:   resourceGraalSystemsWorkflowRead,
		UpdateContext: resourceGraalSystemsWorkflowUpdate,
		DeleteContext: resourceGraalSystemsWorkflowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the workflow",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the workflow",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the project to deploy the workflow on",
			},
			"identity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the identity to use",
			},
			"schedule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Schedule mode of the job. Either `once` or `cron`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type: schema.TypeString,
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								if !slices.Contains(scheduleTypes, val.(string)) {
									errs = append(errs, fmt.Errorf("%q must be one of %q, got: %s", key, scheduleTypes, val))
								}
								return
							},
							Required: true,
						},
						"cron_expression": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cron expression of the schedule. Only used if type is `cron`",
							//TODO: add validate for cron exp ?
						},
						"timezone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Timezone of the schedule. Only used if type is `cron`",
						},
						"infrastructure_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Infrastructure id used for the schedule. Only used if type is `cron`",
						},
						"device_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Device id",
						},
					},
				},
			},
			"job": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of job to chain as a workflow",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ref": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The job ID",
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								if _, err := uuid.ParseUUID(val.(string)); err != nil {
									errs = append(errs, fmt.Errorf("%q must be a valid UUID, got: %q", key, val))
								}
								return
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The job name",
						},
						"depends_on": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of job names to wait for before starting this job",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for every step of the job",
			},
			/* TODO: add the following fields
			"notifications"
			"parameters"
			"metadata"*/
		},
	}
}

func validateUuidList(val []string) (errs []error) {
	for _, elem := range val {
		if _, err := uuid.ParseUUID(elem); err != nil {
			errs = append(errs, fmt.Errorf("%q must contain valid UUIDs, got: %q", val, elem))
		}
	}
	return errs
}

// resourceGraalSystemsWorkflowCreate creates a workflow
func resourceGraalSystemsWorkflowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	projectId := d.Get("project_id").(string)
	identityId := d.Get("identity_id").(string)
	sch := d.Get("schedule").([]interface{})
	lbl := d.Get("labels").(map[string]interface{})

	labels := toStringMap(lbl)
	// Validate and create appropriate schedule
	if diagnostics := validateSchedule(sch[0]); diagnostics != nil {
		return diagnostics
	}
	schedule := defineSchedule(sch[0])

	jobs := d.Get("job").([]interface{})
	if errs := validateJobs(jobs); errs != nil {
		return diag.FromErr(errs)
	}

	workflow := &sdk.Workflow{
		Name:        &name,
		Description: &description,
		ProjectId:   &projectId,
		IdentityId:  &identityId,
		Schedule:    &schedule,
		Tasks:       defineTasks(jobs, "job"),
		Labels:      &labels,
	}

	if registeredWorkflow, _, err := apiClient.ProjectAPI.CreateWorkflowForProject(context.Background(), projectId).XTenant(meta.tenant).Workflow(*workflow).Execute(); err != nil {
		return diag.FromErr(err)
	} else {
		d.SetId(*registeredWorkflow.Id)
	}

	return resourceGraalSystemsWorkflowRead(ctx, d, m)
}

// resourceGraalSystemsWorkflowRead reads the workflow from the GraalSystems API and returns its attributes
func resourceGraalSystemsWorkflowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	// Retrieve the input
	workflowId := d.Id()

	// Retrieve the workflow
	workflow, _, err := apiClient.WorkflowAPI.FindWorkflowById(context.Background(), workflowId).XTenant(meta.tenant).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the resource data
	_ = d.Set("name", workflow.Name)
	if workflow.Description != nil {
		_ = d.Set("description", workflow.Description)
	}
	_ = d.Set("project_id", workflow.ProjectId)
	_ = d.Set("identity_id", workflow.IdentityId)

	if schedule, errors := readSchedule(workflow.Schedule); err != nil {
		return diag.FromErr(errors)
	} else {
		_ = d.Set("schedule", schedule)
	}
	if tasks, er := readTasks(workflow.Tasks); err != nil {
		return diag.FromErr(er)
	} else {
		_ = d.Set("job", tasks)
	}
	if workflow.Notifications != nil {
		_ = d.Set("notifications", *workflow.Notifications)
	}
	_ = d.Set("labels", workflow.Labels)

	return nil
}

// resourceGraalSystemsWorkflowUpdate updates a workflow
func resourceGraalSystemsWorkflowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	workflowId := d.Id()
	if d.HasChange("name") {
		_, _, err := apiClient.WorkflowAPI.UpdateWorkflow(context.Background(), workflowId).XTenant(meta.tenant).Patch(patchFromResourceData(d, "name")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("description") {
		_, _, err := apiClient.WorkflowAPI.UpdateWorkflow(context.Background(), workflowId).XTenant(meta.tenant).Patch(patchFromResourceData(d, "description")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("schedule") {
		return diag.FromErr(fmt.Errorf("cannot change schedule yet, please recreate the workflow"))
		/*sch := d.Get("schedule").([]interface{})
		if diagnostics := validateSchedule(sch[0]); diagnostics != nil {
			return diagnostics
		}
		_, _, err := apiClient.WorkflowAPI.UpdateWorkflow(context.Background(), workflowId).XTenant(meta.tenant).Patch(patchFromResourceData(d, "schedule")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}*/
	}
	if d.HasChange("labels") {
		_, _, err := apiClient.WorkflowAPI.UpdateWorkflow(context.Background(), workflowId).XTenant(meta.tenant).Patch(patchFromResourceData(d, "labels")).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("job") {
		changedJobs, _ := d.GetOk("job")
		if err := validateJobs(changedJobs.([]interface{})); err != nil {
			return diag.FromErr(err)
		}
		if jobPatches := patchJobs(); jobPatches == nil {
			return diag.FromErr(fmt.Errorf("cannot yet update jobs, please recreate the workflow"))
		} else {
			_, _, err := apiClient.WorkflowAPI.UpdateWorkflow(context.Background(), workflowId).XTenant(meta.tenant).Patch(jobPatches).Execute()
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceGraalSystemsWorkflowRead(ctx, d, m)
}

// resourceGraalSystemsWorkflowDelete deletes a workflow
func resourceGraalSystemsWorkflowDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	workflowId := d.Id()
	_, err := apiClient.WorkflowAPI.DeleteWorkflowById(context.Background(), workflowId).XTenant(meta.tenant).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readSchedule(sch sdk.ISchedule) ([]map[string]string, error) {
	if sch == nil {
		return nil, nil
	}

	// Deserialize the schedule into the abstract type
	var schedule sdk.Schedule
	schBytes, err := json.Marshal(sch)
	if err != nil {
		return nil, fmt.Errorf("schedule read marshall error: %s", err)
	}
	if err := json.Unmarshal(schBytes, &schedule); err != nil {
		return nil, fmt.Errorf("schedule read unmarshall error: %s", err)
	}
	// Then, depending on the type, we can deserialize it into the correct type
	if *schedule.Type == "once" {
		return []map[string]string{
			{
				"type": *schedule.Type,
			},
		}, nil
	} else {
		var cron sdk.CronSchedule
		if err := json.Unmarshal(schBytes, &cron); err != nil {
			return nil, fmt.Errorf("cron schedule read unmarshall error: %s", err)
		}

		return []map[string]string{
			{
				"type":              *schedule.Type,
				"cron_expression":   *cron.CronExpression,
				"timezone":          *cron.Timezone,
				"infrastructure_id": *cron.InfrastructureId,
				"device_id":         *cron.DeviceId,
			},
		}, nil
	}
}

func readTasks(tasks []sdk.ITask) ([]map[string]interface{}, error) {
	var taskRefs []map[string]interface{}
	for _, task := range tasks {
		// Deserialize the task into the abstract type
		var sdkTask sdk.Task
		taskBytes, marshErr := json.Marshal(task)
		if marshErr != nil {
			return nil, fmt.Errorf("task read marshall error: %s", marshErr)
		}
		if err := json.Unmarshal(taskBytes, &sdkTask); err != nil {
			return nil, fmt.Errorf("task read unmarshall error: %s", err)
		}
		// Then, depending on the type, we can deserialize it into the correct type
		if *sdkTask.Type == "job" {
			var jobTask sdk.JobTask
			if err := json.Unmarshal(taskBytes, &jobTask); err != nil {
				return nil, fmt.Errorf("job task read unmarshall error: %s", err)
			}
			taskRefs = append(taskRefs, map[string]interface{}{
				"ref":        *jobTask.Ref,
				"name":       *jobTask.Task.Name,
				"depends_on": jobTask.Task.Depends,
			})
		} else {
			return nil, fmt.Errorf("task type %s is not yet supported", *sdkTask.Type)
		}
	}
	return taskRefs, nil
}

func validateJobs(jobs []interface{}) error {
	var jobNames []string
	for i, job := range jobs {
		j := job.(map[string]interface{})
		if val, ok := j["ref"]; !ok {
			if _, err := uuid.ParseUUID(val.(string)); err != nil {
				return err
			}
		}
		if _, ok := j["name"]; !ok {
			return fmt.Errorf("the name field is required for job %d", i)
		} else {
			jobNames = append(jobNames, j["name"].(string))
		}

		if val, ok := j["depends_on"]; ok {
			dependencies := toStringList(val.([]interface{}))
			if i == 0 && len(dependencies) > 0 {
				return fmt.Errorf("the first job must not have a depends_on field")
			}
			if i > 0 {
				for _, dep := range dependencies {
					if !slices.Contains(jobNames, dep) {
						return fmt.Errorf("job %q depends on %q, but this job is not defined", j["name"].(string), dep)
					}
				}
			}
		}
	}
	return nil
}

func defineTasks(tasks []interface{}, taskType string) []sdk.ITask {
	var jobTasks []sdk.ITask
	for _, task := range tasks {
		t := task.(map[string]interface{})
		taskName := t["name"].(string)
		taskDependencies := toStringList(t["depends_on"].([]interface{}))
		taskRef := t["ref"].(string)

		jobTasks = append(jobTasks, sdk.JobTask{Task: sdk.Task{Name: &taskName, Depends: taskDependencies, Type: &taskType}, Ref: &taskRef})
	}
	return jobTasks
}

func patchJobs() []sdk.Patch {
	//TODO: implement patchJobs
	return nil
}
