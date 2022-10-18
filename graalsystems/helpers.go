package graalsystems

import (
	sdk "github.com/graalsystems/sdk/go"
	"net/http"
	"strconv"
)
import "errors"

// // DefaultWaitRetryInterval is used to set the retry interval to 0 during acceptance tests
// var DefaultWaitRetryInterval *time.Duration
//
// // terraformResourceData is an interface for *schema.ResourceData. (used for mock)
//
//	type terraformResourceData interface {
//		HasChange(string) bool
//		GetOk(string) (interface{}, bool)
//		Get(string) interface{}
//		Set(string, interface{}) error
//		SetId(string)
//		Id() string
//	}
//
// isHTTPCodeError returns true if err is an http error with code statusCode
func isHTTPCodeError(err error, statusCode int) bool {
	if err == nil {
		return false
	}

	statusText := strconv.Itoa(statusCode) + " " + http.StatusText(statusCode)

	responseError := &sdk.GenericOpenAPIError{}
	if errors.As(err, &responseError) && responseError.Error() == statusText {
		return true
	}
	return false
}

// is404Error returns true if err is an HTTP 404 error
func is404Error(err error) bool {
	return isHTTPCodeError(err, http.StatusNotFound)
}

//func is412Error(err error) bool {
//	return isHTTPCodeError(err, http.StatusPreconditionFailed)
//}
//
//// organizationIDSchema returns a standard schema for a organization_id
//func organizationIDSchema() *schema.Schema {
//	return &schema.Schema{
//		Type:        schema.TypeString,
//		Description: "The organization_id you want to attach the resource to",
//		Computed:    true,
//	}
//}
//
//// validateStringInSliceWithWarning helps to only returns warnings in case we got a non public locality passed
//func validateStringInSliceWithWarning(correctValues []string, field string) func(i interface{}, path cty.Path) diag.Diagnostics {
//	return func(i interface{}, path cty.Path) diag.Diagnostics {
//		_, rawErr := validation.StringInSlice(correctValues, true)(i, field)
//		var res diag.Diagnostics
//		for _, e := range rawErr {
//			res = append(res, diag.Diagnostic{
//				Severity:      diag.Warning,
//				Summary:       e.Error(),
//				AttributePath: path,
//			})
//		}
//		return res
//	}
//}
//
//// newRandomName returns a random name prefixed for terraform.
//func newRandomName(prefix string) string {
//	return namegenerator.GetRandomName("tf", prefix)
//}
//
//const gb uint64 = 1000 * 1000 * 1000
//
//func flattenTime(date *time.Time) interface{} {
//	if date != nil {
//		return date.Format(time.RFC3339)
//	}
//	return ""
//}
//
//func flattenDuration(duration *time.Duration) interface{} {
//	if duration != nil {
//		return duration.String()
//	}
//	return ""
//}
//
//func expandDuration(data interface{}) (*time.Duration, error) {
//	if data == nil || data == "" {
//		return nil, nil
//	}
//	d, err := time.ParseDuration(data.(string))
//	if err != nil {
//		return nil, err
//	}
//	return &d, nil
//}
//
//func expandOrGenerateString(data interface{}, prefix string) string {
//	if data == nil || data == "" {
//		return newRandomName(prefix)
//	}
//	return data.(string)
//}
//
//func expandStringWithDefault(data interface{}, defaultValue string) string {
//	if data == nil || data.(string) == "" {
//		return defaultValue
//	}
//	return data.(string)
//}
//
//func expandStrings(data interface{}) []string {
//	var stringSlice []string
//	for _, s := range data.([]interface{}) {
//		stringSlice = append(stringSlice, s.(string))
//	}
//	return stringSlice
//}
//
//func expandStringsPtr(data interface{}) *[]string {
//	var stringSlice []string
//	if _, ok := data.([]interface{}); !ok || data == nil {
//		return nil
//	}
//	for _, s := range data.([]interface{}) {
//		stringSlice = append(stringSlice, s.(string))
//	}
//	if stringSlice == nil {
//		return nil
//	}
//	return &stringSlice
//}
//
//// expandUpdatedStringsPtr expands a string slice but will default to an empty list.
//// Should be used on schema update so emptying a list will update resource.
//func expandUpdatedStringsPtr(data interface{}) *[]string {
//	stringSlice := []string{}
//	if _, ok := data.([]interface{}); !ok || data == nil {
//		return &stringSlice
//	}
//	for _, s := range data.([]interface{}) {
//		stringSlice = append(stringSlice, s.(string))
//	}
//	return &stringSlice
//}
//
//func expandSliceIDsPtr(rawIDs interface{}) *[]string {
//	var stringSlice []string
//	if _, ok := rawIDs.([]interface{}); !ok || rawIDs == nil {
//		return &stringSlice
//	}
//	for _, s := range rawIDs.([]interface{}) {
//		stringSlice = append(stringSlice, expandID(s.(string)))
//	}
//	return &stringSlice
//}
//
//func expandStringsOrEmpty(data interface{}) []string {
//	var stringSlice []string
//	if _, ok := data.([]interface{}); !ok || data == nil {
//		return stringSlice
//	}
//	for _, s := range data.([]interface{}) {
//		stringSlice = append(stringSlice, s.(string))
//	}
//	return stringSlice
//}
//
//func expandSliceStringPtr(data interface{}) []*string {
//	if data == nil {
//		return nil
//	}
//	stringSlice := []*string(nil)
//	for _, s := range data.([]interface{}) {
//		stringSlice = append(stringSlice, expandStringPtr(s))
//	}
//	return stringSlice
//}
//
//func flattenIPPtr(ip *net.IP) interface{} {
//	if ip == nil {
//		return ""
//	}
//	return ip.String()
//}
//
//func flattenStringPtr(s *string) interface{} {
//	if s == nil {
//		return ""
//	}
//	return *s
//}
//
//func flattenSliceStringPtr(s []*string) interface{} {
//	res := make([]interface{}, 0, len(s))
//	for _, strPtr := range s {
//		res = append(res, flattenStringPtr(strPtr))
//	}
//	return res
//}
//
//func flattenSliceString(s []string) interface{} {
//	res := make([]interface{}, 0, len(s))
//	for _, strPtr := range s {
//		res = append(res, strPtr)
//	}
//	return res
//}
//
//func flattenBoolPtr(b *bool) interface{} {
//	if b == nil {
//		return nil
//	}
//	return *b
//}
//
//func expandStringPtr(data interface{}) *string {
//	if data == nil || data == "" {
//		return nil
//	}
//	return StringPtr(data.(string))
//}
//
//func expandUpdatedStringPtr(data interface{}) *string {
//	str := ""
//	if data != nil {
//		str = data.(string)
//	}
//	return &str
//}
//
//func expandBoolPtr(data interface{}) *bool {
//	if data == nil {
//		return nil
//	}
//	return scw.BoolPtr(data.(bool))
//}
//
//func flattenInt32Ptr(i *int32) interface{} {
//	if i == nil {
//		return 0
//	}
//	return *i
//}
//
//func expandInt32Ptr(data interface{}) *int32 {
//	if data == nil || data == "" {
//		return nil
//	}
//	return scw.Int32Ptr(int32(data.(int)))
//}
//
//func expandUint32Ptr(data interface{}) *uint32 {
//	if data == nil || data == "" {
//		return nil
//	}
//	return scw.Uint32Ptr(uint32(data.(int)))
//}
//
//func expandIPNet(raw string) (scw.IPNet, error) {
//	if raw == "" {
//		return scw.IPNet{}, nil
//	}
//	var ipNet scw.IPNet
//	err := json.Unmarshal([]byte(strconv.Quote(raw)), &ipNet)
//	if err != nil {
//		return scw.IPNet{}, fmt.Errorf("%s could not be marshaled: %v", raw, err)
//	}
//
//	return ipNet, nil
//}
//
//func flattenIPNet(ipNet scw.IPNet) (string, error) {
//	raw, err := json.Marshal(ipNet)
//	if err != nil {
//		return "", err
//	}
//	return string(raw[1 : len(raw)-1]), nil // remove quotes
//}
//
//func validateDuration() schema.SchemaValidateFunc {
//	return func(i interface{}, s string) (strings []string, errors []error) {
//		str, isStr := i.(string)
//		if !isStr {
//			return nil, []error{fmt.Errorf("%v is not a string", i)}
//		}
//		_, err := time.ParseDuration(str)
//		if err != nil {
//			return nil, []error{fmt.Errorf("cannot parse duration for value %s", str)}
//		}
//		return nil, nil
//	}
//}
//
//func flattenMap(m map[string]string) interface{} {
//	if m == nil {
//		return nil
//	}
//	flattenedMap := make(map[string]interface{})
//	for k, v := range m {
//		flattenedMap[k] = v
//	}
//	return flattenedMap
//}
//
//func flattenMapStringStringPtr(m map[string]*string) interface{} {
//	if m == nil {
//		return nil
//	}
//	flattenedMap := make(map[string]interface{})
//	for k, v := range m {
//		if v != nil {
//			flattenedMap[k] = *v
//		} else {
//			flattenedMap[k] = ""
//		}
//	}
//	return flattenedMap
//}
//
//func diffSuppressFuncDuration(k, oldValue, newValue string, d *schema.ResourceData) bool {
//	if oldValue == newValue {
//		return true
//	}
//	d1, err1 := time.ParseDuration(oldValue)
//	d2, err2 := time.ParseDuration(newValue)
//	if err1 != nil || err2 != nil {
//		return false
//	}
//	return d1 == d2
//}
//
//func diffSuppressFuncIgnoreCase(k, oldValue, newValue string, d *schema.ResourceData) bool {
//	return strings.EqualFold(oldValue, newValue)
//}
//
//func diffSuppressFuncIgnoreCaseAndHyphen(k, oldValue, newValue string, d *schema.ResourceData) bool {
//	return strings.ReplaceAll(strings.ToLower(oldValue), "-", "_") == strings.ReplaceAll(strings.ToLower(newValue), "-", "_")
//}
//
//// diffSuppressFuncLocality is a SuppressDiffFunc to remove the locality from an ID when checking diff.
//// e.g. 2c1a1716-5570-4668-a50a-860c90beabf6 == fr-par-1/2c1a1716-5570-4668-a50a-860c90beabf6
//func diffSuppressFuncLocality(k, oldValue, newValue string, d *schema.ResourceData) bool {
//	return expandID(oldValue) == expandID(newValue)
//}
//
//// TimedOut returns true if the error represents a "wait timed out" condition.
//// Specifically, TimedOut returns true if the error matches all these conditions:
////   - err is of type resource.TimeoutError
////   - TimeoutError.LastError is nil
//func TimedOut(err error) bool {
//	// This explicitly does *not* match wrapped TimeoutErrors
//	timeoutErr, ok := err.(*resource.TimeoutError) //nolint:errorlint // Explicitly does *not* match wrapped TimeoutErrors
//	return ok && timeoutErr.LastError == nil
//}
//
//func expandMapPtrStringString(data interface{}) *map[string]string {
//	if data == nil {
//		return nil
//	}
//	m := make(map[string]string)
//	for k, v := range data.(map[string]interface{}) {
//		m[k] = v.(string)
//	}
//	return &m
//}
//
//func expandMapStringStringPtr(data interface{}) map[string]*string {
//	if data == nil {
//		return nil
//	}
//	m := make(map[string]*string)
//	for k, v := range data.(map[string]interface{}) {
//		m[k] = expandStringPtr(v)
//	}
//	return m
//}
//
//func errorCheck(err error, message string) bool {
//	return strings.Contains(err.Error(), message)
//}
//
//// ErrCodeEquals returns true if the error matches all these conditions:
////   - err is of type scw.Error
////   - Error.Error() equals one of the passed codes
//func ErrCodeEquals(err error, codes ...string) bool {
//	var scwErr scw.SdkError
//	if errors.As(err, &scwErr) {
//		for _, code := range codes {
//			if scwErr.Error() == code {
//				return true
//			}
//		}
//	}
//	return false
//}
//
//func getBool(d *schema.ResourceData, key string) interface{} {
//	val, ok := d.GetOkExists(key)
//	if !ok {
//		return nil
//	}
//	return val
//}
//
//// validateDate will validate that field is a valid ISO 8601
//// It is the same as RFC3339
//func validateDate() schema.SchemaValidateDiagFunc {
//	return func(i interface{}, path cty.Path) diag.Diagnostics {
//		date, isStr := i.(string)
//		if !isStr {
//			return diag.Errorf("%v is not a string", date)
//		}
//		_, err := time.Parse(time.RFC3339, date)
//		if err != nil {
//			return diag.FromErr(err)
//		}
//		return nil
//	}
//}
//
//type ServiceErrorCheckFunc func(*testing.T) resource.ErrorCheckFunc
//
//var serviceErrorCheckFunc map[string]ServiceErrorCheckFunc
//
//func ErrorCheck(t *testing.T, endpointIDs ...string) resource.ErrorCheckFunc {
//	t.Helper()
//	return func(err error) error {
//		if err == nil {
//			return nil
//		}
//
//		for _, endpointID := range endpointIDs {
//			if f, ok := serviceErrorCheckFunc[endpointID]; ok {
//				ef := f(t)
//				err = ef(err)
//			}
//
//			if err == nil {
//				break
//			}
//		}
//
//		return err
//	}
//}
//
//func validateMapKeyLowerCase() schema.SchemaValidateDiagFunc {
//	return func(i interface{}, path cty.Path) diag.Diagnostics {
//		m := expandMapStringStringPtr(i)
//		for k := range m {
//			if strings.ToLower(k) != k {
//				return diag.Diagnostics{diag.Diagnostic{
//					Severity:      diag.Error,
//					AttributePath: cty.IndexStringPath(k),
//					Summary:       "Invalid map content",
//					Detail:        fmt.Sprintf("key (%s) should be lowercase", k),
//				}}
//			}
//		}
//		return nil
//	}
//}
