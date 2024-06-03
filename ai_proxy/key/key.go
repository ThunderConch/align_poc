package key

import (
	"fmt"
	"strings"
)

type KeyInfo struct {
	Name      string
	UpdatedAt int64
	Revoked   *bool
	Key       string
	BpmLimit  *int
	RpmLimit  *int
}

func (k *KeyInfo) Validate() error {
	invalid := []string{}

	if k.UpdatedAt <= 0 {
		invalid = append(invalid, "updatedAt")
	}

	if k.BpmLimit == nil {
		invalid = append(invalid, "bpmLimit")
	}

	if k.RpmLimit == nil {
		invalid = append(invalid, "rpmLimit")
	}

	if len(k.Key) == 0 {
		invalid = append(invalid, "key")
	}

	if len(invalid) > 0 {
		return fmt.Errorf("fields [%s] are invalid", strings.Join(invalid, ", "))
	}

	return nil
}
