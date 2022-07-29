package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"winprockill/internal/entity"
)

type WinCommand struct {
	appName         string
	procNamePattern string
	nssm            []byte
}

func NewWinCommand(appName, procNamePattern string, nssm []byte) *WinCommand {
	return &WinCommand{appName, procNamePattern, nssm}
}

func (w *WinCommand) Processes(ctx context.Context) (processes []entity.WinProcess, err error) {
	script := fmt.Sprintf(`"Get-Process -IncludeUserName | Where-Object {$_.ProcessName -Match "%s"} | Select Name, Id, UserName | ConvertTo-Json"`,
		w.procNamePattern)
	out, err := exec.CommandContext(ctx, "powershell", "-Command", script).Output()
	if err != nil {
		return
	}

	err = json.Unmarshal(out, &processes)
	return
}

func (w *WinCommand) KillProcesses(ctx context.Context) error {
	script := fmt.Sprintf(`"Get-Process | Where-Object {$_.ProcessName -Match "%s"} | Stop-Process -Force"`,
		w.procNamePattern)

	return exec.CommandContext(ctx, "powershell", "-Command", script).Run()
}

func (w *WinCommand) InstallAsService() error {
	f, err := os.Create("nssm.exe")
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Println("close nssm.exe", err)
		}

		err = os.Remove(f.Name())
		if err != nil {
			log.Println("remove nssm.exe", err)
		}
	}()

	if err := f.Chmod(os.FileMode(0766)); err != nil {
		return err
	}

	n, err := f.Write(w.nssm)
	if err != nil {
		return err
	}

	if n != len(w.nssm) {
		return fmt.Errorf("failed unpack nssm, writed only %d of %d", n, len(w.nssm))
	}

	dir := path.Dir(f.Name())
	script := fmt.Sprintf(`"%s install ProcessKiller "%s/%s" AppDirectory "%s" && net start ProcessKiller"`, f.Name(), dir, w.appName, dir)

	return exec.Command("cmd", "/C", script).Run()
}
