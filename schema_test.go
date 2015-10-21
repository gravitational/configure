package configure

import (
	"github.com/gravitational/configure/test"
	"github.com/gravitational/log"
	. "gopkg.in/check.v1"
)

type SchemaSuite struct {
	test.ConfigSuite
}

var _ = Suite(&SchemaSuite{})

func (s *SchemaSuite) SetUpSuite(c *C) {
	log.Initialize("console", "INFO")
}

func (s *SchemaSuite) TestGenerateSchema(c *C) {
	var cfg test.Config
	schema, err := GenerateSchema(&cfg)
	c.Assert(err, IsNil)
	c.Assert(schema.Vars, DeepEquals, nil)
}
