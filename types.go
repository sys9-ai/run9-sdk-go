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
	// OrgKindPersonal identifies one personal organization.
	OrgKindPersonal OrgKind = "personal"
	// OrgKindShared identifies one shared organization.
	OrgKindShared OrgKind = "shared"
)

// Secret scopes accepted by project and box secret APIs.
const (
	// ProjectSecretScopeProject attaches one secret to project scope.
	ProjectSecretScopeProject ProjectSecretScope = "project"
	// ProjectSecretScopeBox attaches one secret to one box.
	ProjectSecretScopeBox ProjectSecretScope = "box"
)

// InvitationState represents one invitation lifecycle state.
type InvitationState string

// BoxState represents the current box state.
type BoxState string

// BoxNetworkMode represents the requested networking mode for a box.
type BoxNetworkMode string

// Box network modes accepted by box create and update APIs.
const (
	// BoxNetworkModeNormal requests the default non-managed network path.
	BoxNetworkModeNormal BoxNetworkMode = "normal"
	// BoxNetworkModeManaged requests the managed network path.
	BoxNetworkModeManaged BoxNetworkMode = "managed"
)

// BoxSecurityMode represents the requested security mode for a box.
type BoxSecurityMode string

// Box security modes accepted by box create and update APIs.
const (
	// BoxSecurityModeRestricted requests the default restricted security mode.
	BoxSecurityModeRestricted BoxSecurityMode = "restricted"
	// BoxSecurityModeUnsafe requests the less restricted unsafe security mode.
	BoxSecurityModeUnsafe BoxSecurityMode = "unsafe"
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
	// DisplayName is the human-readable project name.
	DisplayName string `json:"display_name"`
	// Description is optional free-form project text.
	Description string `json:"description,omitempty"`
}

// CreateProjectSecretRequest creates a new project or box secret.
type CreateProjectSecretRequest struct {
	// Name is an optional human-readable secret name.
	Name string `json:"name,omitempty"`
	// Value is the raw secret value to store.
	Value string `json:"value"`
	// Placeholder is the literal marker that managed egress will replace.
	Placeholder string `json:"placeholder"`
	// AllowedHosts lists the destination hosts where placeholder substitution is allowed.
	AllowedHosts []string `json:"allowed_hosts"`
	// InjectHeaderName selects which HTTP header names are eligible for substitution.
	InjectHeaderName string `json:"inject_header_name"`
}

// UpdateProjectSecretRequest updates one existing project or box secret.
type UpdateProjectSecretRequest struct {
	// Name replaces the stored human-readable secret name when non-nil.
	Name *string `json:"name,omitempty"`
	// Value replaces the stored secret value when non-nil.
	Value *string `json:"value,omitempty"`
	// AllowedHosts replaces the full allowed-host set when non-nil.
	AllowedHosts *[]string `json:"allowed_hosts,omitempty"`
	// InjectHeaderName replaces the header selector when non-nil.
	InjectHeaderName *string `json:"inject_header_name,omitempty"`
}

// ListBoxesRequest describes optional filters for ListBoxes.
type ListBoxesRequest struct {
	// Creator filters by the creator identifier.
	Creator string
	// Label filters by one label selector supported by the control plane.
	Label string
	// State filters by the current box state.
	State BoxState
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
	// Rows is the terminal height in character cells.
	Rows uint32 `json:"rows,omitempty"`
	// Cols is the terminal width in character cells.
	Cols uint32 `json:"cols,omitempty"`
}

// CreateBoxRequest creates a new box from an image or snap.
type CreateBoxRequest struct {
	// BoxID requests one specific box identifier. When empty, the control plane generates one.
	BoxID string `json:"box_id,omitempty"`
	// DesiredShape requests the compute shape for the box.
	DesiredShape string `json:"desired_shape,omitempty"`
	// NetworkMode requests the durable network mode for the box.
	NetworkMode BoxNetworkMode `json:"network_mode,omitempty"`
	// SecurityMode requests the durable security mode for the box.
	SecurityMode BoxSecurityMode `json:"security_mode,omitempty"`
	// Description is optional free-form box text.
	Description string `json:"description,omitempty"`
	// Labels sets durable box labels.
	Labels map[string]string `json:"labels,omitempty"`
	// SourceSnapID creates the box from one existing snap when non-empty.
	SourceSnapID string `json:"source_snap_id,omitempty"`
	// SourceImageRef imports directly from one image reference when non-empty.
	SourceImageRef string `json:"source_image_ref,omitempty"`
}

// CreateBoxFromSharedSnapRequest creates a box from a published shared snap.
type CreateBoxFromSharedSnapRequest struct {
	// Version selects one published version. When nil, the latest version is used.
	Version *int `json:"version,omitempty"`
	// BoxID requests one specific box identifier. When empty, the control plane generates one.
	BoxID string `json:"box_id,omitempty"`
	// DesiredShape requests the compute shape for the box.
	DesiredShape string `json:"desired_shape,omitempty"`
	// NetworkMode requests the durable network mode for the box.
	NetworkMode BoxNetworkMode `json:"network_mode,omitempty"`
	// SecurityMode requests the durable security mode for the box.
	SecurityMode BoxSecurityMode `json:"security_mode,omitempty"`
	// Description is optional free-form box text.
	Description string `json:"description,omitempty"`
	// Labels sets durable box labels.
	Labels map[string]string `json:"labels,omitempty"`
}

// ImportSnapRequest imports a snap from an image reference.
type ImportSnapRequest struct {
	// ImageRef is the OCI-style image reference to import.
	ImageRef string `json:"image_ref"`
}

// ListSnapsRequest describes optional filters for ListSnaps.
type ListSnapsRequest struct {
	// Attached filters by whether a snap is currently attached to a box.
	Attached *bool
}

// UpdateMeRequest updates caller profile fields.
type UpdateMeRequest struct {
	// DisplayName replaces the current account display name when non-nil.
	DisplayName *string `json:"display_name,omitempty"`
}

// CreateSSHKeyRequest creates one account SSH key.
type CreateSSHKeyRequest struct {
	// Label is the human-readable SSH key label.
	Label string `json:"label"`
	// PublicKey is the authorized_keys-formatted SSH public key.
	PublicKey string `json:"public_key"`
}

// UpdateOrgRequest updates mutable organization fields.
type UpdateOrgRequest struct {
	// DisplayName replaces the current organization display name when non-nil.
	DisplayName *string `json:"display_name,omitempty"`
	// OrgCID replaces the durable organization CID when non-nil.
	OrgCID *string `json:"org_cid,omitempty"`
}

// UpdateMembershipRequest updates one organization membership.
type UpdateMembershipRequest struct {
	// Role is the desired organization role.
	Role MembershipRole `json:"role"`
}

// CreateInvitationRequest invites one user into an organization.
type CreateInvitationRequest struct {
	// InviteeEmail is the target email address.
	InviteeEmail string `json:"invitee_email"`
	// Role is the organization role to grant on acceptance.
	Role MembershipRole `json:"role"`
}

// CreateAPIKeyRequest creates one API key for the current organization.
type CreateAPIKeyRequest struct {
	// Description is optional operator-facing API key text.
	Description string `json:"description,omitempty"`
	// ExpiresAt requests one explicit expiry time when non-nil.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// NoExpire requests a non-expiring API key.
	NoExpire bool `json:"no_expire"`
}

// UpdateProjectRequest updates mutable project fields.
type UpdateProjectRequest struct {
	// DisplayName replaces the current project display name when non-nil.
	DisplayName *string `json:"display_name,omitempty"`
	// Description replaces the current project description when non-nil.
	Description *string `json:"description,omitempty"`
}

// UpdateProjectMembershipRequest updates one project membership.
type UpdateProjectMembershipRequest struct {
	// Role is the desired project role.
	Role ProjectRole `json:"role"`
}

// UpdateBoxRequest updates mutable box fields.
type UpdateBoxRequest struct {
	// Description replaces the durable box description when non-nil.
	Description *string `json:"description,omitempty"`
	// Labels replaces the full durable label map when non-nil.
	Labels *map[string]string `json:"labels,omitempty"`
	// DesiredShape replaces the durable target shape when non-nil.
	DesiredShape *string `json:"desired_shape,omitempty"`
	// NetworkMode replaces the durable network mode when non-nil.
	NetworkMode *BoxNetworkMode `json:"network_mode,omitempty"`
	// SecurityMode replaces the durable security mode when non-nil.
	SecurityMode *BoxSecurityMode `json:"security_mode,omitempty"`
}

// PublishSharedSnapRequest publishes one snap into the shared snap catalog.
type PublishSharedSnapRequest struct {
	// Name is the shared snap catalog name.
	Name string `json:"name"`
	// Description is optional catalog text for this publication.
	Description string `json:"description,omitempty"`
	// SourceSnapID identifies the project-local snap to publish.
	SourceSnapID string `json:"source_snap_id"`
}

// CreateSnapFromSharedSnapRequest creates a new snap from a shared snap version.
type CreateSnapFromSharedSnapRequest struct {
	// Version selects one published version. When nil, the latest version is used.
	Version *int `json:"version,omitempty"`
}

// ExecRequest starts one foreground or background exec in a box.
type ExecRequest struct {
	// DeadlineAt requests a hard deadline. When nil, the target exec API applies its default.
	DeadlineAt *time.Time `json:"deadline_at,omitempty"`
	// Command is the argv array executed in the target box.
	Command []string `json:"command"`
	// EnvOverrides replaces or adds environment variables for this exec.
	EnvOverrides map[string]string `json:"env_overrides,omitempty"`
	// User selects the target user account inside the box when non-empty.
	User string `json:"user,omitempty"`
	// Workdir selects the target working directory when non-empty.
	Workdir string `json:"workdir,omitempty"`
	// StdinEnabled keeps stdin available for this exec when true.
	StdinEnabled bool `json:"stdin_enabled,omitempty"`
	// TTY requests one PTY-backed exec session.
	TTY bool `json:"tty,omitempty"`
	// TTYSize provides the initial PTY size when TTY is true.
	TTYSize *TTYSize `json:"tty_size,omitempty"`
}

// ExecStreamEvent describes one event emitted by exec streaming APIs.
type ExecStreamEvent struct {
	// Type identifies the event kind. Common values include started, keepalive, stdout,
	// stderr, exit, error, and cancelled.
	Type string `json:"type"`
	// ExecID carries the durable exec identifier on events that include it.
	ExecID string `json:"exec_id,omitempty"`
	// Data carries stdout or stderr bytes for stream output events.
	Data []byte `json:"data,omitempty"`
	// ExitCode is set on terminal exit events.
	ExitCode int32 `json:"exit_code,omitempty"`
	// FailureReason is set when the stream terminates with an error event.
	FailureReason string `json:"failure_reason,omitempty"`
	// CancelReason is set when the stream terminates with a cancelled event.
	CancelReason string `json:"cancel_reason,omitempty"`
}

// ExecAttachInput describes one input message sent over exec attach.
type ExecAttachInput struct {
	// Type identifies the input frame kind. Supported values are stdin, close_stdin, and resize.
	Type string `json:"type"`
	// Data carries stdin bytes when Type is stdin.
	Data []byte `json:"data,omitempty"`
	// Rows carries the new terminal height when Type is resize.
	Rows uint32 `json:"rows,omitempty"`
	// Cols carries the new terminal width when Type is resize.
	Cols uint32 `json:"cols,omitempty"`
}

// PullBackgroundExecOutputRequest describes one pull-output request.
type PullBackgroundExecOutputRequest struct {
	// Cursor resumes from one previously returned cursor. Leave empty to read from the beginning once.
	Cursor string
	// Wait asks the service to short-poll before returning when the cursor is currently at the live tail.
	Wait time.Duration
}

// WriteBackgroundExecStdinRequest describes one background exec stdin write.
type WriteBackgroundExecStdinRequest struct {
	// Data is the raw stdin chunk to send in this request.
	Data []byte
	// CloseStdin requests EOF immediately after Data is written.
	CloseStdin bool
}

// BackgroundExecPullOutput describes one decoded background exec output window.
type BackgroundExecPullOutput struct {
	// Events carries the decoded output and terminal events returned by this poll.
	Events []BackgroundExecOutputEvent
	// NextCursor is the cursor to use on the next incremental poll.
	NextCursor string
	// State is the current durable exec state.
	State string
	// ExitCode is set after terminal exit when the exec produced one.
	ExitCode *int
	// Reason carries the durable terminal reason when one is available.
	Reason string
	// IdleDeadlineAt is the current idle lease deadline when the control plane reports it.
	IdleDeadlineAt *time.Time
}
