package keeper

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/hypersign-protocol/hid-node/x/ssi/types"
	"github.com/hypersign-protocol/hid-node/x/ssi/utils"
)

// Ref 1: The current implementatition takes in the verification key and checks if that belongs to EITHER
// of the DID controllers. If so, then the signature is valid, which is an approach as opposed to the earlier
// implementation where all the signatures of all DIDs present in DID controller were expected. This needs to be verified.
// Link to DID Controller Spec: https://www.w3.org/TR/did-core/#did-controller
func VerifyIdentitySignature(signer types.Signer, signatures []*types.SignInfo, signingInput []byte) (bool, error) {
	result := false
	// foundOne := false

	for _, info := range signatures {
		did, _ := utils.SplitDidUrlIntoDid(info.VerificationMethodId)
		if did == signer.Signer {
			pubKey, err := utils.FindPublicKey(signer, info.VerificationMethodId)
			if err != nil {
				return false, err
			}

			signature, err := base64.StdEncoding.DecodeString(info.Signature)
			if err != nil {
				return false, err
			}

			result = ed25519.Verify(pubKey, signingInput, signature)
			// foundOne = true
		}
	}

	// if !foundOne {
	// 	return false, fmt.Errorf("signature %s not found", signer.Signer)
	// }

	return result, nil
}

func (k msgServer) VerifySignatureOnDidUpdate(ctx *sdk.Context, oldDIDDoc *types.Did, newDIDDoc *types.Did, signatures []*types.SignInfo) error {
	var signers []types.Signer

	oldController := oldDIDDoc.Controller
	if len(oldController) == 0 {
		oldController = []string{oldDIDDoc.Id}
	}

	for _, controller := range oldController {
		signers = append(signers, types.Signer{Signer: controller})
	}

	for _, oldVM := range oldDIDDoc.VerificationMethod {
		newVM := utils.FindVerificationMethod(newDIDDoc.VerificationMethod, oldVM.Id)

		// Verification Method has been deleted
		if newVM == nil {
			signers = AppendSignerIfNeed(signers, oldVM.Controller, newDIDDoc)
			continue
		}

		// Verification Method has been changed
		if !reflect.DeepEqual(oldVM, newVM) {
			signers = AppendSignerIfNeed(signers, newVM.Controller, newDIDDoc)
		}

		// Verification Method Controller has been changed, need to add old controller
		if newVM.Controller != oldVM.Controller {
			signers = AppendSignerIfNeed(signers, oldVM.Controller, newDIDDoc)
		}
	}

	if err := k.VerifySignature(ctx, newDIDDoc, signers, signatures); err != nil {
		return err
	}

	return nil
}

func AppendSignerIfNeed(signers []types.Signer, controller string, msg *types.Did) []types.Signer {
	for _, signer := range signers {
		if signer.Signer == controller {
			return signers
		}
	}

	signer := types.Signer{
		Signer: controller,
	}

	if controller == msg.Id {
		signer.VerificationMethod = msg.VerificationMethod
		signer.Authentication = msg.Authentication
	}

	return append(signers, signer)
}

// Ref 1: The current implementatition takes in the verification key and checks if that belongs to EITHER
// of the DID controllers. If so, then the signature is valid, which is an approach as opposed to the earlier
// implementation where all the signatures of all DIDs present in DID controller were expected. This needs to be verified.
// Link to DID Controller Spec: https://www.w3.org/TR/did-core/#did-controller
func (k *Keeper) VerifySignature(ctx *sdk.Context, msg *types.Did, signers []types.Signer, signatures []*types.SignInfo) error {
	var validArr []types.ValidDid

	if len(signers) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signer should be present")
	}

	if len(signatures) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signature should be present")
	}

	signingInput := msg.GetSignBytes()

	for _, signer := range signers {
		if signer.VerificationMethod == nil {
			didDoc, err := k.GetDid(ctx, signer.Signer)
			if err != nil {
				return types.ErrDidDocNotFound.Wrap(signer.Signer)
			}

			signer.Authentication = didDoc.Did.Authentication
			signer.VerificationMethod = didDoc.Did.VerificationMethod
		}

		valid, err := VerifyIdentitySignature(signer, signatures, signingInput)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
		}

		validArr = append(validArr, types.ValidDid{Did: signer.Signer, IsValid: valid})
	}

	didFoundTrue := contains(validArr)

	if didFoundTrue == (types.ValidDid{}) {
		return sdkerrors.Wrap(types.ErrInvalidSignature, didFoundTrue.Did)
	}

	return nil
}

// TODO: Look for a better way to do this
func contains(s []types.ValidDid) types.ValidDid {
	for _, v := range s {
		if v.IsValid {
			return v
		}
	}
	return types.ValidDid{}
}

func (k *Keeper) VerifySignatureOnCreateSchema(ctx *sdk.Context, msg *types.Schema, signers []types.Signer, signatures []*types.SignInfo) error {
	if len(signers) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signer should be present")
	}

	if len(signatures) == 0 {
		return types.ErrInvalidSignature.Wrap("At least one signature should be present")
	}

	signingInput := msg.GetSignBytes()

	for _, signer := range signers {
		// TODO: Uncomment when Schema is being implemented properly
		// if signer.PublicKeyStruct == nil {
		// 	state, err := k.GetDid(ctx, signer.Signer)
		// 	if err != nil {
		// 		return types.ErrDidDocNotFound.Wrap(signer.Signer)
		// 	}

		// 	didDoc, err := state.UnpackDataAsDid()
		// 	if err != nil {
		// 		return types.ErrDidDocNotFound.Wrap(signer.Signer)
		// 	}

		// 	signer.Authentication = didDoc.Authentication
		// 	signer.PublicKeyStruct = didDoc.PublicKeyStruct
		// }

		valid, err := VerifyIdentitySignature(signer, signatures, signingInput)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
		}

		if !valid {
			// return sdkerrors.Wrap(types.ErrInvalidSignature, signer.Signer)
			return sdkerrors.Wrap(types.ErrInvalidSignature, string(signingInput))
		}
	}

	return nil
}

func (k *Keeper) ValidateController(ctx *sdk.Context, id string, controller string) error {
	if id == controller {
		return nil
	}
	didDoc, err := k.GetDid(ctx, controller)
	if err != nil {
		return types.ErrDidDocNotFound.Wrap(controller)
	}
	if len(didDoc.Did.Authentication) == 0 {
		return types.ErrBadRequestInvalidVerMethod.Wrap(
			fmt.Sprintf("Verificatition method controller %s doesn't have an authentication keys", controller))
	}
	return nil
}

func (k msgServer) ValidateDidControllers(ctx *sdk.Context, id string, controllers []string, verMethods []*types.VerificationMethod) error {

	for _, verificationMethod := range verMethods {
		if err := k.ValidateController(ctx, id, verificationMethod.Controller); err != nil {
			return err
		}
	}

	for _, didController := range controllers {
		if err := k.ValidateController(ctx, id, didController); err != nil {
			return err
		}
	}
	return nil
}

// Check the Deactivate status of DID
func VerifyDidDeactivate(metadata *types.Metadata, id string) error {
	if metadata.Deactivated {
		return sdkerrors.Wrap(types.ErrDidDocDeactivated, fmt.Sprintf("DidDoc ID: %s", id))
	}
	return nil
}

// Verify Credential Signature
func (k msgServer) VerifyCredentialSignature(msg *types.CredentialStatus, didDoc *types.Did, signature string, verificationMethod string) error {
	signingInput := msg.GetSignBytes()

	signer := types.Signer{
		Signer:             didDoc.GetId(),
		Authentication:     didDoc.GetAuthentication(),
		VerificationMethod: didDoc.GetVerificationMethod(),
	}

	signingInfo := &types.SignInfo{
		VerificationMethodId: verificationMethod,
		Signature:            signature,
	}

	signingInfoList := []*types.SignInfo{
		signingInfo,
	}

	valid, err := VerifyIdentitySignature(signer, signingInfoList, signingInput)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
	}

	if !valid {
		// return sdkerrors.Wrap(types.ErrInvalidSignature, signer.Signer)
		return sdkerrors.Wrap(types.ErrInvalidSignature, string(signer.Signer))
	}
	return nil
}
