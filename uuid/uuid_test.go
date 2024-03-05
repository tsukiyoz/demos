package uuid

import (
	"testing"

	shortUUID "github.com/lithammer/shortuuid"

	googleUUID "github.com/google/uuid"
)

func TestGoogleUUID(t *testing.T) {
	googleUUID.SetNodeID([]byte("instance-0"))
	t.Logf("%v\n", googleUUID.New().String())
}

func TestShortUUID(t *testing.T) {
	t.Logf("%v\n", shortUUID.New())
	t.Logf("%v\n", shortUUID.NewWithNamespace("tsukiyo"))
}
