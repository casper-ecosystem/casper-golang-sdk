package sdk

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/Simplewallethq/casper-golang-sdk/keypair"
	"github.com/Simplewallethq/casper-golang-sdk/keypair/ed25519"
	"github.com/Simplewallethq/casper-golang-sdk/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/blake2b"
)

var source, dest *keypair.PublicKey

var sourceKeyPair keypair.KeyPair

func init() {
	decodedSource, _ := hex.DecodeString("d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061")

	source = &keypair.PublicKey{
		Tag:        keypair.KeyTagEd25519,
		PubKeyData: decodedSource,
	}

	decodedDest, _ := hex.DecodeString("272a2fe949347aa893fdcbb99bfeb4c57e348c5359a45363514c4e15364e5136")

	dest = &keypair.PublicKey{
		Tag:        keypair.KeyTagEd25519,
		PubKeyData: decodedDest,
	}

	sourceKeyPair, _ = ed25519.ParseKeyFiles("../keypair/test_account_keys/account1/public_key.pem", "../keypair/test_account_keys/account1/secret_key.pem")
}

func TestTimeMarshaling(t *testing.T) {
	parse, err := time.Parse("2006-01-02T15:04:05.999Z", "2021-09-13T17:51:59.181Z")
	if err != nil {
		t.Errorf("error : %v", err)
		return
	}

	timeParsed := parse.UnixNano() / 1000000

	assert.Equal(t, "2021-09-13T17:51:59.181Z", time.Unix(0, timeParsed*1000000).UTC().Format("2006-01-02T15:04:05.999Z"))
}

// tests for deploy with hash 48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66 on casper testnet
func TestDeployUtil_HashDeployHeaderCorrectly(t *testing.T) {
	parse, err := time.Parse(time.RFC3339, "2021-09-13T17:51:59.181Z")
	if err != nil {
		t.Errorf("error : %v", err)
		return
	}

	bodyHash, err := hex.DecodeString("f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37")

	if err != nil {
		t.Errorf("error : %v", err)
		return
	}

	accountHash, err := hex.DecodeString("d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061")

	if err != nil {
		t.Errorf("error : %v", err)
		return
	}

	deployHeader := NewDeployHeader(keypair.PublicKey{
		PubKeyData: accountHash,
		Tag:        1,
	}, parse.UnixNano()/1000000, (30 * time.Minute).Milliseconds(), 1, bodyHash, nil, "casper-test")

	assert.Equal(t, "01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061cd8249e07b01000040771b00000000000100000000000000f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37000000000b0000006361737065722d74657374",
		hex.EncodeToString(deployHeader.ToBytes()))

	hashToCompare := blake2b.Sum256(deployHeader.ToBytes())

	assert.Equal(t, "48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66", hex.EncodeToString(hashToCompare[:]))
}

func TestDeployUtil_HashBodyAndMarshalJSONCorrectly(t *testing.T) {
	var payment *ExecutableDeployItem
	var session *ExecutableDeployItem

	payment = StandardPayment(big.NewInt(10000))

	assert.Equal(t, "00000000000100000006000000616d6f756e740300000002102708", hex.EncodeToString(payment.ToBytes()))

	session = NewTransfer(big.NewInt(2500000000), dest, "", uint64(1))

	assert.Equal(t, "050300000006000000616d6f756e74050000000400f90295080600000074617267657420000000a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d0f20000000020000006964090000000101000000000000000d05",
		hex.EncodeToString(session.ToBytes()))

	serializedBody := SerializeBody(payment, session)

	bodyHash := blake2b.Sum256(serializedBody)
	resBodyHash := bodyHash[:]
	assert.Equal(t, "f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37",
		hex.EncodeToString(resBodyHash))

	timeParsed, err := time.Parse(time.RFC3339, "2021-09-13T17:51:59.181Z")
	if err != nil {
		t.Errorf("error : %v", err)
		return
	}

	target, _ := hex.DecodeString("d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061")
	deploy := MakeDeploy(NewDeployParams(keypair.PublicKey{PubKeyData: target, Tag: keypair.KeyTagEd25519}, "casper-test", nil, int64(timeParsed.UnixNano())/1000000), payment, session)

	assert.Equal(t, "48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66", hex.EncodeToString(deploy.Hash))

	marshal, err := json.Marshal(deploy)

	if err != nil {
		t.Errorf("error : %v", err)
		return
	}

	assert.Equal(t,
		"{\"hash\":\"48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66\",\"header\":{\"account\":\"01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061\",\"timestamp\":\"2021-09-13T17:51:59.181Z\",\"ttl\":\"30m0s\",\"gas_price\":1,\"body_hash\":\"f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37\",\"dependencies\":[],\"chain_name\":\"casper-test\"},\"payment\":{\"ModuleBytes\":{\"args\":[[\"amount\",{\"bytes\":\"021027\",\"cl_type\":\"U512\"}]],\"module_bytes\":\"\"}},\"session\":{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"010100000000000000\",\"cl_type\":{\"Option\":\"U64\"}}]]}},\"approvals\":[]}",
		string(marshal))
}

func TestDeployUtil_UnmarshalJSONCorrectly(t *testing.T) {
	deployInJSON := "{\"hash\":\"48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66\",\"header\":{\"account\":\"01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061\",\"timestamp\":\"2021-09-13T17:51:59.181Z\",\"ttl\":\"30m0s\",\"gas_price\":1,\"body_hash\":\"f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37\",\"dependencies\":[],\"chain_name\":\"casper-test\"},\"payment\":{\"ModuleBytes\":{\"args\":[[\"amount\",{\"bytes\":\"021027\",\"cl_type\":\"U512\"}]],\"module_bytes\":\"\"}},\"session\":{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"010100000000000000\",\"cl_type\":{\"Option\":\"U64\"}}]]}},\"approvals\":[]}"

	var deployResult Deploy

	err := json.Unmarshal([]byte(deployInJSON), &deployResult)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.Equal(t, "48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66", hex.EncodeToString(deployResult.Hash))
	assert.Equal(t, "f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37", hex.EncodeToString(deployResult.Header.BodyHash))

	account, err := deployResult.Header.Account.ToBytes()
	if err != nil {
		t.Errorf("can't put already known deploy")
		return
	}
	assert.Equal(t, "01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061", hex.EncodeToString(account))
	assert.Equal(t, 0, len(deployResult.Approvals))
	assert.Equal(t, Duration(1800000), deployResult.Header.TTL)
	assert.Equal(t, 0, len(deployResult.Header.Dependencies))
	assert.Equal(t, Timestamp(1631555519181), deployResult.Header.Timestamp)
	assert.Equal(t, "casper-test", deployResult.Header.ChainName)
	assert.Equal(t, uint64(1), deployResult.Header.GasPrice)
	assert.NotEqual(t, nil, deployResult.Payment.ModuleBytes)
	assert.NotEqual(t, nil, deployResult.Payment.ModuleBytes.Args.Args["amount"])
	assert.NotEqual(t, nil, deployResult.Session.Transfer)
	assert.NotEqual(t, nil, deployResult.Session.Transfer.Args.Args["amount"])
	assert.NotEqual(t, nil, deployResult.Session.Transfer.Args.Args["target"])
	assert.NotEqual(t, nil, deployResult.Session.Transfer.Args.Args["id"])
}

func TestDeployUtil_IsDeploy(t *testing.T) {
	deploy := NewTransferToUniqAddress(*source, UniqAddress{
		PublicKey:  dest,
		TransferId: 10,
	}, big.NewInt(3), big.NewInt(1), "casper-test", "")

	assert.True(t, deploy.IsStandardPayment())
	assert.True(t, deploy.IsTransfer())

	deploy.Payment.ModuleBytes.ModuleBytes = []byte("not standart payment anymore")
	deploy.Session.Type = ExecutableDeployItemTypeModuleBytes

	assert.False(t, deploy.IsStandardPayment())
	assert.False(t, deploy.IsTransfer())
}

func TestDeployUtil_MarshalTransfer(t *testing.T) {
	transfer := NewTransfer(big.NewInt(2500000000), dest, "", 1)

	marshalJSON, err := transfer.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"010100000000000000\",\"cl_type\":{\"Option\":\"U64\"}}]]}}",
		string(marshalJSON))
}

func TestDeployUtil_NewTransferWithoutId(t *testing.T) {
	transfer := NewTransferWithoutId(big.NewInt(2500000000), dest, "")

	marshalJSON, err := transfer.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"00\",\"cl_type\":{\"Option\":\"U64\"}}]]}}",
		string(marshalJSON))
}

func TestDeployUtil_UnmarshalTransfer(t *testing.T) {
	var transfer ExecutableDeployItem

	err := json.Unmarshal([]byte("{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"010100000000000000\",\"cl_type\":{\"Option\":\"U64\"}}]]}}"),
		&transfer)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, transfer.IsTransfer())

	assert.Equal(t, transfer.Transfer.Args.Args["amount"].Tag, types.CLTypeU512)
	assert.Equal(t, transfer.Transfer.Args.Args["amount"].StringBytes, "0400f90295")

	assert.Equal(t, transfer.Transfer.Args.Args["target"].Tag, types.CLTypeByteArray)
	assert.Equal(t, transfer.Transfer.Args.Args["target"].StringBytes, "a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d")

	assert.Equal(t, transfer.Transfer.Args.Args["id"].Tag, types.CLTypeOption)
	assert.True(t, transfer.Transfer.Args.Args["id"].IsOptional)
	assert.Equal(t, transfer.Transfer.Args.Args["id"].Optional.Tag, types.CLTypeU64)
	assert.Equal(t, transfer.Transfer.Args.Args["id"].Optional.StringBytes, "010100000000000000")
}

func TestDeployUtil_UnmarshalIncorrectTransfer(t *testing.T) {
	var transfer ExecutableDeployItem

	err := json.Unmarshal([]byte("{\"Transfer\":{}}"),
		&transfer)
	assert.Error(t, err)

	err = json.Unmarshal([]byte("{\"Transfer\":{\"args\":[[]]}}"),
		&transfer)
	assert.Error(t, err)
}

func TestDeployUtil_StoredContractByHash(t *testing.T) {
	var hash32 [32]byte

	for i := 0; i < 32; i++ {
		hash32[i] = source.PubKeyData[i]
	}

	storedContractByHash := NewStoredContractByHash(hash32, "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	assert.Equal(t, "01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d8303206104000000746573740100000006000000616d6f756e74050000000400f9029508", hex.EncodeToString(storedContractByHash.ToBytes()))
}

func TestDeployUtil_MarshalStoredContractByHash(t *testing.T) {
	var hash32 [32]byte

	for i := 0; i < 32; i++ {
		hash32[i] = source.PubKeyData[i]
	}

	storedContractByHash := NewStoredContractByHash(hash32, "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	marshalJSON, err := storedContractByHash.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"StoredContractByHash\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}]],\"entry_point\":\"test\",\"hash\":\"d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061\"}}",
		string(marshalJSON))
}

func TestDeployUtil_MarshalStoredContractByHashMapValue(t *testing.T) {
	var hash32 [32]byte

	decodedHash, err2 := hex.DecodeString("28ce14c210c53735d43eafef7f4446fb51c0761075c553029a9eb30988a0caa1")
	if err2 != nil {
		return
	}

	for i := 0; i < 32; i++ {
		hash32[i] = decodedHash[i]
	}

	storedContractByHash := NewStoredContractByHash(hash32, "store_batch",
		RuntimeArgs{Args: map[string]Value{"storage_id": {Tag: types.CLTypeString, StringBytes: "050000004142432d32"},
			"storage_data": {
				Tag: types.CLTypeMap,
				Map: &ValueMap{
					KeyType:   types.CLTypeString,
					ValueType: types.CLTypeString,
				},
				StringBytes: "010000000400000032312d324000000033623238633238636664633165386666663030613165346630313931663933383838316330613831396362383435653365366536303037383235336564613662",
			},
		}, KeyOrder: []string{"storage_id", "storage_data"}})

	marshalJSON, err := storedContractByHash.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"StoredContractByHash\":{\"args\":[[\"storage_id\",{\"bytes\":\"050000004142432d32\",\"cl_type\":\"String\"}],[\"storage_data\",{\"bytes\":\"010000000400000032312d324000000033623238633238636664633165386666663030613165346630313931663933383838316330613831396362383435653365366536303037383235336564613662\",\"cl_type\":{\"Map\":{\"key\":\"String\",\"value\":\"String\"}}}]],\"entry_point\":\"store_batch\",\"hash\":\"28ce14c210c53735d43eafef7f4446fb51c0761075c553029a9eb30988a0caa1\"}}",
		string(marshalJSON))
}

func TestDeployUtil_UnmarshalStoredContractByHash(t *testing.T) {
	var storedContractByHash ExecutableDeployItem

	err := json.Unmarshal([]byte("{\"StoredContractByHash\":{\"hash\":\"28ce14c210c53735d43eafef7f4446fb51c0761075c553029a9eb30988a0caa1\",\"entry_point\":\"store_batch\",\"args\":[[\"storage_id\",{\"cl_type\":\"String\",\"bytes\":\"050000004142432d32\",\"parsed\":\"ABC-2\"}],[\"storage_data\",{\"cl_type\":{\"Map\":{\"key\":\"String\",\"value\":\"String\"}},\"bytes\":\"010000000400000032312d324000000033623238633238636664633165386666663030613165346630313931663933383838316330613831396362383435653365366536303037383235336564613662\",\"parsed\":[{\"key\":\"21-2\",\"value\":\"3b28c28cfdc1e8fff00a1e4f0191f938881c0a819cb845e3e6e60078253eda6b\"}]}]]}}"),
		&storedContractByHash)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, storedContractByHash.IsStoredContractByHash())
	assert.Equal(t, "28ce14c210c53735d43eafef7f4446fb51c0761075c553029a9eb30988a0caa1", hex.EncodeToString(storedContractByHash.StoredContractByHash.Hash[:]))
	assert.Equal(t, "store_batch", storedContractByHash.StoredContractByHash.Entrypoint)
	assert.Equal(t, types.CLTypeString, storedContractByHash.StoredContractByHash.Args.Args["storage_id"].Tag)
	assert.Equal(t, "050000004142432d32", storedContractByHash.StoredContractByHash.Args.Args["storage_id"].StringBytes)
	assert.Equal(t, types.CLTypeMap, storedContractByHash.StoredContractByHash.Args.Args["storage_data"].Tag)
	assert.Equal(t, "010000000400000032312d324000000033623238633238636664633165386666663030613165346630313931663933383838316330613831396362383435653365366536303037383235336564613662", storedContractByHash.StoredContractByHash.Args.Args["storage_data"].StringBytes)
	assert.Equal(t, types.CLTypeString, storedContractByHash.StoredContractByHash.Args.Args["storage_data"].Map.KeyType)
	assert.Equal(t, types.CLTypeString, storedContractByHash.StoredContractByHash.Args.Args["storage_data"].Map.ValueType)

	//{"key":"21-2","value":"3b28c28cfdc1e8fff00a1e4f0191f938881c0a819cb845e3e6e60078253eda6b"}

	decodedData, err := hex.DecodeString(storedContractByHash.StoredContractByHash.Args.Args["storage_data"].StringBytes)
	if err != nil {
		return
	}

	clMap := types.CLValue{Type: types.CLTypeMap, Map: &types.CLMap{KeyType: types.CLTypeString, ValueType: types.CLTypeString}}

	n, err := types.UnmarshalCLValue(decodedData, &clMap)
	if err != nil {
		return
	}

	assert.Equal(t, len(decodedData), n)
	assert.Equal(t, types.CLTypeString, clMap.Map.KeyType)
	assert.Equal(t, types.CLTypeString, clMap.Map.ValueType)
	assert.NotNil(t, clMap.Map.Raw["21-2"])
	assert.Equal(t, types.CLTypeString, clMap.Map.Raw["21-2"].Type)
	assert.Equal(t, "3b28c28cfdc1e8fff00a1e4f0191f938881c0a819cb845e3e6e60078253eda6b", *clMap.Map.Raw["21-2"].String)
}

func TestDeployUtil_StoredContractByName(t *testing.T) {
	var hash32 [32]byte

	for i := 0; i < 32; i++ {
		hash32[i] = source.PubKeyData[i]
	}

	storedContractByHash := NewStoredContractByName("test", "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	assert.Equal(t, "02040000007465737404000000746573740100000006000000616d6f756e74050000000400f9029508", hex.EncodeToString(storedContractByHash.ToBytes()))
}

func TestDeployUtil_MarshalStoredContractByName(t *testing.T) {
	storedContractByName := NewStoredContractByName("test", "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	marshalJSON, err := storedContractByName.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"StoredContractByName\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}]],\"entry_point\":\"test\",\"name\":\"test\"}}",
		string(marshalJSON))
}

func TestDeployUtil_UnmarshalStoredContractByName(t *testing.T) {
	var storedContractByName ExecutableDeployItem

	err := json.Unmarshal([]byte("{\"StoredContractByName\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}]],\"entry_point\":\"test\",\"name\":\"test\"}}"),
		&storedContractByName)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, storedContractByName.IsStoredContractByName())
	assert.Equal(t, "test", storedContractByName.StoredContractByName.Name)
	assert.Equal(t, "test", storedContractByName.StoredContractByName.Entrypoint)
	assert.Equal(t, types.CLTypeU512, storedContractByName.StoredContractByName.Args.Args["amount"].Tag)
	assert.Equal(t, "0400f90295", storedContractByName.StoredContractByName.Args.Args["amount"].StringBytes)
}

func TestDeployUtil_StoredVersionedContractByHash(t *testing.T) {
	var hash32 [32]byte

	for i := 0; i < 32; i++ {
		hash32[i] = source.PubKeyData[i]
	}

	storedVersionedContractByHash := NewStoredVersionedContractByHash(hash32, 5, "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	assert.Equal(t, "03d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061010500000004000000746573740100000006000000616d6f756e74050000000400f9029508", hex.EncodeToString(storedVersionedContractByHash.ToBytes()))
}

func TestDeployUtil_MarshalStoredVersionedContractByHash(t *testing.T) {
	var hash32 [32]byte

	for i := 0; i < 32; i++ {
		hash32[i] = source.PubKeyData[i]
	}

	storedVersionedContractByHash := NewStoredVersionedContractByHashWithoutVersion(hash32, "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	marshalJSON, err := storedVersionedContractByHash.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"StoredVersionedContractByHash\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}]],\"entry_point\":\"test\",\"hash\":\"d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061\",\"version\":\"None\"}}",
		string(marshalJSON))
}

func TestDeployUtil_UnmarshalStoredVersionedContractByHash(t *testing.T) {
	var storedVersionedContractByHash ExecutableDeployItem

	err := json.Unmarshal([]byte("{\"StoredVersionedContractByHash\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}]],\"entry_point\":\"test\",\"hash\":\"d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061\",\"version\":\"None\"}}"),
		&storedVersionedContractByHash)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, storedVersionedContractByHash.IsStoredVersionedContractByHash())
	assert.Equal(t, "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061", hex.EncodeToString(storedVersionedContractByHash.StoredVersionedContractByHash.Hash[:]))
	assert.Equal(t, "test", storedVersionedContractByHash.StoredVersionedContractByHash.Entrypoint)
	assert.Equal(t, &types.CLValue{Type: types.CLTypeOption, Option: nil}, storedVersionedContractByHash.StoredVersionedContractByHash.Version)
	assert.Equal(t, types.CLTypeU512, storedVersionedContractByHash.StoredVersionedContractByHash.Args.Args["amount"].Tag)
	assert.Equal(t, "0400f90295", storedVersionedContractByHash.StoredVersionedContractByHash.Args.Args["amount"].StringBytes)
}

func TestDeployUtil_StoredVersionedContractByName(t *testing.T) {
	var hash32 [32]byte

	for i := 0; i < 32; i++ {
		hash32[i] = source.PubKeyData[i]
	}

	storedVersionedContractByHash := NewStoredVersionedContractByName("test", 5, "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	assert.Equal(t, "040400000074657374010500000004000000746573740100000006000000616d6f756e74050000000400f9029508", hex.EncodeToString(storedVersionedContractByHash.ToBytes()))
}

func TestDeployUtil_MarshalStoredVersionedContractByName(t *testing.T) {
	storedVersionedContractByName := NewStoredVersionedContractByName("test", 5, "test",
		RuntimeArgs{Args: map[string]Value{"amount": {Tag: types.CLTypeU512, StringBytes: "0400f90295"}}, KeyOrder: []string{"amount"}})

	marshalJSON, err := storedVersionedContractByName.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"StoredVersionedContractByName\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}]],\"entry_point\":\"test\",\"name\":\"test\",\"version\":\"5\"}}",
		string(marshalJSON))
}

func TestDeployUtil_UnmarshalStoredVersionedContractByName(t *testing.T) {
	var storedVersionedContractByName ExecutableDeployItem

	err := json.Unmarshal([]byte("{\"StoredVersionedContractByName\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}]],\"entry_point\":\"test\",\"name\":\"test\",\"version\":\"5\"}}"),
		&storedVersionedContractByName)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, storedVersionedContractByName.IsStoredVersionedContractByName())
	assert.Equal(t, "test", storedVersionedContractByName.StoredVersionedContractByName.Name)
	assert.NotNil(t, storedVersionedContractByName.StoredVersionedContractByName.Version.Option)
	assert.Equal(t, "5", strconv.Itoa(int(*storedVersionedContractByName.StoredVersionedContractByName.Version.Option.U32)))
	assert.Equal(t, "test", storedVersionedContractByName.StoredVersionedContractByName.Entrypoint)
	assert.Equal(t, types.CLTypeU512, storedVersionedContractByName.StoredVersionedContractByName.Args.Args["amount"].Tag)
	assert.Equal(t, "0400f90295", storedVersionedContractByName.StoredVersionedContractByName.Args.Args["amount"].StringBytes)
}

func TestDeployUtil_MarshalStoredContractByHashPublicKey(t *testing.T) {
	var hash32 [32]byte
	decodedHash, err2 := hex.DecodeString("ccb576d6ce6dec84a551e48f0d0b7af89ddba44c7390b690036257a04a3ae9ea")
	if err2 != nil {
		return
	}

	for i := 0; i < 32; i++ {
		hash32[i] = decodedHash[i]
	}

	storedContractByHash := NewStoredContractByHash(hash32, "delegate",
		RuntimeArgs{Args: map[string]Value{
			"delegator": {
				Tag:         types.CLTypePublicKey,
				StringBytes: "0203b24eb09b295d3122d7abd1aafccb2c899d5db159e73e1c8fc972f017c308a363",
			},
			"validator": {
				Tag:         types.CLTypePublicKey,
				StringBytes: "012bac1d0ff9240ff0b7b06d555815640497861619ca12583ddef434885416e69b",
			},
			"amount": {
				Tag:         types.CLTypeU512,
				StringBytes: "05e0e0c5f4a6",
			},
		}, KeyOrder: []string{"delegator", "validator", "amount"}})

	marshalJSON, err := storedContractByHash.MarshalJSON()
	if err != nil {
		t.Errorf("can't marshal %v", err)
		return
	}

	assert.Equal(t, "{\"StoredContractByHash\":{\"args\":[[\"delegator\",{\"bytes\":\"0203b24eb09b295d3122d7abd1aafccb2c899d5db159e73e1c8fc972f017c308a363\",\"cl_type\":\"PublicKey\"}],[\"validator\",{\"bytes\":\"012bac1d0ff9240ff0b7b06d555815640497861619ca12583ddef434885416e69b\",\"cl_type\":\"PublicKey\"}],[\"amount\",{\"bytes\":\"05e0e0c5f4a6\",\"cl_type\":\"U512\"}]],\"entry_point\":\"delegate\",\"hash\":\"ccb576d6ce6dec84a551e48f0d0b7af89ddba44c7390b690036257a04a3ae9ea\"}}",
		string(marshalJSON))
}

func TestDeployUtil_UnmarshalStoredContractByHashPublicKey(t *testing.T) {
	var storedContractByHash ExecutableDeployItem

	err := json.Unmarshal([]byte("{\"StoredContractByHash\":{\"args\":[[\"delegator\",{\"cl_type\":\"PublicKey\",\"bytes\":\"0203b24eb09b295d3122d7abd1aafccb2c899d5db159e73e1c8fc972f017c308a363\"}],[\"validator\",{\"cl_type\":\"PublicKey\",\"bytes\":\"012bac1d0ff9240ff0b7b06d555815640497861619ca12583ddef434885416e69b\"}],[\"amount\",{\"cl_type\":\"U512\",\"bytes\":\"05e0e0c5f4a6\"}]],\"entry_point\":\"delegate\",\"hash\":\"ccb576d6ce6dec84a551e48f0d0b7af89ddba44c7390b690036257a04a3ae9ea\"}}"),
		&storedContractByHash)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, storedContractByHash.IsStoredContractByHash())
	assert.Equal(t, "ccb576d6ce6dec84a551e48f0d0b7af89ddba44c7390b690036257a04a3ae9ea", hex.EncodeToString(storedContractByHash.StoredContractByHash.Hash[:]))
	assert.Equal(t, "delegate", storedContractByHash.StoredContractByHash.Entrypoint)
	assert.Equal(t, types.CLTypeU512, storedContractByHash.StoredContractByHash.Args.Args["amount"].Tag)
	assert.Equal(t, "05e0e0c5f4a6", storedContractByHash.StoredContractByHash.Args.Args["amount"].StringBytes)
	assert.Equal(t, types.CLTypePublicKey, storedContractByHash.StoredContractByHash.Args.Args["validator"].Tag)
	assert.Equal(t, "012bac1d0ff9240ff0b7b06d555815640497861619ca12583ddef434885416e69b", storedContractByHash.StoredContractByHash.Args.Args["validator"].StringBytes)
	assert.Equal(t, types.CLTypePublicKey, storedContractByHash.StoredContractByHash.Args.Args["delegator"].Tag)
	assert.Equal(t, "0203b24eb09b295d3122d7abd1aafccb2c899d5db159e73e1c8fc972f017c308a363", storedContractByHash.StoredContractByHash.Args.Args["delegator"].StringBytes)
}

func TestDeployUtil_SetArg(t *testing.T) {
	var transfer ExecutableDeployItem
	err := json.Unmarshal([]byte("{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"010100000000000000\",\"cl_type\":{\"Option\":\"U64\"}}]]}}"),
		&transfer)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.Equal(t, transfer.Transfer.Args.Args["amount"].Tag, types.CLTypeU512)
	assert.Equal(t, transfer.Transfer.Args.Args["amount"].StringBytes, "0400f90295")

	amount := uint64(1024)

	err = transfer.SetArg("amount", types.CLValue{Type: types.CLTypeU64, U64: &amount})
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.Equal(t, transfer.Transfer.Args.Args["amount"].Tag, types.CLTypeU64)
	assert.Equal(t, transfer.Transfer.Args.Args["amount"].StringBytes, "0004000000000000")

	one := uint64(1)

	err = transfer.SetArg("option", types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{
		Type: types.CLTypeU64,
		U64:  &one,
	}})
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.Equal(t, transfer.Transfer.Args.Args["option"].Tag, types.CLTypeOption)
	assert.True(t, transfer.Transfer.Args.Args["option"].IsOptional)
	assert.Equal(t, transfer.Transfer.Args.Args["option"].Optional.Tag, types.CLTypeU64)
	assert.Equal(t, transfer.Transfer.Args.Args["option"].Optional.StringBytes, "010100000000000000")

}

func TestDeployUtil_SetArgToDeploy(t *testing.T) {
	deploy := NewTransferToUniqAddress(*source, UniqAddress{
		PublicKey:  dest,
		TransferId: 10,
	}, big.NewInt(3), big.NewInt(1), "casper-test", "")

	assert.True(t, deploy.IsStandardPayment())
	assert.True(t, deploy.IsTransfer())

	amount := uint64(1024)

	err := deploy.AddArgToDeploy("amount", types.CLValue{Type: types.CLTypeU64, U64: &amount})
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.Equal(t, deploy.Session.Transfer.Args.Args["amount"].Tag, types.CLTypeU64)
	assert.Equal(t, deploy.Session.Transfer.Args.Args["amount"].StringBytes, "0004000000000000")

	one := uint64(1)

	err = deploy.Session.SetArg("option", types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{
		Type: types.CLTypeU64,
		U64:  &one,
	}})
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.Equal(t, deploy.Session.Transfer.Args.Args["option"].Tag, types.CLTypeOption)
	assert.True(t, deploy.Session.Transfer.Args.Args["option"].IsOptional)
	assert.Equal(t, deploy.Session.Transfer.Args.Args["option"].Optional.Tag, types.CLTypeU64)
	assert.Equal(t, deploy.Session.Transfer.Args.Args["option"].Optional.StringBytes, "010100000000000000")
}

func TestDeployUtil_ShouldNotSetArgToSignedDeploy(t *testing.T) {
	deploy := NewTransferToUniqAddress(*source, UniqAddress{
		PublicKey:  dest,
		TransferId: 10,
	}, big.NewInt(3), big.NewInt(1), "casper-test", "")

	assert.True(t, deploy.IsStandardPayment())
	assert.True(t, deploy.IsTransfer())

	deploy.Approvals = append(deploy.Approvals, Approval{})

	amount := uint64(1024)

	err := deploy.AddArgToDeploy("amount", types.CLValue{Type: types.CLTypeU64, U64: &amount})
	assert.Error(t, err)
}

func TestDeployUtil_ValidateDeploy(t *testing.T) {
	deploy := NewTransferToUniqAddress(*source, UniqAddress{
		PublicKey:  dest,
		TransferId: 10,
	}, big.NewInt(3), big.NewInt(1), "casper-test", "")

	assert.True(t, deploy.ValidateDeploy())

	deployInJSON := "{\"hash\":\"48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66\",\"header\":{\"account\":\"01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061\",\"timestamp\":\"2021-09-13T17:51:59.181Z\",\"ttl\":\"30m0s\",\"gas_price\":1,\"body_hash\":\"f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37\",\"dependencies\":[],\"chain_name\":\"casper-test\"},\"payment\":{\"ModuleBytes\":{\"args\":[[\"amount\",{\"bytes\":\"021027\",\"cl_type\":\"U512\"}]],\"module_bytes\":\"\"}},\"session\":{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"010100000000000000\",\"cl_type\":{\"Option\":\"U64\"}}]]}},\"approvals\":[]}"

	var deployResult Deploy

	err := json.Unmarshal([]byte(deployInJSON), &deployResult)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, deployResult.ValidateDeploy())

	deployResult.Hash = []byte("1234567")
	assert.False(t, deployResult.ValidateDeploy())
	deploy.Header.BodyHash = []byte("1234567")
	assert.False(t, deploy.ValidateDeploy())
}

func TestDeployUtil_SignDeploy(t *testing.T) {
	deployInJSON := "{\"hash\":\"48b33972cdc075d82363279640490b64bcac26cd540c8cf16da688d400c86b66\",\"header\":{\"account\":\"01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061\",\"timestamp\":\"2021-09-13T17:51:59.181Z\",\"ttl\":\"30m0s\",\"gas_price\":1,\"body_hash\":\"f9608668e24e68cad0c930016e1885d1d82fdb655b254130c32b586c4443af37\",\"dependencies\":[],\"chain_name\":\"casper-test\"},\"payment\":{\"ModuleBytes\":{\"args\":[[\"amount\",{\"bytes\":\"021027\",\"cl_type\":\"U512\"}]],\"module_bytes\":\"\"}},\"session\":{\"Transfer\":{\"args\":[[\"amount\",{\"bytes\":\"0400f90295\",\"cl_type\":\"U512\"}],[\"target\",{\"bytes\":\"a6d3d9fb1044cf5db1b30ad3f8f2c2c69e48ae69ab8aae6f02d69b0d0faa9e3d\",\"cl_type\":{\"ByteArray\":32}}],[\"id\",{\"bytes\":\"010100000000000000\",\"cl_type\":{\"Option\":\"U64\"}}]]}},\"approvals\":[]}"

	var deploy Deploy

	err := json.Unmarshal([]byte(deployInJSON), &deploy)
	if err != nil {
		t.Errorf("can't unmarshal %v", err)
		return
	}

	assert.True(t, deploy.ValidateDeploy())

	deploy.SignDeploy(sourceKeyPair)

	assert.Equal(t, 1, len(deploy.Approvals))
	assert.Equal(t, "d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061", hex.EncodeToString(deploy.Approvals[0].Signer.PubKeyData))
	assert.Equal(t, "4ffe34cf43a62f94181090a9e1bb52db207d37138c927d07d642b81267822f66333b2014fe59cc8aca97a3852b9e66eb2e761cceb4deed2b03776b99bdef0a09", hex.EncodeToString(deploy.Approvals[0].Signature.SignatureData))
}
