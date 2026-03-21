package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var ErrUserNotFound = fmt.Errorf("membership: user not found")

type checkPermissionRequest struct {
	IdentityID     string `json:"identityId"`
	PermissionCode string `json:"permissionCode"`
	OrganizationID string `json:"organizationId"`
}

type checkPermissionResponse struct {
	Allowed bool `json:"allowed"`
}

type setOwnerRequest struct {
	IdentityID     string `json:"ownerIdentityId"`
	OrganizationID string `json:"organizationId"`
}

type UserProfile struct {
	IdentityID string `json:"identityId"`
	Email      string `json:"email"`
	Name       string `json:"name"`
}

type Membership struct {
	MembershipID string   `json:"membershipId"`
	Status       string   `json:"status"`
	JoinedAt     string   `json:"joinedAt"`
	RemovedAt    string   `json:"removedAt"`
	Department   string   `json:"department"`
	Title        string   `json:"title"`
	Roles        []string `json:"roles"`
}

// CheckPermission — POST /api/internal/authorization/check
func (c *Client) CheckPermission(ctx context.Context, identityID, organizationID, permissionCode string) (bool, error) {
	body, err := json.Marshal(checkPermissionRequest{
		IdentityID:     identityID,
		OrganizationID: organizationID,
		PermissionCode: permissionCode,
	})
	if err != nil {
		return false, fmt.Errorf("membership.CheckPermission - marshal: %w", err)
	}

	resp, err := c.post(ctx, "/api/internal/authorization/check", body)
	if err != nil {
		return false, fmt.Errorf("membership.CheckPermission - request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("membership.CheckPermission - unexpected status: %d", resp.StatusCode)
	}

	var result checkPermissionResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("membership.CheckPermission - decode: %w", err)
	}

	return result.Allowed, nil
}

// GetUserByIdentityID — GET /api/internal/users/by-identity/{identityId}
func (c *Client) GetUserByIdentityID(ctx context.Context, identityID string) (*UserProfile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+"/api/internal/users/by-identity/"+identityID, nil)
	if err != nil {
		return nil, fmt.Errorf("membership.GetUserByIdentityID - new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("membership.GetUserByIdentityID - do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrUserNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("membership.GetUserByIdentityID - unexpected status: %d", resp.StatusCode)
	}

	var user UserProfile
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("membership.GetUserByIdentityID - decode: %w", err)
	}

	return &user, nil
}

// SetOrganizationOwner — POST /api/internal/organizations/owner
func (c *Client) SetOrganizationOwner(ctx context.Context, organizationID, identityID string) error {
	body, err := json.Marshal(setOwnerRequest{
		IdentityID:     identityID,
		OrganizationID: organizationID,
	})
	if err != nil {
		return fmt.Errorf("membership.SetOrganizationOwner - marshal: %w", err)
	}

	resp, err := c.post(ctx, "/api/internal/organizations/owner", body)
	if err != nil {
		return fmt.Errorf("membership.SetOrganizationOwner - request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("membership.SetOrganizationOwner - unexpected status: %d. resp: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetMembership — GET /api/internal/organizations/{organizationId}/users/{identityId}/membership
func (c *Client) GetMembership(ctx context.Context, organizationID, identityID string) (*Membership, error) {
	url := fmt.Sprintf("%s/api/internal/organizations/%s/users/%s/membership", c.baseURL, organizationID, identityID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("membership.GetMembership - new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("membership.GetMembership - do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrUserNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("membership.GetMembership - unexpected status: %d", resp.StatusCode)
	}

	var membership Membership
	if err = json.NewDecoder(resp.Body).Decode(&membership); err != nil {
		return nil, fmt.Errorf("membership.GetMembership - decode: %w", err)
	}

	return &membership, nil
}

// RevokeRole — DELETE /api/Organizations/{organizationId}/members/{membershipId}/roles/{roleCode}
func (c *Client) RevokeRole(ctx context.Context, organizationID, membershipID, roleCode, token string) error {
	url := fmt.Sprintf("%s/api/Organizations/%s/members/%s/roles/%s", c.baseURL, organizationID, membershipID, roleCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("membership.RevokeRole - new request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("membership.RevokeRole - do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("membership.RevokeRole - unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}
