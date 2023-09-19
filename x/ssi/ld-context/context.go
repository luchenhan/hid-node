package ldcontext

const didContext string = "https://www.w3.org/ns/did/v1"
const ed25519Context2020 string = "https://w3id.org/security/suites/ed25519-2020/v1"

// As hid-node is not supposed to perform any GET request, the complete Context body of their
// respective Context urls has been maintained below.
var ContextUrlMap map[string]contextObject = map[string]contextObject{
	didContext: {
		"@protected": true,
		"id":         "@id",
		"type":       "@type",
		"alsoKnownAs": map[string]interface{}{
			"@id":   "https://www.w3.org/ns/activitystreams#alsoKnownAs",
			"@type": "@id",
		},
		"assertionMethod": map[string]interface{}{
			"@id":        "https://w3id.org/security#assertionMethod",
			"@type":      "@id",
			"@container": "@set",
		},
		"authentication": map[string]interface{}{
			"@id":        "https://w3id.org/security#authenticationMethod",
			"@type":      "@id",
			"@container": "@set",
		},
		"capabilityDelegation": map[string]interface{}{
			"@id":        "https://w3id.org/security#capabilityDelegationMethod",
			"@type":      "@id",
			"@container": "@set",
		},
		"capabilityInvocation": map[string]interface{}{
			"@id":        "https://w3id.org/security#capabilityInvocationMethod",
			"@type":      "@id",
			"@container": "@set",
		},
		"controller": map[string]interface{}{
			"@id":   "https://w3id.org/security#controller",
			"@type": "@id",
		},
		"keyAgreement": map[string]interface{}{
			"@id":        "https://w3id.org/security#keyAgreementMethod",
			"@type":      "@id",
			"@container": "@set",
		},
		"service": map[string]interface{}{
			"@id":   "https://www.w3.org/ns/did#service",
			"@type": "@id",
			"@context": map[string]interface{}{
				"@protected": true,
				"id":         "@id",
				"type":       "@type",
				"serviceEndpoint": map[string]interface{}{
					"@id":   "https://www.w3.org/ns/did#serviceEndpoint",
					"@type": "@id",
				},
			},
		},
		"verificationMethod": map[string]interface{}{
			"@id":   "https://w3id.org/security#verificationMethod",
			"@type": "@id",
		},
	},
	ed25519Context2020: {
		"id":         "@id",
		"type":       "@type",
		"@protected": true,
		"proof": map[string]interface{}{
			"@id":        "https://w3id.org/security#proof",
			"@type":      "@id",
			"@container": "@graph",
		},
		"Ed25519VerificationKey2020": map[string]interface{}{
			"@id": "https://w3id.org/security#Ed25519VerificationKey2020",
			"@context": map[string]interface{}{
				"@protected": true,
				"id":         "@id",
				"type":       "@type",
				"controller": map[string]interface{}{
					"@id":   "https://w3id.org/security#controller",
					"@type": "@id",
				},
				"revoked": map[string]interface{}{
					"@id":   "https://w3id.org/security#revoked",
					"@type": "http://www.w3.org/2001/XMLSchema#dateTime",
				},
				"publicKeyMultibase": map[string]interface{}{
					"@id":   "https://w3id.org/security#publicKeyMultibase",
					"@type": "https://w3id.org/security#multibase",
				},
			},
		},
		"Ed25519Signature2020": map[string]interface{}{
			"@id": "https://w3id.org/security#Ed25519Signature2020",
			"@context": map[string]interface{}{
				"@protected": true,
				"id":         "@id",
				"type":       "@type",
				"challenge":  "https://w3id.org/security#challenge",
				"created": map[string]interface{}{
					"@id":   "http://purl.org/dc/terms/created",
					"@type": "http://www.w3.org/2001/XMLSchema#dateTime",
				},
				"domain": "https://w3id.org/security#domain",
				"expires": map[string]interface{}{
					"@id":   "https://w3id.org/security#expiration",
					"@type": "http://www.w3.org/2001/XMLSchema#dateTime",
				},
				"nonce": "https://w3id.org/security#nonce",
				"proofPurpose": map[string]interface{}{
					"@id":   "https://w3id.org/security#proofPurpose",
					"@type": "@vocab",
					"@context": map[string]interface{}{
						"@protected": true,
						"id":         "@id",
						"type":       "@type",
						"assertionMethod": map[string]interface{}{
							"@id":        "https://w3id.org/security#assertionMethod",
							"@type":      "@id",
							"@container": "@set",
						},
						"authentication": map[string]interface{}{
							"@id":        "https://w3id.org/security#authenticationMethod",
							"@type":      "@id",
							"@container": "@set",
						},
						"capabilityInvocation": map[string]interface{}{
							"@id":        "https://w3id.org/security#capabilityInvocationMethod",
							"@type":      "@id",
							"@container": "@set",
						},
						"capabilityDelegation": map[string]interface{}{
							"@id":        "https://w3id.org/security#capabilityDelegationMethod",
							"@type":      "@id",
							"@container": "@set",
						},
						"keyAgreement": map[string]interface{}{
							"@id":        "https://w3id.org/security#keyAgreementMethod",
							"@type":      "@id",
							"@container": "@set",
						},
					},
				},
				"proofValue": map[string]interface{}{
					"@id":   "https://w3id.org/security#proofValue",
					"@type": "https://w3id.org/security#multibase",
				},
				"verificationMethod": map[string]interface{}{
					"@id":   "https://w3id.org/security#verificationMethod",
					"@type": "@id",
				},
			},
		},
	},
}
