package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// verify interface at compile time
var _ sdk.Msg = &MsgRevealSolution{}

// MsgRevealSolution - struct for unjailing jailed validator
type MsgRevealSolution struct {
	Scavenger    sdk.AccAddress `json:"scavenger" yaml:"scavenger"`
	SolutionHash string         `json:"solutionHash" yaml:"solutionHash"`
	Solution     string         `json:"solution" yaml:"solution"`
}

func NewMsgRevealSolution(scavenger sdk.AccAddress, solutionHash string, solution string) MsgRevealSolution {
	return MsgRevealSolution{Scavenger: scavenger, SolutionHash: solutionHash, Solution: solution}
}

const RevealSolutionConst = "RevealSolution"

// nolint
func (msg MsgRevealSolution) Route() string { return RouterKey }
func (msg MsgRevealSolution) Type() string  { return RevealSolutionConst }
func (msg MsgRevealSolution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Scavenger)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRevealSolution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRevealSolution) ValidateBasic() error {
	if msg.Scavenger.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "scavenger can't be empty")
	}
	if msg.SolutionHash == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "solutionHash can't be empty")
	}
	if msg.Solution == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "solution can't be empty")
	}

	solutionHash := sha256.Sum256([]byte(msg.Solution))
	solutionHashString := hex.EncodeToString(solutionHash[:])
	if msg.SolutionHash != solutionHashString {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf(
				"Hash of solution (%s) doesn't equal solutionHash (%s)",
				msg.SolutionHash,
				solutionHashString,
			),
		)
	}
	return nil
}
