package tests

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/hypersign-protocol/hid-node/testutil/keeper"
	"github.com/hypersign-protocol/hid-node/x/ssi/keeper"
	"github.com/hypersign-protocol/hid-node/x/ssi/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateDID(t *testing.T) {
	t.Log("Running test for Valid Update DID Tx")
	k, ctx := testkeeper.SsiKeeper(t)
	msgServ := keeper.NewMsgServerImpl(*k)
	goCtx := sdk.WrapSDKContext(ctx)

	keyPair1 := GeneratePublicPrivateKeyPair()
	rpcElements := GenerateDidDocumentRPCElements(keyPair1)
	t.Logf("Registering DID with DID Id: %s", rpcElements.DidDocument.GetId())

	msgCreateDID := &types.MsgCreateDID{
		DidDocString: rpcElements.DidDocument,
		Signatures:   rpcElements.Signatures,
		Creator:      rpcElements.Creator,
	}

	_, err := msgServ.CreateDID(goCtx, msgCreateDID)
	if err != nil {
		t.Error("DID Registeration Failed")
		t.Error(err)
		t.FailNow()
	} else {
		t.Log("Did Registeration Successful")
	}

	// Querying registered did document to get the version ID
	resolvedDidDocument, errResolve := k.GetDid(&ctx, rpcElements.DidDocument.GetId())
	if errResolve != nil {
		t.Error("Error in retrieving registered did document")
		t.Error(errResolve)
		t.FailNow()
	}
	versionId := resolvedDidDocument.GetMetadata().GetVersionId()

	// Updated the existing DID by appending a link
	resolvedDidDocument.Did.Context = append(resolvedDidDocument.Did.Context, "http://www.example.com")

	updateRpcElements := GetModifiedDidDocumentSignature(
		resolvedDidDocument.Did, 
		keyPair1,
		resolvedDidDocument.Did.VerificationMethod[0].Id,
	)
	t.Logf("Updating context field of the registered did with Id: %s", updateRpcElements.DidDocument.Id)

	msgUpdateDID := &types.MsgUpdateDID{
		DidDocString: updateRpcElements.DidDocument,
		Signatures:   updateRpcElements.Signatures,
		VersionId:    versionId,
		Creator:      updateRpcElements.Creator,
	}

	_, errUpdateDID := msgServ.UpdateDID(goCtx, msgUpdateDID)
	if errUpdateDID != nil {
		t.Error("DID Update Failed")
		t.Error(errUpdateDID)
		t.FailNow()
	}
	t.Log("Did Update Successful")

	t.Log("Update DID Tx Test Completed")
}


func TestInvalidChangeByNonControllerDid(t *testing.T) {
	t.Log("Running test for DID Controller (Non-controller DID attempting to make changes)")
	k, ctx := testkeeper.SsiKeeper(t)
	msgServ := keeper.NewMsgServerImpl(*k)
	goCtx := sdk.WrapSDKContext(ctx)

	// Create Two DID: Alice and Charlie
	aliceKeyPair := GeneratePublicPrivateKeyPair()
	aliceDidId, err := CreateDidTx(msgServ, goCtx, aliceKeyPair)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Hello")
	aliceDidDocument := QueryDid(k, ctx, aliceDidId)

	charlieKeyPair := GeneratePublicPrivateKeyPair()
	charlieDidId, err := CreateDidTx(msgServ, goCtx, charlieKeyPair)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	charlieDidDocument := QueryDid(k, ctx, charlieDidId)

	// Charlie will attempt to make changes on Alice's DID
	newContextElement := "somerandomwebiste.earth"
	aliceDidDocument.Did.Context = append(aliceDidDocument.Did.Context, newContextElement)
	versionId := aliceDidDocument.Metadata.VersionId
	updatedRpcElements := GetModifiedDidDocumentSignature(
		aliceDidDocument.Did, 
		charlieKeyPair,
		charlieDidDocument.Did.VerificationMethod[0].Id,
	)
	t.Log("Right here")
	errDidUpdate := UpdateDidTx(msgServ, goCtx, updatedRpcElements, versionId)
	if errDidUpdate == nil {
		t.Error("The test was expected to fail, as charlie's key pair isn't part of alice's control group")
		t.FailNow()
	}

	require.Contains(t, errDidUpdate.Error(), "invalid signature detected")
	t.Log("Cold here")
	t.Log("Test Completed")
}

func TestValidChangeByControllerDid(t *testing.T) {
	t.Log("Running test for DID Controller (Controller DID attempting to make changes)")
	k, ctx := testkeeper.SsiKeeper(t)
	msgServ := keeper.NewMsgServerImpl(*k)
	goCtx := sdk.WrapSDKContext(ctx)

	// Create Two DID: Alice and Charlie
	aliceKeyPair := GeneratePublicPrivateKeyPair()
	aliceDidId, err := CreateDidTx(msgServ, goCtx, aliceKeyPair)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	aliceDidDocument := QueryDid(k, ctx, aliceDidId)

	charlieKeyPair := GeneratePublicPrivateKeyPair()
	charlieDidId, err := CreateDidTx(msgServ, goCtx, charlieKeyPair)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	charlieDidDocument := QueryDid(k, ctx, charlieDidId)

	// Alice adds Charlie to her controller group and signs this transaction her public key
	aliceDidDocument.Did.Controller = append(aliceDidDocument.Did.Controller, charlieDidId)
	versionId := aliceDidDocument.Metadata.VersionId
	updatedRpcElements := GetModifiedDidDocumentSignature(
		aliceDidDocument.Did, 
		aliceKeyPair,
		aliceDidDocument.Did.VerificationMethod[0].Id,
	)

	t.Log("Adding Charlie to Alice's DID Controller Group")
	errDidUpdate := UpdateDidTx(msgServ, goCtx, updatedRpcElements, versionId)
	if errDidUpdate != nil {
		t.Error("Unable to add Charlie to Alice's DID Controller Group")
		t.Error(errDidUpdate)
		t.FailNow()
	}
	t.Log("Charlie has been added to the Controller Group")

	// Charlie attempting to make changes in Alice's Did Document
	t.Log("Charlie will now attempt to make changes in Alice's DID Document")
	aliceDidDocument = QueryDid(k, ctx, aliceDidId)
	newContextElement := "somerandomwebiste.earth"
	aliceDidDocument.Did.Context = append(aliceDidDocument.Did.Context, newContextElement)
	versionId = aliceDidDocument.Metadata.VersionId
	updatedRpcElements = GetModifiedDidDocumentSignature(
		aliceDidDocument.Did, 
		charlieKeyPair,
		charlieDidDocument.Did.VerificationMethod[0].Id,
	)

	err = UpdateDidTx(msgServ, goCtx, updatedRpcElements, versionId)
	if err != nil {
		t.Error("Charlie failed to make changes to Alice's Did Document")
		t.Error(err)
		t.FailNow()
	}

	aliceDidDocument = QueryDid(k, ctx, aliceDidId)
	require.Equal(
		t, 
		aliceDidDocument.Did.Context[len(aliceDidDocument.Did.Context) - 1], 
		newContextElement,
	)

	t.Log("Charlie has succeeded in making the changes.")
	t.Log("Test Completed")

}