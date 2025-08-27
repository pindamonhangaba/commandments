package commandments

import (
	"reflect"
	"testing"
)

type testArgs struct {
	Host string `flag:"host, do host stuff"`
}

func TestMustCMD(t *testing.T) {
	_ = MustCMD("custom", WithConfig(

		func(config testArgs) error {
			return nil
		}), WithDefaultConfig(testArgs{
		Host: "host",
	}))
	t.Log("MustCMD'd")
}
func TestNewCMD(t *testing.T) {
	cmd, err := NewCMD("command", WithConfig(func(config testArgs) error {
		t.Log("ok")
		return nil
	}))
	t.Log(cmd, err)
}

func TestStructToFlags(t *testing.T) {
	res, err := structToFlags[struct {
		Port        int    `flag:"port,Set port number for database"`
		URL         string `flag:"url,Set url for database connection"`
		EnableHTTPS bool   `flag:"enable-https,Enable automatic https"`
	}]()
	if err != nil {
		t.Error(err)
	}
	results := []flag{
		{name: "port", usage: "Set port number for database", kind: reflect.Int},
		{name: "url", usage: "Set url for database connection", kind: reflect.String},
		{name: "enable-https", usage: "Enable automatic https", kind: reflect.Bool},
	}
	t.Log(res)
	for i, v := range res {
		if v.kind != results[i].kind || v.name != results[i].name || v.usage != v.usage {
			t.Errorf("expected %s, %s, %s, got %s, %s, %s", results[i].kind, results[i].name, results[i].usage, v.kind, v.name, v.usage)
		}
	}
}
