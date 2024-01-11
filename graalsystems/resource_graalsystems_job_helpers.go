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
)

const (
	libraryTypeFile = "file"
)

var scheduleTypes = []string{scheduleTypeOnce, scheduleTypeCron}
var optionsTypes = []string{optionTypeBash, optionTypePython}
var libraryTypes = []string{libraryTypeFile}

func validateOptions(input interface{}) diag.Diagnostics {
	var opts sdk.Options
	optBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("validate options marshall error:", err)
	}
	err = json.Unmarshal(optBytes, &opts)
	if err != nil {
		fmt.Println("validate options unmarshall error:", err)
	}

	if *opts.Type == optionTypeBash {
		var opt sdk.BashOptions
		err = json.Unmarshal(optBytes, &opt)
		if err != nil {
			fmt.Println("bash options unmarshall error:", err)
		}
		if len(opt.Lines) == 0 {
			return diag.FromErr(fmt.Errorf("lines parameter is required for options type %s", optionTypeBash))
		}
	}
	if *opts.Type == optionTypePython {
		var opt sdk.PythonOptions
		err = json.Unmarshal(optBytes, &opt)
		if err != nil {
			fmt.Println("python options unmarshall error:", err)
		}
		if *opt.Module == "" {
			return diag.FromErr(fmt.Errorf("module parameter is required for options type %s", optionTypePython))
		}
	}
	return nil
}

func defineOptions(input interface{}) sdk.IOptions {
	var opts sdk.Options
	optBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("options definition marshall error:", err)
	}
	err = json.Unmarshal(optBytes, &opts)
	if err != nil {
		fmt.Println("options definition unmarshall error:", err)
	}
	if *opts.Type == optionTypeBash {
		var opt sdk.BashOptions
		err = json.Unmarshal(optBytes, &opt)
		if err != nil {
			fmt.Println("bash options definition unmarshall error:", err)
		}
		if opt.Lines == nil {

		}
		return opt
	}
	if *opts.Type == optionTypePython {
		var opt sdk.PythonOptions
		err = json.Unmarshal(optBytes, &opt)
		if err != nil {
			fmt.Println("python options definition unmarshall error:", err)
		}
		return opt
	}
	return nil
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
