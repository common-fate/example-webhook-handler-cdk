package audit

import (
	"time"

	"github.com/common-fate/sdk/eid"
)

type Action string

const (
	ActionGrantRequested             Action = "grant.requested"
	ActionGrantApproved              Action = "grant.approved"
	ActionGrantActivated             Action = "grant.activated"
	ActionGrantProvisioned           Action = "grant.provisioned"
	ActionGrantProvisioningAttempted Action = "grant.provisioning_attempted"
	ActionGrantExtended              Action = "grant.extended"
	ActionGrantDeprovisioned         Action = "grant.deprovisioned"
	ActionGrantCancelled             Action = "grant.cancelled"
	ActionGrantRevoked               Action = "grant.revoked"
	ActionGrantProvisioningError     Action = "grant.provisioning_error"
	ActionGrantDeprovisioningError   Action = "grant.deprovisioning_error"
	ActionGrantBreakglassActivated   Action = "grant.breakglass_activated"
	ActionProxySessionStarted        Action = "proxy.session.started"
	ActionProxySessionEnded          Action = "proxy.session.ended"
	ActionForceClosed                Action = "grant.force_closed"
)

// Log are user facing audit trail events relating to a particular grant.
// They carry less information than audit trail events that we emit internally to our customers' SIEMs.
// The use case for these is printing info like 'Access was approved'.
//
// Log is excluded from the Cedar schema.
type Log struct {
	ID string `json:"id" authz:"id"`

	// The resources the event relates to.
	Targets []eid.EID `json:"targets" authz:"targets"`

	Action Action `json:"action" authz:"action"`

	// The identity chain associated with the actor.
	// This will include information about the OIDC token and subject used to make the
	// API call resulting in the action.
	CallerIdentityChain IdentityChain `json:"caller_identity_chain" authz:"caller_identity_chain"`

	// Actor that performed the action
	Actor Actor `json:"actor" authz:"actor"`

	// The principal that the action relates to, for grant actions, this is the principal of the grant.
	Principal *Principal `json:"principal,omitempty" authz:"principal,omitempty"`

	// Message associated with the event
	Message string `json:"message" authz:"message"`

	Context *Context `json:"context,omitempty" authz:"context"`

	// The time the event occurred
	OccurredAt time.Time `json:"occurred_at" authz:"occurred_at"`

	// Justification is not filled in when the audit log event is emitted,
	// but is present for when webhook events are emitted to external
	// webhook integrations (such as a SIEM)
	Justification *Justification `json:"justification,omitempty" authz:"justification"`

	// Index is used to represent the intended order of logs which are recorded with the same timestamp
	// For example, audit logs which are emitted from the same API operation
	// The index should be used as a secondary sort order when OccurredAt are equal
	// The index should not be treated as an ID and it is not unique
	Index int `json:"index" authz:"index"`
}

type IdentityChain []IdentityLink

type IdentityLink struct {
	ID    eid.EID `json:"id" authz:"id"`
	Label *string `json:"label,omitempty" authz:"label"`
}

// Context is context to include in the audit log events.
type Context struct {
	Request      *RequestContext      `json:"request,omitempty" authz:"request"`
	Authz        *AuthzContext        `json:"authz,omitempty" authz:"authz"`
	ProxySession *ProxySessionContext `json:"proxy_session,omitempty" authz:"proxy_session"`
}

type RequestContext struct {
	ClientAddr string `json:"client_addr" authz:"client_addr"`
	UserAgent  string `json:"user_agent" authz:"user_addr"`
}

type AuthzContext struct {
	EvalID string `json:"eval" authz:"eval"`
}
type ProxySessionContext struct {
	SessionID string `json:"session_id" authz:"session_id"`
}

type Justification struct {
	Reason string `json:"reason" authz:"reason"`
}

type Principal struct {
	Type  string `json:"type" authz:"type"`
	ID    string `json:"id" authz:"id"`
	Name  string `json:"name" authz:"name"`
	Email string `json:"email" authz:"email"`
}

func (a Principal) EID() eid.EID {
	return eid.New(a.Type, a.ID)
}

type Actor struct {
	Type string `json:"type" authz:"type"`
	ID   string `json:"id" authz:"id"`
	// Name is not filled in when the audit log event is emitted,
	// but is present for when webhook events are emitted to external
	// webhook integrations (such as a SIEM)
	Name string `json:"name,omitempty" authz:"name"`
	// Email is not filled in when the audit log event is emitted,
	// but is present for when webhook events are emitted to external
	// webhook integrations (such as a SIEM)
	Email string `json:"email,omitempty" authz:"email"`
}

func (a Actor) EID() eid.EID {
	return eid.New(a.Type, a.ID)
}

func ActorFromEID(e eid.EID) Actor {
	return Actor{
		ID:   e.ID,
		Type: e.Type,
	}
}
