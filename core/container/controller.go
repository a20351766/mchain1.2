/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/context"

	"github.com/hyperledger/mchain/common/flogging"
	"github.com/hyperledger/mchain/core/chaincode/platforms"
	"github.com/hyperledger/mchain/core/container/ccintf"
	pb "github.com/hyperledger/mchain/protos/peer"
)

type VMProvider interface {
	NewVM() VM
}

type Builder interface {
	Build() (io.Reader, error)
}

//VM is an abstract virtual image for supporting arbitrary virual machines
type VM interface {
	Start(ctxt context.Context, ccid ccintf.CCID, args []string, env []string, filesToUpload map[string][]byte, builder Builder) error
	Stop(ctxt context.Context, ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error
}

type refCountedLock struct {
	refCount int
	lock     *sync.RWMutex
}

//VMController - manages VMs
//   . abstract construction of different types of VMs (we only care about Docker for now)
//   . manage lifecycle of VM (start with build, start, stop ...
//     eventually probably need fine grained management)
type VMController struct {
	sync.RWMutex
	containerLocks map[string]*refCountedLock
	vmProviders    map[string]VMProvider
}

var vmLogger = flogging.MustGetLogger("container")

// NewVMController creates a new instance of VMController
func NewVMController(vmProviders map[string]VMProvider) *VMController {
	return &VMController{
		containerLocks: make(map[string]*refCountedLock),
		vmProviders:    vmProviders,
	}
}

func (vmc *VMController) newVM(typ string) VM {
	v, ok := vmc.vmProviders[typ]
	if !ok {
		vmLogger.Panicf("Programming error: unsupported VM type: %s", typ)
	}
	return v.NewVM()
}

func (vmc *VMController) lockContainer(id string) {
	//get the container lock under global lock
	vmc.Lock()
	var refLck *refCountedLock
	var ok bool
	if refLck, ok = vmc.containerLocks[id]; !ok {
		refLck = &refCountedLock{refCount: 1, lock: &sync.RWMutex{}}
		vmc.containerLocks[id] = refLck
	} else {
		refLck.refCount++
		vmLogger.Debugf("refcount %d (%s)", refLck.refCount, id)
	}
	vmc.Unlock()
	vmLogger.Debugf("waiting for container(%s) lock", id)
	refLck.lock.Lock()
	vmLogger.Debugf("got container (%s) lock", id)
}

func (vmc *VMController) unlockContainer(id string) {
	vmc.Lock()
	if refLck, ok := vmc.containerLocks[id]; ok {
		if refLck.refCount <= 0 {
			panic("refcnt <= 0")
		}
		refLck.lock.Unlock()
		if refLck.refCount--; refLck.refCount == 0 {
			vmLogger.Debugf("container lock deleted(%s)", id)
			delete(vmc.containerLocks, id)
		}
	} else {
		vmLogger.Debugf("no lock to unlock(%s)!!", id)
	}
	vmc.Unlock()
}

//VMCReq - all requests should implement this interface.
//The context should be passed and tested at each layer till we stop
//note that we'd stop on the first method on the stack that does not
//take context
type VMCReq interface {
	Do(ctxt context.Context, v VM) error
	GetCCID() ccintf.CCID
}

//StartContainerReq - properties for starting a container.
type StartContainerReq struct {
	ccintf.CCID
	Builder       Builder
	Args          []string
	Env           []string
	FilesToUpload map[string][]byte
}

// PlatformBuilder implements the Build interface using
// the platforms package GenerateDockerBuild function.
// XXX This is a pretty awkward spot for the builder, it should
// really probably be pushed into the dockercontroller, as it only
// builds docker images, but, doing so would require contaminating
// the dockercontroller package with the CDS, which is also
// undesirable.
type PlatformBuilder struct {
	DeploymentSpec *pb.ChaincodeDeploymentSpec
}

// Build a tar stream based on the CDS
func (b *PlatformBuilder) Build() (io.Reader, error) {
	return platforms.GenerateDockerBuild(b.DeploymentSpec)
}

func (si StartContainerReq) Do(ctxt context.Context, v VM) error {
	return v.Start(ctxt, si.CCID, si.Args, si.Env, si.FilesToUpload, si.Builder)
}

func (si StartContainerReq) GetCCID() ccintf.CCID {
	return si.CCID
}

//StopContainerReq - properties for stopping a container.
type StopContainerReq struct {
	ccintf.CCID
	Timeout uint
	//by default we will kill the container after stopping
	Dontkill bool
	//by default we will remove the container after killing
	Dontremove bool
}

func (si StopContainerReq) Do(ctxt context.Context, v VM) error {
	return v.Stop(ctxt, si.CCID, si.Timeout, si.Dontkill, si.Dontremove)
}

func (si StopContainerReq) GetCCID() ccintf.CCID {
	return si.CCID
}

//Process should be used as follows
//   . construct a context
//   . construct req of the right type (e.g., CreateImageReq)
//   . call it in a go routine
//   . process response in the go routing
//context can be cancelled. VMCProcess will try to cancel calling functions if it can
//For instance docker clients api's such as BuildImage are not cancelable.
//In all cases VMCProcess will wait for the called go routine to return
func (vmc *VMController) Process(ctxt context.Context, vmtype string, req VMCReq) error {
	v := vmc.newVM(vmtype)

	c := make(chan error)
	go func() {
		ccid := req.GetCCID()
		id := ccid.GetName()
		vmc.lockContainer(id)
		err := req.Do(ctxt, v)
		vmc.unlockContainer(id)
		c <- err
	}()

	select {
	case err := <-c:
		return err
	case <-ctxt.Done():
		//TODO cancel req.do ... (needed) ?
		// XXX This logic doesn't make much sense, why return the context error if it's canceled,
		// but still wait for the request to complete, and ignore its error
		<-c
		return ctxt.Err()
	}
}

// GetChaincodePackageBytes creates bytes for docker container generation using the supplied chaincode specification
func GetChaincodePackageBytes(spec *pb.ChaincodeSpec) ([]byte, error) {
	if spec == nil || spec.ChaincodeId == nil {
		return nil, fmt.Errorf("invalid chaincode spec")
	}

	return platforms.GetDeploymentPayload(spec)
}
