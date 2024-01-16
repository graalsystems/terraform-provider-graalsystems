package graalsystems

import (
	"context"
	"fmt"
	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"slices"
)

func resourceGraalSystemsJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraalSystemsJobCreate,
		ReadContext:   resourceGraalSystemsJobRead,
		UpdateContext: resourceGraalSystemsJobUpdate,
		DeleteContext: resourceGraalSystemsJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the job to create",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the job",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project id of the project the job belongs to",
			},
			"identity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The identity id of the identity used to run the job",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum duration of the job",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be a positive integer, got: %d", key, v))
					}
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q cannot be null, got: %d", key, v))
					}
					return
				},
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum retries in case of failure",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)
					if v <= 0 {
						errs = append(errs, fmt.Errorf("%q must be greater or equal 0, got: %d", key, v))
					}
					return
				},
			},
			"options": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Job definition options",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"env": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Key value pairs of environment variables for the job",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"docker_image": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Docker image to use for the job",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Compute instance type to use for the job. Check which instance types are available for your project",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Type of the job. Possible values in %q.", optionsTypes),
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								if !slices.Contains(optionsTypes, val.(string)) {
									errs = append(errs, fmt.Errorf("%q must be one of %q, got: %s", key, optionsTypes, val))
								}
								return
							},
						},
						"lines": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of bash lines to execute. Only used if type is `bash`",
						},
						"module": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Python module to execute. Only used if type is `python`",
						},
					},
				},
			},
			"secrets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of secret ids",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"library": {
				// TODO: Create a resource & data source for libraries it will allow for easy key retrieval of existing libraries
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of libraries to use for the job run",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Library type, should be file by default",
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								if !slices.Contains(libraryTypes, val.(string)) {
									errs = append(errs, fmt.Errorf("%q must be one of %q, got: %s", key, libraryTypes, val))
								}
								return
							},
							Default: "file",
						},
						"dep": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dependency of the library. Only used if type is `pypi`",
						},
						"repo": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Repository url of the library. Only used if type is `maven`",
						},
						"dependency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dependency of the library. Only used if type is `maven`",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Repository url of the library. Only used if type is `git`",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path of the library. Only used if type is `git`",
						},
						"revision": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Revision of the library. Only used if type is `git`",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Username to use to connect to git. Only used if type is `git`",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Password to use to connect to git. Only used if type is `git`",
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Id of the library to use in the job. Only used if type is `file`",
						},
						"ref": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Reference of the library. Only used if type is `cran`",
						},
					},
				},
			},
			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of parameters",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for every step of the job",
			},
			"schedule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
			/* TODO: add the following fields
			"notifications"
			"secrets"
			"metadata"*/

		},
	}
}

func resourceGraalSystemsJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	projectId := d.Get("project_id").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	identityId := d.Get("identity_id").(string)
	ts := d.Get("timeout_seconds").(int)
	mr := d.Get("max_retries").(int)
	parameters := d.Get("parameters").([]interface{})
	lbl := d.Get("labels").(map[string]interface{})

	// Convert the definition to the sdk.Job expected types
	timeoutSeconds := int32(ts)
	maxRetries := int32(mr)
	labels := toStringMap(lbl)

	opts := d.Get("options").([]interface{})
	if diagnostics := validateOptions(opts[0]); diagnostics != nil {
		return diagnostics
	}
	currentOptions := defineOptions(opts[0])

	sch := d.Get("schedule").([]interface{})
	if diagnostics := validateSchedule(sch[0]); diagnostics != nil {
		return diagnostics
	}
	schedule := defineSchedule(sch[0])

	libs := d.Get("library").([]interface{})
	if diagnostics := validateLibraries(libs); diagnostics != nil {
		return diagnostics
	}
	libraries := defineLibraries(libs)

	job := &sdk.Job{
		Name:           &name,
		Description:    &description,
		ProjectId:      &projectId,
		IdentityId:     &identityId,
		Options:        &currentOptions,
		TimeoutSeconds: &timeoutSeconds,
		MaxRetries:     &maxRetries,
		Parameters:     toStringList(parameters),
		Labels:         &labels,
		Schedule:       &schedule,
		Libraries:      libraries,
	}
	result, response, err := apiClient.ProjectAPI.CreateJobForProject(context.Background(), projectId).XTenant(meta.tenant).Job(*job).Execute()
	if err != nil {
		fmt.Printf("Error %+v", err)
		return diag.FromErr(err)
	}
	if response.StatusCode == 200 {
		return diag.FromErr(fmt.Errorf("Job created, but could not retrieve its info. Check that every parameter you entered is valid."))
	}

	d.SetId(*result.Id)

	return resourceGraalSystemsJobRead(ctx, d, meta)
}

func resourceGraalSystemsJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	res, _, err := apiClient.JobAPI.FindJobByJobId(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("description", res.Description)

	return nil
}

func resourceGraalSystemsJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient
	//TODO: update when other parameters are updated
	if d.HasChange("name") {
		path := "/name"

		value := d.Get("name").(string)

		patch := &sdk.Patch{
			Op:    nil,
			Path:  &path,
			Value: &value,
		}
		patchs := &[]sdk.Patch{*patch}
		_, _, err := apiClient.JobAPI.UpdateJob(context.Background(), d.Id()).XTenant(meta.tenant).Patch(*patchs).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGraalSystemsJobRead(ctx, d, meta)
}

func resourceGraalSystemsJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)
	apiClient := meta.apiClient

	_, err := apiClient.JobAPI.DeleteJobById(context.Background(), d.Id()).XTenant(meta.tenant).Execute()
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
