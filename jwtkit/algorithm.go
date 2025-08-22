package jwtkit

import (
	"slices"
	
	"github.com/lestrrat-go/jwx/v3/jwa"
)

type Algorithm string

const (
	secretType = iota
	rsaType
	eddsaType
	ecdsaType
	noneType
	
	// AlgorithmHS256 means HMAC using SHA-256.
	// This is the most common and recommended symmetric algorithm for general use cases.
	AlgorithmHS256 Algorithm = "HS256"
	
	// AlgorithmHS384 means HMAC using SHA-384.
	// A stronger, but slightly slower alternative to HS256.
	AlgorithmHS384 Algorithm = "HS384"
	
	// AlgorithmHS512 means HMAC using SHA-512.
	// The strongest symmetric option, but also the slowest.
	AlgorithmHS512 Algorithm = "HS512"
	
	// AlgorithmRS256 means RSASSA-PKCS1-v1_5 using SHA-256.
	// This is the most common and widely supported asymmetric algorithm.
	AlgorithmRS256 Algorithm = "RS256"
	
	// AlgorithmRS384 means RSASSA-PKCS1-v1_5 using SHA-384.
	AlgorithmRS384 Algorithm = "RS384"
	
	// AlgorithmRS512 means RSASSA-PKCS1-v1_5 using SHA-512.
	AlgorithmRS512 Algorithm = "RS512"
	
	// AlgorithmES256 means ECDSA using P-256 and SHA-256.
	// A modern and more efficient alternative to RS256 that produces a smaller token size.
	AlgorithmES256 Algorithm = "ES256"
	
	// AlgorithmES256K means ECDSA using secp256k1 and SHA-256.
	// It is a modern and efficient asymmetric algorithm,
	// particularly common in blockchain-related applications.
	AlgorithmES256K Algorithm = "ES256K"
	
	// AlgorithmES384 means ECDSA using P-384 and SHA-384.
	AlgorithmES384 Algorithm = "ES384"
	
	// AlgorithmES512 means ECDSA using P-512 and SHA-512.
	AlgorithmES512 Algorithm = "ES512"
	
	// AlgorithmEdDSA means EdDSA using Ed25519.
	// A newer and highly efficient asymmetric algorithm.
	AlgorithmEdDSA Algorithm = "EdDSA"
	
	// AlgorithmPS256 means RSASSA-PSS using SHA-256 and MGF1 with SHA-256.
	AlgorithmPS256 Algorithm = "PS256"
	
	// AlgorithmPS384 means RSASSA-PSS using SHA-384 and MGF1 with SHA-384.
	AlgorithmPS384 Algorithm = "PS384"
	
	// AlgorithmPS512 means RSASSA-PSS using SHA-512 and MGF1 with SHA-512.
	AlgorithmPS512 Algorithm = "PS512"
	
	// AlgorithmNone means no signature.
	// This algorithm is used for UNSIGNED tokens.
	// It is highly insecure and SHOULD NEVER be used in production.
	AlgorithmNone Algorithm = "none"
)

var (
	symmetricAlgorithms = []Algorithm{
		AlgorithmHS256, AlgorithmHS384, AlgorithmHS512,
	}
	
	asymmetricAlgorithms = []Algorithm{
		AlgorithmRS256, AlgorithmRS384, AlgorithmRS512,
		AlgorithmES256, AlgorithmES256K, AlgorithmES384, AlgorithmES512,
		AlgorithmEdDSA,
		AlgorithmPS256, AlgorithmPS384, AlgorithmPS512,
	}
	
	rsaAlgorithms = []Algorithm{
		AlgorithmRS256, AlgorithmRS384, AlgorithmRS512,
		AlgorithmPS256, AlgorithmPS384, AlgorithmPS512,
	}
	
	ecdsaAlgorithms = []Algorithm{
		AlgorithmES256, AlgorithmES256K, AlgorithmES384, AlgorithmES512,
	}
	
	eddsaAlgorithms = []Algorithm{
		AlgorithmEdDSA,
	}
)

// IsSymmetric returns true if the algorithm is symmetric
func (algorithm Algorithm) IsSymmetric() bool {
	return slices.Contains(symmetricAlgorithms, algorithm)
}

// IsRSA returns true if the algorithm is RSA-based
func (algorithm Algorithm) IsRSA() bool {
	return slices.Contains(rsaAlgorithms, algorithm)
}

// IsECDSA returns true if the algorithm is ECDSA-based
func (algorithm Algorithm) IsECDSA() bool {
	return slices.Contains(ecdsaAlgorithms, algorithm)
}

// IsEdDSA returns true if the algorithm is EdDSA-based
func (algorithm Algorithm) IsEdDSA() bool {
	return slices.Contains(eddsaAlgorithms, algorithm)
}

// IsAsymmetric returns true if the algorithm is asymmetric
func (algorithm Algorithm) IsAsymmetric() bool {
	return slices.Contains(asymmetricAlgorithms, algorithm)
}

func (algorithm Algorithm) toJwaAlgo() jwa.KeyAlgorithm {
	switch algorithm {
	case AlgorithmHS256:
		return jwa.HS256()
	case AlgorithmHS384:
		return jwa.HS384()
	case AlgorithmHS512:
		return jwa.HS512()
	case AlgorithmRS256:
		return jwa.RS256()
	case AlgorithmRS384:
		return jwa.RS384()
	case AlgorithmRS512:
		return jwa.RS512()
	case AlgorithmES256:
		return jwa.ES256()
	case AlgorithmES256K:
		return jwa.ES256K()
	case AlgorithmES384:
		return jwa.ES384()
	case AlgorithmES512:
		return jwa.ES512()
	case AlgorithmEdDSA:
		return jwa.EdDSA()
	case AlgorithmPS256:
		return jwa.PS256()
	case AlgorithmPS384:
		return jwa.PS384()
	case AlgorithmPS512:
		return jwa.PS512()
	default:
		return jwa.NoSignature()
	}
}
