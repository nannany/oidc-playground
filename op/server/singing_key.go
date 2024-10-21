package server

import (
	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var MySigningKey SigningKey

type SigningKey struct {
	privateKey jose.JSONWebKey
}

// init処理をする
func init() {
	var jwk jose.JSONWebKey
	_ = jwk.UnmarshalJSON([]byte(`
{
  "kty": "RSA",
  "n": "5srSQZgeolXrjpTvw1OuHxBrHiBKnBEOxeeOgBDaB_61Dm1nr39rnbjd7CmuVel9o1CQof26741AoqxFxDAzc1KtnG2pysT32kcKVLBYQYSyXl860jrMXBgs-eR2Gz_YJl5UmmMvexYmnJ1CAhDUxMK23MeeR0_llTUIRDPrE1JFgE033gvUF8PfNxSUzeI5FHu6PjbLrwiatg3sOhUAkxQhC5IPGJoSVuS0_taU72lRoSEKT2Ij32HnhLWx7dAZ_PXcSZGU3L86AGksenF-bfDes6_OXIWkCBtlcpXGo51WNWzmVgX1KBVe48SCWwO9qIr8F6oRNe0zxcIvaSWHKpfMw711uF8OT8XpF9jOvlMxXGOASpAJ8eDVh4DK4YHfG4GFg4mlzQ6wr7_MHl8yXLj5v_-03XS3-AzskLs86haHi91U6zoA2zGkQ6f_KsBa5Mi7Yn9XkjT3LqIdE2Eq6PzLkXa0_BPyoA4yu1AQiZ0UneNCZpxqD_1UVzU2ZmoyvNprAd1Y5RK7pimWx8NAEkcZfLg3OjsQvxho4l0YeyqZPrnmYy2G61BvCWkgzpjoHIxn9IJgdXsS80ugJKOWF-hfKUYwyW5iWuO285WvZbF_jSoqfGvKk21bsyf5_4Pj0i_5lY5OmTrYnDHhGKNcO_FrZXKHvEVTLFC7h1FJOo8",
  "e": "AQAB",
  "d": "v6HvXXnDDgYLtnNidix4auDGINitOtNLEfbIMvxLHdJXYihjrKw_nlHPK3YchMeAg-e0gF1fJ-OApNIoh9OCFSK0dyvzjSNprbvJDzLxxU06kfO-sqsR3vJP_hq-GrgzcQBofV0YIrrsyqQlF_QDx3DqR34UpZ_owLCtMXxul7d7cCIt3lix_h2srV6QhDdwL9rgdSIolCiv1bTPzhZ8SKOxIJthn4rm6XGz-9oPBo_LsN_vYKTbpnLzuxTnLdDak1lOVWPlWDoMHiU2QmEUmpMjR84VqmtpHTYqUS-eJWLSH0BSLWF7A5JjHPqhRpNfmp_03G1hJwbTM0Se-06m4JpLcVYTbNZD15qQDHl9mcPHziWcxgkoSDo3pryh8jLFne0h8k8-03bgt-nr7zN3q8V_UYdEmTPar0Gr-jkc87FX4qvuWNgpNz6Ptiz0zW2xIh1wnwMDnY58IRqbT_Ulhctm5diis2IQ7q3CZtLv_amUqfCYqgCPxm3PNcSuvMoFkyIGz06lOZiy3fiMit8HeoNQWjkNO1oc0M9_7i41NjT2StenNZLga04yJEh1-A1ggqQuTMhM6ih3fqnVGHwnjmbFJFtTKLgG4SG4CWhjJmnWyRJbgfjfaHm0qw_sExWucvQXtzGNs8kSAY4i_yeD0Wjd29MqhdkDF77onXFgaDE",
  "p": "8-8SuVznEBOGzgeeKZtvu3bvCJpavn3K4ZUujDeslxO27jdUpFbyL4zMTSzs8QVK_kG6C_zV2ne-C3NoC3atEV163shgXs4sH49QhFOsjFkkAuWYi_A1VbOfK1XJPLLizmv-O1YdTRh22D7JyUDebytaIUCnXG3khYU_5-GWK7_O9qPNPbDHSn82ea-N5j8wUjtJuTbOIj7VSQvMQRH8hV27LPOyaKKFUvte5uPdWOkbH7ySYpmnQBXnbMjqQyvhtAjHRpXbt21F2eH8Yw6x6oXne3Rh-OKdoDLiw-JlpisDcW5ikk7tzuR_U6LcdEa86m6Kqfn3u4YbB614L8bsyQ",
  "q": "8jVWHp8daJZWSD90h89dgjK6jqWoSJiIVDpMPiJ8LlFbbS-neArvHIXCkf8Ms33QjbwdsA81q7Rsz82Q5MTMuqrP3NPiyLWYqiDu4ibScVcf-SlXgOoLQTrllY2DF4eKIX-yGFWRy3wLs5CmPri-9qH_8nzXA_xrTw-fmUu3tPTBleH8eDsxpjYqtroXBQzrLs7iQc_DOv8qWS47X1pi9FbrAWyfg07ZVt8F76gjsGaT69zpWYF_ft68N_GE_fv5r3MINmnzYyejxuC-r7XyxPX0ZC02eKOUc9TU-s6Knsik9QlYvGtIDSy8oRoalSV1aIN8qEx3OTSQ06GPQW8Qlw",
  "dp": "ivewWxmqKWZ2bfm5CUscJFhlZSlKeSuA4XLzyb4N_SOmG6A6AEXoQ16bJXxqoAS77I3VR-8KhiOhiTR-GcnKXxI6ZaESBfC_AlvLKxfgPTSrZ1sVxONb_y8Nhsqgkov22lJ7y6ILn1hInHloy9bA4eR4vGjw759LiSWYeqnxu4rShYBb9ME5SB-hEUIKPRnSlYZhQbsPREM3jx3Bh3CPxAraD4nmIeY0vkmmjzNRbs2ePj5XxwRyW_gel8L-crYVJ1O482V8fssp-C7ecjB0-369mX98MSYpVpKmzaG09y2aEI0qat-8axmR0DwAC94g0g2Xwa4-i_6id4VD6zQQwQ",
  "dq": "mP8kFVfBRe2hNyYOQDO3B6VvvufZs3HWvA3PV3iFJOTzcbcmfGe7vzKnQ33u0frpoH5x6mLRLlNcYK-jlT7TbB62CvI1UT-U3fLLs0N-r2Wvpr58rcwpq9y-ZYfshRJzKNH29_QlknobEiSPPxOIXVbbzJsbX0M1rc1arYTs2Hu7RShpSLFrnWid9qr6G5CbhwfcWGIbGstQPuE7U_JOi7XT4SRyZomKNJvQriBx2t4RFm6HYKiylruv_U6tCZr_j9qEF4s52SqTA--3xFqUreIHrLf6rp32Cm0o0_1gBwkLWlW7HdpOuzPFO6a3n_r3fTuBpPYYYaRYz7ZAPb8gpQ",
  "qi": "Hr7GYLltqdPHn8RUJb7iNuORlsWY-Fi8HJoI0uSaLejoCNWTWK3C0TkecXOUCU98BCIFmeKCIMKSUyFx7dlKPJLz_BQEfO7M-dl_2JGU-MfxEebkXdXKGbf0C_IflVScOMKW4BclNCp7eRiyWzyF6NBGLHuaUU4u0VQoGIsnvQiXxw0HMX13JHcaxQiSPWqFHXZSWG6guBFlcWAz-aZUHlFp9a0v7IY6Hw812PEyga8cfs8tqjRLdOvsYPC2G2aULbNvOh94LGWZB3EqWIblxenNIG9UHuC8pDWIpWM4-cOUb2QSBJ5eeUocUvry3d1xZED3i97KQWWVhmsMozTGcA",
  "alg": "RS256",
  "use": "sig",
  "kid": "avUja_OmJ6soJ6KUnmM_IWoPLxny3Ph-uWLZnFxrpuE"
}
    `))

	MySigningKey = SigningKey{
		privateKey: jwk,
	}

}

func (s SigningKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(s.privateKey.Algorithm)
}

func (s SigningKey) Key() any {
	return s.privateKey.Key

}

func (s SigningKey) ID() string {
	return s.privateKey.KeyID
}

// SigningKey は op.SigningKey インターフェースを実装しています。
var _ op.SigningKey = (*SigningKey)(nil)
