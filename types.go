package run9

import "time"

// MembershipRole represents one organization membership role.
type MembershipRole string

// OrgKind represents one organization kind.
type OrgKind string

// ProjectRole represents one project membership role.
type ProjectRole string

// ProjectSecretScope represents where a secret is attached.
type ProjectSecretScope string

// Organization kinds returned by the control plane.
const (
	OrgKindPersonal OrgKind = "personal"
	OrgKindShared   OrgKind = "shared"
)

// Secret scopes accepted by project and box secret APIs.
const (
	ProjectSecretScopeProject ProjectSecretScope = "project"
	ProjectSecretScopeBox     ProjectSecretScope = "box"
)

// InvitationState represents one invitation lifecycle state.
type InvitationState string

// BoxState represents the current box state.
type BoxState string

// BoxNetworkMode represents the requested networking mode for a box.
type BoxNetworkMode string

// Box network modes accepted by box create and update APIs.
const (
	BoxNetworkModeNormal  BoxNetworkMode = "normal"
	BoxNetworkModeManaged BoxNetworkMode = "managed"
)

// BoxSecurityMode represents the requested security mode for a box.
type BoxSecurityMode string

// Box security modes accepted by box create and update APIs.
const (
	BoxSecurityModeRestricted BoxSecurityMode = "restricted"
	BoxSecurityModeUnsafe     BoxSecurityMode = "unsafe"
)

// SnapState represents the current snap state.
type SnapState string

// CurrentSubscriptionView describes the caller's current subscription snapshot.
type CurrentSubscriptionView struct {
	Tier      string    `json:"tier"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

// MeView describes one user account.
type MeView struct {
	UserID          string    `json:"user_id"`
	PrimaryEmail    string    `json:"primary_email"`
	DisplayName     string    `json:"display_name,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	IsSystemManager bool      `json:"is_system_manager"`
}

// OrgView describes one organization visible to the caller.
type OrgView struct {
	OrgID               string                  `json:"org_id"`
	OrgCID              string                  `json:"org_cid"`
	DisplayName         string                  `json:"display_name"`
	Kind                OrgKind                 `json:"kind"`
	Role                MembershipRole          `json:"role"`
	CreatedBy           string                  `json:"created_by"`
	CreatedAt           time.Time               `json:"created_at"`
	CurrentSubscription CurrentSubscriptionView `json:"current_subscription"`
}

// CurrentOrgIdentityView describes the current authenticated user and org.
type CurrentOrgIdentityView struct {
	User     MeView  `json:"user"`
	Org      OrgView `json:"org"`
	AuthKind string  `json:"auth_kind"`
}

// DeleteOrgResult describes the accepted result of deleting an organization.
type DeleteOrgResult struct {
	OrgID  string `json:"org_id"`
	Status string `json:"status"`
}

// MembershipView describes one organization member.
type MembershipView struct {
	OrgID        string         `json:"org_id"`
	UserID       string         `json:"user_id"`
	PrimaryEmail string         `json:"primary_email"`
	DisplayName  string         `json:"display_name,omitempty"`
	Role         MembershipRole `json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	LastLoginAt  time.Time      `json:"last_login_at,omitempty"`
}

// InvitationView describes one organization invitation.
type InvitationView struct {
	InvitationID string          `json:"invitation_id"`
	OrgID        string          `json:"org_id"`
	InviteeEmail string          `json:"invitee_email"`
	Role         MembershipRole  `json:"role"`
	InvitedBy    string          `json:"invited_by"`
	State        InvitationState `json:"state"`
	ExpiresAt    time.Time       `json:"expires_at"`
	AcceptedBy   string          `json:"accepted_by,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// DeleteInvitationResult describes the accepted result of revoking an invitation.
type DeleteInvitationResult struct {
	InvitationID string `json:"invitation_id"`
	Status       string `json:"status"`
}

// APIKeyView describes one API key without secret key material.
type APIKeyView struct {
	APIKeyID          string    `json:"api_key_id"`
	OrgID             string    `json:"org_id"`
	UserID            string    `json:"user_id"`
	OwnerPrimaryEmail string    `json:"owner_primary_email"`
	OwnerDisplayName  string    `json:"owner_display_name,omitempty"`
	Description       string    `json:"description,omitempty"`
	AK                string    `json:"ak"`
	DisplayPrefix     string    `json:"display_prefix"`
	DisplaySuffix     string    `json:"display_suffix"`
	CreatedAt         time.Time `json:"created_at"`
	ExpiresAt         time.Time `json:"expires_at,omitempty"`
	NoExpire          bool      `json:"no_expire"`
}

// CreatedAPIKeyView describes a newly created API key, including its secret key.
type CreatedAPIKeyView struct {
	APIKeyID          string    `json:"api_key_id"`
	OrgID             string    `json:"org_id"`
	UserID            string    `json:"user_id"`
	OwnerPrimaryEmail string    `json:"owner_primary_email"`
	OwnerDisplayName  string    `json:"owner_display_name,omitempty"`
	Description       string    `json:"description,omitempty"`
	AK                string    `json:"ak"`
	SK                string    `json:"sk"`
	DisplayPrefix     string    `json:"display_prefix"`
	DisplaySuffix     string    `json:"display_suffix"`
	CreatedAt         time.Time `json:"created_at"`
	ExpiresAt         time.Time `json:"expires_at,omitempty"`
	NoExpire          bool      `json:"no_expire"`
}

// SSHKeyView describes one SSH key on the caller's account.
type SSHKeyView struct {
	SSHKeyID    string     `json:"ssh_key_id"`
	Label       string     `json:"label"`
	Fingerprint string     `json:"fingerprint"`
	CreatedAt   time.Time  `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
}

// ProjectView describes one project visible to the caller.
type ProjectView struct {
	ProjectID   string      `json:"project_id"`
	OrgID       string      `json:"org_id"`
	ProjectCID  string      `json:"project_cid"`
	DisplayName string      `json:"display_name"`
	Description string      `json:"description,omitempty"`
	Role        ProjectRole `json:"role"`
	CreatedBy   string      `json:"created_by"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// DeleteProjectResult describes the accepted result of deleting a project.
type DeleteProjectResult struct {
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
}

// ProjectMembershipView describes one project member.
type ProjectMembershipView struct {
	ProjectID    string      `json:"project_id"`
	UserID       string      `json:"user_id"`
	PrimaryEmail string      `json:"primary_email"`
	DisplayName  string      `json:"display_name,omitempty"`
	Role         ProjectRole `json:"role"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// ProjectSecretView describes one secret attached to project or box scope.
type ProjectSecretView struct {
	SecretID         string             `json:"secret_id"`
	OrgID            string             `json:"org_id"`
	ProjectID        string             `json:"project_id"`
	Scope            ProjectSecretScope `json:"scope"`
	BoxID            string             `json:"box_id,omitempty"`
	Name             string             `json:"name,omitempty"`
	Placeholder      string             `json:"placeholder"`
	AllowedHosts     []string           `json:"allowed_hosts"`
	InjectHeaderName string             `json:"inject_header_name"`
	CreatedBy        string             `json:"created_by"`
	UpdatedBy        string             `json:"updated_by"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// CreateProjectRequest creates a new project.
type CreateProjectRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description,omitempty"`
}

// CreateProjectSecretRequest creates a new project or box secret.
type CreateProjectSecretRequest struct {
	Name             string   `json:"name,omitempty"`
	Value            string   `json:"value"`
	Placeholder      string   `json:"placeholder"`
	AllowedHosts     []string `json:"allowed_hosts"`
	InjectHeaderName string   `json:"inject_header_name"`
}

// UpdateProjectSecretRequest updates one existing project or box secret.
type UpdateProjectSecretRequest struct {
	Name             *string   `json:"name,omitempty"`
	Value            *string   `json:"value,omitempty"`
	AllowedHosts     *[]string `json:"allowed_hosts,omitempty"`
	InjectHeaderName *string   `json:"inject_header_name,omitempty"`
}

// BoxView describes one box.
type BoxView struct {
	BoxID                      string            `json:"box_id"`
	OrgID                      string            `json:"org_id"`
	ProjectID                  string            `json:"project_id"`
	Creator                    string            `json:"creator"`
	CreatedAt                  time.Time         `json:"created_at"`
	LastUsedAt                 time.Time         `json:"last_used_at"`
	Description                string            `json:"description,omitempty"`
	Labels                     map[string]string `json:"labels,omitempty"`
	State                      BoxState          `json:"state"`
	Reason                     string            `json:"reason,omitempty"`
	BoxSnapID                  string            `json:"box_snap_id"`
	DesiredShape               string            `json:"desired_shape"`
	NetworkMode                BoxNetworkMode    `json:"network_mode"`
	SecurityMode               BoxSecurityMode   `json:"security_mode"`
	CurrentHostID              string            `json:"current_host_id,omitempty"`
	CurrentRuntimeShape        string            `json:"current_runtime_shape,omitempty"`
	CurrentRuntimeNetworkMode  BoxNetworkMode    `json:"current_runtime_network_mode,omitempty"`
	CurrentRuntimeSecurityMode BoxSecurityMode   `json:"current_runtime_security_mode,omitempty"`
	PendingShapeChange         bool              `json:"pending_shape_change"`
	PendingNetworkModeChange   bool              `json:"pending_network_mode_change"`
	PendingSecurityModeChange  bool              `json:"pending_security_mode_change"`
}

// SnapView describes one snap.
type SnapView struct {
	SnapID              string            `json:"snap_id"`
	OrgID               string            `json:"org_id"`
	ProjectID           string            `json:"project_id"`
	Creator             string            `json:"creator"`
	CreatedAt           time.Time         `json:"created_at"`
	LastUsedAt          time.Time         `json:"last_used_at"`
	State               SnapState         `json:"state"`
	InUseReason         string            `json:"inuse_reason,omitempty"`
	Reason              string            `json:"reason,omitempty"`
	ParentChain         []string          `json:"parent_chain,omitempty"`
	SourceImageRef      string            `json:"source_image_ref,omitempty"`
	SourceImageDigest   string            `json:"source_image_digest,omitempty"`
	SourceImagePlatform string            `json:"source_image_platform,omitempty"`
	Attached            bool              `json:"attached"`
	AttachedBoxID       string            `json:"attached_box_id,omitempty"`
	Size                *SnapSize         `json:"size,omitempty"`
	OwnedStorage        *SnapOwnedStorage `json:"owned_storage,omitempty"`
}

// SnapTreeView describes the ancestry tree returned for one snap.
type SnapTreeView struct {
	Supported  bool               `json:"supported"`
	RootSnapID string             `json:"root_snap_id,omitempty"`
	Nodes      []SnapTreeNodeView `json:"nodes,omitempty"`
}

// SnapTreeNodeView describes one node inside a snap ancestry tree.
type SnapTreeNodeView struct {
	SnapID              string                   `json:"snap_id"`
	ParentSnapID        string                   `json:"parent_snap_id,omitempty"`
	Creator             string                   `json:"creator"`
	CreatorDisplayName  string                   `json:"creator_display_name,omitempty"`
	CreatorPrimaryEmail string                   `json:"creator_primary_email,omitempty"`
	CreatedAt           time.Time                `json:"created_at"`
	State               SnapState                `json:"state"`
	Reason              string                   `json:"reason,omitempty"`
	SourceImageRef      string                   `json:"source_image_ref,omitempty"`
	SourceImageDigest   string                   `json:"source_image_digest,omitempty"`
	SourceImagePlatform string                   `json:"source_image_platform,omitempty"`
	AttachedBox         *SnapTreeAttachedBoxView `json:"attached_box,omitempty"`
}

// SnapTreeAttachedBoxView describes the box attached to one snap tree node.
type SnapTreeAttachedBoxView struct {
	BoxID        string   `json:"box_id"`
	DesiredShape string   `json:"desired_shape"`
	State        BoxState `json:"state"`
}

// SnapSize is the cached filesystem usage measured for one snap.
type SnapSize struct {
	UsedBytes  uint64    `json:"used_bytes"`
	UsedInodes uint64    `json:"used_inodes"`
	MeasuredAt time.Time `json:"measured_at"`
}

// SnapOwnedStorage is the cached first-owner object storage attribution for one snap.
type SnapOwnedStorage struct {
	Bytes      uint64    `json:"bytes"`
	Objects    uint64    `json:"objects"`
	MeasuredAt time.Time `json:"measured_at"`
}

// RuntimeRequestView describes one asynchronous runtime request.
type RuntimeRequestView struct {
	RuntimeRequestID string `json:"runtime_request_id"`
	State            string `json:"state"`
	SessionID        string `json:"session_id,omitempty"`
	HostID           string `json:"host_id,omitempty"`
}

// ExecView describes one exec request and its current state.
type ExecView struct {
	ExecID         string         `json:"exec_id"`
	BoxID          string         `json:"box_id"`
	OrgID          string         `json:"org_id"`
	ProjectID      string         `json:"project_id"`
	Creator        string         `json:"creator"`
	AcceptedAt     time.Time      `json:"accepted_at"`
	Source         string         `json:"source"`
	CommandSummary string         `json:"command_summary"`
	ShapeSnapshot  string         `json:"shape_snapshot"`
	Mode           string         `json:"mode,omitempty"`
	State          string         `json:"state"`
	ExitCode       *int           `json:"exit_code,omitempty"`
	OutputSummary  string         `json:"output_summary,omitempty"`
	Reason         string         `json:"reason,omitempty"`
	HardDeadlineAt *time.Time     `json:"hard_deadline_at,omitempty"`
	IdleDeadlineAt *time.Time     `json:"idle_deadline_at,omitempty"`
	StdinEnabled   bool           `json:"stdin_enabled,omitempty"`
	AttachURL      string         `json:"attach_url,omitempty"`
	Diagnostics    map[string]any `json:"diagnostics,omitempty"`
}

// SharedSnapLineView describes one shared snap in list results.
type SharedSnapLineView struct {
	OrgID                    string    `json:"org_id"`
	Name                     string    `json:"name"`
	LatestVersion            int       `json:"latest_version"`
	Publisher                string    `json:"publisher"`
	PublishedBy              string    `json:"published_by"`
	PublishedAt              time.Time `json:"published_at"`
	Description              string    `json:"description,omitempty"`
	SourceProjectID          string    `json:"source_project_id"`
	SourceProjectCID         string    `json:"source_project_cid"`
	SourceProjectDisplayName string    `json:"source_project_display_name"`
}

// SharedSnapVersionView describes one published shared snap version.
type SharedSnapVersionView struct {
	OrgID                    string    `json:"org_id"`
	Name                     string    `json:"name"`
	Version                  int       `json:"version"`
	Description              string    `json:"description,omitempty"`
	PublishedBy              string    `json:"published_by"`
	PublishedAt              time.Time `json:"published_at"`
	SourceProjectID          string    `json:"source_project_id"`
	SourceProjectCID         string    `json:"source_project_cid"`
	SourceProjectDisplayName string    `json:"source_project_display_name"`
}

// SharedSnapDetailView describes one shared snap and its version history.
type SharedSnapDetailView struct {
	OrgID         string                  `json:"org_id"`
	Name          string                  `json:"name"`
	LatestVersion int                     `json:"latest_version"`
	Publisher     string                  `json:"publisher"`
	Versions      []SharedSnapVersionView `json:"versions"`
}

// HostIssue describes one issue reported for a runtime host.
type HostIssue struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HostView describes one runtime host assigned to an organization.
type HostView struct {
	HostID string `json:"host_id"`

	HostClass      string `json:"host_class,omitempty"`
	LifecycleState string `json:"lifecycle_state,omitempty"`
	InstanceID     string `json:"instance_id,omitempty"`
	OwnerOrgID     string `json:"owner_org_id,omitempty"`

	MachineID                 string `json:"machine_id,omitempty"`
	BootID                    string `json:"boot_id,omitempty"`
	Hostname                  string `json:"hostname,omitempty"`
	VmaVersion                string `json:"vma_version,omitempty"`
	ChVersion                 string `json:"ch_version,omitempty"`
	RtVersion                 string `json:"rt_version,omitempty"`
	ForegroundRelayConfigured bool   `json:"foreground_relay_configured,omitempty"`

	Connected       bool      `json:"connected"`
	Ready           bool      `json:"ready"`
	LastHeartbeatAt time.Time `json:"last_heartbeat_at,omitempty"`

	ActiveBoxes                                uint32 `json:"active_boxes"`
	ActiveExecs                                uint32 `json:"active_execs"`
	ActiveTransfers                            uint32 `json:"active_transfers"`
	ActiveBackgroundOwnerExecs                 uint32 `json:"active_background_owner_execs"`
	VMAServiceRestartPreservesBackgroundOwners bool   `json:"vma_service_restart_preserves_background_owners"`

	PlanningReservedCPUMillis   uint32 `json:"planning_reserved_cpu_millis"`
	PlanningReservedMemoryBytes uint64 `json:"planning_reserved_memory_bytes"`

	CPUTotalCores              float64 `json:"cpu_total_cores"`
	MemoryTotalBytes           uint64  `json:"memory_total_bytes"`
	CPUUsedCores               float64 `json:"cpu_used_cores"`
	MemoryUsedBytes            uint64  `json:"memory_used_bytes"`
	RuntimeReservedCPUMillis   uint32  `json:"runtime_reserved_cpu_millis"`
	RuntimeReservedMemoryBytes uint64  `json:"runtime_reserved_memory_bytes"`

	LastIssueSummary string      `json:"last_issue_summary,omitempty"`
	Issues           []HostIssue `json:"issues,omitempty"`
}

// OrgHostsView describes runtime host assignment for one organization.
type OrgHostsView struct {
	OrgID         string     `json:"org_id"`
	AssignedHosts int        `json:"assigned_hosts"`
	Hosts         []HostView `json:"hosts"`
}

// TTYSize describes terminal size in rows and columns.
type TTYSize struct {
	Rows uint32 `json:"rows,omitempty"`
	Cols uint32 `json:"cols,omitempty"`
}

// CreateBoxRequest creates a new box from an image or snap.
type CreateBoxRequest struct {
	BoxID          string            `json:"box_id,omitempty"`
	DesiredShape   string            `json:"desired_shape,omitempty"`
	NetworkMode    BoxNetworkMode    `json:"network_mode,omitempty"`
	SecurityMode   BoxSecurityMode   `json:"security_mode,omitempty"`
	Description    string            `json:"description,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	SourceSnapID   string            `json:"source_snap_id,omitempty"`
	SourceImageRef string            `json:"source_image_ref,omitempty"`
}

// CreateBoxFromSharedSnapRequest creates a box from a published shared snap.
type CreateBoxFromSharedSnapRequest struct {
	Version      *int              `json:"version,omitempty"`
	BoxID        string            `json:"box_id,omitempty"`
	DesiredShape string            `json:"desired_shape,omitempty"`
	NetworkMode  BoxNetworkMode    `json:"network_mode,omitempty"`
	SecurityMode BoxSecurityMode   `json:"security_mode,omitempty"`
	Description  string            `json:"description,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
}

// ImportSnapRequest imports a snap from an image reference.
type ImportSnapRequest struct {
	ImageRef string `json:"image_ref"`
}

// UpdateMeRequest updates caller profile fields.
type UpdateMeRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
}

// CreateSSHKeyRequest creates one account SSH key.
type CreateSSHKeyRequest struct {
	Label     string `json:"label"`
	PublicKey string `json:"public_key"`
}

// UpdateOrgRequest updates mutable organization fields.
type UpdateOrgRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	OrgCID      *string `json:"org_cid,omitempty"`
}

// UpdateMembershipRequest updates one organization membership.
type UpdateMembershipRequest struct {
	Role MembershipRole `json:"role"`
}

// CreateInvitationRequest invites one user into an organization.
type CreateInvitationRequest struct {
	InviteeEmail string         `json:"invitee_email"`
	Role         MembershipRole `json:"role"`
}

// CreateAPIKeyRequest creates one API key for the current organization.
type CreateAPIKeyRequest struct {
	Description string     `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	NoExpire    bool       `json:"no_expire"`
}

// UpdateProjectRequest updates mutable project fields.
type UpdateProjectRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateProjectMembershipRequest updates one project membership.
type UpdateProjectMembershipRequest struct {
	Role ProjectRole `json:"role"`
}

// UpdateBoxRequest updates mutable box fields.
type UpdateBoxRequest struct {
	Description  *string            `json:"description,omitempty"`
	Labels       *map[string]string `json:"labels,omitempty"`
	DesiredShape *string            `json:"desired_shape,omitempty"`
	NetworkMode  *BoxNetworkMode    `json:"network_mode,omitempty"`
	SecurityMode *BoxSecurityMode   `json:"security_mode,omitempty"`
}

// PublishSharedSnapRequest publishes one snap into the shared snap catalog.
type PublishSharedSnapRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	SourceSnapID string `json:"source_snap_id"`
}

// CreateSnapFromSharedSnapRequest creates a new snap from a shared snap version.
type CreateSnapFromSharedSnapRequest struct {
	Version *int `json:"version,omitempty"`
}

// ExecBoxRequest starts one foreground or background exec in a box.
type ExecBoxRequest struct {
	DeadlineAt   *time.Time        `json:"deadline_at,omitempty"`
	Command      []string          `json:"command"`
	EnvOverrides map[string]string `json:"env_overrides,omitempty"`
	User         string            `json:"user,omitempty"`
	Workdir      string            `json:"workdir,omitempty"`
	StdinEnabled bool              `json:"stdin_enabled,omitempty"`
	TTY          bool              `json:"tty,omitempty"`
	TTYSize      *TTYSize          `json:"tty_size,omitempty"`
}

// ExecStreamEvent describes one event emitted by exec streaming APIs.
type ExecStreamEvent struct {
	Type          string `json:"type"`
	ExecID        string `json:"exec_id,omitempty"`
	Data          []byte `json:"data,omitempty"`
	ExitCode      int32  `json:"exit_code,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	CancelReason  string `json:"cancel_reason,omitempty"`
}

// ExecAttachInput describes one input message sent over exec attach.
type ExecAttachInput struct {
	Type string `json:"type"`
	Data []byte `json:"data,omitempty"`
	Rows uint32 `json:"rows,omitempty"`
	Cols uint32 `json:"cols,omitempty"`
}

// BackgroundExecPullOutput describes one poll result from PullBackgroundExecOutput.
type BackgroundExecPullOutput struct {
	Body           []byte
	NextCursor     string
	State          string
	ExitCode       *int
	Reason         string
	IdleDeadlineAt *time.Time
}
