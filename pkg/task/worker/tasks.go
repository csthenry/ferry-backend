package worker

import (
	"context"
	"errors"
	"ferry/pkg/logger"
	"os/exec"
	"syscall"

	"github.com/RichardKnop/machinery/v1/tasks"
)

var asyncTaskMap map[string]interface{}

func executeTaskBase(envPath string, scriptPath string, params string) (err error) {
	command := exec.Command(envPath, scriptPath, params) //初始化Cmd
	out, err := command.CombinedOutput()
	logger.Info("envPath: ", envPath)
	// 解决 Windows 下 python3 可能不可用的问题
	if err != nil && envPath == "python3" {
		logger.Info("Python3 -> Python")
		envPath = "python"
		command = exec.Command(envPath, scriptPath, params)
		out, err = command.CombinedOutput()
	}
	if err != nil {
		logger.Info("Output: ", string(out))
		logger.Errorf("Task exec failed: %v", err.Error())
		return
	}
	logger.Info("Output: ", string(out))
	logger.Info("ProcessState PID: ", command.ProcessState.Pid())
	logger.Info("Exit Code ", command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
	return
}

// ExecCommand 异步任务
func ExecCommand(classify string, scriptPath string, params string) (err error) {
	if classify == "shell" {
		logger.Info("start exec shell - ", scriptPath)
		err = executeTaskBase("/bin/sh", scriptPath, params)
		if err != nil {
			return
		}
	} else if classify == "python" {
		logger.Info("start exec python - ", scriptPath)
		err = executeTaskBase("python3", scriptPath, params)
		if err != nil {
			return
		}
	} else {
		err = errors.New("目前仅支持Python与Shell脚本的执行，请知悉。")
		return
	}
	return
}

func SendTask(ctx context.Context, classify string, scriptPath string, params string) {
	args := make([]tasks.Arg, 0)
	args = append(args, tasks.Arg{
		Name:  "classify",
		Type:  "string",
		Value: classify,
	})
	args = append(args, tasks.Arg{
		Name:  "scriptPath",
		Type:  "string",
		Value: scriptPath,
	})
	args = append(args, tasks.Arg{
		Name:  "params",
		Type:  "string",
		Value: params,
	})
	task, _ := tasks.NewSignature("ExecCommandTask", args)
	task.RetryCount = 5
	_, err := AsyncTaskCenter.SendTaskWithContext(ctx, task)
	if err != nil {
		logger.Error(err.Error())
	}
}

func initAsyncTaskMap() {
	asyncTaskMap = make(map[string]interface{})
	asyncTaskMap["ExecCommandTask"] = ExecCommand
}
