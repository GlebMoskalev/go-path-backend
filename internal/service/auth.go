package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/GlebMoskalev/go-path-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	log          *zap.Logger
	userRepo     repository.UserRepository
	stateRepo    repository.StateRepository
	oauth2Config *oauth2.Config
	jwtSecret    []byte
	accessTTl    time.Duration
	refreshTTL   time.Duration
	userInfoURL  string
}

func NewAuthService(
	log *zap.Logger, userRepo repository.UserRepository, stateRepo repository.StateRepository,
	clientID, clientSecret, redirectURL, jwtSecret, userInfoURL string,
	accessTTl, refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		log:       log,
		userRepo:  userRepo,
		stateRepo: stateRepo,
		jwtSecret: []byte(jwtSecret),
		oauth2Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "email", "profile"},
		},
		userInfoURL: userInfoURL,
		accessTTl:   accessTTl,
		refreshTTL:  refreshTTL,
	}
}

func (s *AuthService) GetGoogleLoginURL(ctx context.Context) (string, error) {
	state, err := generateState()
	if err != nil {
		s.log.Error("failed to generate state", zap.Error(err))
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	if err := s.stateRepo.Save(ctx, state, 10*time.Minute); err != nil {
		s.log.Error("failed to save state", zap.Error(err))
		return "", fmt.Errorf("failed to save state: %w", err)
	}
	return s.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code, state string) (*model.TokenPair, error) {
	valid, err := s.stateRepo.Validate(ctx, state)
	if err != nil {
		s.log.Error("failed to validate state", zap.Error(err))
		return nil, fmt.Errorf("validate state: %w", err)
	}
	if !valid {
		s.log.Warn("invalid or expired oauth state", zap.String("state", state))
		return nil, errors.New("invalid or expired state")
	}

	token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		s.log.Error("failed to exchange code", zap.Error(err))
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	userInfo, err := s.getGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		s.log.Error("failed to get google user info", zap.Error(err))
		return nil, fmt.Errorf("get google user info: %w", err)
	}

	user, err := s.userRepo.GetByGoogleID(ctx, userInfo.ID)
	if err != nil {
		if !errors.Is(err, repository.UserNotFound) {
			s.log.Error("failed to get user by google id", zap.Error(err))
			return nil, fmt.Errorf("get user by google id: %w", err)
		}
		user, err = s.createUserFromGoogle(ctx, userInfo)
		if err != nil {
			s.log.Error("failed to create user", zap.Error(err))
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	if err = s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.log.Error("failed to update last login", zap.String("user_id", user.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("update last login: %w", err)
	}

	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		s.log.Error("failed to generate token pair", zap.Error(err))
		return nil, fmt.Errorf("generate token pair: %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		s.log.Error("failed to parse refresh token", zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		err = errors.New("invalid refresh token")
		s.log.Error("failed to parse refresh token", zap.Error(err))
		return nil, err
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		err = errors.New("invalid token type")
		s.log.Error("failed to parse refresh token", zap.Error(err))
		return nil, err
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		err = errors.New("invalid user id in token")
		s.log.Error("failed to parse refresh token", zap.Error(err))
		return nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		err = fmt.Errorf("invalid user_id format: %w", err)
		s.log.Error("failed to parse refresh token", zap.Error(err))
		return nil, err
	}

	tokenVersion, ok := claims["token_version"].(float64)
	if !ok {
		err = errors.New("invalid token version in token")
		s.log.Error("failed to parse refresh token", zap.Error(err))
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repository.UserNotFound) {
			err = fmt.Errorf("failed to get user: %w", err)
			s.log.Error("failed to get user in repository", zap.Error(err))
			return nil, err
		}
		s.log.Warn("user not found")
		return nil, err
	}

	if int(tokenVersion) != user.TokenVersion {
		err = errors.New("token has been invalidated")
		s.log.Error(err.Error())
		return nil, err
	}

	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		err = fmt.Errorf("failed to generate tokens: %w", err)
		s.log.Error(err.Error())
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	err := s.userRepo.IncrementTokenVersion(ctx, userID)
	if err != nil {
		s.log.Error("failed to increment token version", zap.Error(err))
	}
	return err
}

// ValidateAccessToken проверяет JWT access token и возвращает user ID и token version
func (s *AuthService) ValidateAccessToken(ctx context.Context, tokenStr string) (uuid.UUID, int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		s.log.Error("failed to parse access token", zap.Error(err))
		return uuid.Nil, 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			err = errors.New("invalid token type")
			s.log.Error("invalid token", zap.Error(err))
			return uuid.Nil, 0, err
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			err = errors.New("invalid user id in token")
			s.log.Error("invalid token", zap.Error(err))
			return uuid.Nil, 0, err
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			err = fmt.Errorf("invalid user id format: %w", err)
			s.log.Error("invalid token", zap.Error(err))
			return uuid.Nil, 0, err
		}

		tokenVersion, ok := claims["token_version"].(float64)
		if !ok {
			err = errors.New("invalid token version in token")
			s.log.Error("invalid token", zap.Error(err))
			return uuid.Nil, 0, err
		}

		return userID, int(tokenVersion), nil
	}

	s.log.Error("access token invalid")
	return uuid.Nil, 0, errors.New("invalid token")
}

func (s *AuthService) generateTokenPair(user *model.User) (*model.TokenPair, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":          "access",
		"user_id":       user.ID.String(),
		"token_version": user.TokenVersion,
		"exp":           time.Now().Add(s.accessTTl).Unix(),
		"iat":           time.Now().Unix(),
	})
	accessTokenStr, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to signed token: %w", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":          "refresh",
		"user_id":       user.ID.String(),
		"token_version": user.TokenVersion,
		"exp":           time.Now().Add(s.refreshTTL).Unix(),
		"iat":           time.Now().Unix(),
	})
	refreshTokenStr, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to signed token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int(s.accessTTl.Seconds()),
	}, nil
}

func (s *AuthService) createUserFromGoogle(ctx context.Context, googleInfo *model.GoogleUserInfo) (*model.User, error) {
	if googleInfo == nil {
		return nil, errors.New("google info is nil")
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        googleInfo.Email,
		Name:         googleInfo.Name,
		Picture:      googleInfo.Picture,
		GoogleID:     googleInfo.ID,
		TokenVersion: 0,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) getGoogleUserInfo(ctx context.Context, token string) (*model.GoogleUserInfo, error) {
	resp, err := http.Get(s.userInfoURL + "?access_token=" + token)
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info from Google and status code equal %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var userInfo model.GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal userInfo: %w", err)
	}

	return &userInfo, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
