package alert

import (
	"fmt"
	"gosecurity/sys"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

var AlertPath string

var LoadedAlerts []Alert

type Alert struct {
	Source     string      `yaml:"source"`
	Query      string      `yaml:"query"`
	Conditions []Condition `yaml:"conditions"`
	Index      string      `yaml:"index"`
	Results    chan ResultsChannel
}

type Condition struct {
	Field    string `yaml:"field"`
	Operator string `yaml:"operator"`
	Value    string `yaml:"value"`
}

type ResultsChannel struct {
	Result bool
	Event  map[string]interface{}
	Reason string
}

func ExecuteAll() {
	for _, alert := range LoadedAlerts {
		go alert.Run()
	}
}

func SetAlertPath(path string) {
	AlertPath = path
}

func (a Alert) Fire(event map[string]interface{}, reason string) {
	fmt.Println("--------------------------------")
	fmt.Println("Firing alert")
	fmt.Println("Event", event)
	fmt.Println("Reason", reason)
	fmt.Println("--------------------------------")
}

func Load(path string) (Alert, error) {
	fmt.Println("Loading alert config")
	file, err := os.Open(AlertPath + path)
	if err != nil {
		return Alert{}, fmt.Errorf("error reading YAML file: %w", err)
	}
	defer file.Close()

	var config Alert

	// Unmarshal the YAML data into the Config struct
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return Alert{}, err
	}
	LoadedAlerts = append(LoadedAlerts, config)
	return config, err
}

// func (c Condition) CompareInt(first int, second int) bool {
// 	switch c.Operator {
// 	case "GREATER":
// 		return GreaterInt(first, second)
// 	default:
// 		return false
// 	}
// }

func (a Alert) Run() (results []map[string]interface{}, triggered bool) {
	res := sys.System.Db.Query(a.Index, a.Query)
	results = sys.System.Db.GetResults(res)

	// Print the results
	fmt.Println("Query results:")
	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}

	// Check all the conditions against the results
	triggered = a.Check(results)

	return results, triggered
}

func (a Alert) Check(r []map[string]interface{}) bool {
	// Create initial channels for comms
	// When we are fully done
	done := make(chan bool)
	// When we get a result from the checks
	res := make(chan bool)
	// The number of checks we have done
	count := make(chan int)

	// Printing for logs
	fmt.Println("Starting the checks")
	fmt.Println("Number of conditions:", len(a.Conditions))
	fmt.Println("Number of rows:", len(r))

	// Create the results channel for this specific alert
	a.Results = make(chan ResultsChannel, len(a.Conditions))

	// Start the process checks to gather results using the channels
	go a.processChecks(len(r), res, done, count)

	// Loop through each Row we are given from results and check all conditions on it
	for _, row := range r {
		go a.CheckConditions(&a.Results, row)
	}

	// Wait for the results to come back
	// final = Did we have an alert trigger for any event
	final := <-res
	// The number of checks is the number of conditions that were met
	num_checks := <-count
	// Wait for the process checks to finish
	<-done
	fmt.Println("Done processing all checks")
	fmt.Println("Number of FIRES:", num_checks)
	return final
}

func (a Alert) processChecks(size int, res chan bool, done chan bool, num chan int) {
	// The rolling bool is the result of all the conditions
	rollingBool := false
	// The count is the number of rows that were met
	count := 0
	for i := 0; i < size; i++ {
		// a.Results has res, the event that triggered it and the reason
		finishedCheck := <-a.Results // Read from the channel
		if finishedCheck.Result {
			// Fire the alert
			a.Fire(finishedCheck.Event, finishedCheck.Reason)
			count++
		}
		rollingBool = rollingBool || finishedCheck.Result
	}
	res <- rollingBool
	num <- count
	done <- true
}

func (a Alert) CheckConditions(aResults *chan ResultsChannel, r map[string]interface{}) {
	rollingBool := true
	reason := "The following conditions were met:\n"
	for _, cond := range a.Conditions {
		res := cond.Check(r)
		rollingBool = res && rollingBool
		if res {
			reason += fmt.Sprintf("Condition: %v, Result: %v\n", cond, res)
		}
	}
	*aResults <- ResultsChannel{Result: rollingBool, Event: r, Reason: reason}
}

// func (c Condition) CheckCondition(ch *chan ResultsChannel, r []map[string]interface{}) {
// 	rollingBool := true
// 	for _, row := range r {
// 		res := c.Check(row)
// 		rollingBool = res && rollingBool
// 		*ch <- ResultsChannel{Result: rollingBool, Event: row}
// 	}
// }

// func (a Alert) CheckRow(r map[string]interface{}) bool {
// 	for _, cond := range a.Conditions {
// 		cond.Check(r)
// 	}
// }

func checkYAMLValueIsInt(strValue string) (int, bool) {
	// Attempt to convert to integer
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		fmt.Printf("Value '%s' cannot be converted to an integer\n", strValue)
		return 0, false
	}

	return intValue, true
}

func (c Condition) Check(r map[string]interface{}) bool {
	intvalue, boolRes := checkYAMLValueIsInt(c.Value)
	fmt.Println("Int value:", intvalue)
	if boolRes {
		return c.CheckInt(r)
	}
	return c.CheckString(r)
}

func (c Condition) CheckInt(r map[string]interface{}) bool {
	fmt.Printf("Checking Int - Field: %s, Value: %v, Operator: %s\n", c.Field, r[c.Field], c.Operator)

	var intValue int
	var err error

	switch v := r[c.Field].(type) {
	case int:
		intValue = v
	case int64:
		intValue = int(v)
	case float64:
		intValue = int(v)
	case string:
		intValue, err = strconv.Atoi(v)
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)
			return false
		}
	default:
		fmt.Printf("Unsupported type for field %s: %T\n", c.Field, v)
		return false
	}

	conditionValue, err := strconv.Atoi(c.Value)
	if err != nil {
		fmt.Printf("Error converting condition value to int: %v\n", err)
		return false
	}

	switch c.Operator {
	case "GREATER":
		return GreaterInt(intValue, conditionValue)
	case "LESS":
		return LessInt(intValue, conditionValue)
	case "EQUALS":
		return EqualsInt(intValue, conditionValue)
	default:
		fmt.Printf("Unknown operator for int comparison: %s\n", c.Operator)
		return false
	}
}

func (c Condition) CheckString(r map[string]interface{}) bool {
	fmt.Printf("Checking String - Field: %s, Value: %v, Operator: %s\n", c.Field, r[c.Field], c.Operator)

	value, ok := r[c.Field].(string)
	if !ok {
		fmt.Printf("Type assertion failed for field %s\n", c.Field)
		return false
	}

	conditionValue := c.Value

	switch c.Operator {
	case "EQUALS":
		return EqualsString(value, conditionValue)
	case "NOT_EQUALS":
		return NotEqualsString(value, conditionValue)
	default:
		fmt.Printf("Unknown operator for string comparison: %s\n", c.Operator)
		return false
	}
}
