// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package loader

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
	"github.com/sirupsen/logrus"

	"github.com/cilium/cilium/pkg/command/exec"
	"github.com/cilium/cilium/pkg/common"
	"github.com/cilium/cilium/pkg/datapath/linux/probes"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/option"
)

// OutputType determines the type to be generated by the compilation steps.
type OutputType string

const (
	outputObject   = OutputType("obj")
	outputAssembly = OutputType("asm")
	outputSource   = OutputType("c")

	compiler = "clang"
	linker   = "llc"

	endpointPrefix   = "bpf_lxc"
	endpointProg     = endpointPrefix + "." + string(outputSource)
	endpointObj      = endpointPrefix + ".o"
	endpointObjDebug = endpointPrefix + ".dbg.o"
	endpointAsm      = endpointPrefix + "." + string(outputAssembly)

	hostEndpointPrefix       = "bpf_host"
	hostEndpointNetdevPrefix = "bpf_netdev_"
	hostEndpointProg         = hostEndpointPrefix + "." + string(outputSource)
	hostEndpointObj          = hostEndpointPrefix + ".o"
	hostEndpointObjDebug     = hostEndpointPrefix + ".dbg.o"
	hostEndpointAsm          = hostEndpointPrefix + "." + string(outputAssembly)

	networkPrefix = "bpf_network"
	networkProg   = networkPrefix + "." + string(outputSource)
	networkObj    = networkPrefix + ".o"

	xdpPrefix = "bpf_xdp"
	xdpProg   = xdpPrefix + "." + string(outputSource)
	xdpObj    = xdpPrefix + ".o"
)

var (
	probeCPUOnce sync.Once

	// default fallback
	nameBPFCPU = "v1"
)

// progInfo describes a program to be compiled with the expected output format
type progInfo struct {
	// Source is the program source (base) filename to be compiled
	Source string
	// Output is the expected (base) filename produced from the source
	Output string
	// OutputType to be created by LLVM
	OutputType OutputType
	// Options are passed directly to LLVM as individual parameters
	Options []string
}

// directoryInfo includes relevant directories for compilation and linking
type directoryInfo struct {
	// Library contains the library code to be used for compilation
	Library string
	// Runtime contains headers for compilation
	Runtime string
	// State contains node, lxc, and features headers for templatization
	State string
	// Output is the directory where the files will be stored
	Output string
}

var (
	standardCFlags = []string{"-O2", "-target", "bpf", "-std=gnu89",
		"-nostdinc", fmt.Sprintf("-D__NR_CPUS__=%d", common.GetNumPossibleCPUs(log)),
		"-Wall", "-Wextra", "-Werror", "-Wshadow",
		"-Wno-address-of-packed-member",
		"-Wno-unknown-warning-option",
		"-Wno-gnu-variable-sized-type-not-at-end",
		"-Wdeclaration-after-statement",
		"-Wimplicit-int-conversion",
		"-Wenum-conversion"}
	standardLDFlags = []string{"-march=bpf"}

	// testIncludes allows the unit tests to inject additional include
	// paths into the compile command at test time. It is usually nil.
	testIncludes []string

	debugProgs = []*progInfo{
		{
			Source:     endpointProg,
			Output:     endpointObjDebug,
			OutputType: outputObject,
		},
		{
			Source:     endpointProg,
			Output:     endpointAsm,
			OutputType: outputAssembly,
		},
		{
			Source:     endpointProg,
			Output:     endpointProg,
			OutputType: outputSource,
		},
	}
	debugHostProgs = []*progInfo{
		{
			Source:     hostEndpointProg,
			Output:     hostEndpointObjDebug,
			OutputType: outputObject,
		},
		{
			Source:     hostEndpointProg,
			Output:     hostEndpointAsm,
			OutputType: outputAssembly,
		},
		{
			Source:     hostEndpointProg,
			Output:     hostEndpointProg,
			OutputType: outputSource,
		},
	}
	epProg = &progInfo{
		Source:     endpointProg,
		Output:     endpointObj,
		OutputType: outputObject,
	}
	hostEpProg = &progInfo{
		Source:     hostEndpointProg,
		Output:     hostEndpointObj,
		OutputType: outputObject,
	}
	networkTcProg = &progInfo{
		Source:     networkProg,
		Output:     networkObj,
		OutputType: outputObject,
	}
)

// GetBPFCPU returns the BPF CPU for this host.
func GetBPFCPU() string {
	probeCPUOnce.Do(func() {
		if !option.Config.DryMode {
			// We can probe the availability of BPF instructions indirectly
			// based on what kernel helpers are available when both were
			// added in the same release.
			// We want to enable v3 only on kernels 5.10+ where we have
			// tested it and need it to work around complexity issues.
			if probes.HaveV3ISA() == nil {
				if probes.HaveProgramHelper(ebpf.SchedCLS, asm.FnRedirectNeigh) == nil {
					nameBPFCPU = "v3"
					return
				}
			}
			// We want to enable v2 on all kernels that support it, that is,
			// kernels 4.14+.
			if probes.HaveV2ISA() == nil {
				nameBPFCPU = "v2"
			}
		}
	})
	return nameBPFCPU
}

// progLDFlags determines the loader flags for the specified prog and paths.
func progLDFlags(prog *progInfo, dir *directoryInfo) []string {
	return []string{
		fmt.Sprintf("-filetype=%s", prog.OutputType),
		"-o", path.Join(dir.Output, prog.Output),
	}
}

// prepareCmdPipes attaches pipes to the stdout and stderr of the specified
// command, and returns the stdout, stderr, and any error that may have
// occurred while creating the pipes.
func prepareCmdPipes(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get stdout pipe: %s", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdout.Close()
		return nil, nil, fmt.Errorf("Failed to get stderr pipe: %s", err)
	}

	return stdout, stderr, nil
}

func pidFromProcess(proc *os.Process) string {
	result := "not-started"
	if proc != nil {
		result = fmt.Sprintf("%d", proc.Pid)
	}
	return result
}

// compileAndLink links the specified program from the specified path to the
// intermediate representation, to the output specified in the prog's info.
func compileAndLink(ctx context.Context, prog *progInfo, dir *directoryInfo, compileArgs ...string) error {
	compileCmd, cancelCompile := exec.WithCancel(ctx, compiler, compileArgs...)
	defer cancelCompile()
	compilerStdout, compilerStderr, err := prepareCmdPipes(compileCmd)
	if err != nil {
		return err
	}

	linkArgs := make([]string, 0, 8)
	linkArgs = append(linkArgs, standardLDFlags...)
	linkArgs = append(linkArgs, "-mcpu="+GetBPFCPU())
	linkArgs = append(linkArgs, progLDFlags(prog, dir)...)

	linkCmd := exec.CommandContext(ctx, linker, linkArgs...)
	linkCmd.Stdin = compilerStdout
	if err := compileCmd.Start(); err != nil {
		return fmt.Errorf("Failed to start command %s: %s", compileCmd.Args, err)
	}

	var compileOut []byte
	/* Ignoring the output here because pkg/command/exec will log it. */
	_, err = linkCmd.CombinedOutput(log, true)
	if err == nil {
		compileOut, _ = io.ReadAll(compilerStderr)
		err = compileCmd.Wait()
	} else {
		cancelCompile()
	}
	if err != nil {
		err = fmt.Errorf("Failed to compile %s: %w", prog.Output, err)

		// In linux/unix based implementations, cancelling the context for a cmd.Run() will
		// return errors: "context cancelled" if the context is cancelled prior to the process
		// starting and "signal: killed" if it is already running.
		// This can mess up calling logging logic which expects the returned error to have
		// context.Cancelled so we join this error in to fix that.
		if errors.Is(ctx.Err(), context.Canceled) &&
			compileCmd.ProcessState != nil &&
			!compileCmd.ProcessState.Exited() &&
			strings.HasSuffix(err.Error(), syscall.SIGKILL.String()) {
			err = errors.Join(err, ctx.Err())
		}

		if !errors.Is(err, context.Canceled) {
			log.WithFields(logrus.Fields{
				"compiler-pid": pidFromProcess(compileCmd.Process),
				"linker-pid":   pidFromProcess(linkCmd.Process),
			}).Error(err)
		}

		if compileOut != nil {
			scanner := bufio.NewScanner(bytes.NewReader(compileOut))
			for scanner.Scan() {
				log.Warn(scanner.Text())
			}
		}
	}

	return err
}

// progCFlags determines the compiler flags for the specified prog and paths.
func progCFlags(prog *progInfo, dir *directoryInfo) []string {
	var output string

	if prog.OutputType == outputSource {
		output = path.Join(dir.Output, prog.Output)
	} else {
		output = "-" // stdout
	}

	return append(testIncludes,
		fmt.Sprintf("-I%s", path.Join(dir.Runtime, "globals")),
		fmt.Sprintf("-I%s", dir.State),
		fmt.Sprintf("-I%s", dir.Library),
		fmt.Sprintf("-I%s", path.Join(dir.Library, "include")),
		"-c", path.Join(dir.Library, prog.Source),
		"-o", output,
	)
}

// compile and link a program.
func compile(ctx context.Context, prog *progInfo, dir *directoryInfo) (err error) {
	args := make([]string, 0, 16)
	if prog.OutputType == outputSource {
		args = append(args, "-E") // Preprocessor
	} else {
		args = append(args, "-emit-llvm")
		args = append(args, "-g")
	}

	args = append(args, standardCFlags...)
	args = append(args, prog.Options...)
	args = append(args, progCFlags(prog, dir)...)

	// Compilation is split between two exec calls. First clang generates
	// LLVM bitcode and then later llc compiles it to byte-code.
	log.WithFields(logrus.Fields{
		"target": compiler,
		"args":   args,
	}).Debug("Launching compiler")
	switch prog.OutputType {
	case outputSource:
		compileCmd := exec.CommandContext(ctx, compiler, args...)
		_, err = compileCmd.CombinedOutput(log, true)
	case outputObject, outputAssembly:
		err = compileAndLink(ctx, prog, dir, args...)
	default:
		log.Fatalf("Unhandled progInfo.OutputType %s", prog.OutputType)
	}

	return err
}

// compileDatapath invokes the compiler and linker to create all state files for
// the BPF datapath, with the primary target being the BPF ELF binary.
//
// It also creates the following output files:
// * Preprocessed C
// * Assembly
// * Object compiled with debug symbols
func compileDatapath(ctx context.Context, dirs *directoryInfo, isHost bool, logger *logrus.Entry) error {
	scopedLog := logger.WithField(logfields.Debug, true)

	versionCmd := exec.CommandContext(ctx, compiler, "--version")
	compilerVersion, err := versionCmd.CombinedOutput(scopedLog, true)
	if err != nil {
		return err
	}
	versionCmd = exec.CommandContext(ctx, linker, "--version")
	linkerVersion, err := versionCmd.CombinedOutput(scopedLog, true)
	if err != nil {
		return err
	}
	scopedLog.WithFields(logrus.Fields{
		compiler: string(compilerVersion),
		linker:   string(linkerVersion),
	}).Debug("Compiling datapath")

	if option.Config.Debug {
		// Write out assembly and preprocessing files for debugging purposes
		progs := debugProgs
		if isHost {
			progs = debugHostProgs
		}
		for _, p := range progs {
			if err := compile(ctx, p, dirs); err != nil {
				// Only log an error here if the context was not canceled. This log message
				// should only represent failures with respect to compiling the program.
				if !errors.Is(err, context.Canceled) {
					scopedLog.WithField(logfields.Params, logfields.Repr(p)).WithError(err).Debug("JoinEP: Failed to compile")
				}
				return err
			}
		}
	}

	// Compile the new program
	prog := epProg
	if isHost {
		prog = hostEpProg
	}
	if err := compile(ctx, prog, dirs); err != nil {
		// Only log an error here if the context was not canceled. This log message
		// should only represent failures with respect to compiling the program.
		if !errors.Is(err, context.Canceled) {
			scopedLog.WithField(logfields.Params, logfields.Repr(prog)).WithError(err).Warn("JoinEP: Failed to compile")
		}
		return err
	}

	return nil
}

// CompileWithOptions compiles a BPF program generating an object file,
// using a set of provided compiler options.
func CompileWithOptions(ctx context.Context, src string, out string, opts []string) error {
	prog := progInfo{
		Source:     src,
		Options:    opts,
		Output:     out,
		OutputType: outputObject,
	}
	dirs := directoryInfo{
		Library: option.Config.BpfDir,
		Runtime: option.Config.StateDir,
		Output:  option.Config.StateDir,
		State:   option.Config.StateDir,
	}
	return compile(ctx, &prog, &dirs)
}

// Compile compiles a BPF program generating an object file.
func Compile(ctx context.Context, src string, out string) error {
	return CompileWithOptions(ctx, src, out, nil)
}

// compileTemplate compiles a BPF program generating a template object file.
func compileTemplate(ctx context.Context, out string, isHost bool) error {
	dirs := directoryInfo{
		Library: option.Config.BpfDir,
		Runtime: option.Config.StateDir,
		Output:  out,
		State:   out,
	}
	return compileDatapath(ctx, &dirs, isHost, log)
}

// compileNetwork compiles a BPF program attached to network
func compileNetwork(ctx context.Context) error {
	dirs := directoryInfo{
		Library: option.Config.BpfDir,
		Runtime: option.Config.StateDir,
		Output:  option.Config.StateDir,
		State:   option.Config.StateDir,
	}
	scopedLog := log.WithField(logfields.Debug, true)

	versionCmd := exec.CommandContext(ctx, compiler, "--version")
	compilerVersion, err := versionCmd.CombinedOutput(scopedLog, true)
	if err != nil {
		return err
	}
	versionCmd = exec.CommandContext(ctx, linker, "--version")
	linkerVersion, err := versionCmd.CombinedOutput(scopedLog, true)
	if err != nil {
		return err
	}
	scopedLog.WithFields(logrus.Fields{
		compiler: string(compilerVersion),
		linker:   string(linkerVersion),
	}).Debug("Compiling network programs")

	// Write out assembly and preprocessing files for debugging purposes
	if err := compile(ctx, networkTcProg, &dirs); err != nil {
		scopedLog.WithField(logfields.Params, logfields.Repr(networkTcProg)).
			WithError(err).Warn("Failed to compile")
		return err
	}
	return nil
}
