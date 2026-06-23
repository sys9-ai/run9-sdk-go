package run9

import "time"

type MembershipRole string

type OrgKind string

type ProjectRole string

type ProjectSecretScope string

const (
	OrgKindPersonal OrgKind = "personal"
	OrgKindShared   OrgKind = "shared"
)

const (
	ProjectSecretScopeProject ProjectSecretScope = "project"
	ProjectSecretScopeBox     ProjectSecretScope = "box"
)

type InvitationState string

type BoxState string

type BoxNetworkMode string

const (
	BoxNetworkModeNormal  BoxNetworkMode = "normal"
	BoxNetworkModeManaged BoxNetworkMode = "managed"
)

type BoxSecurityMode string

const (
	BoxSecurityModeRestricted BoxSecurityMode = "restricted"
	BoxSecurityModeUnsafe     BoxSecurityMode = "unsafe"
)

type SnapState string

type CurrentSubscriptionView struct {
	Tier      string    `json:"tier"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

type MeView struct {
	UserID          string    `json:"user_id"`
	PrimaryEmail    string    `json:"primary_email"`
	DisplayName     string    `json:"display_name,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	IsSystemManager bool      `json:"is_system_manager"`
}

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

type CurrentOrgIdentityView struct {
	User     MeView  `json:"user"`
	Org      OrgView `json:"org"`
	AuthKind string  `json:"auth_kind"`
}

type DeleteOrgResult struct {
	OrgID  string `json:"org_id"`
	Status string `json:"status"`
}

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

type DeleteInvitationResult struct {
	InvitationID string `json:"invitation_id"`
	Status       string `json:"status"`
}

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

type SSHKeyView struct {
	SSHKeyID    string     `json:"ssh_key_id"`
	Label       string     `json:"label"`
	Fingerprint string     `json:"fingerprint"`
	CreatedAt   time.Time  `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
}

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

type DeleteProjectResult struct {
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
}

type ProjectMembershipView struct {
	ProjectID    string      `json:"project_id"`
	UserID       string      `json:"user_id"`
	PrimaryEmail string      `json:"primary_email"`
	DisplayName  string      `json:"display_name,omitempty"`
	Role         ProjectRole `json:"role"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

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

type CreateProjectRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description,omitempty"`
}

type CreateProjectSecretRequest struct {
	Name             string   `json:"name,omitempty"`
	Value            string   `json:"value"`
	Placeholder      string   `json:"placeholder"`
	AllowedHosts     []string `json:"allowed_hosts"`
	InjectHeaderName string   `json:"inject_header_name"`
}

type UpdateProjectSecretRequest struct {
	Name             *string   `json:"name,omitempty"`
	Value            *string   `json:"value,omitempty"`
	AllowedHosts     *[]string `json:"allowed_hosts,omitempty"`
	InjectHeaderName *string   `json:"inject_header_name,omitempty"`
}

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

type SnapTreeView struct {
	Supported  bool               `json:"supported"`
	RootSnapID string             `json:"root_snap_id,omitempty"`
	Nodes      []SnapTreeNodeView `json:"nodes,omitempty"`
}

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

type RuntimeRequestView struct {
	RuntimeRequestID string `json:"runtime_request_id"`
	State            string `json:"state"`
	SessionID        string `json:"session_id,omitempty"`
	HostID           string `json:"host_id,omitempty"`
}

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

type SharedSnapDetailView struct {
	OrgID         string                  `json:"org_id"`
	Name          string                  `json:"name"`
	LatestVersion int                     `json:"latest_version"`
	Publisher     string                  `json:"publisher"`
	Versions      []SharedSnapVersionView `json:"versions"`
}

type HostIssue struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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

type OrgHostsView struct {
	OrgID         string     `json:"org_id"`
	AssignedHosts int        `json:"assigned_hosts"`
	Hosts         []HostView `json:"hosts"`
}

type TTYSize struct {
	Rows uint32 `json:"rows,omitempty"`
	Cols uint32 `json:"cols,omitempty"`
}

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

type CreateBoxFromSharedSnapRequest struct {
	Version      *int              `json:"version,omitempty"`
	BoxID        string            `json:"box_id,omitempty"`
	DesiredShape string            `json:"desired_shape,omitempty"`
	NetworkMode  BoxNetworkMode    `json:"network_mode,omitempty"`
	SecurityMode BoxSecurityMode   `json:"security_mode,omitempty"`
	Description  string            `json:"description,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
}

type ImportSnapRequest struct {
	ImageRef string `json:"image_ref"`
}

type UpdateMeRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
}

type CreateSSHKeyRequest struct {
	Label     string `json:"label"`
	PublicKey string `json:"public_key"`
}

type UpdateOrgRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	OrgCID      *string `json:"org_cid,omitempty"`
}

type UpdateMembershipRequest struct {
	Role MembershipRole `json:"role"`
}

type CreateInvitationRequest struct {
	InviteeEmail string         `json:"invitee_email"`
	Role         MembershipRole `json:"role"`
}

type CreateAPIKeyRequest struct {
	Description string     `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	NoExpire    bool       `json:"no_expire"`
}

type UpdateProjectRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateProjectMembershipRequest struct {
	Role ProjectRole `json:"role"`
}

type UpdateBoxRequest struct {
	Description  *string            `json:"description,omitempty"`
	Labels       *map[string]string `json:"labels,omitempty"`
	DesiredShape *string            `json:"desired_shape,omitempty"`
	NetworkMode  *BoxNetworkMode    `json:"network_mode,omitempty"`
	SecurityMode *BoxSecurityMode   `json:"security_mode,omitempty"`
}

type PublishSharedSnapRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	SourceSnapID string `json:"source_snap_id"`
}

type CreateSnapFromSharedSnapRequest struct {
	Version *int `json:"version,omitempty"`
}

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

type ExecStreamEvent struct {
	Type          string `json:"type"`
	ExecID        string `json:"exec_id,omitempty"`
	Data          []byte `json:"data,omitempty"`
	ExitCode      int32  `json:"exit_code,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	CancelReason  string `json:"cancel_reason,omitempty"`
}

type ExecAttachInput struct {
	Type string `json:"type"`
	Data []byte `json:"data,omitempty"`
	Rows uint32 `json:"rows,omitempty"`
	Cols uint32 `json:"cols,omitempty"`
}

type BackgroundExecPullOutput struct {
	Body           []byte
	NextCursor     string
	State          string
	ExitCode       *int
	Reason         string
	IdleDeadlineAt *time.Time
}
