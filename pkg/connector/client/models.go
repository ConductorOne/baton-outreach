package client

import (
	"time"
)

type Pagination struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

type UserAttributes struct {
	AccountsViewId                       int       `json:"accountsViewId"`
	ActivityNotificationsDisabled        bool      `json:"activityNotificationsDisabled"`
	BounceWarningEmailEnabled            bool      `json:"bounceWarningEmailEnabled"`
	BridgePhone                          string    `json:"bridgePhone"`
	BridgePhoneExtension                 string    `json:"bridgePhoneExtension"`
	CallsViewId                          int       `json:"callsViewId"`
	ControlledTabDefault                 string    `json:"controlledTabDefault"`
	CreatedAt                            time.Time `json:"createdAt"`
	CurrentSignInAt                      time.Time `json:"currentSignInAt"`
	DailyDigestEmailEnabled              bool      `json:"dailyDigestEmailEnabled"`
	DefaultRulesetId                     int       `json:"defaultRulesetId"`
	Email                                string    `json:"email"`
	EnableVoiceRecordings                bool      `json:"enableVoiceRecordings"`
	EngagementEmailsEnabled              bool      `json:"engagementEmailsEnabled"`
	FirstName                            string    `json:"firstName"`
	InboundBridgePhone                   string    `json:"inboundBridgePhone"`
	InboundBridgePhoneExtension          string    `json:"inboundBridgePhoneExtension"`
	InboundCallBehavior                  string    `json:"inboundCallBehavior"`
	InboundPhoneType                     string    `json:"inboundPhoneType"`
	InboundVoicemailCustomMessageText    string    `json:"inboundVoicemailCustomMessageText"`
	InboundVoicemailMessageTextVoice     string    `json:"inboundVoicemailMessageTextVoice"`
	InboundVoicemailPromptType           string    `json:"inboundVoicemailPromptType"`
	KaiaRecordingsViewId                 int       `json:"kaiaRecordingsViewId"`
	KeepBridgePhoneConnected             bool      `json:"keepBridgePhoneConnected"`
	LastName                             string    `json:"lastName"`
	LastSignInAt                         time.Time `json:"lastSignInAt"`
	Locked                               bool      `json:"locked"`
	MailboxErrorEmailEnabled             bool      `json:"mailboxErrorEmailEnabled"`
	MeetingEngagementNotificationEnabled bool      `json:"meetingEngagementNotificationEnabled"`
	Name                                 string    `json:"name"`
	NotificationsEnabled                 bool      `json:"notificationsEnabled"`
	OceClickToDialEverywhere             bool      `json:"oceClickToDialEverywhere"`
	OceGmailToolbar                      bool      `json:"oceGmailToolbar"`
	OceGmailTrackingState                string    `json:"oceGmailTrackingState"`
	OceSalesforceEmailDecorating         bool      `json:"oceSalesforceEmailDecorating"`
	OceSalesforcePhoneDecorating         bool      `json:"oceSalesforcePhoneDecorating"`
	OceUniversalTaskFlow                 bool      `json:"oceUniversalTaskFlow"`
	OceWindowMode                        bool      `json:"oceWindowMode"`
	OpportunitiesViewId                  int       `json:"opportunitiesViewId"`
	PasswordExpiresAt                    time.Time `json:"passwordExpiresAt"`
	PhoneCountryCode                     string    `json:"phoneCountryCode"`
	PhoneNumber                          string    `json:"phoneNumber"`
	PhoneType                            string    `json:"phoneType"`
	PluginAlertNotificationEnabled       bool      `json:"pluginAlertNotificationEnabled"`
	PreferredVoiceRegion                 string    `json:"preferredVoiceRegion"`
	PrefersLocalPresence                 bool      `json:"prefersLocalPresence"`
	PrimaryTimezone                      string    `json:"primaryTimezone"`
	ProspectsViewId                      int       `json:"prospectsViewId"`
	ReportsTeamPerfViewId                int       `json:"reportsTeamPerfViewId"`
	ReportsViewId                        int       `json:"reportsViewId"`
	ScimExternalId                       string    `json:"scimExternalId"`
	ScimSource                           string    `json:"scimSource"`
	SecondaryTimezone                    string    `json:"secondaryTimezone"`
	SenderNotificationsExcluded          bool      `json:"senderNotificationsExcluded"`
	TasksViewId                          int       `json:"tasksViewId"`
	TeamsViewId                          int       `json:"teamsViewId"`
	TertiaryTimezone                     string    `json:"tertiaryTimezone"`
	TextingEmailNotifications            bool      `json:"textingEmailNotifications"`
	Title                                string    `json:"title"`
	UnknownReplyEmailEnabled             bool      `json:"unknownReplyEmailEnabled"`
	UpdatedAt                            time.Time `json:"updatedAt"`
	UserGUID                             string    `json:"userGuid"`
	Username                             string    `json:"username"`
	UsersViewId                          int       `json:"usersViewId"`
	VoicemailNotificationEnabled         bool      `json:"voicemailNotificationEnabled"`
	WeeklyDigestEmailEnabled             bool      `json:"weeklyDigestEmailEnabled"`
}

type UserRelationships struct {
	Profile *struct {
		Data *DataDetailPair `json:"data,omitempty"`
	} `json:"profile,omitempty"`
	Teams *struct {
		Data *[]DataDetailPair `json:"data,omitempty"`
	} `json:"teams,omitempty"`
}

type User struct {
	Id            int                `json:"id"`
	Attributes    UserAttributes     `json:"attributes"`
	Relationships *UserRelationships `json:"relationships,omitempty"`
	Type          string             `json:"type"`
}

type UsersResponse struct {
	Links   *Pagination `json:"links,omitempty"`
	Results []*User     `json:"data"`
}

type ProfileAttributes struct {
	CreatedAt string `json:"createdAt"`
	IsAdmin   bool   `json:"isAdmin"`
	Name      string `json:"name"`
	SpecialId string `json:"specialId"`
	UpdatedAt string `json:"updatedAt"`
}

type Profile struct {
	Attributes ProfileAttributes `json:"attributes"`
	Id         int               `json:"id"`
	Type       string            `json:"type"`
}

type ProfilesResponse struct {
	Links   *Pagination `json:"links,omitempty"`
	Results []*Profile  `json:"data"`
}

type TeamAttributes struct {
	Color          string `json:"color"`
	CreatedAt      string `json:"createdAt"`
	Name           string `json:"name"`
	ScimExternalId string `json:"scimExternalId"`
	ScimSource     string `json:"scimSource"`
	UpdatedAt      string `json:"updatedAt"`
}

type TeamRelationships struct {
	Users *struct {
		Data *[]DataDetailPair `json:"data,omitempty"`
	} `json:"users,omitempty"`
}

type Team struct {
	Attributes    TeamAttributes     `json:"attributes"`
	Id            int                `json:"id"`
	Relationships *TeamRelationships `json:"relationships,omitempty"`
	Type          string             `json:"type"`
}

type TeamsResponse struct {
	Links   *Pagination `json:"links,omitempty"`
	Results []*Team     `json:"data"`
}

type UpdateTeamBody struct {
	Id            int                     `json:"id"`
	Type          string                  `json:"type"` // Type should always be 'team'.
	Relationships UpdateTeamRelationships `json:"relationships"`
}

type UpdateTeamRelationships struct {
	Users struct {
		Data []DataDetailPair `json:"data"`
	} `json:"users"`
}

type DataDetailPair struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

type UpdateUsersProfileBody struct {
	Id            int                      `json:"id"`
	Type          string                   `json:"type"`
	Relationships UserProfileRelationships `json:"relationships"`
}

type UserProfileRelationships struct {
	Profile struct {
		Data DataDetailPair `json:"data"`
	} `json:"profile"`
}

type NewUserBody struct {
	Data struct {
		Type       string            `json:"type"` // The type should always be 'user'.
		Attributes NewUserAttributes `json:"attributes"`
	} `json:"data"`
}

type NewUserAttributes struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserLockStatusUpdate struct {
	Id         int    `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Locked bool `json:"locked"`
	} `json:"attributes"`
}
