package bootstrap

import (
	"fmt"
	"testing"

	"github.com/bufbuild/protovalidate-go"

	"origadmin/basic-layout/internal/configs"
)

func TestValidate(t *testing.T) {
	msg := &configs.Bootstrap{
		Mode: "cluster",
	}

	v, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate(msg); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}

	msg = &configs.Bootstrap{
		Mode: "cluster",
	}

	v, err = protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate(msg); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}

	msg = &configs.Bootstrap{
		Mode: "other",
	}

	v, err = protovalidate.New(protovalidate.WithFailFast(true))
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate(msg); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}
}
