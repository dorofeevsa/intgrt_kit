package ololo

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func Hello() {
	logrus.Info("logrus dep")
	fmt.Println("testing packege hosting")
}

func HelloTwo() {
	logrus.Info("logrus dep 2 ")
	fmt.Println("testing packege hosting 2 ")
}
