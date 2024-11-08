package configs

import (
	"fmt"
	"testing"

	"github.com/bufbuild/protovalidate-go"
)

func TestValidate(t *testing.T) {
	msg := &Bootstrap{
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

	msg = &Bootstrap{
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

	msg = &Bootstrap{
		Mode: "other",
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
}
