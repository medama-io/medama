// Code generated by ogen, DO NOT EDIT.

package api

import (
	"github.com/go-faster/errors"
)

// Request body for logging in.
// Ref: #/components/schemas/AuthLogin
type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetEmail returns the value of Email.
func (s *AuthLogin) GetEmail() string {
	return s.Email
}

// GetPassword returns the value of Password.
func (s *AuthLogin) GetPassword() string {
	return s.Password
}

// SetEmail sets the value of Email.
func (s *AuthLogin) SetEmail(val string) {
	s.Email = val
}

// SetPassword sets the value of Password.
func (s *AuthLogin) SetPassword(val string) {
	s.Password = val
}

type BadRequestError struct {
	Error BadRequestErrorError `json:"error"`
}

// GetError returns the value of Error.
func (s *BadRequestError) GetError() BadRequestErrorError {
	return s.Error
}

// SetError sets the value of Error.
func (s *BadRequestError) SetError(val BadRequestErrorError) {
	s.Error = val
}

func (*BadRequestError) deleteUserRes()          {}
func (*BadRequestError) deleteWebsitesIDRes()    {}
func (*BadRequestError) getUserRes()             {}
func (*BadRequestError) getWebsiteIDSummaryRes() {}
func (*BadRequestError) getWebsitesIDActiveRes() {}
func (*BadRequestError) getWebsitesIDRes()       {}
func (*BadRequestError) getWebsitesRes()         {}
func (*BadRequestError) patchUserRes()           {}
func (*BadRequestError) patchWebsitesIDRes()     {}
func (*BadRequestError) postAuthLoginRes()       {}
func (*BadRequestError) postUserRes()            {}
func (*BadRequestError) postWebsitesRes()        {}

type BadRequestErrorError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *BadRequestErrorError) GetCode() int32 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *BadRequestErrorError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *BadRequestErrorError) SetCode(val int32) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *BadRequestErrorError) SetMessage(val string) {
	s.Message = val
}

type ConflictError struct {
	Error ConflictErrorError `json:"error"`
}

// GetError returns the value of Error.
func (s *ConflictError) GetError() ConflictErrorError {
	return s.Error
}

// SetError sets the value of Error.
func (s *ConflictError) SetError(val ConflictErrorError) {
	s.Error = val
}

func (*ConflictError) deleteUserRes()   {}
func (*ConflictError) patchUserRes()    {}
func (*ConflictError) postUserRes()     {}
func (*ConflictError) postWebsitesRes() {}

type ConflictErrorError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *ConflictErrorError) GetCode() int32 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *ConflictErrorError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *ConflictErrorError) SetCode(val int32) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *ConflictErrorError) SetMessage(val string) {
	s.Message = val
}

type CookieAuth struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *CookieAuth) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *CookieAuth) SetAPIKey(val string) {
	s.APIKey = val
}

// DeleteUserOK is response for DeleteUser operation.
type DeleteUserOK struct{}

func (*DeleteUserOK) deleteUserRes() {}

// DeleteWebsitesIDOK is response for DeleteWebsitesID operation.
type DeleteWebsitesIDOK struct{}

func (*DeleteWebsitesIDOK) deleteWebsitesIDRes() {}

// Website hit event.
// Ref: #/components/schemas/EventHit
type EventHit struct {
	// UUID generated for each user to link multiple events on the same page together.
	B string `json:"b"`
	// Page URL including query parameters.
	U string `json:"u"`
	// Referrer URL.
	R string `json:"r"`
	// Event type consisting of either 'pagehide', 'unload', 'load', 'hidden' or 'visible'.
	E string `json:"e"`
	// Title of page.
	T OptInt `json:"t"`
	// Timezone of the user.
	D OptString `json:"d"`
	// Screen width.
	W OptInt `json:"w"`
	// Screen height.
	H OptInt `json:"h"`
	// Time spent on page. Only sent on unload.
	M OptInt `json:"m"`
}

// GetB returns the value of B.
func (s *EventHit) GetB() string {
	return s.B
}

// GetU returns the value of U.
func (s *EventHit) GetU() string {
	return s.U
}

// GetR returns the value of R.
func (s *EventHit) GetR() string {
	return s.R
}

// GetE returns the value of E.
func (s *EventHit) GetE() string {
	return s.E
}

// GetT returns the value of T.
func (s *EventHit) GetT() OptInt {
	return s.T
}

// GetD returns the value of D.
func (s *EventHit) GetD() OptString {
	return s.D
}

// GetW returns the value of W.
func (s *EventHit) GetW() OptInt {
	return s.W
}

// GetH returns the value of H.
func (s *EventHit) GetH() OptInt {
	return s.H
}

// GetM returns the value of M.
func (s *EventHit) GetM() OptInt {
	return s.M
}

// SetB sets the value of B.
func (s *EventHit) SetB(val string) {
	s.B = val
}

// SetU sets the value of U.
func (s *EventHit) SetU(val string) {
	s.U = val
}

// SetR sets the value of R.
func (s *EventHit) SetR(val string) {
	s.R = val
}

// SetE sets the value of E.
func (s *EventHit) SetE(val string) {
	s.E = val
}

// SetT sets the value of T.
func (s *EventHit) SetT(val OptInt) {
	s.T = val
}

// SetD sets the value of D.
func (s *EventHit) SetD(val OptString) {
	s.D = val
}

// SetW sets the value of W.
func (s *EventHit) SetW(val OptInt) {
	s.W = val
}

// SetH sets the value of H.
func (s *EventHit) SetH(val OptInt) {
	s.H = val
}

// SetM sets the value of M.
func (s *EventHit) SetM(val OptInt) {
	s.M = val
}

type ForbiddenError struct {
	Error ForbiddenErrorError `json:"error"`
}

// GetError returns the value of Error.
func (s *ForbiddenError) GetError() ForbiddenErrorError {
	return s.Error
}

// SetError sets the value of Error.
func (s *ForbiddenError) SetError(val ForbiddenErrorError) {
	s.Error = val
}

func (*ForbiddenError) postUserRes() {}

type ForbiddenErrorError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *ForbiddenErrorError) GetCode() int32 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *ForbiddenErrorError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *ForbiddenErrorError) SetCode(val int32) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *ForbiddenErrorError) SetMessage(val string) {
	s.Message = val
}

// GetEventPingOK is response for GetEventPing operation.
type GetEventPingOK struct {
	LastModified string
}

// GetLastModified returns the value of LastModified.
func (s *GetEventPingOK) GetLastModified() string {
	return s.LastModified
}

// SetLastModified sets the value of LastModified.
func (s *GetEventPingOK) SetLastModified(val string) {
	s.LastModified = val
}

func (*GetEventPingOK) getEventPingRes() {}

type GetWebsitesOKApplicationJSON []WebsiteGet

func (*GetWebsitesOKApplicationJSON) getWebsitesRes() {}

type InternalServerError struct {
	Error InternalServerErrorError `json:"error"`
}

// GetError returns the value of Error.
func (s *InternalServerError) GetError() InternalServerErrorError {
	return s.Error
}

// SetError sets the value of Error.
func (s *InternalServerError) SetError(val InternalServerErrorError) {
	s.Error = val
}

func (*InternalServerError) deleteUserRes()          {}
func (*InternalServerError) deleteWebsitesIDRes()    {}
func (*InternalServerError) getEventPingRes()        {}
func (*InternalServerError) getUserRes()             {}
func (*InternalServerError) getWebsiteIDSummaryRes() {}
func (*InternalServerError) getWebsitesIDActiveRes() {}
func (*InternalServerError) getWebsitesIDRes()       {}
func (*InternalServerError) getWebsitesRes()         {}
func (*InternalServerError) patchUserRes()           {}
func (*InternalServerError) patchWebsitesIDRes()     {}
func (*InternalServerError) postAuthLoginRes()       {}
func (*InternalServerError) postEventHitRes()        {}
func (*InternalServerError) postUserRes()            {}
func (*InternalServerError) postWebsitesRes()        {}

type InternalServerErrorError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *InternalServerErrorError) GetCode() int32 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *InternalServerErrorError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *InternalServerErrorError) SetCode(val int32) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *InternalServerErrorError) SetMessage(val string) {
	s.Message = val
}

type NotFoundError struct {
	Error NotFoundErrorError `json:"error"`
}

// GetError returns the value of Error.
func (s *NotFoundError) GetError() NotFoundErrorError {
	return s.Error
}

// SetError sets the value of Error.
func (s *NotFoundError) SetError(val NotFoundErrorError) {
	s.Error = val
}

func (*NotFoundError) deleteUserRes()          {}
func (*NotFoundError) deleteWebsitesIDRes()    {}
func (*NotFoundError) getUserRes()             {}
func (*NotFoundError) getWebsiteIDSummaryRes() {}
func (*NotFoundError) getWebsitesIDActiveRes() {}
func (*NotFoundError) getWebsitesIDRes()       {}
func (*NotFoundError) getWebsitesRes()         {}
func (*NotFoundError) patchUserRes()           {}
func (*NotFoundError) patchWebsitesIDRes()     {}

type NotFoundErrorError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *NotFoundErrorError) GetCode() int32 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *NotFoundErrorError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *NotFoundErrorError) SetCode(val int32) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *NotFoundErrorError) SetMessage(val string) {
	s.Message = val
}

// NewOptFloat32 returns new OptFloat32 with value set to v.
func NewOptFloat32(v float32) OptFloat32 {
	return OptFloat32{
		Value: v,
		Set:   true,
	}
}

// OptFloat32 is optional float32.
type OptFloat32 struct {
	Value float32
	Set   bool
}

// IsSet returns true if OptFloat32 was set.
func (o OptFloat32) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptFloat32) Reset() {
	var v float32
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptFloat32) SetTo(v float32) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptFloat32) Get() (v float32, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptFloat32) Or(d float32) float32 {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptInt returns new OptInt with value set to v.
func NewOptInt(v int) OptInt {
	return OptInt{
		Value: v,
		Set:   true,
	}
}

// OptInt is optional int.
type OptInt struct {
	Value int
	Set   bool
}

// IsSet returns true if OptInt was set.
func (o OptInt) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt) Reset() {
	var v int
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt) SetTo(v int) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt) Get() (v int, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt) Or(d int) int {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptUserCreateLanguage returns new OptUserCreateLanguage with value set to v.
func NewOptUserCreateLanguage(v UserCreateLanguage) OptUserCreateLanguage {
	return OptUserCreateLanguage{
		Value: v,
		Set:   true,
	}
}

// OptUserCreateLanguage is optional UserCreateLanguage.
type OptUserCreateLanguage struct {
	Value UserCreateLanguage
	Set   bool
}

// IsSet returns true if OptUserCreateLanguage was set.
func (o OptUserCreateLanguage) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptUserCreateLanguage) Reset() {
	var v UserCreateLanguage
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptUserCreateLanguage) SetTo(v UserCreateLanguage) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptUserCreateLanguage) Get() (v UserCreateLanguage, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptUserCreateLanguage) Or(d UserCreateLanguage) UserCreateLanguage {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptUserPatchLanguage returns new OptUserPatchLanguage with value set to v.
func NewOptUserPatchLanguage(v UserPatchLanguage) OptUserPatchLanguage {
	return OptUserPatchLanguage{
		Value: v,
		Set:   true,
	}
}

// OptUserPatchLanguage is optional UserPatchLanguage.
type OptUserPatchLanguage struct {
	Value UserPatchLanguage
	Set   bool
}

// IsSet returns true if OptUserPatchLanguage was set.
func (o OptUserPatchLanguage) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptUserPatchLanguage) Reset() {
	var v UserPatchLanguage
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptUserPatchLanguage) SetTo(v UserPatchLanguage) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptUserPatchLanguage) Get() (v UserPatchLanguage, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptUserPatchLanguage) Or(d UserPatchLanguage) UserPatchLanguage {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// PostAuthLoginOK is response for PostAuthLogin operation.
type PostAuthLoginOK struct {
	SetCookie string
}

// GetSetCookie returns the value of SetCookie.
func (s *PostAuthLoginOK) GetSetCookie() string {
	return s.SetCookie
}

// SetSetCookie sets the value of SetCookie.
func (s *PostAuthLoginOK) SetSetCookie(val string) {
	s.SetCookie = val
}

func (*PostAuthLoginOK) postAuthLoginRes() {}

// PostEventHitNotFound is response for PostEventHit operation.
type PostEventHitNotFound struct{}

func (*PostEventHitNotFound) postEventHitRes() {}

// PostEventHitOK is response for PostEventHit operation.
type PostEventHitOK struct{}

func (*PostEventHitOK) postEventHitRes() {}

// Return the number of active realtime users.
// Ref: #/components/schemas/StatsActive
type StatsActive struct {
	Visitors int `json:"visitors"`
}

// GetVisitors returns the value of Visitors.
func (s *StatsActive) GetVisitors() int {
	return s.Visitors
}

// SetVisitors sets the value of Visitors.
func (s *StatsActive) SetVisitors(val int) {
	s.Visitors = val
}

func (*StatsActive) getWebsitesIDActiveRes() {}

// Ref: #/components/schemas/StatsSummary
type StatsSummary struct {
	Uniques   OptInt     `json:"uniques"`
	Pageviews OptInt     `json:"pageviews"`
	Bounces   OptFloat32 `json:"bounces"`
	Duration  OptInt     `json:"duration"`
}

// GetUniques returns the value of Uniques.
func (s *StatsSummary) GetUniques() OptInt {
	return s.Uniques
}

// GetPageviews returns the value of Pageviews.
func (s *StatsSummary) GetPageviews() OptInt {
	return s.Pageviews
}

// GetBounces returns the value of Bounces.
func (s *StatsSummary) GetBounces() OptFloat32 {
	return s.Bounces
}

// GetDuration returns the value of Duration.
func (s *StatsSummary) GetDuration() OptInt {
	return s.Duration
}

// SetUniques sets the value of Uniques.
func (s *StatsSummary) SetUniques(val OptInt) {
	s.Uniques = val
}

// SetPageviews sets the value of Pageviews.
func (s *StatsSummary) SetPageviews(val OptInt) {
	s.Pageviews = val
}

// SetBounces sets the value of Bounces.
func (s *StatsSummary) SetBounces(val OptFloat32) {
	s.Bounces = val
}

// SetDuration sets the value of Duration.
func (s *StatsSummary) SetDuration(val OptInt) {
	s.Duration = val
}

func (*StatsSummary) getWebsiteIDSummaryRes() {}

type UnauthorisedError struct {
	Error UnauthorisedErrorError `json:"error"`
}

// GetError returns the value of Error.
func (s *UnauthorisedError) GetError() UnauthorisedErrorError {
	return s.Error
}

// SetError sets the value of Error.
func (s *UnauthorisedError) SetError(val UnauthorisedErrorError) {
	s.Error = val
}

func (*UnauthorisedError) deleteUserRes()          {}
func (*UnauthorisedError) deleteWebsitesIDRes()    {}
func (*UnauthorisedError) getUserRes()             {}
func (*UnauthorisedError) getWebsiteIDSummaryRes() {}
func (*UnauthorisedError) getWebsitesIDActiveRes() {}
func (*UnauthorisedError) getWebsitesIDRes()       {}
func (*UnauthorisedError) getWebsitesRes()         {}
func (*UnauthorisedError) patchUserRes()           {}
func (*UnauthorisedError) patchWebsitesIDRes()     {}
func (*UnauthorisedError) postAuthLoginRes()       {}
func (*UnauthorisedError) postUserRes()            {}
func (*UnauthorisedError) postWebsitesRes()        {}

type UnauthorisedErrorError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *UnauthorisedErrorError) GetCode() int32 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *UnauthorisedErrorError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *UnauthorisedErrorError) SetCode(val int32) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *UnauthorisedErrorError) SetMessage(val string) {
	s.Message = val
}

// Request body for creating a user.
// Ref: #/components/schemas/UserCreate
type UserCreate struct {
	Email    string                `json:"email"`
	Password string                `json:"password"`
	Language OptUserCreateLanguage `json:"language"`
}

// GetEmail returns the value of Email.
func (s *UserCreate) GetEmail() string {
	return s.Email
}

// GetPassword returns the value of Password.
func (s *UserCreate) GetPassword() string {
	return s.Password
}

// GetLanguage returns the value of Language.
func (s *UserCreate) GetLanguage() OptUserCreateLanguage {
	return s.Language
}

// SetEmail sets the value of Email.
func (s *UserCreate) SetEmail(val string) {
	s.Email = val
}

// SetPassword sets the value of Password.
func (s *UserCreate) SetPassword(val string) {
	s.Password = val
}

// SetLanguage sets the value of Language.
func (s *UserCreate) SetLanguage(val OptUserCreateLanguage) {
	s.Language = val
}

type UserCreateLanguage string

const (
	UserCreateLanguageEn UserCreateLanguage = "en"
)

// AllValues returns all UserCreateLanguage values.
func (UserCreateLanguage) AllValues() []UserCreateLanguage {
	return []UserCreateLanguage{
		UserCreateLanguageEn,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s UserCreateLanguage) MarshalText() ([]byte, error) {
	switch s {
	case UserCreateLanguageEn:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *UserCreateLanguage) UnmarshalText(data []byte) error {
	switch UserCreateLanguage(data) {
	case UserCreateLanguageEn:
		*s = UserCreateLanguageEn
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Response body for getting a user.
// Ref: #/components/schemas/UserGet
type UserGet struct {
	Email       string          `json:"email"`
	Language    UserGetLanguage `json:"language"`
	DateCreated int64           `json:"dateCreated"`
	DateUpdated int64           `json:"dateUpdated"`
}

// GetEmail returns the value of Email.
func (s *UserGet) GetEmail() string {
	return s.Email
}

// GetLanguage returns the value of Language.
func (s *UserGet) GetLanguage() UserGetLanguage {
	return s.Language
}

// GetDateCreated returns the value of DateCreated.
func (s *UserGet) GetDateCreated() int64 {
	return s.DateCreated
}

// GetDateUpdated returns the value of DateUpdated.
func (s *UserGet) GetDateUpdated() int64 {
	return s.DateUpdated
}

// SetEmail sets the value of Email.
func (s *UserGet) SetEmail(val string) {
	s.Email = val
}

// SetLanguage sets the value of Language.
func (s *UserGet) SetLanguage(val UserGetLanguage) {
	s.Language = val
}

// SetDateCreated sets the value of DateCreated.
func (s *UserGet) SetDateCreated(val int64) {
	s.DateCreated = val
}

// SetDateUpdated sets the value of DateUpdated.
func (s *UserGet) SetDateUpdated(val int64) {
	s.DateUpdated = val
}

func (*UserGet) getUserRes()   {}
func (*UserGet) patchUserRes() {}

// UserGetHeaders wraps UserGet with response headers.
type UserGetHeaders struct {
	SetCookie string
	Response  UserGet
}

// GetSetCookie returns the value of SetCookie.
func (s *UserGetHeaders) GetSetCookie() string {
	return s.SetCookie
}

// GetResponse returns the value of Response.
func (s *UserGetHeaders) GetResponse() UserGet {
	return s.Response
}

// SetSetCookie sets the value of SetCookie.
func (s *UserGetHeaders) SetSetCookie(val string) {
	s.SetCookie = val
}

// SetResponse sets the value of Response.
func (s *UserGetHeaders) SetResponse(val UserGet) {
	s.Response = val
}

func (*UserGetHeaders) postUserRes() {}

type UserGetLanguage string

const (
	UserGetLanguageEn UserGetLanguage = "en"
)

// AllValues returns all UserGetLanguage values.
func (UserGetLanguage) AllValues() []UserGetLanguage {
	return []UserGetLanguage{
		UserGetLanguageEn,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s UserGetLanguage) MarshalText() ([]byte, error) {
	switch s {
	case UserGetLanguageEn:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *UserGetLanguage) UnmarshalText(data []byte) error {
	switch UserGetLanguage(data) {
	case UserGetLanguageEn:
		*s = UserGetLanguageEn
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Request body for updating a user.
// Ref: #/components/schemas/UserPatch
type UserPatch struct {
	Email    OptString            `json:"email"`
	Password OptString            `json:"password"`
	Language OptUserPatchLanguage `json:"language"`
}

// GetEmail returns the value of Email.
func (s *UserPatch) GetEmail() OptString {
	return s.Email
}

// GetPassword returns the value of Password.
func (s *UserPatch) GetPassword() OptString {
	return s.Password
}

// GetLanguage returns the value of Language.
func (s *UserPatch) GetLanguage() OptUserPatchLanguage {
	return s.Language
}

// SetEmail sets the value of Email.
func (s *UserPatch) SetEmail(val OptString) {
	s.Email = val
}

// SetPassword sets the value of Password.
func (s *UserPatch) SetPassword(val OptString) {
	s.Password = val
}

// SetLanguage sets the value of Language.
func (s *UserPatch) SetLanguage(val OptUserPatchLanguage) {
	s.Language = val
}

type UserPatchLanguage string

const (
	UserPatchLanguageEn UserPatchLanguage = "en"
)

// AllValues returns all UserPatchLanguage values.
func (UserPatchLanguage) AllValues() []UserPatchLanguage {
	return []UserPatchLanguage{
		UserPatchLanguageEn,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s UserPatchLanguage) MarshalText() ([]byte, error) {
	switch s {
	case UserPatchLanguageEn:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *UserPatchLanguage) UnmarshalText(data []byte) error {
	switch UserPatchLanguage(data) {
	case UserPatchLanguageEn:
		*s = UserPatchLanguageEn
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Request body for creating a website.
// Ref: #/components/schemas/WebsiteCreate
type WebsiteCreate struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

// GetName returns the value of Name.
func (s *WebsiteCreate) GetName() string {
	return s.Name
}

// GetHostname returns the value of Hostname.
func (s *WebsiteCreate) GetHostname() string {
	return s.Hostname
}

// SetName sets the value of Name.
func (s *WebsiteCreate) SetName(val string) {
	s.Name = val
}

// SetHostname sets the value of Hostname.
func (s *WebsiteCreate) SetHostname(val string) {
	s.Hostname = val
}

// Response body for getting a website.
// Ref: #/components/schemas/WebsiteGet
type WebsiteGet struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

// GetName returns the value of Name.
func (s *WebsiteGet) GetName() string {
	return s.Name
}

// GetHostname returns the value of Hostname.
func (s *WebsiteGet) GetHostname() string {
	return s.Hostname
}

// SetName sets the value of Name.
func (s *WebsiteGet) SetName(val string) {
	s.Name = val
}

// SetHostname sets the value of Hostname.
func (s *WebsiteGet) SetHostname(val string) {
	s.Hostname = val
}

func (*WebsiteGet) getWebsitesIDRes()   {}
func (*WebsiteGet) patchWebsitesIDRes() {}
func (*WebsiteGet) postWebsitesRes()    {}

// Request body for updating a website.
// Ref: #/components/schemas/WebsitePatch
type WebsitePatch struct {
	Name     OptString `json:"name"`
	Hostname OptString `json:"hostname"`
}

// GetName returns the value of Name.
func (s *WebsitePatch) GetName() OptString {
	return s.Name
}

// GetHostname returns the value of Hostname.
func (s *WebsitePatch) GetHostname() OptString {
	return s.Hostname
}

// SetName sets the value of Name.
func (s *WebsitePatch) SetName(val OptString) {
	s.Name = val
}

// SetHostname sets the value of Hostname.
func (s *WebsitePatch) SetHostname(val OptString) {
	s.Hostname = val
}
