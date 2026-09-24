package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/grpc/fake"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	fake *fake.FakeGRPCServer
	srv  *Server
}

func setupAnchorHTTP(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	gRPCServer := fake.NewFakeGRPCServer(t, ts, &grpc.Config{
		Enable: true,
		Listen: "127.0.0.1:0",
	})
	require.NoError(t, gRPCServer.Server.StartServer())

	conf := DefaultConfig()
	conf.Enable = true
	conf.Listen = "127.0.0.1:0"
	srv := NewServer(t.Context(), conf)
	require.NoError(t, srv.StartServer(gRPCServer.Server.Address()))
	t.Cleanup(srv.StopServer)

	return &testData{
		TestSuite: ts,
		fake:      gRPCServer,
		srv:       srv,
	}
}

func TestGatewayGetAnchor(t *testing.T) {
	td := setupAnchorHTTP(t)
	root := bytesFill(0x11)
	address := plantAnchor(t, td, 1, root)

	status, body := getJSON(t, td, "/pactus/blockchain/get_anchor?address="+url.QueryEscape(address))
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, true, body["found"])
	anchor := objectField(t, body, "anchor")
	require.Equal(t, base64.StdEncoding.EncodeToString(root), stringField(t, anchor, "rootHash", "root_hash"))

	missing := td.RandAccAddress().String()
	status, body = getJSON(t, td, "/pactus/blockchain/get_anchor?address="+url.QueryEscape(missing))
	require.Equal(t, http.StatusOK, status)
	require.NotEqual(t, true, body["found"])
	require.Nil(t, body["anchor"])

	for _, address := range []string{crypto.TreasuryAddress.String(), td.RandValAddress().String()} {
		status, _ = getJSON(t, td, "/pactus/blockchain/get_anchor?address="+url.QueryEscape(address))
		require.Equal(t, http.StatusBadRequest, status, address)
	}
}

// The HTTP gateway writes the protobuf message with protojson: int64 fields
// are JSON strings, so a locked deposit above 2^53 keeps its exact digits;
// uint32 fields are JSON numbers and bytes are base64.
func TestGatewayAnchorBoundaryValues(t *testing.T) {
	td := setupAnchorHTTP(t)
	want := testsuite.BoundaryAnchor()
	addr, acc := td.GenerateTestAccount(testsuite.AccountWithNumber(1), testsuite.AccountWithAnchor(want))
	td.fake.FakeState.AddTestAccount(addr, acc)
	query := "?address=" + url.QueryEscape(addr.String())

	body := getJSONExact(t, td, "/pactus/blockchain/get_anchor"+query)
	require.Equal(t, true, body["found"])
	anchor := objectField(t, body, "anchor")

	digits := func(value int64) json.Number {
		return json.Number(strconv.FormatInt(value, 10))
	}
	require.Equal(t, strconv.FormatInt(int64(want.LockedDeposit), 10),
		anyField(t, anchor, "lockedDeposit", "locked_deposit"), "int64 is a JSON string")
	require.Equal(t, digits(int64(want.AnchorType)), anyField(t, anchor, "anchorType", "anchor_type"))
	require.Equal(t, digits(int64(want.CreatedAtHeight)), anyField(t, anchor, "createdAtHeight", "created_at_height"))
	require.Equal(t, digits(int64(want.CreatedAtTime)), anyField(t, anchor, "createdAtTime", "created_at_time"))
	require.Equal(t, digits(int64(want.UpdatedAtHeight)), anyField(t, anchor, "updatedAtHeight", "updated_at_height"))
	require.Equal(t, digits(int64(want.UpdatedAtTime)), anyField(t, anchor, "updatedAtTime", "updated_at_time"))
	require.Equal(t, base64.StdEncoding.EncodeToString(want.RootHash), stringField(t, anchor, "rootHash", "root_hash"))
	require.Equal(t, want.ManifestURI, stringField(t, anchor, "manifestUri", "manifest_uri"))

	account := getJSONExact(t, td, "/pactus/blockchain/get_account"+query)
	require.Equal(t, anchor, objectField(t, objectField(t, account, "account"), "anchor"))

	list := getJSONExact(t, td, "/pactus/blockchain/list_anchors?count=100")
	items, ok := list["items"].([]any)
	require.True(t, ok)
	found := false
	for _, raw := range items {
		item, isObject := raw.(map[string]any)
		require.True(t, isObject)
		if item["address"] == addr.String() {
			require.Equal(t, anchor, item["anchor"])
			found = true
		}
	}
	require.True(t, found)
}

func TestGatewayGetAccount(t *testing.T) {
	td := setupAnchorHTTP(t)
	root := bytesFill(0x21)
	anchored := plantAnchor(t, td, 2, root)
	plainAddr, plain := td.GenerateTestAccount(testsuite.AccountWithBalance(4))
	td.fake.FakeState.AddTestAccount(plainAddr, plain)

	treasury := account.NewAccount(0)
	treasury.AddToBalance(5)
	td.fake.FakeState.AddTestAccount(crypto.TreasuryAddress, treasury)

	status, body := getJSON(t, td, "/pactus/blockchain/get_account?address="+url.QueryEscape(anchored))
	require.Equal(t, http.StatusOK, status)
	accountBody := objectField(t, body, "account")
	anchor := objectField(t, accountBody, "anchor")
	require.Equal(t, base64.StdEncoding.EncodeToString(root), stringField(t, anchor, "rootHash", "root_hash"))

	status, body = getJSON(t, td, "/pactus/blockchain/get_account?address="+url.QueryEscape(plainAddr.String()))
	require.Equal(t, http.StatusOK, status)
	accountBody = objectField(t, body, "account")
	require.Nil(t, accountBody["anchor"])

	treasuryURL := "/pactus/blockchain/get_account?address=" + url.QueryEscape(crypto.TreasuryAddress.String())
	status, body = getJSON(t, td, treasuryURL)
	require.Equal(t, http.StatusOK, status)
	accountBody = objectField(t, body, "account")
	require.Equal(t, int64(5), numberField(t, accountBody, "balance"))
	require.Nil(t, accountBody["anchor"])
}

func TestGatewayListAnchors(t *testing.T) {
	td := setupAnchorHTTP(t)
	first := plantAnchor(t, td, 1, bytesFill(0x11))
	second := plantAnchor(t, td, 2, bytesFill(0x22))

	status, body := getJSON(t, td, "/pactus/blockchain/list_anchors?count=1")
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, int64(2), numberField(t, body, "total"))
	items, ok := body["items"].([]any)
	require.True(t, ok)
	require.Len(t, items, 1)
	require.Equal(t, first, stringField(t, items[0].(map[string]any), "address"))

	status, _ = getJSON(t, td, "/pactus/blockchain/list_anchors?count=101")
	require.Equal(t, http.StatusBadRequest, status)

	status, body = getJSON(t, td, "/pactus/blockchain/list_anchors?skip=2&count=20")
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, int64(2), numberField(t, body, "total"))
	_, hasItems := body["items"]
	if hasItems {
		require.Empty(t, body["items"])
	}
	require.NotEqual(t, first, second)
}

func TestGatewayRawAnchor(t *testing.T) {
	td := setupAnchorHTTP(t)
	from := td.RandAccAddress().String()
	root := bytesFill(0x31)

	status, body := postJSON(t, td, "/pactus/transaction/get_raw_anchor_transaction", map[string]any{
		"from":     from,
		"action":   0,
		"rootHash": base64.StdEncoding.EncodeToString(root),
		"deposit":  5,
		"fee":      7,
		"lockTime": 1,
	})
	require.Equal(t, http.StatusOK, status)
	raw := stringField(t, body, "rawTransaction", "raw_transaction")
	parsed, err := tx.FromString(raw)
	require.NoError(t, err)
	pld := parsed.Payload().(*payload.AnchorPayload)
	require.Equal(t, amount.Amount(5), pld.Deposit)
	require.Equal(t, amount.Amount(7), parsed.Fee())
	require.NotEqual(t, pld.Deposit, parsed.Fee())

	status, _ = postJSON(t, td, "/pactus/transaction/get_raw_anchor_transaction", map[string]any{
		"from":     from,
		"action":   1,
		"deposit":  1,
		"rootHash": base64.StdEncoding.EncodeToString(root),
	})
	require.Equal(t, http.StatusBadRequest, status)

	td.fake.FakeState.EXPECT().CalculateFee(amount.Amount(5), payload.TypeAnchor).Return(amount.Amount(4242))
	status, body = postJSON(t, td, "/pactus/transaction/get_raw_anchor_transaction", map[string]any{
		"from":     from,
		"action":   0,
		"rootHash": base64.StdEncoding.EncodeToString(root),
		"deposit":  5,
		"fee":      0,
		"lockTime": 1,
	})
	require.Equal(t, http.StatusOK, status)
	raw = stringField(t, body, "rawTransaction", "raw_transaction")
	parsed, err = tx.FromString(raw)
	require.NoError(t, err)
	require.Equal(t, amount.Amount(4242), parsed.Fee())
	require.Equal(t, amount.Amount(5), parsed.Payload().(*payload.AnchorPayload).Deposit)
}

func plantAnchor(t *testing.T, td *testData, number int32, root []byte) string {
	t.Helper()

	addr, acc := td.GenerateTestAccount(testsuite.AccountWithNumber(number))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        root,
		LockedDeposit:   9,
		CreatedAtHeight: 3,
		CreatedAtTime:   111,
	}))
	td.fake.FakeState.AddTestAccount(addr, acc)

	return addr.String()
}

func (td *testData) endpoint(path string) string {
	return "http://" + td.srv.listener.Addr().String() + "/http/api" + path
}

func getJSON(t *testing.T, td *testData, path string) (int, map[string]any) {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, td.endpoint(path), http.NoBody)
	require.NoError(t, err)

	return doJSON(t, req)
}

func postJSON(t *testing.T, td *testData, path string, payload map[string]any) (int, map[string]any) {
	t.Helper()

	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, td.endpoint(path), bytes.NewReader(raw))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	return doJSON(t, req)
}

func doJSON(t *testing.T, req *http.Request) (int, map[string]any) {
	t.Helper()

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, nil
	}

	var body map[string]any
	require.NoError(t, json.Unmarshal(raw, &body), string(raw))

	return resp.StatusCode, body
}

// getJSONExact decodes numbers as json.Number, so no digit is lost to float64.
func getJSONExact(t *testing.T, td *testData, path string) map[string]any {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, td.endpoint(path), http.NoBody)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(raw))

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var body map[string]any
	require.NoError(t, decoder.Decode(&body), string(raw))

	return body
}

func anyField(t *testing.T, body map[string]any, keys ...string) any {
	t.Helper()

	for _, key := range keys {
		if value, ok := body[key]; ok {
			return value
		}
	}
	require.Failf(t, "missing field", "%v in %v", keys, body)

	return nil
}

func objectField(t *testing.T, body map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := body[key].(map[string]any)
	require.Truef(t, ok, "missing %s in %v", key, body)

	return value
}

func stringField(t *testing.T, body map[string]any, keys ...string) string {
	t.Helper()

	for _, key := range keys {
		if value, ok := body[key].(string); ok {
			return value
		}
	}
	require.Failf(t, "missing string", "%v in %v", keys, body)

	return ""
}

func numberField(t *testing.T, body map[string]any, key string) int64 {
	t.Helper()

	switch value := body[key].(type) {
	case float64:
		return int64(value)
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		require.NoError(t, err)

		return parsed
	default:
		require.Failf(t, "missing number", "%s in %v", key, body)

		return 0
	}
}

func bytesFill(fill byte) []byte {
	return bytes.Repeat([]byte{fill}, 32)
}
