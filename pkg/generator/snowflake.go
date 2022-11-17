// Package generator
package generator

import (
	"fmt"

	"github.com/brainlabs/snowflake"

	"sharefood/pkg/logger"
)

var (
	snowFlakeGenerator *snowflake.Node
)

// Setup initiated snowflake
func Setup(node uint64) {
	s, err := snowflake.NewNode(int64(node))

	if err != nil {
		logger.Fatal(fmt.Sprintf("snowflake generator error %s", err.Error()), logger.EventName("generator"))
	}

	snowFlakeGenerator = s
}

// GenerateInt64 generate id int64
func GenerateInt64() int64 {
	return snowFlakeGenerator.Generate().Int64()
}

// GenerateString generate id string number
func GenerateString() string {
	return snowFlakeGenerator.Generate().String()
}
