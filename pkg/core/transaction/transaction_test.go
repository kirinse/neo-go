package transaction

import (
	"encoding/hex"
	"testing"

	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/internal/testserdes"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWitnessEncodeDecode(t *testing.T) {
	verif, err := hex.DecodeString("552102486fd15702c4490a26703112a5cc1d0923fd697a33406bd5a1c00e0013b09a7021024c7b7fb6c310fccf1ba33b082519d82964ea93868d676662d4a59ad548df0e7d2102aaec38470f6aad0042c6e877cfd8087d2676b0f516fddd362801b9bd3936399e2103b209fd4f53a7170ea4444e0cb0a6bb6a53c2bd016926989cf85f9b0fba17a70c2103b8d9d5771d8f513aa0869b9cc8d50986403b78c6da36890638c3d46a5adce04a2102ca0e27697b9c248f6f16e085fd0061e26f44da85b58ee835c110caa5ec3ba5542102df48f60e8f3e01c48ff40b9b7f1310d7a8b2a193188befe1c2e3df740e89509357ae")
	assert.Nil(t, err)
	invoc, err := hex.DecodeString("404edf5005771de04619235d5a4c7a9a11bb78e008541f1da7725f654c33380a3c87e2959a025da706d7255cb3a3fa07ebe9c6559d0d9e6213c68049168eb1056f4038a338f879930c8adc168983f60aae6f8542365d844f004976346b70fb0dd31aa1dbd4abd81e4a4aeef9941ecd4e2dd2c1a5b05e1cc74454d0403edaee6d7a4d4099d33c0b889bf6f3e6d87ab1b11140282e9a3265b0b9b918d6020b2c62d5a040c7e0c2c7c1dae3af9b19b178c71552ebd0b596e401c175067c70ea75717c8c00404e0ebd369e81093866fe29406dbf6b402c003774541799d08bf9bb0fc6070ec0f6bad908ab95f05fa64e682b485800b3c12102a8596e6c715ec76f4564d5eff34070e0521979fcd2cbbfa1456d97cc18d9b4a6ad87a97a2a0bcdedbf71b6c9676c645886056821b6f3fec8694894c66f41b762bc4e29e46ad15aee47f05d27d822")
	assert.Nil(t, err)

	lenInvoc := len(invoc)
	lenVerif := len(verif)
	t.Log(lenInvoc)
	t.Log(lenVerif)

	wit := &Witness{
		InvocationScript:   invoc,
		VerificationScript: verif,
	}

	witDecode := &Witness{}
	testserdes.EncodeDecodeBinary(t, wit, witDecode)

	t.Log(len(witDecode.VerificationScript))
	t.Log(len(witDecode.InvocationScript))
}

// TODO NEO3.0: update binary
/*
func TestDecodeEncodeClaimTX(t *testing.T) {
	tx := decodeTransaction(rawClaimTX, t)
	assert.Equal(t, tx.Type, ClaimType)
	assert.IsType(t, tx.Data, &ClaimTX{})
	claimTX := tx.Data.(*ClaimTX)
	assert.Equal(t, 4, len(claimTX.Claims))
	assert.Equal(t, 0, len(tx.Attributes))
	assert.Equal(t, 0, len(tx.Inputs))
	assert.Equal(t, 1, len(tx.Outputs))
	assert.Equal(t, "AQJseD8iBmCD4sgfHRhMahmoi9zvopG6yz", address.Uint160ToString(tx.Outputs[0].ScriptHash))
	assert.Equal(t, "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7", tx.Outputs[0].AssetID.StringLE())
	assert.Equal(t, tx.Outputs[0].Amount.String(), "0.06247739")
	invoc := "40456349cec43053009accdb7781b0799c6b591c812768804ab0a0b56b5eae7a97694227fcd33e70899c075848b2cee8fae733faac6865b484d3f7df8949e2aadb"
	verif := "2103945fae1ed3c31d778f149192b76734fcc951b400ba3598faa81ff92ebe477eacac"
	assert.Equal(t, 1, len(tx.Scripts))
	assert.Equal(t, invoc, hex.EncodeToString(tx.Scripts[0].InvocationScript))
	assert.Equal(t, verif, hex.EncodeToString(tx.Scripts[0].VerificationScript))

	data, err := testserdes.EncodeBinary(tx)
	assert.NoError(t, err)
	assert.Equal(t, rawClaimTX, hex.EncodeToString(data))

	hash := "2c6a45547b3898318e400e541628990a07acb00f3b9a15a8e966ae49525304da"
	assert.Equal(t, hash, tx.hash.StringLE())
}

func TestDecodeEncodeInvocationTX(t *testing.T) {
	tx := decodeTransaction(rawInvocationTX, t)
	assert.Equal(t, tx.Type, InvocationType)
	assert.IsType(t, tx.Data, &InvocationTX{})

	invocTX := tx.Data.(*InvocationTX)
	script := "0400b33f7114839c33710da24cf8e7d536b8d244f3991cf565c8146063795d3b9b3cd55aef026eae992b91063db0db53c1087472616e7366657267c5cc1cb5392019e2cc4e6d6b5ea54c8d4b6d11acf166cb072961424c54f6"
	assert.Equal(t, script, hex.EncodeToString(invocTX.Script))
	assert.Equal(t, util.Fixed8(0), invocTX.Gas)

	assert.Equal(t, 1, len(tx.Attributes))
	assert.Equal(t, 0, len(tx.Inputs))
	assert.Equal(t, 0, len(tx.Outputs))
	invoc := "40c6a131c55ca38995402dff8e92ac55d89cbed4b98dfebbcb01acbc01bd78fa2ce2061be921b8999a9ab79c2958875bccfafe7ce1bbbaf1f56580815ea3a4feed"
	verif := "2102d41ddce2c97be4c9aa571b8a32cbc305aa29afffbcae71b0ef568db0e93929aaac"
	assert.Equal(t, 1, len(tx.Scripts))
	assert.Equal(t, invoc, hex.EncodeToString(tx.Scripts[0].InvocationScript))
	assert.Equal(t, verif, hex.EncodeToString(tx.Scripts[0].VerificationScript))

	data, err := testserdes.EncodeBinary(tx)
	assert.NoError(t, err)
	assert.Equal(t, rawInvocationTX, hex.EncodeToString(data))
}
*/

func TestNewInvocationTX(t *testing.T) {
	script := []byte{0x51}
	tx := NewInvocationTX(script, 1)
	txData := tx.Data.(*InvocationTX)
	assert.Equal(t, InvocationType, tx.Type)
	assert.Equal(t, tx.Version, txData.Version)
	assert.Equal(t, script, txData.Script)
	// Update hash fields to match tx2 that is gonna autoupdate them on decode.
	_ = tx.Hash()
	testserdes.EncodeDecodeBinary(t, tx, new(Transaction))
}

func TestEncodingTXWithNoData(t *testing.T) {
	_, err := testserdes.EncodeBinary(new(Transaction))
	require.Error(t, err)
}

func TestMarshalUnmarshalJSONContractTX(t *testing.T) {
	tx := NewContractTX()
	tx.Outputs = []Output{{
		AssetID:    util.Uint256{1, 2, 3, 4},
		Amount:     567,
		ScriptHash: util.Uint160{7, 8, 9, 10},
		Position:   13,
	}}
	tx.Scripts = []Witness{{
		InvocationScript:   []byte{5, 3, 1},
		VerificationScript: []byte{2, 4, 6},
	}}
	tx.Data = &ContractTX{}
	testserdes.MarshalUnmarshalJSON(t, tx, new(Transaction))
}

func TestMarshalUnmarshalJSONClaimTX(t *testing.T) {
	tx := &Transaction{
		Type:    ClaimType,
		Version: 0,
		Data: &ClaimTX{Claims: []Input{
			{
				PrevHash:  util.Uint256{1, 2, 3, 4},
				PrevIndex: uint16(56),
			},
		}},
		Attributes: []Attribute{},
		Inputs: []Input{{
			PrevHash:  util.Uint256{5, 6, 7, 8},
			PrevIndex: uint16(12),
		}},
		Outputs: []Output{{
			AssetID:    util.Uint256{1, 2, 3},
			Amount:     util.Fixed8FromInt64(1),
			ScriptHash: util.Uint160{1, 2, 3},
			Position:   0,
		}},
		Scripts: []Witness{},
		Trimmed: false,
	}

	testserdes.MarshalUnmarshalJSON(t, tx, new(Transaction))
}

func TestMarshalUnmarshalJSONInvocationTX(t *testing.T) {
	tx := &Transaction{
		Type:    InvocationType,
		Version: 3,
		Data: &InvocationTX{
			Script:  []byte{1, 2, 3, 4},
			Gas:     util.Fixed8FromFloat(100),
			Version: 3,
		},
		Attributes: []Attribute{},
		Inputs: []Input{{
			PrevHash:  util.Uint256{5, 6, 7, 8},
			PrevIndex: uint16(12),
		}},
		Outputs: []Output{{
			AssetID:    util.Uint256{1, 2, 3},
			Amount:     util.Fixed8FromInt64(1),
			ScriptHash: util.Uint160{1, 2, 3},
			Position:   0,
		}},
		Scripts: []Witness{},
		Trimmed: false,
	}

	testserdes.MarshalUnmarshalJSON(t, tx, new(Transaction))
}

func TestMarshalUnmarshalJSONRegisterTX(t *testing.T) {
	tx := &Transaction{
		Type:    RegisterType,
		Version: 5,
		Data: &RegisterTX{
			AssetType: 0,
			Name:      `[{"lang":"zh-CN","name":"小蚁股"},{"lang":"en","name":"AntShare"}]`,
			Amount:    1000000,
			Precision: 0,
			Owner:     keys.PublicKey{},
			Admin:     util.Uint160{},
		},
		Attributes: []Attribute{},
		Inputs: []Input{{
			PrevHash:  util.Uint256{5, 6, 7, 8},
			PrevIndex: uint16(12),
		}},
		Outputs: []Output{{
			AssetID:    util.Uint256{1, 2, 3},
			Amount:     util.Fixed8FromInt64(1),
			ScriptHash: util.Uint160{1, 2, 3},
			Position:   0,
		}},
		Scripts: []Witness{},
		Trimmed: false,
	}

	testserdes.MarshalUnmarshalJSON(t, tx, new(Transaction))
}
