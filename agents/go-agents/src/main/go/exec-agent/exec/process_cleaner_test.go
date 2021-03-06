//
// Copyright (c) 2012-2017 Codenvy, S.A.
// All rights reserved. This program and the accompanying materials
// are made available under the terms of the Eclipse Public License v1.0
// which accompanies this distribution, and is available at
// http://www.eclipse.org/legal/epl-v10.html
//
// Contributors:
//   Codenvy, S.A. - initial API and implementation
//

package exec_test

import (
	"testing"
	"time"

	"github.com/eclipse/che/agents/go-agents/src/main/go/exec-agent/exec"
)

func TestCleanWithZeroThreshold(t *testing.T) {
	p := startAndWaitTestProcess(testCmd, t)
	defer cleanupLogsDir()

	exec.NewCleaner(0, 0).CleanOnce()

	_, err := exec.Get(p.Pid)
	if err == nil {
		t.Fatal("Must not exist")
	}
	if _, ok := err.(*exec.NoProcessError); !ok {
		t.Fatal(err)
	}
}

func TestCleansOnlyUnusedProcesses(t *testing.T) {
	p1 := startAndWaitTestProcess(testCmd, t)
	p2 := startAndWaitTestProcess(testCmd, t)

	time.Sleep(500 * time.Millisecond)

	// use one of the processes, so it is used now
	exec.Get(p1.Pid)

	// cleanup immediately
	(&exec.Cleaner{CleanupPeriod: 0, CleanupThreshold: 500 * time.Millisecond}).CleanOnce()

	_, err1 := exec.Get(p1.Pid)
	_, err2 := exec.Get(p2.Pid)

	// process 1 must be cleaned
	if err1 != nil {
		t.Fatalf("Expected process 2 to exist, but got an error: %s", err1.Error())
	}

	// process 2 must exist
	if _, ok := err2.(*exec.NoProcessError); !ok {
		t.Fatal("Expected process 2 to be cleaned")
	}
}
