package certificate_test

// func TestVoteCertificateSignBytes(t *testing.T) {
// 	ts := testsuite.NewTestSuite(t)

// 	h := ts.RandHash()
// 	height := ts.RandHeight()
// 	round := ts.RandRound()
// 	cpRound := ts.RandRound()
// 	cpValue := ts.RandInt8(3)

// 	cert1 := certificate.NewBlockCertificate(height, round, true)
// 	cert2 := certificate.NewVoteCertificate(height, round)

// 	sb2 := cert2.SignBytes(h)
// 	sb3 := cert3.SignBytes(h)
// 	sb4 := cert4.SignBytes(h)
// 	sb5 := cert5.SignBytes(h)

// 	assert.NotEqual(t, sb2, sb3)
// 	assert.NotEqual(t, sb2, sb4)
// 	assert.NotEqual(t, sb3, sb4)
// 	assert.NotEqual(t, sb4, sb5)

// 	// BlockCertificate (fast path) has same sign bytes as Prepare certificate
// 	assert.Equal(t, cert1.SignBytes(h), cert2.SignBytes(h))

// 	assert.Contains(t, string(sb2), "PREPARE")
// 	assert.Contains(t, string(sb3), "PRE-VOTE")
// 	assert.Contains(t, string(sb4), "MAIN-VOTE")
// 	assert.Contains(t, string(sb5), "DECIDED")
// }
