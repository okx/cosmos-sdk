package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	tmcrypto "github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"

	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestLedgerErrorHandling(t *testing.T) {
	// first, try to generate a key, must return an error
	// (no panic)
	path := *hd.NewParams(44, 555, 0, false, 0)
	_, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
	require.Error(t, err)
}

func TestPublicKeyUnsafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, priv)

	require.Equal(t, "eb5ae987210216a102323bd08b25b5e8f0cf573c7715a16a5ed907d7298fc21ee3aceef29628",
		fmt.Sprintf("%x", priv.PubKey().Bytes()),
		"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "okexchainpub1addwnpepqgt2zq3j80ggkfd4arcv74euwu26z6j7myraw2v0cg0w8t8w72tzsw0zse0",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	addr := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, "okexchain1sjjdwhadqgt4e4507ps54zwc3zvx9c5hfxtysq",
		addr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)
}

func TestPublicKeyUnsafeHDPath(t *testing.T) {
	expectedAnswers := []string{
		"okexchainpub1addwnpepqgt2zq3j80ggkfd4arcv74euwu26z6j7myraw2v0cg0w8t8w72tzsw0zse0",
		"okexchainpub1addwnpepqdvr3cuk0yvr6nct8zf6l7pxtyxxgtnc8ncxjnjjs4z2j3s2xlffxja4pkf",
		"okexchainpub1addwnpepqwmkd5ng58gt8mpt3j5hekp5jnq2p7a55p7qut5tx9xdkyguqgvgjvxdxjq",
		"okexchainpub1addwnpepq0pk8w5se4p352jm0ga36mrdtxc62rwppthncs9h632xhvxlw8ft6ygdfvm",
		"okexchainpub1addwnpepq2e43ypcapx6d03rctews7qkjtcxkuav58dazkesrwe5nl77gucvspwql4x",
		"okexchainpub1addwnpepqdyarkhe2kh56l5j5kde55qq29mff5u49y0qsthww34ymffsps0jv4q9rnz",
		"okexchainpub1addwnpepqgtg3xaplrf8a79e6j09cxp0hwaq40yrfn90udwz09dwg50asacpceup5r3",
		"okexchainpub1addwnpepq2wvn84fky0chza9zrje6qkg34rr7fwwlan9npln7j5vay6w3uh9qvzctya",
		"okexchainpub1addwnpepqg4u7rz8uksafzdgjxvepssnxtn897ettdpkfqcn26ygh7gsp6fygl3nlzl",
		"okexchainpub1addwnpepqvesztxmxvh8t2hl5q6l0m3vgcd26dnuvsct6fptzt3sruflun255mjuhn3",
	}

	const numIters = 10

	privKeys := make([]tmcrypto.PrivKey, numIters)

	// Check with device
	for i := uint32(0); i < 10; i++ {
		path := *hd.NewFundraiserParams(0, sdk.CoinType, i)
		fmt.Printf("Checking keys at %v\n", path)

		priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
		require.Nil(t, err, "%s", err)
		require.NotNil(t, priv)

		// Check other methods
		require.NoError(t, priv.(PrivKeyLedgerSecp256k1).ValidateKey())
		tmp := priv.(PrivKeyLedgerSecp256k1)
		(&tmp).AssertIsPrivKeyInner()

		pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedAnswers[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func TestPublicKeySafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, addr, err := NewPrivKeyLedgerSecp256k1(path, "cosmos")

	require.Nil(t, err, "%s", err)
	require.NotNil(t, priv)

	require.Equal(t, "eb5ae987210216a102323bd08b25b5e8f0cf573c7715a16a5ed907d7298fc21ee3aceef29628",
		fmt.Sprintf("%x", priv.PubKey().Bytes()),
		"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "okexchainpub1addwnpepqgt2zq3j80ggkfd4arcv74euwu26z6j7myraw2v0cg0w8t8w72tzsw0zse0",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	require.Equal(t, "okexchain1sjjdwhadqgt4e4507ps54zwc3zvx9c5hfxtysq",
		addr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, addr, addr2)
}

func TestPublicKeyHDPath(t *testing.T) {
	expectedPubKeys := []string{
		"okexchainpub1addwnpepqgt2zq3j80ggkfd4arcv74euwu26z6j7myraw2v0cg0w8t8w72tzsw0zse0",
		"okexchainpub1addwnpepqdvr3cuk0yvr6nct8zf6l7pxtyxxgtnc8ncxjnjjs4z2j3s2xlffxja4pkf",
		"okexchainpub1addwnpepqwmkd5ng58gt8mpt3j5hekp5jnq2p7a55p7qut5tx9xdkyguqgvgjvxdxjq",
		"okexchainpub1addwnpepq0pk8w5se4p352jm0ga36mrdtxc62rwppthncs9h632xhvxlw8ft6ygdfvm",
		"okexchainpub1addwnpepq2e43ypcapx6d03rctews7qkjtcxkuav58dazkesrwe5nl77gucvspwql4x",
		"okexchainpub1addwnpepqdyarkhe2kh56l5j5kde55qq29mff5u49y0qsthww34ymffsps0jv4q9rnz",
		"okexchainpub1addwnpepqgtg3xaplrf8a79e6j09cxp0hwaq40yrfn90udwz09dwg50asacpceup5r3",
		"okexchainpub1addwnpepq2wvn84fky0chza9zrje6qkg34rr7fwwlan9npln7j5vay6w3uh9qvzctya",
		"okexchainpub1addwnpepqg4u7rz8uksafzdgjxvepssnxtn897ettdpkfqcn26ygh7gsp6fygl3nlzl",
		"okexchainpub1addwnpepqvesztxmxvh8t2hl5q6l0m3vgcd26dnuvsct6fptzt3sruflun255mjuhn3",
	}

	expectedAddrs := []string{
		"okexchain1sjjdwhadqgt4e4507ps54zwc3zvx9c5hfxtysq",
		"okexchain1kywusqu9t5ppe9nl2rq26tg0cnzwygdctmmz5q",
		"okexchain1w9jw7273g3n8l0fyjtmlxmqhrctcn2afnlv9a2",
		"okexchain1lpjkyd2esun6y3l6qkjt377szaszn7ndte9sl7",
		"okexchain1m4hehxcu4nnq48wh6wummhdqtsnatm2cu8nlvc",
		"okexchain1chkys65xeekky8rjsr492n9ghytqprm6n3a2nt",
		"okexchain166tecxhd52la8w33g4kvpnnhp33wdfy0v6uda7",
		"okexchain1u8nwm4vk3q025azsujfqxwnftzzffjhkm0l7kr",
		"okexchain1m79f2vy0xm756v7evs2dz9umkw4swfn7j3e2a9",
		"okexchain13aex5w5w3hyg6zyvfq367g7alu47nk25msf2pj",
	}

	const numIters = 10

	privKeys := make([]tmcrypto.PrivKey, numIters)

	// Check with device
	for i := uint32(0); i < 10; i++ {
		path := *hd.NewFundraiserParams(0, sdk.CoinType, i)
		fmt.Printf("Checking keys at %v\n", path)

		priv, addr, err := NewPrivKeyLedgerSecp256k1(path, "cosmos")
		require.Nil(t, err, "%s", err)
		require.NotNil(t, addr)
		require.NotNil(t, priv)

		addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
		require.Equal(t, addr2, addr)
		require.Equal(t,
			expectedAddrs[i], addr,
			"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

		// Check other methods
		require.NoError(t, priv.(PrivKeyLedgerSecp256k1).ValidateKey())
		tmp := priv.(PrivKeyLedgerSecp256k1)
		(&tmp).AssertIsPrivKeyInner()

		pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedPubKeys[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func getFakeTx(accountNumber uint32) []byte {
	tmp := fmt.Sprintf(
		`{"account_number":"%d","chain_id":"1234","fee":{"amount":[{"amount":"150","denom":"atom"}],"gas":"5000"},"memo":"memo","msgs":[[""]],"sequence":"6"}`,
		accountNumber)

	return []byte(tmp)
}

func TestSignaturesHD(t *testing.T) {
	for account := uint32(0); account < 100; account += 30 {
		msg := getFakeTx(account)

		path := *hd.NewFundraiserParams(account, sdk.CoinType, account/5)
		fmt.Printf("Checking signature at %v    ---   PLEASE REVIEW AND ACCEPT IN THE DEVICE\n", path)

		priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
		require.Nil(t, err, "%s", err)

		pub := priv.PubKey()
		sig, err := priv.Sign(msg)
		require.Nil(t, err)

		valid := pub.VerifyBytes(msg, sig)
		require.True(t, valid, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)
	}
}

func TestRealLedgerSecp256k1(t *testing.T) {
	msg := getFakeTx(50)
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
	require.Nil(t, err, "%s", err)

	pub := priv.PubKey()
	sig, err := priv.Sign(msg)
	require.Nil(t, err)

	valid := pub.VerifyBytes(msg, sig)
	require.True(t, valid)

	// now, let's serialize the public key and make sure it still works
	bs := priv.PubKey().Bytes()
	pub2, err := cryptoAmino.PubKeyFromBytes(bs)
	require.Nil(t, err, "%+v", err)

	// make sure we get the same pubkey when we load from disk
	require.Equal(t, pub, pub2)

	// signing with the loaded key should match the original pubkey
	sig, err = priv.Sign(msg)
	require.Nil(t, err)
	valid = pub.VerifyBytes(msg, sig)
	require.True(t, valid)

	// make sure pubkeys serialize properly as well
	bs = pub.Bytes()
	bpub, err := cryptoAmino.PubKeyFromBytes(bs)
	require.NoError(t, err)
	require.Equal(t, pub, bpub)
}
