package buildinfo_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func compile(code string, args ...string) (string, error) {
	if err := os.Mkdir("./test", 0755); err != nil {
		return "", err
	}

	defer os.RemoveAll("./test")

	if err := os.WriteFile("./test/main.go", []byte(code), 0644); err != nil {
		return "", err
	}

	os.Chdir("./test")
	defer os.Chdir("..")

	args = append([]string{"build"}, args...)
	out, err := exec.Command("go", args...).Output()

	if err != nil {
		return string(out), err
	}

	out, err = exec.Command("./test").Output()

	if err != nil {
		return string(out), err
	}

	return strings.TrimSpace(string(out)), nil
}

func TestInspect(t *testing.T) {
	output, err := compile(`
		package main

		import (
			"fmt"
			"github.com/quartercastle/buildinfo"
		)

		type Test struct {}

		func main() {
			fmt.Println(buildinfo.Inspect(Test{}))
		}
	`)

	if err != nil {
		t.Error(err)
	}

	expected := "&{main dev  <nil>}"

	if output != expected {
		t.Errorf("expected %s; got %s", expected, output)
	}
}

func TestVersionVCS(t *testing.T) {
	output, err := compile(`
		package main

		import (
			"fmt"
			"github.com/quartercastle/buildinfo"
		)

		func main() {
			fmt.Println(buildinfo.Version())
		}
	`)

	if err != nil {
		t.Error(err)
	}

	sha, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	expected := strings.TrimSpace(string(sha))

	if output != expected {
		t.Errorf("expected %s; got %s", expected, output)
	}
}

func TestVersionInjected(t *testing.T) {
	output, err := compile(`
		package main

		import (
			"fmt"
			"github.com/quartercastle/buildinfo"
		)

		func main() {
			fmt.Println(buildinfo.Version())
		}
	`, "-ldflags", "-X github.com/quartercastle/buildinfo.version=v1.0.0")

	if err != nil {
		t.Error(err)
	}

	expected := "v1.0.0"

	if output != expected {
		t.Errorf("expected %s; got %s", expected, output)
	}
}
