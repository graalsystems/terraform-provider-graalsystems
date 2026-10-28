package graalsystems

import (
	"encoding/json"
	"fmt"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const (
	scheduleTypeOnce = "once"
	scheduleTypeCron = "cron"
)

const (
	optionTypeBash   = "bash"
	optionTypePython = "python"
	optionTypeSpark  = "spark"
)

const (
	libraryTypeFile = "file"
)

var scheduleTypes = []string{scheduleTypeOnce, scheduleTypeCron}
var optionsTypes = []string{optionTypeBash, optionTypePython, optionTypeSpark}
var libraryTypes = []string{libraryTypeFile}

func validateOptions(input interface{}) diag.Diagnostics {
	opts, _ := input.(map[string]interface{})
	str := func(k string) string { v, _ := opts[k].(string); return v }
	switch str("type") {
	case optionTypeBash:
		if lines, _ := opts["lines"].([]interface{}); len(lines) == 0 {
			return diag.FromErr(fmt.Errorf("lines is required for options type %s", optionTypeBash))
		}
	case optionTypePython:
		if str("file") == "" && str("module") == "" {
			return diag.FromErr(fmt.Errorf("file or module is required for options type %s", optionTypePython))
		}
	case optionTypeSpark:
		if str("file") == "" && str("main_class_name") == "" {
			return diag.FromErr(fmt.Errorf("file or main_class_name is required for options type %s", optionTypeSpark))
		}
	}
	return nil
}

// defineOptions builds the options object of a job from the schema block,
// leaving out unset attributes.
func defineOptions(input interface{}) sdk.IOptions {
	opts, _ := input.(map[string]interface{})
	out := map[string]interface{}{}
	for _, k := range []string{"type", "instance_type", "docker_image", "file", "module", "main_class_name"} {
		if v, _ := opts[k].(string); v != "" {
			out[k] = v
		}
	}
	if lines, _ := opts["lines"].([]interface{}); len(lines) > 0 {
		out["lines"] = toStringList(lines)
	}
	if env, _ := opts["env"].(map[string]interface{}); len(env) > 0 {
		out["env"] = toStringMap(env)
	}
	return out
}

func validateSchedule(input interface{}) diag.Diagnostics {
	convertedInput := toStringMap(input.(map[string]interface{}))

	if convertedInput["type"] == scheduleTypeOnce {
		if vCron := convertedInput["cron_expression"]; len(vCron) > 0 {
			return diag.FromErr(fmt.Errorf("cron_expression is not allowed for schedule type %s", scheduleTypeOnce))
		}
		if vTime := convertedInput["timezone"]; len(vTime) > 0 {
			return diag.FromErr(fmt.Errorf("timezone is not allowed for schedule type %s", scheduleTypeOnce))
		}
		if vInfra := convertedInput["infrastructure_id"]; len(vInfra) > 0 {
			return diag.FromErr(fmt.Errorf("infrastructure_id is not allowed for schedule type %s", scheduleTypeOnce))
		}
		if vDevice := convertedInput["device_id"]; len(vDevice) > 0 {
			return diag.FromErr(fmt.Errorf("device_id is not allowed for schedule type %s", scheduleTypeOnce))
		}
	}
	if convertedInput["type"] == scheduleTypeCron {
		if vCron := convertedInput["cron_expression"]; len(vCron) == 0 {
			return diag.FromErr(fmt.Errorf("cron_expression is required for schedule type %s", scheduleTypeCron))
		}
		if vTime := convertedInput["timezone"]; len(vTime) == 0 {
			return diag.FromErr(fmt.Errorf("timezone is required for schedule type %s", scheduleTypeCron))
		}
		if vInfra := convertedInput["infrastructure_id"]; len(vInfra) == 0 {
			return diag.FromErr(fmt.Errorf("infrastructure_id is required for schedule type %s", scheduleTypeCron))
		}
	}
	return nil
}

func defineSchedule(input interface{}) sdk.ISchedule {
	convertedInput := toStringMap(input.(map[string]interface{}))

	if convertedInput["type"] == scheduleTypeCron {
		var sch sdk.CronSchedule
		bytes, err := json.Marshal(convertedInput)
		if err != nil {
			fmt.Println("Cron schedule marshall error:", err)
		}
		err = json.Unmarshal(bytes, &sch)
		if err != nil {
			fmt.Println("Cron schedule unmarshall error:", err)
		}
		return sch
	}
	return *sdk.NewRunOnceSchedule()
}

func validateLibraries(input []interface{}) diag.Diagnostics {
	for _, lib := range input {
		convertedInput := toStringMap(lib.(map[string]interface{}))
		if convertedInput["type"] == libraryTypeFile {
			if vKey := convertedInput["key"]; len(vKey) == 0 {
				return diag.FromErr(fmt.Errorf("key is required for library type %s", libraryTypeFile))
			}
		}
	}
	return nil
}

func defineLibraries(input []interface{}) []sdk.ILibrary {
	var libs []sdk.ILibrary
	for _, lib := range input {
		convertedInput := toStringMap(lib.(map[string]interface{}))
		if convertedInput["type"] == libraryTypeFile {
			var lib sdk.FileLibrary
			bytes, err := json.Marshal(convertedInput)
			if err != nil {
				fmt.Println("file library marshall error:", err)
			}
			err = json.Unmarshal(bytes, &lib)
			if err != nil {
				fmt.Println("file library unmarshall error:", err)
			}
			libs = append(libs, lib)
		}
	}
	return libs
}
