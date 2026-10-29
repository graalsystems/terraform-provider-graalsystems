package graalsystems

import (
	"context"
	"fmt"
	"slices"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeInt,
				Optional: true,
				// The API sets a default (3600 s) when the field is omitted: without
				// Computed, every plan would try to set it back to null.
				Computed:    true,
				Description: "Maximum duration of the job, in seconds. Defaults to the platform's value (3600).",
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
				Computed:    true,
				Description: "Maximum retries in case of failure. Defaults to the platform's value.",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q must be greater than or equal to 0, got: %d", key, v))
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
							Optional:    true,
							Description: "Container image to use for the job. Defaults to the runtime image configured on the installation for the job type",
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
						"file": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File to run, relative to the run's /workspace where libraries are copied. Used by `python` and `spark` jobs",
						},
						"main_class_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Main class of a JVM Spark application. Only used if type is `spark`",
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
	parameters := d.Get("parameters").([]interface{})
	labels := toStringMap(d.Get("labels").(map[string]interface{}))

	opts := d.Get("options").([]interface{})
	if diagnostics := validateOptions(opts[0]); diagnostics != nil {
		return diagnostics
	}
	currentOptions := defineOptions(opts[0])

	libs := d.Get("library").([]interface{})
	if diagnostics := validateLibraries(libs); diagnostics != nil {
		return diagnostics
	}

	job := &sdk.Job{
		Name:        &name,
		Description: &description,
		ProjectId:   &projectId,
		IdentityId:  &identityId,
		Options:     &currentOptions,
		Parameters:  toStringList(parameters),
		Secrets:     toStringList(d.Get("secrets").([]interface{})),
		Labels:      &labels,
		Libraries:   defineLibraries(libs),
	}
	if v, ok := d.GetOk("timeout_seconds"); ok {
		timeoutSeconds := int32(v.(int))
		job.TimeoutSeconds = &timeoutSeconds
	}
	if v, ok := d.GetOk("max_retries"); ok {
		maxRetries := int32(v.(int))
		job.MaxRetries = &maxRetries
	}
	// Without a schedule block, the job only runs when started explicitly.
	if sch := d.Get("schedule").([]interface{}); len(sch) > 0 && sch[0] != nil {
		if diagnostics := validateSchedule(sch[0]); diagnostics != nil {
			return diagnostics
		}
		schedule := defineSchedule(sch[0])
		job.Schedule = &schedule
	}

	result, _, err := apiClient.ProjectAPI.CreateJobForProject(context.Background(), projectId).XTenant(meta.tenant).Job(*job).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil || result.Id == nil {
		return diag.FromErr(fmt.Errorf("the job was created but the API did not return its id"))
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
	_ = d.Set("project_id", res.ProjectId)
	_ = d.Set("identity_id", res.IdentityId)
	if res.TimeoutSeconds != nil {
		_ = d.Set("timeout_seconds", int(*res.TimeoutSeconds))
	}
	if res.MaxRetries != nil {
		_ = d.Set("max_retries", int(*res.MaxRetries))
	}
	_ = d.Set("parameters", res.Parameters)
	_ = d.Set("secrets", res.Secrets)
	if res.Labels != nil {
		_ = d.Set("labels", *res.Labels)
	}

	return nil
}

func resourceGraalSystemsJobUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*Meta)

	ops := replaceChanged(d, "name", "description", "labels", "parameters", "secrets")
	for _, attr := range []string{"timeout_seconds", "max_retries"} {
		if !d.HasChange(attr) {
			continue
		}
		if v, ok := d.GetOk(attr); ok {
			ops = append(ops, patchOperation{Op: "add", Path: "/" + attr, Value: v})
		} else {
			ops = append(ops, patchOperation{Op: "remove", Path: "/" + attr})
		}
	}
	if d.HasChange("options") {
		opts := d.Get("options").([]interface{})
		if diagnostics := validateOptions(opts[0]); diagnostics != nil {
			return diagnostics
		}
		ops = append(ops, patchOperation{Op: "add", Path: "/options", Value: defineOptions(opts[0])})
	}
	if d.HasChange("schedule") {
		sch := d.Get("schedule").([]interface{})
		if len(sch) == 0 || sch[0] == nil {
			ops = append(ops, patchOperation{Op: "remove", Path: "/schedule"})
		} else {
			if diagnostics := validateSchedule(sch[0]); diagnostics != nil {
				return diagnostics
			}
			ops = append(ops, patchOperation{Op: "add", Path: "/schedule", Value: defineSchedule(sch[0])})
		}
	}
	if d.HasChange("library") {
		libs := d.Get("library").([]interface{})
		if diagnostics := validateLibraries(libs); diagnostics != nil {
			return diagnostics
		}
		ops = append(ops, patchOperation{Op: "add", Path: "/libraries", Value: defineLibraries(libs)})
	}
	if err := applyPatch(ctx, meta, "job", "/jobs/"+d.Id(), ops); err != nil {
		return diag.FromErr(err)
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
