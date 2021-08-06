package slide_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hello-slide/slide-manager/slide"
	"github.com/hello-slide/slide-manager/storage"
)

func TestSlideManager(t *testing.T) {
	filePath := "helloslide-82be1a83e598-key.json"
	_, err := os.Stat(filePath)

	if err == nil {
		ctx := context.Background()
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatal(err)
		}
		storage.Key = bytes
		client, err := storage.CreateClient(ctx)
		if err != nil {
			t.Fatal(err)
		}

		slideManger := slide.NewSlideManager(ctx, *client, "helloslide-test", "test_user")

		// id, err := slideManger.Create("test_title")
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// fmt.Println(id)

		info, err := slideManger.GetInfo()
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(info)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
