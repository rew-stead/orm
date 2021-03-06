package orm_test

import (
	"testing"

	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/phogolabs/orm/fixture"
)

func TestOAK(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OAK Suite")
}

type ObjP struct {
	Id int `db:"id"`
}

type ObjM struct {
	Id int
}

func (m *ObjM) Map() map[string]interface{} {
	param := make(map[string]interface{})
	param["id"] = m.Id
	param["name"] = "jack"
	return param
}
