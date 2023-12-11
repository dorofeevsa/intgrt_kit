package intgrt_afick

import (
	"fmt"
	"github.com/dorofeevsa/intgrt_kit/pkg/common"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type AfickOption struct {
	OptName  string
	OptValue string
}

type AfickCheckRes struct {
	Scanned        int
	New            int
	Delete         int
	Changed        int
	Dangling       int
	Exclude_suffix int
	Exclude_prefix int
	Exclude_re     int
	Degraded       int
}

func (a *AfickOption) Name() string {
	return a.OptName
}

func (a *AfickOption) Value() interface{} {
	return a.OptValue
}

const (
	OptViolationNew     = "intgrt-new"
	OptViolationChanged = "intgrt-change"
	OptViolationDelete  = "intgrt-delete"
)

const (
	OptAfickSecAlias = "afk-alias"
)

type AfickIC struct {
	afickCmdPath string
	configFile   string
}

func (a *AfickIC) RefreshIntegrityDatabase() error {

	cmd := exec.Command(a.afickCmdPath, "-u", "-c", a.configFile)

	// Output() will run command!
	out, err := cmd.Output()
	if err != nil {
		err2 := processAfickInternalError(out)
		if err2 != nil {
			return err2
		}
	}

	fmt.Printf("Refresh integrity database result:\n %s\n", out)
	return nil
}

// WARNING! Afick will return non-zero result
// even after SUCCESSFULLY operations,
// if files was changed >_<
// check this via parsing - if ok, no errors from  operation
func processAfickInternalError(out []byte) error {

	_, err := parseCheckOutput([]byte(out))
	if err != nil {
		return fmt.Errorf("RefreshIntegrityDatabase parsing error: %s", err)
	}

	return nil
}

func (a *AfickIC) HasIntegrityViolation(filepath string, checkOpts ...string) (bool, map[string]interface{}, error) {
	res, err := a.CheckFileByControl(filepath)
	if err != nil {
		return false, nil, err
	}
	if len(checkOpts) == 0 {
		return false, nil, fmt.Errorf("integrity violation options are empty")
	}

	wasChanged := false
	violationDetails := make(map[string]interface{})
	for _, opt := range checkOpts {
		switch opt {
		case OptViolationNew:
			if res.New > 0 {
				wasChanged = true
				violationDetails[OptViolationNew] = res.New
			}
		case OptViolationDelete:
			if res.Delete > 0 {
				wasChanged = true
				violationDetails[OptViolationDelete] = res.Delete
			}
		case OptViolationChanged:
			if res.Changed > 0 {
				wasChanged = true
				violationDetails[OptViolationChanged] = res.Changed
			}
		}
	}

	return wasChanged, violationDetails, nil
}

func (a *AfickIC) CheckFileByControl(file string) (*AfickCheckRes, error) {
	s, err := filepath.Abs(file)
	if err != nil {
		return nil, fmt.Errorf("CheckFileByControl can't resolve absolute path for file: %s", s)
	}
	cmd := exec.Command(a.afickCmdPath, "-l", s, "-c", a.configFile)
	out, err := cmd.Output()

	if err != nil {
		err2 := processAfickInternalError(out)
		if err2 != nil {
			return nil, err2
		}
	}

	result, err := parseCheckOutput([]byte(out))
	if err != nil {

		return nil, fmt.Errorf("integrity check error: %#v", result)
	}
	fmt.Printf("integrity check: %#v", result)
	return result, nil
}

func (a *AfickIC) InitDatabase() error {
	cmd := exec.Command(a.afickCmdPath, "-i", a.configFile)
	out, err := cmd.Output()

	if err != nil {
		err2 := processAfickInternalError(out)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func parseCheckOutput(output []byte) (*AfickCheckRes, error) {
	res := AfickCheckRes{}

	var cont = []struct {
		reg    string
		target *int
	}{
		{"([0-9]+) files", &res.Scanned},
		{"new : ([0-9]+);", &res.New},
		{"delete : ([0-9]+);", &res.Delete},
		{"changed : ([0-9]+);", &res.Changed},
		{"dangling : ([0-9]+);", &res.Dangling},
		{"exclude_suffix : ([0-9]+);", &res.Exclude_suffix},
		{"exclude_prefix : ([0-9]+);", &res.Exclude_prefix},
		{"exclude_re : ([0-9]+);", &res.Exclude_re},
		{"degraded : ([0-9]+)", &res.Degraded},
	}
	findedCount := 0
	for _, s := range cont {
		re := regexp.MustCompile(s.reg)
		val := re.FindAllSubmatchIndex(output, -1)

		if len(val) > 0 {
			i, err := strconv.Atoi(string(output[val[0][2]:val[0][3]]))
			if err != nil {
				return nil, err
			}
			*s.target = i
			findedCount++
		}
	}

	if findedCount == 0 {
		return &res, fmt.Errorf("output result check was failed: ")
	}

	return &res, nil
}

func (a *AfickIC) AddFileToIc(s string, options ...common.ICOption) error {
	data, err := prepareConfigData(a.configFile)
	if err != nil {
		return err
	}
	filename := s
	for _, o := range options {
		switch o.Name() {
		case OptAfickSecAlias:
			filename = filename + " " + o.Value().(string)

		}
	}
	data = append(data, filename)
	err = writeConfigData(a.configFile, data)
	if err != nil {
		return err
	}
	return nil
}

func NewAfickIC(afickPath string, afickConf string) (*AfickIC, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("wrapper for AfickIC designed to use in linux system!")
	}
	return &AfickIC{afickCmdPath: afickPath,
		configFile: afickConf}, nil
}

func prepareConfigData(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	configData := strings.Split(string(data), "\n")
	return configData, nil
}

func writeConfigData(filname string, data []string) error {
	out := strings.Join(data, "\n")
	err := os.WriteFile(filname, []byte(out), 0644)
	if err != nil {
		return err
	}
	return nil
}
