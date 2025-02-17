package ice_cream_maker_test

import (
	"testing"

	"github.com/jaehoonkim/sentinel/pkg/manager/database/vanilla/ice_cream_maker"
)

func TestColumnValue(t *testing.T) {

	objs := []interface{}{
		ServiceStep_essential{},
		ServiceStep{},
	}

	s, err := ice_cream_maker.ColumnValues(objs...)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}
