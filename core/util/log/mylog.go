/*
	This source you can use!
*/
package log

import "github.com/hunterhug/fafacms/core/util"

func New(filename string) {
	logsconf, err := util.ReadfromFile(filename)
	if err != nil {
		panic(err)
	}
	err = Init(string(logsconf))
	if err != nil {
		panic(err)
	}
}
