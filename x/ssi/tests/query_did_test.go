package tests

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hypersign-protocol/hid-node/x/ssi/keeper"
	"github.com/hypersign-protocol/hid-node/x/ssi/types"
	"github.com/stretchr/testify/assert"
)

func TestDidResolve(t *testing.T) {
	t.Log("Running test for DidResolve (Query)")
	k, ctx := TestKeeper(t)
	msgServer := keeper.NewMsgServerImpl(*k)
	goCtx := sdk.WrapSDKContext(ctx)

	k.SetDidNamespace(&ctx, "devnet")

	keyPair1 := GeneratePublicPrivateKeyPair()
	rpcElements := GenerateDidDocumentRPCElements(keyPair1)
	didId := rpcElements.DidDocument.GetId()
	t.Logf("Registering DID with DID Id: %s", didId)

	msgCreateDID := &types.MsgCreateDID{
		DidDocString: rpcElements.DidDocument,
		Signatures:   rpcElements.Signatures,
		Creator:      rpcElements.Creator,
	}

	_, errCreateDID := msgServer.CreateDID(goCtx, msgCreateDID)
	if errCreateDID != nil {
		t.Error("DID Registeration Failed")
		t.Log(rpcElements.DidDocument.Id)
		t.Error(errCreateDID)
		t.FailNow()
	}

	t.Log("Did Registeration Successful")
	t.Log("Querying the DID from store")

	req := &types.QueryGetDidDocByIdRequest{
		DidId: didId,
	}

	res, errResponse := k.ResolveDid(goCtx, req)
	if errResponse != nil {
		t.Error("Did Resolve Failed")
		t.Error(errResponse)
		t.FailNow()
	}
	t.Log("Querying successful")
	// To check if queried Did Document is not nil
	assert.NotNil(t, res.DidDocument)
	t.Log("Did Resolve Test Completed")
}

func TestDidParam(t *testing.T) {
	t.Log("Running test for DidParam (Query)")
	k, ctx := TestKeeper(t)
	msgServer := keeper.NewMsgServerImpl(*k)
	goCtx := sdk.WrapSDKContext(ctx)

	k.SetDidNamespace(&ctx, "devnet")
	
	keyPair1 := GeneratePublicPrivateKeyPair()
	rpcElements := GenerateDidDocumentRPCElements(keyPair1)
	didId := rpcElements.DidDocument.GetId()
	t.Logf("Registering DID with DID Id: %s", didId)

	msgCreateDID := &types.MsgCreateDID{
		DidDocString: rpcElements.DidDocument,
		Signatures:   rpcElements.Signatures,
		Creator:      rpcElements.Creator,
	}

	_, errCreateDID := msgServer.CreateDID(goCtx, msgCreateDID)
	if errCreateDID != nil {
		t.Error("DID Registeration Failed")
		t.Log(rpcElements.DidDocument.Id)
		t.Error(errCreateDID)
		t.FailNow()
	}
	
	t.Log("Did Registeration Successful")
	t.Log("Querying the list of Did Documents")

	req := &types.QueryDidParamRequest{}

	res, errResponse := k.DidParam(goCtx, req)
	if errResponse != nil {
		t.Error("Did Resolve Failed")
		t.Error(errResponse)
		t.FailNow()
	}

	t.Log("Querying successful")

	// Did Document Count should't be zero 
	assert.NotEqual(t, "0", res.TotalDidCount)
	// List should be populated with a single Did Document
	assert.Equal(t, 1, len(res.DidDocList))
	// Did Document shouldnt be nil
	assert.NotNil(t, res.DidDocList[0].DidDocument)

	t.Log("Did Param Test Completed")
}