package types_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/types"
)

var invalidStrs = []string{
	"hello, world!",
	"AAA",
	types.Bech32PrefixAccAddr + "AB0C",
	types.Bech32PrefixAccPub + "1234",
	types.Bech32PrefixValAddr + "5678",
	types.Bech32PrefixValPub + "BBAB",
	types.Bech32PrefixConsAddr + "FF04",
	types.Bech32PrefixConsPub + "6789",
}

func testMarshal(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}

func TestEmptyAddresses(t *testing.T) {
	require.Equal(t, (types.AccAddress{}).String(), "")
	require.Equal(t, (types.ValAddress{}).String(), "")
	require.Equal(t, (types.ConsAddress{}).String(), "")

	accAddr, err := types.AccAddressFromBech32("")
	require.True(t, accAddr.Empty())
	require.Nil(t, err)

	valAddr, err := types.ValAddressFromBech32("")
	require.True(t, valAddr.Empty())
	require.Nil(t, err)

	consAddr, err := types.ConsAddressFromBech32("")
	require.True(t, consAddr.Empty())
	require.Nil(t, err)
}

func TestRandBech32PubkeyConsistency(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 1000; i++ {
		rand.Read(pub[:])

		mustBech32AccPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pub)
		bech32AccPub, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pub)
		require.Nil(t, err)
		require.Equal(t, bech32AccPub, mustBech32AccPub)

		mustBech32ValPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeValPub, pub)
		bech32ValPub, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeValPub, pub)
		require.Nil(t, err)
		require.Equal(t, bech32ValPub, mustBech32ValPub)

		mustBech32ConsPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeConsPub, pub)
		bech32ConsPub, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeConsPub, pub)
		require.Nil(t, err)
		require.Equal(t, bech32ConsPub, mustBech32ConsPub)

		mustAccPub := types.MustGetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, bech32AccPub)
		accPub, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, bech32AccPub)
		require.Nil(t, err)
		require.Equal(t, accPub, mustAccPub)

		mustValPub := types.MustGetPubKeyFromBech32(types.Bech32PubKeyTypeValPub, bech32ValPub)
		valPub, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeValPub, bech32ValPub)
		require.Nil(t, err)
		require.Equal(t, valPub, mustValPub)

		mustConsPub := types.MustGetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, bech32ConsPub)
		consPub, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, bech32ConsPub)
		require.Nil(t, err)
		require.Equal(t, consPub, mustConsPub)

		require.Equal(t, valPub, accPub)
		require.Equal(t, valPub, consPub)
	}
}

func TestYAMLMarshalers(t *testing.T) {
	addr := secp256k1.GenPrivKey().PubKey().Address()

	acc := types.AccAddress(addr)
	val := types.ValAddress(addr)
	cons := types.ConsAddress(addr)

	got, _ := yaml.Marshal(&acc)
	require.Equal(t, acc.String()+"\n", string(got))

	got, _ = yaml.Marshal(&val)
	require.Equal(t, val.String()+"\n", string(got))

	got, _ = yaml.Marshal(&cons)
	require.Equal(t, cons.String()+"\n", string(got))
}

func TestRandBech32AccAddrConsistency(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 1000; i++ {
		rand.Read(pub[:])

		acc := types.AccAddress(pub.Address())
		res := types.AccAddress{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.AccAddressFromBech32(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.AccAddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.AccAddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.AccAddressFromBech32(str)
		require.NotNil(t, err)

		err = (*types.AccAddress)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

func TestValAddr(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 20; i++ {
		rand.Read(pub[:])

		acc := types.ValAddress(pub.Address())
		res := types.ValAddress{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.ValAddressFromBech32(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.ValAddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.ValAddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.ValAddressFromBech32(str)
		require.NotNil(t, err)

		err = (*types.ValAddress)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

func TestConsAddress(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 20; i++ {
		rand.Read(pub[:])

		acc := types.ConsAddress(pub.Address())
		res := types.ConsAddress{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.ConsAddressFromBech32(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.ConsAddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.ConsAddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.ConsAddressFromBech32(str)
		require.NotNil(t, err)

		err = (*types.ConsAddress)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestConfiguredPrefix(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	for length := 1; length < 10; length++ {
		for times := 1; times < 20; times++ {
			rand.Read(pub[:])
			// Test if randomly generated prefix of a given length works
			prefix := RandString(length)

			// Assuming that GetConfig is not sealed.
			config := types.GetConfig()
			config.SetBech32PrefixForAccount(
				prefix+types.PrefixAccount,
				prefix+types.PrefixPublic)

			acc := types.AccAddress(pub.Address())
			require.True(t, strings.HasPrefix(
				acc.String(),
				prefix+types.PrefixAccount), acc.String())

			bech32Pub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pub)
			require.True(t, strings.HasPrefix(
				bech32Pub,
				prefix+types.PrefixPublic))

			config.SetBech32PrefixForValidator(
				prefix+types.PrefixValidator+types.PrefixAddress,
				prefix+types.PrefixValidator+types.PrefixPublic)

			val := types.ValAddress(pub.Address())
			require.True(t, strings.HasPrefix(
				val.String(),
				prefix+types.PrefixValidator+types.PrefixAddress))

			bech32ValPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeValPub, pub)
			require.True(t, strings.HasPrefix(
				bech32ValPub,
				prefix+types.PrefixValidator+types.PrefixPublic))

			config.SetBech32PrefixForConsensusNode(
				prefix+types.PrefixConsensus+types.PrefixAddress,
				prefix+types.PrefixConsensus+types.PrefixPublic)

			cons := types.ConsAddress(pub.Address())
			require.True(t, strings.HasPrefix(
				cons.String(),
				prefix+types.PrefixConsensus+types.PrefixAddress))

			bech32ConsPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeConsPub, pub)
			require.True(t, strings.HasPrefix(
				bech32ConsPub,
				prefix+types.PrefixConsensus+types.PrefixPublic))
		}
	}
}

func TestAddressInterface(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	addrs := []types.Address{
		types.ConsAddress(pub.Address()),
		types.ValAddress(pub.Address()),
		types.AccAddress(pub.Address()),
	}

	for _, addr := range addrs {
		switch addr := addr.(type) {
		case types.AccAddress:
			_, err := types.AccAddressFromBech32(addr.String())
			require.Nil(t, err)
		case types.ValAddress:
			_, err := types.ValAddressFromBech32(addr.String())
			require.Nil(t, err)
		case types.ConsAddress:
			_, err := types.ConsAddressFromBech32(addr.String())
			require.Nil(t, err)
		default:
			t.Fail()
		}
	}

}

func TestCustomAddressVerifier(t *testing.T) {
	// Create a 10 byte address
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	accBech := types.AccAddress(addr).String()
	valBech := types.ValAddress(addr).String()
	consBech := types.ConsAddress(addr).String()
	// Verifiy that the default logic rejects this 10 byte address
	err := types.VerifyAddressFormat(addr)
	require.NotNil(t, err)
	_, err = types.AccAddressFromBech32(accBech)
	require.NotNil(t, err)
	_, err = types.ValAddressFromBech32(valBech)
	require.NotNil(t, err)
	_, err = types.ConsAddressFromBech32(consBech)
	require.NotNil(t, err)

	// Set a custom address verifier that accepts 10 or 20 byte addresses
	types.GetConfig().SetAddressVerifier(func(bz []byte) error {
		n := len(bz)
		if n == 10 || n == types.AddrLen {
			return nil
		}
		return fmt.Errorf("incorrect address length %d", n)
	})

	// Verifiy that the custom logic accepts this 10 byte address
	err = types.VerifyAddressFormat(addr)
	require.Nil(t, err)
	_, err = types.AccAddressFromBech32(accBech)
	require.Nil(t, err)
	_, err = types.ValAddressFromBech32(valBech)
	require.Nil(t, err)
	_, err = types.ConsAddressFromBech32(consBech)
	require.Nil(t, err)
}

var (
	first = `ok.scf.debug 018f8526dbfd6e38e3d8307702ca8469bae6c56c15 4c96dfce0a3c0a148f8526dbfd6e38e3d8307702ca8469bae6c56c15121f0a036f6b741218353731373438313035323031353132333135363534383535208a70280112207b1ba3c740ec72c185154730cb2106a272097785a74a1f28cc127e44194fde4c
ok.scf.debug 01d9a2e75b371f01d629b619dba9d4c890a7963a4a 4c96dfce0a620a14d9a2e75b371f01d629b619dba9d4c890a7963a4a121b0a036f6b74121431323736363137363135373939393935343636361a26f3b3cd032102b52317ef3392acb4ec9fc2b304002dda9a638478f424706e563b2187d0e40a8420eb2028e49c081220c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470
ok.scf.debug 01f1829676db577682e944fc3493d451b67ff3e29f 0bfd7d9a0a300a14f1829676db577682e944fc3493d451b67ff3e29f12160a036f6b74120f3936393432333030303132343238352042120d6665655f636f6c6c6563746f72
ok.scf.debug 01f193f0bb30a4ba945dd8d23fada8f2f32c45eb9a 4c96dfce0a1c0a14f193f0bb30a4ba945dd8d23fada8f2f32c45eb9a20bcd15c28011220af3d11c57c8a1bf1e2edd8ad108c7803515b08a6445468c71a28b97c9d5656da
ok.scf.debug 01df54b6c6195ea4d948d03bfd818d365cf175cfc2 4c96dfce0a1b0a14df54b6c6195ea4d948d03bfd818d365cf175cfc2209d66280112209f15a21bb6644b381aeea4b081d751308c1a0e8234a7a190ff8168a0effd80d3
ok.scf.debug 05f193f0bb30a4ba945dd8d23fada8f2f32c45eb9a277837436ce9384d5ac2563360db8d9f710e96e89d9758ffe719f263deeada93 61516c5b00000000535b6ebd6590366057c500000000152c0254a82bda1f0a9a
ok.scf.debug 05df54b6c6195ea4d948d03bfd818d365cf175cfc2556d6e2c4b9da71a601dd40b0cff622b44d811d0d60becc545886f83d7119b8c ffffffffffffffffffffffffffffffffffffffffffffbdccb06707d8dc7b2c1d
ok.scf.debug 05f193f0bb30a4ba945dd8d23fada8f2f32c45eb9a952507c2efcdacdd6b303d9fe5e44a234d96d16291599e04e9d5551b2a83af9f 0000000000000000000000000000000dd037701cb79cd3e9396e0946df997ff9
ok.scf.debug 05f193f0bb30a4ba945dd8d23fada8f2f32c45eb9af0b0c5d71d5a660b0091970f85dbd0a76223ae7151d6ef240aff74209c0981a0 000000000000000000000000000000f213d9355d23d57e0c4c68d50efa096d98
ok.scf.debug 058f8526dbfd6e38e3d8307702ca8469bae6c56c153ce253bab56f3082f5476eae8d5cb106f7dcba3005b018aa4e71696747cdd1cd 0000000000000000000000000000000000000000000000c1d8f1843f9a184967
ok.scf.debug 058f8526dbfd6e38e3d8307702ca8469bae6c56c159ff08bde739d18e834d71dedc28eddec9e32378f80ca2bb9fe471f5e50e2f2bf 00000000000000000000000000000000000000000000152c0254a82bda1f0a9a
ok.scf.debug 05df54b6c6195ea4d948d03bfd818d365cf175cfc24156aede4259e94fd8c9ad42c8e53a64ce2c0b27a0390710c0af8982055d3a11 00000000000000000000000000000000000000000000535b6ebd6590366057c5
ok.scf.debug 05df54b6c6195ea4d948d03bfd818d365cf175cfc29bd25dd7057a5e8ca5fb476e41c428be65c5f89517855e9a3492630a12b177f2 0000000000000000000000000000000000000000000000721b4f014af4546f36
`
)

func TestAsd(t *testing.T) {
	fmt.Println("first", first)

}
