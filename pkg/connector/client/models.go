package client

import "time"

type Pagination struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

type UserAttributes struct {
	AccountsViewId                int       `json:"accountsViewId"`
	ActivityNotificationsDisabled bool      `json:"activityNotificationsDisabled"`
	BounceWarningEmailEnabled     bool      `json:"bounceWarningEmailEnabled"`
	BridgePhone                   string    `json:"bridgePhone"`
	BridgePhoneExtension          string    `json:"bridgePhoneExtension"`
	CallsViewId                   int       `json:"callsViewId"`
	ControlledTabDefault          string    `json:"controlledTabDefault"`
	CreatedAt                     time.Time `json:"createdAt"`
	CurrentSignInAt               time.Time `json:"currentSignInAt"`
	//Custom1                       string    `json:"custom1"`
	//Custom2                       string    `json:"custom2"`
	//Custom3                       string    `json:"custom3"`
	//Custom4                       string    `json:"custom4"`
	//Custom5                       string    `json:"custom5"`
	DailyDigestEmailEnabled bool `json:"dailyDigestEmailEnabled"`
	DefaultRulesetId        int  `json:"defaultRulesetId"`
	Duties                  struct {
	} `json:"duties"`
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
	UserGuid                             string    `json:"userGuid"`
	Username                             string    `json:"username"`
	UsersViewId                          int       `json:"usersViewId"`
	VoicemailNotificationEnabled         bool      `json:"voicemailNotificationEnabled"`
	WeeklyDigestEmailEnabled             bool      `json:"weeklyDigestEmailEnabled"`
}

type User struct {
	Id            int            `json:"id"`
	Attributes    UserAttributes `json:"attributes"`
	Relationships struct {
		Batches []struct {
			Links struct {
				Related string `json:"related"`
			} `json:"links"`
		} `json:"batches"`
		ContentCategories []struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"contentCategories"`
		Creator struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"creator"`
		Mailbox struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"mailbox"`
		Mailboxes []struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"mailboxes"`
		Profile struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"profile"`
		Recipients []struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"recipients"`
		Role struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"role"`
		Teams []struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"teams"`
		Updater struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"updater"`
	} `json:"relationships"`
	Type string `json:"type"`
}

type UsersResponse struct {
	Links   *Pagination `json:"links,omitempty"`
	Results []*User     `json:"data"`
}

type RoleAttributes struct {
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Role struct {
	Attributes    RoleAttributes `json:"attributes"`
	Id            int            `json:"id"`
	Relationships *struct {
		ParentRole *struct {
			Data *struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data,omitempty"`
		} `json:"parentRole,omitempty"`
	} `json:"relationships,omitempty"`
	Type string `json:"type"`
}

type RolesResponse struct {
	Links   *Pagination `json:"links,omitempty"`
	Results []*Role     `json:"data"`
}

type TeamAttributes struct {
	Color          string    `json:"color"`
	CreatedAt      time.Time `json:"createdAt"`
	Name           string    `json:"name"`
	ScimExternalId string    `json:"scimExternalId"`
	ScimSource     string    `json:"scimSource"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Team struct {
	Attributes    TeamAttributes `json:"attributes"`
	Id            int            `json:"id"`
	Relationships *struct {
		Batches []struct {
			Links struct {
				Related string `json:"related"`
			} `json:"links"`
		} `json:"batches"`
		ContentCategories []struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"contentCategories"`
		Creator struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"creator"`
		Updater struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"updater"`
		Users []struct {
			Data struct {
				Id   int    `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"users,omitempty"`
	} `json:"relationships,omitempty"`
	Type string `json:"type"`
}

type TeamsResponse struct {
	Links   *Pagination `json:"links,omitempty"`
	Results []*Team     `json:"data"`
}
