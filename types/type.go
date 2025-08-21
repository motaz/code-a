package types

type SetUp struct {
	Server         string
	Database       string
	User           string
	Password       string
	IsConfigExists bool
	Version        string
}

type DomainType struct {
	DomainID      int
	DomainName    string
	IsLocal       bool
	DefaultDomain bool
	IsEnabled     bool
	RemoteURL     string
	DefaultPage   string
}

type Users struct {
	Userid   int
	Login    string
	Password string
	Fullname string
	// Email            string
	Info      string
	Isenabled int
	Isadmin   int
	DomainID  int
	// ResetToken       string
	// ResetTokenExpiry *time.Time
}
type UserInfo struct {
	Userid        int
	Login         string
	Password      string
	Fullname      string
	Info          string
	Isadmin       bool
	DomainID      int
	DomainName    string
	IsLocal       bool
	DefaultDomain bool
	IsEnabled     bool
	RemoteURL     string
	DefaultPage   string
}

type AccessResponseType struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"errorcode"`
	UserID    interface{} `json:"user_id"`
}

type Login struct {
	Success     bool
	ErrorCode   int
	Message     string
	Id          int
	SessionID   string
	SessionInfo string
	Username    string
	Domain      string
}

type OperationResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorcode"`
	ID        int    `json:"user_id"`
}

type Home struct {
	Username        string
	UserID          int
	Page            string
	Url             string
	Domain          string
	IsAdmin         bool
	DomainID        int
	Domains         []DomainType
	UserInfo        []UserInfo
	User            UserInfo
	DomainInfo      DomainType
	AlertType       string
	ResponseMessage string
	ResponseStatus  bool
	View            string
	Showall         string
	Search          string
	SearchButton    string
	Modify          string
	UserType        int
}

type SessionResult struct {
	Success           bool
	ErrorCode         int
	Message           string
	Id                int
	SessionID         string
	SessionInfo       string
	SessionTime       string
	SessionExpiration string
	UserID            int
	Username          string
	DomainName        string
}

type AuthenticateTemplate struct {
	IsAuthenticate bool
	Success        bool
	Key            string
	NewSessionID   string
	Returnto       string
	Version        string
}
