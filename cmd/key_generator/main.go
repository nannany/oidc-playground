package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
)

type JWK struct {
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
	D   string `json:"d,omitempty"`  // 公開鍵にはDは不要
	P   string `json:"p,omitempty"`  // 公開鍵にはPは不要
	Q   string `json:"q,omitempty"`  // 公開鍵にはQは不要
	Dp  string `json:"dp,omitempty"` // 公開鍵にはDpは不要
	Dq  string `json:"dq,omitempty"` // 公開鍵にはDqは不要
	Qi  string `json:"qi,omitempty"` // 公開鍵にはQiは不要
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid,omitempty"`
}

func main() {
	// RSA秘密鍵を生成
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("秘密鍵の生成に失敗しました: %s", err)
	}

	// 公開指数 (e) をバイト配列に変換
	eBytes := big.NewInt(int64(privateKey.PublicKey.E)).Bytes()

	// 各RSAコンポーネントをBase64URLエンコーディングでエンコード
	n := base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(eBytes)
	d := base64.RawURLEncoding.EncodeToString(privateKey.D.Bytes())
	p := base64.RawURLEncoding.EncodeToString(privateKey.Primes[0].Bytes())
	q := base64.RawURLEncoding.EncodeToString(privateKey.Primes[1].Bytes())
	dp := base64.RawURLEncoding.EncodeToString(privateKey.Precomputed.Dp.Bytes())
	dq := base64.RawURLEncoding.EncodeToString(privateKey.Precomputed.Dq.Bytes())
	qi := base64.RawURLEncoding.EncodeToString(privateKey.Precomputed.Qinv.Bytes())

	// 公開鍵のSHA-256ハッシュを計算してkidとして使用
	publicKeyHash := sha256.Sum256(privateKey.PublicKey.N.Bytes())
	kid := base64.RawURLEncoding.EncodeToString(publicKeyHash[:])

	// 秘密鍵JWKオブジェクトを作成
	privateJWK := JWK{
		Kty: "RSA",
		N:   n,
		E:   e,
		D:   d,
		P:   p,
		Q:   q,
		Dp:  dp,
		Dq:  dq,
		Qi:  qi,
		Alg: "RS256",
		Use: "sig",
		Kid: kid, // kidを追加
	}

	// 公開鍵JWKオブジェクトを作成（DやP, Qなどは含めない）
	publicJWK := JWK{
		Kty: "RSA",
		N:   n,
		E:   e,
		Alg: "RS256",
		Use: "sig",
		Kid: kid, // kidを追加
	}

	// 秘密JWKをJSONに変換して表示
	privateKeyJSON, err := json.MarshalIndent(privateJWK, "", "  ")
	if err != nil {
		log.Fatalf("秘密JWKのJSON変換に失敗しました: %s", err)
	}
	fmt.Printf("Private JWK:\n%s\n\n", privateKeyJSON)

	// 公開JWKをJSONに変換して表示
	publicKeyJSON, err := json.MarshalIndent(publicJWK, "", "  ")
	if err != nil {
		log.Fatalf("公開JWKのJSON変換に失敗しました: %s", err)
	}
	fmt.Printf("Public JWK:\n%s\n", publicKeyJSON)
}
