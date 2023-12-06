package os_plugin

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const query = "ps -ef | grep rp1 | grep -v grep | grep -v rp2 | awk '{print $2}'"

func StopProcess(name string) {
	text := strings.Replace(query, "rp1", name, -1)
	text = strings.Replace(text, "rp2", strconv.Itoa(os.Getpid()), -1)
	os.WriteFile("temp.sh", []byte(text), 0777)
	defer os.Remove("temp.sh")

	out, err := exec.Command("bash", "temp.sh").Output()
	if err != nil {
		return
	}
	exec.Command("kill", "-2", strings.TrimSpace(string(out))).Run()
}
