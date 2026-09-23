package jsonrpc_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestJSONRPCAccountAnchor(t *testing.T) {
	td := setup(t)
	root := bytesFill(0x11)
	addr, acc := td.GenerateTestAccount(testsuite.AccountWithNumber(3), testsuite.AccountWithBalance(5))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        root,
		LockedDeposit:   9,
		CreatedAtHeight: 4,
		CreatedAtTime:   111,
		UpdatedAtHeight: 4,
		UpdatedAtTime:   111,
	}))
	td.gRPCServer.FakeState.AddTestAccount(addr, acc)

	result := callJSONRPC(t, td, "pactus.blockchain.get_account", map[string]any{
		"address": addr.String(),
	})
	accountBody := objectField(t, result, "account")
	anchor := objectField(t, accountBody, "anchor")
	require.Equal(t, base64.StdEncoding.EncodeToString(root), stringField(t, anchor, "root_hash", "rootHash"))
	require.Equal(t, int64(4), numberField(t, anchor, "created_at_height", "createdAtHeight"))
	require.Equal(t, int64(9), numberField(t, anchor, "locked_deposit", "lockedDeposit"))

	plainAddr, plain := td.GenerateTestAccount(testsuite.AccountWithNumber(8))
	td.gRPCServer.FakeState.AddTestAccount(plainAddr, plain)
	plainResult := callJSONRPC(t, td, "pactus.blockchain.get_account", map[string]any{
		"address": plainAddr.String(),
	})
	plainAccount := objectField(t, plainResult, "account")
	require.Nil(t, plainAccount["anchor"])

	acc.ClearAnchor()
	cleared := callJSONRPC(t, td, "pactus.blockchain.get_account", map[string]any{
		"address": addr.String(),
	})
	clearedAccount := objectField(t, cleared, "account")
	require.Nil(t, clearedAccount["anchor"])
}

func TestJSONRPCGetAnchorMissing(t *testing.T) {
	td := setup(t)
	result := callJSONRPC(t, td, "pactus.blockchain.get_anchor", map[string]any{
		"address": td.RandAccAddress().String(),
	})
	require.NotEqual(t, true, result["found"])
	require.Nil(t, result["anchor"])
}

func TestJSONRPCListAnchors(t *testing.T) {
	td := setup(t)
	first, firstAcc := td.GenerateTestAccount(testsuite.AccountWithNumber(1))
	second, secondAcc := td.GenerateTestAccount(testsuite.AccountWithNumber(2))
	require.NoError(t, firstAcc.SetAnchor(account.AnchorData{
		RootHash:      bytesFill(0x11),
		LockedDeposit: 1,
	}))
	require.NoError(t, secondAcc.SetAnchor(account.AnchorData{
		RootHash:      bytesFill(0x22),
		LockedDeposit: 1,
	}))
	td.gRPCServer.FakeState.AddTestAccount(first, firstAcc)
	td.gRPCServer.FakeState.AddTestAccount(second, secondAcc)

	result := callJSONRPC(t, td, "pactus.blockchain.list_anchors", map[string]any{
		"count": 20,
	})
	require.Equal(t, int64(2), numberField(t, result, "total"))
	items, ok := result["items"].([]any)
	require.True(t, ok)
	require.Len(t, items, 2)
	require.Equal(t, first.String(), stringField(t, items[0].(map[string]any), "address"))
	require.Equal(t, second.String(), stringField(t, items[1].(map[string]any), "address"))
}

func TestJSONRPCRawAnchor(t *testing.T) {
	td := setup(t)
	from := td.RandAccAddress().String()
	root := base64.StdEncoding.EncodeToString(bytesFill(0x31))

	rawResult := callJSONRPC(t, td, "pactus.transaction.get_raw_anchor_transaction", map[string]any{
		"from":      from,
		"action":    0,
		"root_hash": root,
		"deposit":   5,
		"fee":       7,
		"lock_time": 1,
	})
	raw := stringField(t, rawResult, "raw_transaction", "rawTransaction")

	decoded := callJSONRPC(t, td, "pactus.transaction.decode_raw_transaction", map[string]any{
		"raw_transaction": raw,
	})
	trx := objectField(t, decoded, "transaction")
	require.True(t, payloadIsAnchor(field(t, trx, "payload_type", "payloadType")))
	anchor := anchorObject(t, trx)
	action := int64(0)
	if _, ok := anchor["action"]; ok {
		action = numberField(t, anchor, "action")
	}
	require.Equal(t, int64(0), action)
	require.Equal(t, int64(5), numberField(t, anchor, "deposit"))
	require.Equal(t, int64(7), numberField(t, trx, "fee"))
	require.NotEqual(t, numberField(t, anchor, "deposit"), numberField(t, trx, "fee"))
}

func callJSONRPC(t *testing.T, td *testData, method string, params map[string]any) map[string]any {
	t.Helper()

	requestBody, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  method,
		"params":  params,
	})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		t.Context(), http.MethodPost,
		"http://"+td.jsonrpcServer.Address(),
		bytes.NewBuffer(requestBody),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var response map[string]any
	require.NoError(t, json.Unmarshal(raw, &response), string(raw))
	require.NotContains(t, response, "error", string(raw))
	result, ok := response["result"].(map[string]any)
	require.Truef(t, ok, "body %s", raw)

	return result
}

func payloadIsAnchor(value any) bool {
	switch typed := value.(type) {
	case float64:
		return typed == float64(payload.TypeAnchor)
	case string:
		return typed == "PAYLOAD_TYPE_ANCHOR" || typed == "7"
	default:
		return false
	}
}

func anchorObject(t *testing.T, trx map[string]any) map[string]any {
	t.Helper()

	if anchor, ok := trx["anchor"].(map[string]any); ok {
		return anchor
	}
	payloadBody, ok := trx["Payload"].(map[string]any)
	require.Truef(t, ok, "anchor payload missing in %v", trx)
	anchor, ok := payloadBody["Anchor"].(map[string]any)
	require.Truef(t, ok, "anchor payload missing in %v", trx)

	return anchor
}

func objectField(t *testing.T, body map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := body[key].(map[string]any)
	require.Truef(t, ok, "missing %s in %v", key, body)

	return value
}

func field(t *testing.T, body map[string]any, keys ...string) any {
	t.Helper()

	for _, key := range keys {
		if value, ok := body[key]; ok {
			return value
		}
	}
	require.Failf(t, "missing field", "%v in %v", keys, body)

	return nil
}

func stringField(t *testing.T, body map[string]any, keys ...string) string {
	t.Helper()

	value, ok := field(t, body, keys...).(string)
	require.Truef(t, ok, "%v in %v", keys, body)

	return value
}

func numberField(t *testing.T, body map[string]any, keys ...string) int64 {
	t.Helper()

	switch value := field(t, body, keys...).(type) {
	case float64:
		return int64(value)
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		require.NoError(t, err)

		return parsed
	default:
		require.Failf(t, "not a number", "%v in %v", keys, body)

		return 0
	}
}

func bytesFill(fill byte) []byte {
	return bytes.Repeat([]byte{fill}, 32)
}
