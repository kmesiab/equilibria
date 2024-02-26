package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/jwt"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
)

const (
	TokenTypeString   = "TOKEN"
	PrincipleID       = "user"
	PolicyEffectDeny  = "Deny"
	PolicyEffectAllow = "Allow"
)

type AuthorizerLambdaHandler struct {
	lib.LambdaHandler
	KeyRotator   jwt.KeyRotatorInterface
	TokenService jwt.TokenServiceInterface
}

func (h AuthorizerLambdaHandler) HandleRequest(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	publicKey, err := h.TokenService.GetPublicKey(h.KeyRotator)

	if err != nil {
		log.New("Error getting public key from remote").
			AddError(err).Log()

		return generatePolicy(PrincipleID, PolicyEffectDeny, request.MethodArn), nil
	}

	token := strings.TrimPrefix(request.AuthorizationToken, "Bearer ")
	claims, err := h.TokenService.Validate(token, publicKey)

	if err != nil {
		log.New("Error validating JWT: %s", err).Add("token", token).Log()

		return generatePolicy(PrincipleID, PolicyEffectDeny, request.MethodArn), nil
	}

	// Log this JWT
	log.New("Authorized %s", claims.PhoneNumber).
		Add("user_id", strconv.FormatInt(claims.UserID, 10)).
		Add("phone_number", claims.PhoneNumber).
		Add("lastname", claims.Lastname).
		Add("firstname", claims.Firstname).
		Add("email", claims.Email).
		Add("account_status", claims.AccountStatus).
		Add("phone_verified", log.FormatBool(claims.PhoneVerified)).
		Add("jwt_issued_at", strconv.FormatInt(claims.IssuedAt, 10)).
		Add("jwt_expires_at", strconv.FormatInt(claims.ExpiresAt, 10)).
		Add("jwt_issuer", claims.Issuer).
		Add("token", token).
		Log()

	// If the token is invalid, deny access
	return generatePolicy(PrincipleID, PolicyEffectAllow, request.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	log.New("Issuing %s policy for %s resource", effect, resource).Log()

	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}
	if effect != "" && resource != "" {
		policyDocument := events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17", // Default version
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
		authResponse.PolicyDocument = policyDocument
	}
	return authResponse
}

func main() {
	log.New("Authorizer Lambda booting...").Log()

	cfg := config.Get()
	if cfg == nil {
		log.New("Could not load config").Log()
		os.Exit(1)
	}

	handler := &AuthorizerLambdaHandler{
		KeyRotator:   &jwt.KeyRotator{},
		TokenService: &jwt.TokenService{},
	}

	log.New("Authorizer Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}
