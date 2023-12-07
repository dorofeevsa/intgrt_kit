package main

import (
	"fmt"
	"github.com/dorofeevsa/intgrt_kit/pkg/common"
	afick "github.com/dorofeevsa/intgrt_kit/pkg/intgrt_afick"

	"os"
	"time"
)

func main() {

	cntr, err := afick.NewAfickIC("/usr/bin/afick", "/home/drofa/work/intgrt_kit/example.conf")
	interfaceCheck(cntr)
	err = cntr.InitDatabase()
	if err != nil {
		fmt.Printf("Problem witj afick init: %s", err)
	}

	if err != nil {
		panic("Afick wrapper init fail")
	}
	err = cntr.AddFileToIc("/home/drofa/work/intgrt_kit/example_folder/", &afick.AfickOption{OptName: afick.OptAfickSecAlias, OptValue: "PARSEC"})
	if err != nil {
		fmt.Printf("Problem: %s", err)
	}

	//initial file
	tempData := time.Now().String()

	filename := "example_folder/example_log"
	_ = os.WriteFile(filename, []byte(tempData), os.ModeAppend)
	_ = cntr.RefreshIntegrityDatabase()

	go func() {
		for {
			time.Sleep(time.Second * 2)
			file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				continue
			}

			tempData = time.Now().String()
			_, err = file.WriteString("ololo")
			err = file.Close()
			violation, res, err := cntr.HasIntegrityViolation("/home/drofa/work/intgrt_kit/example_folder/", afick.OptViolationNew, afick.OptViolationChanged, afick.OptViolationDelete)
			if violation {
				fmt.Printf("Integrity violation: %#v", res)
			}

			fmt.Printf("CheckedRes: %#v", res)

			err = cntr.RefreshIntegrityDatabase()
			if err != nil {
				continue
			}
		}

	}()

	select {}
}

func interfaceCheck(cntr common.IntegrityController) {
	fmt.Printf("IntegrityController interface matching done: #%v", cntr)
}
